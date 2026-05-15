package instance

import (
	"context"
	"encoding/json"
	"net/http"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/security"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const WriteWait = 10 * time.Second
const PongWait = 60 * time.Second
const PingPeriod = (PongWait * 9) / 10
const MaxMessageSize = 128 * 1024

var wsOriginPolicy = security.NewOriginPolicyFromEnv()

var upgrader = websocket.Upgrader{
	EnableCompression: true,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}

		return wsOriginPolicy.HasAllowedOrigins() && wsOriginPolicy.IsAllowed(origin)
	},
}

type Client struct {
	ID       uuid.UUID          `json:"ID"`
	User     *game.User         `json:"user,omitempty"`
	conn     *websocket.Conn    `json:"-"`
	ctx      context.Context    `json:"-"`
	cancel   context.CancelFunc `json:"-"`
	instance *Instance          `json:"-"`
	inbox    chan Response      `json:"-"`
	once     sync.Once          `json:"-"`
}

func (c *Client) AttachUser(user *game.User) {
	c.User = user
	if user != nil {
		c.ID = user.ID
	}
}

func NewClient(instance *Instance) *Client {
	ctx, cancel := context.WithCancel(instance.ctx)
	return &Client{
		ID:       uuid.New(),
		ctx:      ctx,
		cancel:   cancel,
		instance: instance,

		inbox: make(chan Response, 100),
	}
}

func (c *Client) Connect(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	if c.User != nil {
		c.ID = c.User.ID
	}
	c.instance.Register <- c
	return nil
}

func (c *Client) Close() {
	c.once.Do(func() {
		if c.conn != nil {
			c.conn.Close()
		}
		c.instance.Unregister <- c
		c.cancel()
	})
}

func (c *Client) WriteResponse(response *Response) error {
	json, err := json.Marshal(response)
	if err != nil {
		return err
	}
	if err = c.conn.WriteMessage(websocket.TextMessage, json); err != nil {
		return err
	}

	return nil
}

func (c *Client) TryWriteResponse(response Response) bool {
	select {
	case c.inbox <- response:
		return true
	default:
		return false
	}
}

func (c *Client) listenForRequest(request *Request) error {
	_, raw, err := c.conn.ReadMessage()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, request); err != nil {
		return err
	}
	return nil
}

func (c *Client) listenIn() {
	defer c.Close()

	pongHandler := func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	}
	c.conn.SetReadLimit(MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(pongHandler)

	for {
		var request Request
		if err := c.listenForRequest(&request); err != nil {
			break
		}

		select {
		case c.instance.ReadRequest <- request:
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) listenOut() {
	clock := time.NewTicker(PingPeriod)
	defer func() {
		clock.Stop()
		c.Close()
	}()

	for {
		select {
		case response := <-c.inbox:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.WriteResponse(&response); err != nil {
				return
			}
		case <-clock.C:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) Run() {
	go c.listenIn()
	go c.listenOut()
}

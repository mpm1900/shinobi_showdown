package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"shinobi_showdown/internal/auth"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/instance"

	"github.com/google/uuid"
)

type InstancesHandler struct {
	mux         *http.ServeMux
	instances   map[uuid.UUID]*instance.Instance
	instancesMu sync.RWMutex
}

func NewInstancesHandler(ctx context.Context) *InstancesHandler {
	handler := &InstancesHandler{
		mux:       http.NewServeMux(),
		instances: map[uuid.UUID]*instance.Instance{},
	}

	handler.mux.HandleFunc("GET /{instanceID}/connect", handler.handleGameConnection(ctx))
	return handler
}

func (ih *InstancesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ih.mux.ServeHTTP(w, r)
}

func (ih *InstancesHandler) NewInstance(instanceID uuid.UUID, ctx context.Context) *instance.Instance {
	instance := instance.NewInstance(ctx, instanceID, ih.RemoveInstance)
	ih.instancesMu.Lock()
	ih.instances[instance.ID] = instance
	ih.instancesMu.Unlock()
	go instance.Run()

	return instance
}

func (ih *InstancesHandler) GetOrCreateInstance(instanceID uuid.UUID, ctx context.Context) *instance.Instance {
	ih.instancesMu.Lock()
	defer ih.instancesMu.Unlock()

	if existing, ok := ih.instances[instanceID]; ok && existing != nil {
		return existing
	}

	created := instance.NewInstance(ctx, instanceID, ih.RemoveInstance)
	ih.instances[created.ID] = created
	go created.Run()
	return created
}

func (ih *InstancesHandler) RemoveInstance(instanceID uuid.UUID) {
	ih.instancesMu.Lock()
	delete(ih.instances, instanceID)
	ih.instancesMu.Unlock()
}

func (ih *InstancesHandler) GetInstance(instanceID uuid.UUID) (*instance.Instance, bool) {
	ih.instancesMu.RLock()
	defer ih.instancesMu.RUnlock()
	instance, ok := ih.instances[instanceID]

	return instance, ok
}

func (ih *InstancesHandler) GetAllInstances() []instance.InstanceJSON {
	ih.instancesMu.RLock()
	defer ih.instancesMu.RUnlock()
	games := make([]instance.InstanceJSON, 0, len(ih.instances))
	for _, instance := range ih.instances {
		if instance != nil {
			games = append(games, instance.ToJSON())
		}
	}
	return games
}

func (ih *InstancesHandler) HandleGetGames(w http.ResponseWriter, r *http.Request) {
	games := ih.GetAllInstances()
	if len(games) == 0 {
		games = make([]instance.InstanceJSON, 0)
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(games)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ih *InstancesHandler) handleGameConnection(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.AuthenticatedUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		instanceID, err := uuid.Parse(r.PathValue("instanceID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		i := ih.GetOrCreateInstance(instanceID, ctx)

		client := instance.NewClient(i)
		client.AttachUser(&game.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
		if err := client.Connect(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		client.Run()
	}
}

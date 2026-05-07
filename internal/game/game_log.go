package game

import (
	"strings"

	"github.com/google/uuid"
)

type GameLog struct {
	ID      uuid.UUID `json:"ID"`
	Text    string    `json:"text"`
	Context Context   `json:"context"`
}

const GameLogIndent = "→"

func NewLog(text string) GameLog {
	return GameLog{
		ID:      uuid.New(),
		Text:    text,
		Context: NewContext(),
	}
}
func NewLogContext(text string, context Context) GameLog {
	return GameLog{
		ID:      uuid.New(),
		Text:    text,
		Context: context,
	}
}

func MakeGameLog(text string, context Context, depth int) GameLog {
	var sb strings.Builder
	for range depth {
		sb.WriteString(GameLogIndent)
	}
	if depth > 0 {
		sb.WriteByte(' ')
	}
	sb.WriteString(text)
	log := NewLogContext(sb.String(), context)
	return log
}

package modifiers

import (
	"fmt"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var addNatureID = uuid.MustParse("0cac062a-b0be-5ad3-8f82-e51433b66bed")

func AddNature(nature game.NatureSet, duration int) game.Modifier {
	return game.Modifier{
		ID:       addNatureID,
		Name:     fmt.Sprintf("Add Nature: %s", nature),
		Show:     false,
		Duration: duration,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				nil,
				game.MutPriorityDefault,
				game.ComposeAF(game.SourceFilter, game.ActiveFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					a.Natures[nature] = game.NATURES[nature]
					return a
				},
			),
		},
	}
}
func SetNature(nature game.NatureSet, duration int) game.Modifier {
	return game.Modifier{
		ID:       addNatureID,
		Name:     fmt.Sprintf("Set Nature: %s", nature),
		Show:     false,
		Duration: duration,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				nil,
				game.MutPriorityPostStagedStats,
				game.ComposeAF(game.SourceFilter, game.ActiveFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					a.Natures = make(map[game.NatureSet][]game.Nature)
					a.Natures[nature] = game.NATURES[nature]
					return a
				},
			),
		},
	}
}

var removeNatureID = uuid.MustParse("be5ce975-cda5-59e2-b03c-7bdcb3db2e33")

func RemoveNature(nature game.NatureSet) game.Modifier {
	return game.Modifier{
		ID:       removeNatureID,
		Name:     fmt.Sprintf("Remove Nature: %s", nature),
		Show:     false,
		Duration: 0,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				nil,
				game.MutPriorityDefault,
				game.ComposeAF(game.SourceFilter, game.ActiveFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					delete(a.Natures, nature)
					return a
				},
			),
		},
	}
}

package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Mifune = game.ActorDef{
	ActorID:      uuid.MustParse("9a87963d-05f0-42b0-89cc-84f04b367312"),
	SpriteURL:    "/sprites/mifune_64.png",
	Name:         "Mifune",
	Affiliations: []string{},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        135,
		game.StatDefense:       120,
		game.StatChakraAttack:  60,
		game.StatChakraDefense: 85,
		game.StatSpeed:         50,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsTai,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.IaiSlash.ID,
		actions.IaiBlock.ID,
		actions.IaiTrueStrike.ID,
		actions.Disable.ID,
		actions.SwordsStance.ID,
		actions.NoEscape.ID,
	}, GlobalActionsExcept(actions.Substitution.ID)...),
}

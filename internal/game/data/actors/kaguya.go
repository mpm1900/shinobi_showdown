package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Kaguya = game.ActorDef{
	ActorID:      uuid.MustParse("47f76e4c-edaa-5cd7-ba31-8a1803213d9e"),
	Name:         "Kaguya Ōtsutsuki",
	Affiliations: []string{},

	Stats: map[game.ActorStat]int{
		game.StatHP:            255,
		game.StatStamina:       255,
		game.StatAttack:        135,
		game.StatDefense:       255,
		game.StatChakraAttack:  135,
		game.StatChakraDefense: 255,
		game.StatSpeed:         125,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Chidori.ID,
		actions.DragonStance.ID,
	}, GlobalActions...),
}

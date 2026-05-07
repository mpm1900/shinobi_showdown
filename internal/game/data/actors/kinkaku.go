package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Kinkaku = game.ActorDef{
	ActorID:      uuid.MustParse("b94e4ed5-64dc-4265-9b81-572c023082e5"),
	SpriteURL:    "/sprites/kinkaku_64.png",
	Name:         "Kinkaku",
	Affiliations: []string{game.AffKumo},

	Stats: map[game.ActorStat]int{
		game.StatHP:            91,
		game.StatStamina:       130,
		game.StatAttack:        134,
		game.StatDefense:       95,
		game.StatChakraAttack:  100,
		game.StatChakraDefense: 100,
		game.StatSpeed:         80,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsLightning,
		game.NsWater,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CollidingWave.ID,
		actions.WaterWall.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
	}, GlobalActions...),
}

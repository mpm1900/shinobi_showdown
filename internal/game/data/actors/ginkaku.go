package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Ginkaku = game.ActorDef{
	ActorID:      uuid.MustParse("490b0f40-b921-44eb-9cc4-d6f375fd20c9"),
	SpriteURL:    "/sprites/ginkaku_64.png",
	Name:         "Ginkaku",
	Affiliations: []string{game.AffKumo},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       130,
		game.StatAttack:        134,
		game.StatDefense:       110,
		game.StatChakraAttack:  95,
		game.StatChakraDefense: 100,
		game.StatSpeed:         61,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsWind,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CollidingWave.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.VacuumBlast.ID,
	}, GlobalActions...),
}

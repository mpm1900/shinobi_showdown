package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var KillerBee = game.ActorDef{
	ActorID:      uuid.MustParse("e92a1675-7ea0-52c7-a43c-f2bb2068e365"),
	SpriteURL:    "/sprites/bee_64.png",
	Name:         "Killer Bee",
	Affiliations: []string{game.AffKumo},

	Stats: map[game.ActorStat]int{
		game.StatHP:            74,
		game.StatStamina:       100,
		game.StatAttack:        130,
		game.StatDefense:       100,
		game.StatChakraAttack:  110,
		game.StatChakraDefense: 70,
		game.StatSpeed:         116,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsLightning,
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.ChakraArmor,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CollidingWave.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

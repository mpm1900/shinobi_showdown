package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Haku = game.ActorDef{
	ActorID:      uuid.MustParse("af4a0622-69f9-4e05-a3c1-5579c4d8d3e2"),
	SpriteURL:    "/sprites/haku_64.png",
	Name:         "Haku",
	Affiliations: []string{game.AffKuri},

	Stats: map[game.ActorStat]int{
		game.StatHP:            55,
		game.StatStamina:       100,
		game.StatAttack:        55,
		game.StatDefense:       75,
		game.StatChakraAttack:  135,
		game.StatChakraDefense: 120,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
		game.NsIce,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Haze.ID,
		actions.CollidingWave.ID,
		actions.HiddenMist.ID,
		actions.CreateRain.ID,
		actions.OneThousandNeedles.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

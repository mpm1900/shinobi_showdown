package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Yamato = game.ActorDef{
	ActorID:      uuid.MustParse("aec1edee-f092-5422-a669-fe8eb73b556f"),
	SpriteURL:    "/sprites/yamato_64.png",
	Name:         "Yamato",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            80,
		game.StatStamina:       101,
		game.StatAttack:        80,
		game.StatDefense:       100,
		game.StatChakraAttack:  95,
		game.StatChakraDefense: 135,
		game.StatSpeed:         70,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
		game.NsWood,
	}),
	Abilities: []game.Modifier{
		modifiers.NeutralizingChakra,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Recover.ID,
		actions.WhirlwindKick.ID,
		actions.GreatTreeSpear.ID,
		actions.SummonAlly.ID,
		actions.Yawn.ID,
		actions.Taunt.ID,
		actions.HealthSplit.ID,
		actions.GreatWaterfall.ID,
		actions.MudWall.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

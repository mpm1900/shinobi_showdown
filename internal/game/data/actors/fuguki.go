package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Fuguki = game.ActorDef{
	ActorID:      uuid.MustParse("e9c58fd3-ca37-4af5-8715-da628a1343f7"),
	SpriteURL:    "/sprites/fuguki_64.png",
	Name:         "Fuguki Suikazan",
	Affiliations: []string{game.AffKuri},

	Stats: map[game.ActorStat]int{
		game.StatHP:            130,
		game.StatStamina:       100,
		game.StatAttack:        140,
		game.StatDefense:       105,
		game.StatChakraAttack:  45,
		game.StatChakraDefense: 70,
		game.StatSpeed:         50,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
		game.NsYang,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.KebariSenbon.ID,
		actions.DragonStance.ID,
		actions.CollidingWave.ID,
		actions.HiddenMist.ID,
		actions.CreateRain.ID,
		actions.WaterDragon.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

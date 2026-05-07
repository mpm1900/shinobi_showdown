package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Kimimaro = game.ActorDef{
	ActorID:   uuid.MustParse("ced6b910-2207-5982-9202-fc41804c5071"),
	SpriteURL: "/sprites/kimimaro_64.png",
	Name:      "Kimimaro Kaguya",
	Affiliations: []string{
		game.AffKonoha,
	},

	Stats: map[game.ActorStat]int{
		game.StatHP:            76,
		game.StatStamina:       70,
		game.StatAttack:        147,
		game.StatDefense:       90,
		game.StatChakraAttack:  60,
		game.StatChakraDefense: 70,
		game.StatSpeed:         97,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsTai,
		game.NsEarth,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Taunt.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.CamelliaDance.ID,
		actions.ClematisDanceFlower.ID,
		actions.BodyFlicker.ID,
		actions.Earthquake.ID,
	}, GlobalActions...),
}

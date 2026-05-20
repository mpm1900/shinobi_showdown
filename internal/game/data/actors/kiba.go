package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Kiba = game.ActorDef{
	ActorID:      uuid.MustParse("eec3b447-b6f6-4632-8262-ca0bcb38c2ce"),
	SpriteURL:    "/sprites/kiba_64.png",
	Name:         "Kiba Inuzuka",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            60,
		game.StatStamina:       100,
		game.StatAttack:        80,
		game.StatDefense:       110,
		game.StatChakraAttack:  50,
		game.StatChakraDefense: 80,
		game.StatSpeed:         45,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsYang,
	}),
	Abilities: []game.Modifier{
		modifiers.HardHeaded,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Expansion.ID,
		actions.WhirlwindKick.ID,
		actions.RockFist.ID,
		actions.Earthquake.ID,
		actions.FangOverFang.ID,
	}, GlobalActions...),
	DefaultModifiers: []game.Modifier{
		modifiers.MansBestFriend,
	},
}

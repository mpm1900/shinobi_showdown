package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Ino = game.ActorDef{
	ActorID:   uuid.MustParse("62a307d4-006a-4a14-9122-784909a37ec1"),
	SpriteURL: "/sprites/ino_64.png",
	Name:      "Ino Yamanaka",
	Affiliations: []string{
		game.AffKonoha,
	},

	Stats: map[game.ActorStat]int{
		game.StatHP:            90,
		game.StatStamina:       100,
		game.StatAttack:        60,
		game.StatDefense:       80,
		game.StatChakraAttack:  90,
		game.StatChakraDefense: 110,
		game.StatSpeed:         60,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsYin,
	}),
	Abilities: []game.Modifier{
		modifiers.MindLink,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Barrier.ID,
		actions.Command.ID,
		actions.Coercion.ID,
		actions.Disable.ID,
		actions.SageMode.ID,
		actions.MudWall.ID,
		actions.TradeOffer.ID,
		actions.Yawn.ID,
		actions.WaterWall.ID,
		actions.Taunt.ID,
	}, GlobalActions...),
}

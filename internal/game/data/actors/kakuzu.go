package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Kakuzu = game.ActorDef{
	ActorID:      uuid.MustParse("9a273c6c-d268-5d54-9667-4c264d5192d8"),
	SpriteURL:    "/sprites/kakuzu_64.png",
	Name:         "Kakuzu",
	Affiliations: []string{game.AffAkatsuki, game.AffTaki},

	Stats: map[game.ActorStat]int{
		game.StatHP:            80,
		game.StatStamina:       100,
		game.StatAttack:        110,
		game.StatDefense:       125,
		game.StatChakraAttack:  125,
		game.StatChakraDefense: 80,
		game.StatSpeed:         80,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
	}),
	Abilities: []game.Modifier{
		modifiers.ConsumeChakra,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.IronSkin.ID,
		actions.IronBody.ID,
		actions.DragonStance.ID,
		actions.SearingMigraine.ID,
		actions.FalseDarkness.ID,
		actions.PressureDamage.ID,
		actions.WhirlwindKick.ID,
		actions.Earthquake.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

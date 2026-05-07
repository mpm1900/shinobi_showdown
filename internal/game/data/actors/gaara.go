package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Gaara = game.ActorDef{
	ActorID:   uuid.MustParse("b32fadc0-a6c4-4cf3-8f34-f91f67eb1b77"),
	SpriteURL: "/sprites/gaara_64.png",
	Name:      "Gaara",
	Affiliations: []string{
		game.AffSun,
	},

	Stats: map[game.ActorStat]int{
		game.StatHP:            90,
		game.StatStamina:       100,
		game.StatAttack:        66,
		game.StatDefense:       116,
		game.StatChakraAttack:  89,
		game.StatChakraDefense: 126,
		game.StatSpeed:         33,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsMagnet,
	}),
	Abilities: []game.Modifier{
		modifiers.SandAura,
		modifiers.AllyGuard,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.EarthDomePrison.ID,
		actions.Earthquake.ID,
	}, GlobalActions...),
}

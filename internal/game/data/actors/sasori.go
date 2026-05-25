package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Sasori = game.ActorDef{
	ActorID:   uuid.MustParse("21ce56e6-9463-434f-ad1d-cd6fcc72a46e"),
	SpriteURL: "/sprites/sasori_64.png",
	Name:      "Sasori",
	Affiliations: []string{
		game.AffAkatsuki,
		game.AffSun,
	},
	Stats: map[game.ActorStat]int{
		game.StatHP:            80,
		game.StatStamina:       100,
		game.StatAttack:        130,
		game.StatDefense:       60,
		game.StatChakraAttack:  40,
		game.StatChakraDefense: 80,
		game.StatSpeed:         120,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWind,
		game.NsEarth,
	}),
	Abilities: []game.Modifier{
		modifiers.PoisonInfused,
	},
	DefaultModifiers: []game.Modifier{
		modifiers.Immune,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.SnakeStrike.ID,
		actions.Flash.ID,
		actions.Taunt.ID,
		actions.RetreatingStrike.ID,
		actions.PuppetAssault.ID,
	}, GlobalActions...),
}

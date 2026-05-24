package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Jiraiya = game.ActorDef{
	ActorID:      uuid.MustParse("fbf88375-8f7c-5380-97b9-2f2748581e76"),
	SpriteURL:    "/sprites/jiraiya_64.png",
	Name:         "Jiraiya",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            91,
		game.StatStamina:       100,
		game.StatAttack:        90,
		game.StatDefense:       106,
		game.StatChakraAttack:  130,
		game.StatChakraDefense: 106,
		game.StatSpeed:         77,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsEarth,
	}),
	Abilities: []game.Modifier{
		modifiers.SagesBlessing,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Rasengan.ID,
		actions.GiantRasengan.ID,
		actions.Haze.ID,
		actions.SageMode.ID,
		actions.FlameBullet.ID,
		actions.GreatFireball.ID,
		actions.SummonGamabunta.ID,
		actions.KebariSenbon.ID,
		actions.Yawn.ID,
		actions.Earthquake.ID,
		actions.DarkSwamp.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

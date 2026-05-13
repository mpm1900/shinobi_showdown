package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Kabuto = game.ActorDef{
	ActorID:      uuid.MustParse("a94c7c47-4211-53e4-8821-49a12e308840"),
	SpriteURL:    "/sprites/kabuto_64.png",
	Name:         "Kabuto Yakushi",
	Affiliations: []string{game.AffOto},

	Stats: map[game.ActorStat]int{
		game.StatHP:            55,
		game.StatStamina:       100,
		game.StatAttack:        20,
		game.StatDefense:       35,
		game.StatChakraAttack:  20,
		game.StatChakraDefense: 45,
		game.StatSpeed:         75,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYin,
		game.NsYang,
	}),
	Abilities: []game.Modifier{
		modifiers.AllyGuard,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Haze.ID,
		actions.SageMode.ID,
		actions.Tailwind.ID,
		actions.Redirect.ID,
		actions.TempleOfNirvana.ID,
		actions.PerishSong.ID,
		actions.ThirtyTwoPalms.ID,
		actions.Flash.ID,
		actions.Barrier.ID,
		actions.MudWall.ID,
		actions.WaterWall.ID,
		actions.HealthSplit.ID,
	}, GlobalActions...),
}

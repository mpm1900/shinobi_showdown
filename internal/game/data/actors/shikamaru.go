package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Shikamaru = game.ActorDef{
	ActorID:      uuid.MustParse("4d0ad337-2685-5ac1-be21-d62a2e23f678"),
	SpriteURL:    "/sprites/shikamaru_64.png",
	Name:         "Shikamaru Nara",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            70,
		game.StatStamina:       100,
		game.StatAttack:        55,
		game.StatDefense:       65,
		game.StatChakraAttack:  95,
		game.StatChakraDefense: 105,
		game.StatSpeed:         85,
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
		modifiers.FastThinking,
		modifiers.ShadowCage,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.PatternBreak.ID,
		actions.Haze.ID,
		actions.Tailwind.ID,
		actions.Redirect.ID,
		actions.ShadowPossession.ID,
		actions.MudWall.ID,
		actions.WaterWall.ID,
		actions.Yawn.ID,
		actions.TradeOffer.ID,
	}, GlobalActions...),
}

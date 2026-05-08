package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Obito = game.ActorDef{
	ActorID:      uuid.MustParse("e01398df-0a65-45a4-a4d6-878af3fa9d4c"),
	SpriteURL:    "/sprites/obito_64.png",
	Name:         "Masked Obito",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffAkatsuki, game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        135,
		game.StatDefense:       120,
		game.StatChakraAttack:  60,
		game.StatChakraDefense: 85,
		game.StatSpeed:         50,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYin,
	}),
	Abilities: []game.Modifier{
		modifiers.Focused,
		modifiers.NatureSpecialist,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.NoEscape.ID,
		actions.Kamui.ID,
		actions.KamuiCounter.ID,
		actions.KamuiSlash.ID,
		actions.SageMode.ID,
		actions.PatternBreak.ID,
		actions.KusariChains.ID,
		actions.SwordsStance.ID,
		actions.DesperateStrike.ID,
	}, GlobalActionsExcept(actions.BodyReplacement.ID)...),
}

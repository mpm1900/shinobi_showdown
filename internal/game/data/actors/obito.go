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
		game.StatHP:            120,
		game.StatStamina:       100,
		game.StatAttack:        112,
		game.StatDefense:       65,
		game.StatChakraAttack:  80,
		game.StatChakraDefense: 75,
		game.StatSpeed:         78,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYin,
	}),
	Abilities: []game.Modifier{
		modifiers.NatureSpecialist,
		modifiers.Focused,
		modifiers.Defiant,
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
		actions.WillOfTheFallen.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActionsExcept(actions.Substitution.ID)...),
}

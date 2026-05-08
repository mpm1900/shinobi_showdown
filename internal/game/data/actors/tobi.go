package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Tobi = game.ActorDef{
	ActorID:      uuid.MustParse("e02084f2-2545-468d-8395-b7b5c84fc63d"),
	SpriteURL:    "/sprites/tobi_64.png",
	Name:         "Tobi",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffAkatsuki, game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            95,
		game.StatStamina:       100,
		game.StatAttack:        115,
		game.StatDefense:       90,
		game.StatChakraAttack:  80,
		game.StatChakraDefense: 90,
		game.StatSpeed:         60,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsYin,
	}),
	Abilities: []game.Modifier{
		modifiers.MoldBreaker,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Kamui.ID,
		actions.KamuiEscape.ID,
		actions.KamuiCounter.ID,
		actions.KamuiSlash.ID,
		actions.GreatFireball.ID,
		actions.Flash.ID,
		actions.RetreatingStrike.ID,
		actions.TradeOffer.ID,
	}, GlobalActionsExcept(actions.BodyReplacement.ID)...),
}

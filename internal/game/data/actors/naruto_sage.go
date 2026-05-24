package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var NarutoSage = game.ActorDef{
	ActorID:      uuid.MustParse("a0dd319e-aaa0-49c0-9ba0-eba14cba39d1"),
	SpriteURL:    "/sprites/naruto_sage_64.png",
	Name:         "Naruto (Sage)",
	Clan:         game.ClanUzumaki,
	Affiliations: []string{game.AffKonoha},
	Restricted:   true,

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       130,
		game.StatAttack:        100,
		game.StatDefense:       120,
		game.StatChakraAttack:  150,
		game.StatChakraDefense: 100,
		game.StatSpeed:         80,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsSage,
		game.NsWind,
	}),
	Abilities: []game.Modifier{
		modifiers.PriorityFailure,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Rasengan.ID,
		actions.GiantRasengan.ID,
		actions.PowerBoost.ID,
		actions.SageMode.ID,
		actions.WhirlwindKick.ID,
		actions.SummonGamabunta.ID,
		actions.Rasenshuriken.ID,
		actions.VacuumBlast.ID,
		actions.TradeOffer.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

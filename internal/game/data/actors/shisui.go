package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Shisui = game.ActorDef{
	ActorID:      uuid.MustParse("3b9a402f-9ac1-5e47-8d56-37ec1ed287ff"),
	SpriteURL:    "/sprites/shisui_64.png",
	Name:         "Shisui Uchiha",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        100,
		game.StatDefense:       100,
		game.StatChakraAttack:  100,
		game.StatChakraDefense: 100,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsWind,
	}),
	Abilities: []game.Modifier{
		modifiers.Guts,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.NoEscape.ID,
		actions.Taunt.ID,
		actions.SharinganGlare.ID,
		actions.Redirect.ID,
		actions.PatternBreak.ID,
		actions.Chidori.ID,
		actions.DragonStance.ID,
		actions.GreatFireball.ID,
		actions.Firestorm.ID,
		actions.BodyFlicker.ID,
		actions.Recover.ID,
		actions.TradeOffer.ID,
	}, GlobalActions...),
}

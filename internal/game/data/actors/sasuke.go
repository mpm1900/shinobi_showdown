package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Sasuke = game.ActorDef{
	ActorID:      uuid.MustParse("2c7025a3-77c8-57b7-9120-e99b473f669f"),
	Name:         "Sasuke Uchiha",
	SpriteURL:    "/sprites/sasuke_64.png",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffAkatsuki, game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            80,
		game.StatStamina:       100,
		game.StatAttack:        95,
		game.StatDefense:       70,
		game.StatChakraAttack:  120,
		game.StatChakraDefense: 85,
		game.StatSpeed:         120,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsLightning,
	}),
	Abilities: []game.Modifier{
		modifiers.InnerFocus,
		modifiers.Unburden,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.NoEscape.ID,
		actions.Chidori.ID,
		actions.ChidoriSpear.ID,
		actions.ChidoriStream.ID,
		actions.Kirin.ID,
		actions.DragonStance.ID,
		actions.DragonFire.ID,
		actions.GreatFireball.ID,
		actions.Amaterasu.ID,
	}, GlobalActions...),
}

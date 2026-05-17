package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var SasukeEMS = game.ActorDef{
	ActorID:      uuid.MustParse("d5b318c0-5cb7-4bae-9c24-1fc5d3ab9214"),
	Name:         "Sasuke (EMS)",
	SpriteURL:    "/sprites/sasuke_ems_64.png",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffAkatsuki, game.AffKonoha},
	Restricted:   true,

	Stats: map[game.ActorStat]int{
		game.StatHP:            98,
		game.StatStamina:       100,
		game.StatAttack:        110,
		game.StatDefense:       100,
		game.StatChakraAttack:  130,
		game.StatChakraDefense: 90,
		game.StatSpeed:         142,
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
		actions.GreatFireball.ID,
		actions.Amaterasu.ID,
	}, GlobalActions...),
}

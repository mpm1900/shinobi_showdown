package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Itachi = game.ActorDef{
	ActorID:      uuid.MustParse("52713900-6202-50f5-9b3d-23dc408e3c63"),
	SpriteURL:    "/sprites/itachi_64.png",
	Name:         "Itachi Uchiha",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffAkatsuki, game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            55,
		game.StatStamina:       100,
		game.StatAttack:        55,
		game.StatDefense:       55,
		game.StatChakraAttack:  135,
		game.StatChakraDefense: 135,
		game.StatSpeed:         135,
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
		modifiers.Intimidate,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.NoEscape.ID,
		actions.Disable.ID,
		actions.PatternBreak.ID,
		actions.Flash.ID,
		actions.MirageCrow.ID,
		actions.GreatFireball.ID,
		actions.PhoenixFlower.ID,
		actions.PunishingFire.ID,
		actions.Amaterasu.ID,
		actions.Coercion.ID,
		actions.InstilFear.ID,
		actions.PerishSong.ID,
		actions.Taunt.ID,
		actions.DisarmingStrike.ID,
		actions.SharinganGlare.ID,
		actions.Caltrops.ID,
	}, GlobalActions...),
}

package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Madara = game.ActorDef{
	ActorID:      uuid.MustParse("87096d92-2694-5262-bfa3-59f23600be6b"),
	SpriteURL:    "/sprites/madara_64.png",
	Name:         "Madara Uchiha",
	Clan:         game.ClanUchiha,
	Affiliations: []string{game.AffAkatsuki, game.AffKonoha},
	Restricted:   true,

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       110,
		game.StatAttack:        110,
		game.StatDefense:       80,
		game.StatChakraAttack:  165,
		game.StatChakraDefense: 105,
		game.StatSpeed:         130,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsYinYang,
	}),
	Abilities: []game.Modifier{
		modifiers.WillOfFire,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.NoEscape.ID,
		actions.Flash.ID,
		actions.Taunt.ID,
		actions.DragonStance.ID,
		actions.GreatFireball.ID,
		actions.WhirlwindKick.ID,
		actions.GreatFireAnnihilation.ID,
		actions.MajesticFlame.ID,
		actions.Firestorm.ID,
		actions.SharinganGlare.ID,
		actions.ChibakuTensei.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Shino = game.ActorDef{
	ActorID:      uuid.MustParse("4bc617c5-2df1-4ae4-8e6e-dd086d017104"),
	SpriteURL:    "/sprites/shino_64.png",
	Name:         "Shino Aburame",
	Affiliations: []string{game.AffKonoha},
	Clan:         game.ClanAburame,

	Stats: map[game.ActorStat]int{
		game.StatHP:            80,
		game.StatStamina:       100,
		game.StatAttack:        80,
		game.StatDefense:       80,
		game.StatChakraAttack:  100,
		game.StatChakraDefense: 80,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsYang,
	}),
	Abilities: []game.Modifier{},
	DefaultModifiers: []game.Modifier{
		modifiers.Immune,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
	}, GlobalActions...),
}

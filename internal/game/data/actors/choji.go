package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Choji = game.ActorDef{
	ActorID:      uuid.MustParse("00f5a224-63ce-5cd7-87d9-5623dca59e92"),
	SpriteURL:    "/sprites/choji_64.png",
	Name:         "Choji Akimichi",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            150,
		game.StatStamina:       100,
		game.StatAttack:        100,
		game.StatDefense:       115,
		game.StatChakraAttack:  65,
		game.StatChakraDefense: 65,
		game.StatSpeed:         35,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsYang,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Expansion.ID,
		actions.WhirlwindKick.ID,
		actions.HumanBoulder.ID,
		actions.HeavyPunch.ID,
		actions.ChilliPill.ID,
		actions.RockFist.ID,
		actions.Earthquake.ID,
	}, GlobalActions...),
}

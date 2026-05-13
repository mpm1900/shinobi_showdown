package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Onoki = game.ActorDef{
	ActorID:   uuid.MustParse("94b7b9d0-bf5d-5b4f-b29c-61779635177e"),
	SpriteURL: "/sprites/onoki_64.png",
	Name:      "Ōnoki",
	Affiliations: []string{
		game.AffIwa,
	},

	Stats: map[game.ActorStat]int{
		game.StatHP:            56,
		game.StatStamina:       100,
		game.StatAttack:        80,
		game.StatDefense:       114,
		game.StatChakraAttack:  124,
		game.StatChakraDefense: 80,
		game.StatSpeed:         116,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsParticle,
	}),
	Abilities: []game.Modifier{},
	DefaultModifiers: []game.Modifier{
		modifiers.Flying,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.RockFist.ID,
		actions.MudWall.ID,
		actions.Barrier.ID,
		actions.DisarmingStrike.ID,
		actions.AtomicDismantlingCharge.ID,
		actions.Earthquake.ID,
	}, GlobalActions...),
}

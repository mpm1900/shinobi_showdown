package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Mu = game.ActorDef{
	ActorID:   uuid.MustParse("1198d088-4f2a-46e1-b320-0d0361f3dc33"),
	SpriteURL: "/sprites/mu_64.png",
	Name:      "Mū",
	Affiliations: []string{
		game.AffIwa,
	},

	Stats: map[game.ActorStat]int{
		game.StatHP:            89,
		game.StatStamina:       100,
		game.StatAttack:        115,
		game.StatDefense:       70,
		game.StatChakraAttack:  135,
		game.StatChakraDefense: 90,
		game.StatSpeed:         101,
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
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

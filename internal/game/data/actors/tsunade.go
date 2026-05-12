package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Tsunade = game.ActorDef{
	ActorID:      uuid.MustParse("aad15064-9fcb-4bdb-a66b-36e56694316e"),
	SpriteURL:    "/sprites/tsunade_64.png",
	Name:         "Tsunade Senju",
	Clan:         game.ClanSenju,
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            110,
		game.StatStamina:       100,
		game.StatAttack:        110,
		game.StatDefense:       110,
		game.StatChakraAttack:  75,
		game.StatChakraDefense: 100,
		game.StatSpeed:         75,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
		game.NsYang,
	}),
	Abilities: []game.Modifier{
		modifiers.Regeneration,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Haze.ID,
		actions.HeavyPunch.ID,
		actions.SageMode.ID,
		actions.TeamHeal.ID,
		actions.OneHundredHealingsStatus.ID,
		actions.HealthSplit.ID,
	}, GlobalActions...),
}

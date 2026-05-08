package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Raikage4 = game.ActorDef{
	ActorID:      uuid.MustParse("aecc5cde-93ae-42f7-8e0f-0e9db4294c6a"),
	SpriteURL:    "/sprites/4_raikage_64.png",
	Name:         "Ay (4th Raikage)",
	Affiliations: []string{game.AffKumo},
	Restricted:   true,

	Stats: map[game.ActorStat]int{
		game.StatHP:            110,
		game.StatStamina:       100,
		game.StatAttack:        145,
		game.StatDefense:       130,
		game.StatChakraAttack:  85,
		game.StatChakraDefense: 85,
		game.StatSpeed:         105,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsLightning,
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.LightningMode,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.DragonStance.ID,
		actions.RockFist.ID,
		actions.LightningLariat.ID,
		actions.LightningArmor.ID,
		actions.LightningLigerBomb.ID,
	}, GlobalActions...),
}

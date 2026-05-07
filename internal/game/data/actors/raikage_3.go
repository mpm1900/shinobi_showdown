package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Raikage3 = game.ActorDef{
	ActorID:      uuid.MustParse("1a2a1fa8-d2b2-4bb8-b748-b03c5160b8bd"),
	SpriteURL:    "/sprites/3_raikage_64.png",
	Name:         "Ay (3rd Raikage)",
	Affiliations: []string{game.AffKumo},

	Stats: map[game.ActorStat]int{
		game.StatHP:            90,
		game.StatStamina:       100,
		game.StatAttack:        112,
		game.StatDefense:       120,
		game.StatChakraAttack:  72,
		game.StatChakraDefense: 70,
		game.StatSpeed:         106,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsLightning,
		game.NsEarth,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.DragonStance.ID,
		actions.RockFist.ID,
		actions.LightningLariat.ID,
		actions.LightningArmor.ID,
		actions.Earthquake.ID,
	}, GlobalActions...),
}

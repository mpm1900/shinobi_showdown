package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Sakura = game.ActorDef{
	ActorID:      uuid.MustParse("dfa5199c-4d51-4db8-955b-a260fc01fc24"),
	SpriteURL:    "/sprites/sakura_64.png",
	Name:         "Sakura Haruno",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            71,
		game.StatStamina:       100,
		game.StatAttack:        121,
		game.StatDefense:       106,
		game.StatChakraAttack:  60,
		game.StatChakraDefense: 80,
		game.StatSpeed:         70,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsEarth,
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.HealingTactics,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CreateRain.ID,
		actions.Haze.ID,
		actions.Redirect.ID,
		actions.GreatWaterfall.ID,
		actions.RockFist.ID,
		actions.CherryBlossomImpact.ID,
		actions.SageMode.ID,
		actions.TeamHeal.ID,
		actions.MudWall.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

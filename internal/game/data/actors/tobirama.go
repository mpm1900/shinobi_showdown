package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Tobirama = game.ActorDef{
	ActorID:      uuid.MustParse("057322de-bf24-5b4a-b6df-53a0817e93cf"),
	SpriteURL:    "/sprites/tobirama_64.png",
	Name:         "Tobirama Senju",
	Clan:         game.ClanSenju,
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            70,
		game.StatStamina:       100,
		game.StatAttack:        110,
		game.StatDefense:       75,
		game.StatChakraAttack:  125,
		game.StatChakraDefense: 85,
		game.StatSpeed:         135,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
		game.NsLightning,
	}),
	Abilities: []game.Modifier{
		modifiers.ConsumeChakra,
		modifiers.MoldBreaker,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CreateRain.ID,
		actions.Tailwind.ID,
		actions.LeechSeed.ID,
		actions.Haze.ID,
		actions.GreatTreeSpear.ID,
		actions.WhirlwindKick.ID,
		actions.VacuumBlast.ID,
		actions.WaterDragon.ID,
		actions.GreatWaterfall.ID,
		actions.WaterSlicer.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

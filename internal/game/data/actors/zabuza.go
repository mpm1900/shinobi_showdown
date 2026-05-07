package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Zabuza = game.ActorDef{
	ActorID:      uuid.MustParse("7fec97b6-26ed-460f-becc-e368b98ba98d"),
	SpriteURL:    "/sprites/zabuza_64.png",
	Name:         "Zabuza Momochi",
	Affiliations: []string{game.AffKuri},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        125,
		game.StatDefense:       90,
		game.StatChakraAttack:  90,
		game.StatChakraDefense: 70,
		game.StatSpeed:         85,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Chidori.ID,
		actions.DragonStance.ID,
		actions.CollidingWave.ID,
		actions.WhirlwindKick.ID,
		actions.HiddenMist.ID,
		actions.CreateRain.ID,
		actions.WaterDragon.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

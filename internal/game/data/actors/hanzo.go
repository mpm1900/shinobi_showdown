package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Hanzo = game.ActorDef{
	ActorID:      uuid.MustParse("6ceedb86-2e9c-4be2-9a27-14f888b50eda"),
	SpriteURL:    "/sprites/hanzo_64.png",
	Name:         "Hanzō",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Stats: map[game.ActorStat]int{
		game.StatHP:            108,
		game.StatStamina:       100,
		game.StatAttack:        130,
		game.StatDefense:       95,
		game.StatChakraAttack:  80,
		game.StatChakraDefense: 85,
		game.StatSpeed:         102,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.PoisonInfused,
		modifiers.Raincaller,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Barrier.ID,
		actions.Caltrops.ID,
		actions.SwordsStance.ID,
		actions.DragonStance.ID,
		actions.GreatFireball.ID,
		actions.WaterDragon.ID,
		actions.WaterWall.ID,
		actions.CreateRain.ID,
		actions.Haze.ID,
		actions.HiddenMist.ID,
		actions.WhirlwindKick.ID,
		actions.BurningAsh.ID,
		actions.DragonFire.ID,
		actions.DarkSwamp.ID,
		actions.NoEscape.ID,
	}, GlobalActions...),
	Immunities: map[uuid.UUID]struct{}{},
}

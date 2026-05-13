package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Orochimaru = game.ActorDef{
	ActorID:      uuid.MustParse("2e9d220b-be84-524b-b7a3-b078c226fa2d"),
	SpriteURL:    "/sprites/orochimaru_64.png",
	Name:         "Orochimaru",
	Affiliations: []string{game.AffAkatsuki, game.AffOto},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        95,
		game.StatDefense:       92,
		game.StatChakraAttack:  119,
		game.StatChakraDefense: 90,
		game.StatSpeed:         94,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWind,
		game.NsYin,
	}),
	DefaultModifiers: []game.Modifier{
		modifiers.Insomnia,
	},
	Abilities: []game.Modifier{
		modifiers.Focused,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Recover.ID,
		actions.SageMode.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.MindTransfer.ID,
		actions.PatternBreak.ID,
		actions.SnakeStrike.ID,
		actions.MudWall.ID,
		actions.WaterWall.ID,
		actions.Devistation.ID,
		actions.VacuumBlast.ID,
	}, GlobalActions...),
}

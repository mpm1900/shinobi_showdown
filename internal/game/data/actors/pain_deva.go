package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PainDeva = game.ActorDef{
	ActorID:      uuid.MustParse("58e03e29-77bd-494f-a3d3-6a555519eced"),
	SpriteURL:    "/sprites/pain_deva_64.png",
	Name:         "Pain (Deva Path)",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Stats: map[game.ActorStat]int{
		game.StatHP:            65,
		game.StatStamina:       100,
		game.StatAttack:        115,
		game.StatDefense:       80,
		game.StatChakraAttack:  110,
		game.StatChakraDefense: 70,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
	}),
	Abilities: []game.Modifier{
		modifiers.Raincaller,
		modifiers.VoiceOfPain,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CreateRain.ID,
		actions.Sekiryoku.ID,
		actions.ShinraTensei.ID,
		actions.NegateJutsu.ID,
		actions.Tailwind.ID,
		actions.MindTransfer.ID,
		actions.RetreatingStrike.ID,
		actions.BlackNeedle.ID,
		actions.WaterWall.ID,
	}, GlobalActionsExcept(actions.BodyReplacement.ID)...),
	Immunities: map[uuid.UUID]struct{}{
		modifiers.BurdenOfPain.ID:    {},
		modifiers.ChainsOfPain.ID:    {},
		modifiers.JudgementOfPain.ID: {},
		modifiers.VoiceOfPain.ID:     {},
	},
}

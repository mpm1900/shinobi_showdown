package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PainPreta = game.ActorDef{
	ActorID:      uuid.MustParse("276059f9-61e2-440c-8e0b-2b2b29e88268"),
	SpriteURL:    "/sprites/pain_preta_64.png",
	Name:         "Pain (Preta Path)",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Stats: map[game.ActorStat]int{
		game.StatHP:            65,
		game.StatStamina:       100,
		game.StatAttack:        85,
		game.StatDefense:       70,
		game.StatChakraAttack:  90,
		game.StatChakraDefense: 120,
		game.StatSpeed:         60,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
	}),
	Abilities: []game.Modifier{
		modifiers.ChainsOfPain,
		modifiers.NeutralizingChakra,
		modifiers.PowerNullification,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Recover.ID,
		actions.RetreatingStrike.ID,
		actions.NegateJutsu.ID,
		actions.CopyJutsu.ID,
		actions.ChannelChakra.ID,
		actions.Expansion.ID,
		actions.BlackNeedle.ID,
	}, GlobalActionsExcept(actions.Substitution.ID)...),
	Immunities: map[uuid.UUID]struct{}{
		modifiers.BurdenOfPain.ID:    {},
		modifiers.ChainsOfPain.ID:    {},
		modifiers.JudgementOfPain.ID: {},
		modifiers.VoiceOfPain.ID:     {},
	},
}

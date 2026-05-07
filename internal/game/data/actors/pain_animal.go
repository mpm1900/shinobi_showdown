package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PainAnimal = game.ActorDef{
	ActorID:      uuid.MustParse("f6ac5c74-c27b-4879-8fbb-e61699d6049a"),
	SpriteURL:    "/sprites/pain_animal_64.png",
	Name:         "Pain (Animal Path)",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Stats: map[game.ActorStat]int{
		game.StatHP:            65,
		game.StatStamina:       100,
		game.StatAttack:        81,
		game.StatDefense:       71,
		game.StatChakraAttack:  74,
		game.StatChakraDefense: 69,
		game.StatSpeed:         124,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.RetreatingStrike.ID,
		actions.BlackNeedle.ID,
		actions.SummonAlly.ID,
		actions.DragonStance.ID,
		actions.RetreatingStrike.ID,
	}, GlobalActions...),
	Immunities: map[uuid.UUID]struct{}{
		modifiers.BurdenOfPain.ID:    {},
		modifiers.ChainsOfPain.ID:    {},
		modifiers.JudgementOfPain.ID: {},
		modifiers.VoiceOfPain.ID:     {},
	},
}

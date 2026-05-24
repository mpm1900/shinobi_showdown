package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PainHuman = game.ActorDef{
	ActorID:      uuid.MustParse("81803b0a-fb3a-4a79-bd8a-b4531503bbed"),
	SpriteURL:    "/sprites/pain_human_64.png",
	Name:         "Pain (Human Path)",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Stats: map[game.ActorStat]int{
		game.StatHP:            65,
		game.StatStamina:       100,
		game.StatAttack:        75,
		game.StatDefense:       100,
		game.StatChakraAttack:  75,
		game.StatChakraDefense: 100,
		game.StatSpeed:         95,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
	}),
	Abilities: []game.Modifier{
		modifiers.BurdenOfPain,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Redirect.ID,
		actions.PerishSong.ID,
		actions.MindTransfer.ID,
		actions.InstilFear.ID,
		actions.Flash.ID,
		actions.RetreatingStrike.ID,
		actions.SummonAlly.ID,
		actions.BlackNeedle.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
	Immunities: map[uuid.UUID]struct{}{
		modifiers.BurdenOfPain.ID:    {},
		modifiers.ChainsOfPain.ID:    {},
		modifiers.JudgementOfPain.ID: {},
		modifiers.VoiceOfPain.ID:     {},
	},
}

package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PainNaraka = game.ActorDef{
	ActorID:      uuid.MustParse("e8906381-8f07-4034-a235-1e7153f2da54"),
	SpriteURL:    "/sprites/pain_naraka_64.png",
	Name:         "Pain (Naraka Path)",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Stats: map[game.ActorStat]int{
		game.StatHP:            98,
		game.StatStamina:       100,
		game.StatAttack:        53,
		game.StatDefense:       85,
		game.StatChakraAttack:  87,
		game.StatChakraDefense: 105,
		game.StatSpeed:         67,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
	}),
	Abilities: []game.Modifier{
		modifiers.JudgementOfPain,
		modifiers.StatusReflection,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.PerishSong.ID,
		actions.TempleOfNirvana.ID,
		actions.Flash.ID,
		actions.Revival.ID,
		actions.RetreatingStrike.ID,
		actions.SummonAlly.ID,
		actions.BlackNeedle.ID,
		actions.InstilFear.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
	Immunities: map[uuid.UUID]struct{}{
		modifiers.BurdenOfPain.ID:    {},
		modifiers.ChainsOfPain.ID:    {},
		modifiers.JudgementOfPain.ID: {},
		modifiers.VoiceOfPain.ID:     {},
	},
}

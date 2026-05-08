package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Nagato = game.ActorDef{
	ActorID:      uuid.MustParse("f344914c-26d5-5e7b-8673-ad73d5b8f334"),
	SpriteURL:    "/sprites/nagato_64.png",
	Name:         "Nagato Uzumaki",
	Clan:         game.ClanUzumaki,
	Affiliations: []string{game.AffAkatsuki, game.AffAme},
	Restricted:   true,

	Stats: map[game.ActorStat]int{
		game.StatHP:            106,
		game.StatStamina:       100,
		game.StatAttack:        90,
		game.StatDefense:       99,
		game.StatChakraAttack:  130,
		game.StatChakraDefense: 134,
		game.StatSpeed:         91,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsYinYang,
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.BurdenOfPain,
		modifiers.VoiceOfPain,
		modifiers.ChainsOfPain,
		modifiers.JudgementOfPain,
		modifiers.SpeedBoost,
		modifiers.Raincaller,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Sekiryoku.ID,
		actions.Tailwind.ID,
		actions.MindTransfer.ID,
		actions.CollidingWave.ID,
	}, GlobalActions...),
	Immunities: map[uuid.UUID]struct{}{
		modifiers.BurdenOfPain.ID:    {},
		modifiers.ChainsOfPain.ID:    {},
		modifiers.JudgementOfPain.ID: {},
		modifiers.VoiceOfPain.ID:     {},
	},
}

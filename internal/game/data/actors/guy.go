package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Guy = game.ActorDef{
	ActorID:   uuid.MustParse("bf98da03-0afb-54c9-9b0a-0636552cb32d"),
	SpriteURL: "/sprites/guy_64.png",
	Name:      "Might Guy",
	Affiliations: []string{
		game.AffKonoha,
	},

	Stats: map[game.ActorStat]int{
		game.StatHP:            87,
		game.StatStamina:       80,
		game.StatAttack:        145,
		game.StatDefense:       92,
		game.StatChakraAttack:  55,
		game.StatChakraDefense: 86,
		game.StatSpeed:         115,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsTai,
		game.NsYang,
	}),
	Abilities: []game.Modifier{
		modifiers.Rage,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.LuckyStrikes.ID,
		actions.DragonStance.ID,
		actions.FlyingLotus.ID,
		actions.Asakujaku.ID,
		actions.Hirudora.ID,
		actions.WhirlwindKick.ID,
		actions.HeavyPunch.ID,
	}, GlobalActions...),
}

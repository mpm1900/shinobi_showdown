package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Gamabunta game.ActorDef = game.ActorDef{
	ActorID:   uuid.New(),
	SpriteURL: "/sprites/sub_64.png",
	Name:      "Gamabunta",
	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        100,
		game.StatDefense:       90,
		game.StatChakraAttack:  150,
		game.StatChakraDefense: 140,
		game.StatSpeed:         90,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
	}),
}

var GamabuntaActions []game.Action = []game.Action{
	CollidingWave,
}

var SummonGamabunta = MakeSummonGamabunta()

func MakeSummonGamabunta() game.Action {
	nature := game.NsSage
	config := game.ActionConfig{
		Name:        "Summon Gamabunta",
		Nature:      &nature,
		Jutsu:       game.Ninjutsu,
		Description: "",
		TargetType:  game.TargetActorID,
		Summon:      true,
	}

	return game.Action{
		ID:              uuid.MustParse("17967cc9-5b82-43d4-9cd6-3a4a2c70ca39"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.ComposeGF(game.SourceIsAlive, game.SourceHasHpRatio(0.25)),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}

				transactions = append(
					transactions,
					applySummon(context, Gamabunta, GamabuntaActions)...,
				)

				return transactions
			},
		},
	}
}

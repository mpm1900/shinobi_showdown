package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Kisame = game.ActorDef{
	ActorID:      uuid.MustParse("e34e2dec-6b2b-59f5-92c4-afb7e473f3e9"),
	Name:         "Kisame Hoshigaki",
	SpriteURL:    "/sprites/kisame_64.png",
	Affiliations: []string{game.AffAkatsuki, game.AffKuri},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       130,
		game.StatAttack:        110,
		game.StatDefense:       90,
		game.StatChakraAttack:  110,
		game.StatChakraDefense: 90,
		game.StatSpeed:         80,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.NatureSpecialist,
		modifiers.SwiftSwim,
		modifiers.WaterAbsorb,
		samehadaTransform,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.CollidingWave.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.HiddenMist.ID,
		actions.WaterDragon.ID,
		actions.SharkBomb.ID,
		actions.GreatWaterfall.ID,
		actions.WaterSlicer.ID,
	}, GlobalActions...),
}

var KisameSamehadaFusion = game.ActorDef{
	ActorID:      uuid.MustParse("f62c3886-93fd-4749-8a88-6774ffdb93d1"),
	Name:         "Kisame-Samehada Fusion",
	SpriteURL:    "/sprites/kisame_shark_64.png",
	Affiliations: []string{game.AffAkatsuki, game.AffKuri},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       130,
		game.StatAttack:        130,
		game.StatDefense:       100,
		game.StatChakraAttack:  140,
		game.StatChakraDefense: 120,
		game.StatSpeed:         90,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.SwiftSwim,
	},
}

var samehadaTransformID = uuid.MustParse("ca64b95f-4c34-4e8b-aad3-8d95f39ef924")
var samehadaTransformTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: samehadaTransformID,
	On:         game.OnActorEnter,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			state, _ := g.GetState(context)
			if state.Terrain != game.GameTerrainFlooded {
				return transactions
			}

			context.TargetPositionIDs = []uuid.UUID{}
			context.TargetActorIDs = []uuid.UUID{*context.SourceActorID}
			mut := mutations.Transform(KisameSamehadaFusion)
			transactions = append(transactions, game.MakeTransaction(mut, context))
			return transactions
		},
	},
}

var samehadaTransform = game.Modifier{
	ID:          samehadaTransformID,
	GroupID:     &samehadaTransformID,
	Icon:        "samehada_transform",
	Name:        "Samehada Fusion",
	Description: "On switch in: Transform if there's flooded terrain.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&samehadaTransformID),
	},
	Triggers: []game.Trigger{
		samehadaTransformTrigger,
	},
}

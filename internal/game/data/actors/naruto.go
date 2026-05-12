package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Naruto = game.ActorDef{
	ActorID:      uuid.MustParse("7b8d8818-ebb3-5c79-8d67-20c5df3d026d"),
	SpriteURL:    "/sprites/naruto_64.png",
	Name:         "Naruto Uzumaki",
	Clan:         game.ClanUzumaki,
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            110,
		game.StatStamina:       130,
		game.StatAttack:        60,
		game.StatDefense:       80,
		game.StatChakraAttack:  80,
		game.StatChakraDefense: 60,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsPure,
		game.NsWind,
	}),
	Abilities: []game.Modifier{
		narutoTransform,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Rasengan.ID,
		actions.GiantRasengan.ID,
		actions.PowerBoost.ID,
		actions.ActionBoost.ID,
		actions.SageMode.ID,
		actions.WhirlwindKick.ID,
		actions.SummonGamabunta.ID,
		actions.Rasenshuriken.ID,
		actions.VacuumBlast.ID,
		actions.Tailwind.ID,
		actions.TradeOffer.ID,
	}, GlobalActions...),
}

var KCMNaurto = game.ActorDef{
	ActorID:      uuid.MustParse("cb34b284-efab-40df-a15c-81930f46064c"),
	SpriteURL:    "/sprites/naruto_kcm_64.png",
	Name:         "Kurama Chakra Naruto",
	Clan:         game.ClanUzumaki,
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            110,
		game.StatStamina:       130,
		game.StatAttack:        116,
		game.StatDefense:       107,
		game.StatChakraAttack:  150,
		game.StatChakraDefense: 97,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsPure,
		game.NsWind,
	}),
	Abilities: []game.Modifier{
		kcmTransformed,
	},
}

var narutoTransformID = uuid.MustParse("d9be7d8f-55cf-4877-ace6-85f97f05a4f2")
var narutoTransformTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: narutoTransformID,
	On:         game.OnActorEnter,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			if source.Enters <= 1 {
				return transactions
			}

			context.TargetPositionIDs = []uuid.UUID{}
			context.TargetActorIDs = []uuid.UUID{*context.SourceActorID}
			mut := mutations.Transform(KCMNaurto)
			transactions = append(transactions, game.MakeTransaction(mut, context))
			return transactions
		},
	},
}

var narutoTransform = game.Modifier{
	ID:          narutoTransformID,
	GroupID:     &narutoTransformID,
	Icon:        "naruto_transform",
	Name:        "Kurama Transformation",
	Description: "On switch in: Transform if it's the second time.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&narutoTransformID),
	},
	Triggers: []game.Trigger{
		narutoTransformTrigger,
	},
}

var kcmTransformedID = uuid.MustParse("072bf59b-00dc-44ed-8f8b-a1486155c02d")

var kcmTransformed = game.Modifier{
	ID:          kcmTransformedID,
	GroupID:     &kcmTransformedID,
	Icon:        "kcm_transformed",
	Name:        "Kurama Transformation",
	Description: "On transform: Apply Chakra Terrain.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&kcmTransformedID),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: kcmTransformedID,
			On:         game.OnActorTransform,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					return modifiers.ApplyTerrain(g, context, game.GameTerrainChakra, modifiers.ChakraTerrain())
				},
			},
		},
	},
}

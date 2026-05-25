package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

func NewStageDelta(
	stat game.ActorStat,
	groupID *uuid.UUID,
	filter func(game.Game, game.Actor, game.Context) bool,
	priority int,
	delta int,
) game.Modifier {
	mut := game.MakeActorMutation(
		groupID,
		priority,
		filter,
		func(g game.Game, actor game.Actor, context game.Context) game.Actor {
			actor.Stages[stat] = actor.Stages[stat] + delta
			return actor
		},
	)

	return game.Modifier{
		ID:       uuid.New(),
		GroupID:  groupID,
		Duration: game.ModifierDurationInf,
		ActorMutations: []game.ActorMutation{
			mut,
		},
	}
}

func NewStageDeltaFromMap(
	stages map[game.ActorStat]int,
	groupID *uuid.UUID,
	filter func(game.Game, game.Actor, game.Context) bool,
	priority int,
) game.Modifier {
	mut := game.MakeActorMutation(
		groupID,
		priority,
		filter,
		func(g game.Game, actor game.Actor, context game.Context) game.Actor {
			for stat, delta := range stages {
				actor.Stages[stat] = actor.Stages[stat] + delta
			}
			return actor
		},
	)

	return game.Modifier{
		ID:       uuid.New(),
		GroupID:  groupID,
		Duration: game.ModifierDurationInf,
		ActorMutations: []game.ActorMutation{
			mut,
		},
	}
}

func NewStatMult(
	stat game.ActorStat,
	groupID *uuid.UUID,
	filter func(game.Game, game.Actor, game.Context) bool,
	priority int,
	mult float64,
) game.Modifier {
	mut := game.MakeActorMutation(
		groupID,
		priority,
		filter,
		func(g game.Game, actor game.Actor, context game.Context) game.Actor {
			actor.Stats[stat] = game.Round(float64(actor.Stats[stat]) * mult)
			return actor
		},
	)

	return game.Modifier{
		ID:       uuid.New(),
		GroupID:  groupID,
		Duration: game.ModifierDurationInf,
		ActorMutations: []game.ActorMutation{
			mut,
		},
	}
}

func MakeStatDeltaSource(stat game.ActorStat, groupID *uuid.UUID, delta int, name string) game.Modifier {
	modifier := NewStageDelta(stat, groupID, game.ComposeAF(game.ActiveFilter, game.SourceFilter), game.MutPriorityDefault, delta)
	modifier.Name = name
	return modifier
}
func MakeStatDeltaSourceWithShow(stat game.ActorStat, groupID *uuid.UUID, delta int, name string) game.Modifier {
	modifier := MakeStatDeltaSource(stat, groupID, delta, name)
	modifier.Show = true
	return modifier
}
func MakeStatDeltaTarget(stat game.ActorStat, groupID *uuid.UUID, delta int, name string) game.Modifier {
	modifier := NewStageDelta(stat, groupID, game.ComposeAF(game.ActiveFilter, game.TargetFilter), game.MutPriorityDefault, delta)
	modifier.Name = name
	return modifier
}
func MakeStatDeltaTargetWithShow(stat game.ActorStat, groupID *uuid.UUID, delta int, name string) game.Modifier {
	modifier := MakeStatDeltaTarget(stat, groupID, delta, name)
	modifier.Show = true
	return modifier
}

func MakeStatDeltaTeam(stat game.ActorStat, name string, groupID *uuid.UUID, delta int) game.Modifier {
	modifier := NewStageDelta(stat, groupID, game.ComposeAF(game.ActiveFilter, game.TeamFilter), game.MutPriorityDefault, delta)
	modifier.Name = name
	return modifier
}
func MakeStatDeltaOtherTeam(stat game.ActorStat, name string, groupID *uuid.UUID, delta int) game.Modifier {
	modifier := NewStageDelta(stat, groupID, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter), game.MutPriorityDefault, delta)
	modifier.Name = name
	return modifier
}

func MakeStatDeltaAll(stat game.ActorStat, name string, groupID *uuid.UUID, delta int) game.Modifier {
	modifier := NewStageDelta(stat, groupID, game.ActiveFilter, game.MutPriorityDefault, delta)
	modifier.Name = name
	return modifier
}

func MakeStatMultTeam(stat game.ActorStat, name string, groupID *uuid.UUID, mult float64, priority int) game.Modifier {
	modifier := NewStatMult(stat, groupID, game.ComposeAF(game.ActiveFilter, game.TeamFilter), priority, mult)
	modifier.Name = name
	modifier.Show = true
	return modifier
}
func MakeStatMultAll(stat game.ActorStat, name string, groupID *uuid.UUID, mult float64, priority int) game.Modifier {
	modifier := NewStatMult(stat, groupID, game.ComposeAF(game.ActiveFilter), priority, mult)
	modifier.Name = name
	modifier.Show = true
	return modifier
}

var AttackUpID = uuid.MustParse("ad74e7f2-2b26-50f7-a27a-24b24cb403bb")
var AttackDownID = uuid.MustParse("ec4383a0-e6c2-5302-a63b-cfe988e88291")
var AttackUpSource = MakeStatDeltaSource(game.StatAttack, &AttackUpID, 1, "Attack Up")
var AttackUp2Source = MakeStatDeltaSource(game.StatAttack, &AttackUpID, 2, "Attack Up (2)")
var AttackDownSource = MakeStatDeltaSource(game.StatAttack, &AttackDownID, -1, "Attack Down")
var AttackDown2Source = MakeStatDeltaSource(game.StatAttack, &AttackDownID, -2, "Attack Down (2)")
var AttackUpTarget = MakeStatDeltaTarget(game.StatAttack, &AttackUpID, 1, "Attack Up")
var AttackUp2Target = MakeStatDeltaTarget(game.StatAttack, &AttackUpID, 2, "Attack Up (2)")
var AttackDownTarget = MakeStatDeltaTarget(game.StatAttack, &AttackDownID, -1, "Attack Down")
var AttackDown2Target = MakeStatDeltaTarget(game.StatAttack, &AttackDownID, -2, "Attack Down (2)")
var MaxAttackSource = MakeStatDeltaSource(game.StatAttack, &AttackUpID, +6, "Max Attack")
var DefenseUpID = uuid.MustParse("215a8ec8-3b79-528e-89a7-f99759652315")
var DefenseDownID = uuid.MustParse("dd82e0f3-43e0-581e-a471-894801fb3d47")
var DefenseUpSource = MakeStatDeltaSource(game.StatDefense, &DefenseUpID, 1, "Defense Up")
var DefenseUp2Source = MakeStatDeltaSource(game.StatDefense, &DefenseUpID, 2, "Defense Up (2)")
var DefenseDownSource = MakeStatDeltaSource(game.StatDefense, &DefenseDownID, -1, "Defense Down")
var DefenseDown2Source = MakeStatDeltaSource(game.StatDefense, &DefenseDownID, -2, "Defense Down (2)")
var DefenseUpTarget = MakeStatDeltaTarget(game.StatDefense, &DefenseUpID, 1, "Defense Up")
var DefenseUp2Target = MakeStatDeltaTarget(game.StatDefense, &DefenseUpID, 2, "Defense Up (2)")
var DefenseDownTarget = MakeStatDeltaTarget(game.StatDefense, &DefenseDownID, -1, "Defense Down")
var DefenseDown2Target = MakeStatDeltaTarget(game.StatDefense, &DefenseDownID, -2, "Defense Down (2)")
var ChakraAttackUpID = uuid.MustParse("f88afbce-a618-5f5d-a557-aa359287cfc9")
var ChakraAttackDownID = uuid.MustParse("870c93e5-121d-5bb2-869b-a39e443dc4b7")
var ChakraAttackUpSource = MakeStatDeltaSource(game.StatChakraAttack, &ChakraAttackUpID, 1, "Chakra Attack Up")
var ChakraAttackUp2Source = MakeStatDeltaSource(game.StatChakraAttack, &ChakraAttackUpID, 2, "Chakra Attack Up (2)")
var ChakraAttackDownSource = MakeStatDeltaSource(game.StatChakraAttack, &ChakraAttackDownID, -1, "Chakra Attack Down")
var ChakraAttackDown2Source = MakeStatDeltaSource(game.StatChakraAttack, &ChakraAttackDownID, -2, "Chakra Attack Down (2)")
var ChakraAttackUpTarget = MakeStatDeltaTarget(game.StatChakraAttack, &ChakraAttackUpID, 1, "Chakra Attack Up")
var ChakraAttackUp2Target = MakeStatDeltaTarget(game.StatChakraAttack, &ChakraAttackUpID, 2, "Chakra Attack Up (2)")
var ChakraAttackDownTarget = MakeStatDeltaTarget(game.StatChakraAttack, &ChakraAttackDownID, -1, "Chakra Attack Down")
var ChakraAttackDown2Target = MakeStatDeltaTarget(game.StatChakraAttack, &ChakraAttackDownID, -2, "Chakra Attack Down (2)")
var ChakraDefenseUpID = uuid.MustParse("5160744e-9d6c-50e6-836c-9ad7f4e7b9e5")
var ChakraDefenseDownID = uuid.MustParse("73d3ea68-4b6d-52f8-9f83-5123abe824bb")
var ChakraDefenseUpSource = MakeStatDeltaSource(game.StatChakraDefense, &ChakraDefenseUpID, 1, "Chakra Defense Up")
var ChakraDefenseUp2Source = MakeStatDeltaSource(game.StatChakraDefense, &ChakraDefenseUpID, 2, "Chakra Defense Up (2)")
var ChakraDefenseDownSource = MakeStatDeltaSource(game.StatChakraDefense, &ChakraDefenseDownID, -1, "Chakra Defense Down")
var ChakraDefenseDown2Source = MakeStatDeltaSource(game.StatChakraDefense, &ChakraDefenseDownID, -2, "Chakra Defense Down (2)")
var ChakraDefenseUpTarget = MakeStatDeltaTarget(game.StatChakraDefense, &ChakraDefenseUpID, 1, "Chakra Defense Up")
var ChakraDefenseUp2Target = MakeStatDeltaTarget(game.StatChakraDefense, &ChakraDefenseUpID, 2, "Chakra Defense Up (2)")
var ChakraDefenseDownTarget = MakeStatDeltaTarget(game.StatChakraDefense, &ChakraDefenseDownID, -1, "Chakra Defense Down")
var ChakraDefenseDown2Target = MakeStatDeltaTarget(game.StatChakraDefense, &ChakraDefenseDownID, -2, "Chakra Defense Down (2)")
var SpeedUpID = uuid.MustParse("d7622698-9639-563e-b4d1-e7411247453f")
var SpeedDownID = uuid.MustParse("a843a341-f9f2-564a-8b8c-1dd2db7555f3")
var SpeedUpSource = MakeStatDeltaSource(game.StatSpeed, &SpeedUpID, 1, "Speed Up")
var SpeedUp2Source = MakeStatDeltaSource(game.StatSpeed, &SpeedUpID, 2, "Speed Up (2)")
var SpeedDownSource = MakeStatDeltaSource(game.StatSpeed, &SpeedDownID, -1, "Speed Down")
var SpeedDown2Source = MakeStatDeltaSource(game.StatSpeed, &SpeedDownID, -2, "Speed Down (2)")
var SpeedUpTarget = MakeStatDeltaTarget(game.StatSpeed, &SpeedUpID, 1, "Speed Up")
var SpeedUp2Target = MakeStatDeltaTarget(game.StatSpeed, &SpeedUpID, 2, "Speed Up (2)")
var SpeedDownTarget = MakeStatDeltaTarget(game.StatSpeed, &SpeedDownID, -1, "Speed Down")
var SpeedDown2Target = MakeStatDeltaTarget(game.StatSpeed, &SpeedDownID, -2, "Speed Down (2)")

var EvasionUpID = uuid.MustParse("4067fc87-e8a8-58bb-bfa1-5e852b2296b6")
var EvasionDownID = uuid.MustParse("160398fa-03b7-5370-a7a1-b7780c2f943d")
var EvasionUpSource = MakeStatDeltaSource(game.StatEvasion, &EvasionUpID, 1, "Evasion Up")
var EvasionUpTarget = MakeStatDeltaTarget(game.StatEvasion, &EvasionUpID, 1, "Evasion Up")
var AccuracyUpID = uuid.MustParse("79566281-477e-5bce-b2b4-b878b11f2882")
var AccuracyDownID = uuid.MustParse("0915300e-72d3-5b0c-b5d1-41f894a7d394")
var AccuracyUpSource = MakeStatDeltaSource(game.StatAccuracy, &AccuracyUpID, 1, "Accuracy Up")
var AccuracyUpTarget = MakeStatDeltaTarget(game.StatAccuracy, &AccuracyUpID, 1, "Accuracy Up")

// NAMED STAT UPS
var TailwindID = uuid.MustParse("cd2010e6-90d8-530f-be90-79634690e33d")
var Tailwind = MakeStatMultTeam(game.StatSpeed, "Tailwind", &TailwindID, 2.0, game.MutPriorityPostStagedStats)

var SageModeID = uuid.MustParse("764b5ee9-9136-5994-b598-40c3881e79dc")
var SageMode = MakeStatMultAll(game.StatSpeed, "Sage Mode", &SageModeID, -1, game.MutPriorityPostSet)

var HiddenMistID = uuid.MustParse("2d007fb7-270b-492c-b114-f752f9c7a17a")
var HiddenMist = MakeStatDeltaAll(game.StatAccuracy, "Hidden Mist", &HiddenMistID, -2)

// HAZE
var hazeID = uuid.MustParse("1f9dc965-2621-5e04-aa5e-6484bcf9a723")
var Haze game.Modifier = game.Modifier{
	ID:          hazeID,
	GroupID:     &hazeID,
	Icon:        "haze",
	Name:        "Haze",
	Description: "Reset all stat stages.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&hazeID,
			game.MutPriorityStagedZero,
			game.ActiveFilter,
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Stages[game.StatAccuracy] = 0
				actor.Stages[game.StatAttack] = 0
				actor.Stages[game.StatChakraAttack] = 0
				actor.Stages[game.StatChakraDefense] = 0
				actor.Stages[game.StatDefense] = 0
				actor.Stages[game.StatEvasion] = 0
				actor.Stages[game.StatHP] = 0
				actor.Stages[game.StatSpeed] = 0
				actor.Stages[game.StatStamina] = 0

				return actor
			},
		),
	},
}

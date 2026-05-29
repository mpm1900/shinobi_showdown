package game

import (
	"testing"
)

func TestResolveBaseStatsSnapshot(t *testing.T) {
	// Setup a dummy game and actor
	g := Game{}
	actor := Actor{
		ActorDef: &ActorDef{
			Stats: map[ActorStat]int{
				StatAttack: 100,
			},
		},
		Stats: map[ActorStat]int{
			StatAttack: 100,
		},
		Level: 100,
		Stages: map[ActorStat]int{
			StatAttack: 0,
		},
	}
	actor.ActorDef.Stats = map[ActorStat]int{
		StatAttack: 100,
	}
	actor.Stats = map[ActorStat]int{
		StatAttack: 100,
	}

	// Resolve the actor
	resolved := actor.Resolve(g)

	// Modify the resolved stats to simulate what a mutation might do in-place
	resolved.Stats[StatAttack] = 999

	// Verify BaseStats still has the original value
	if resolved.BaseStats[StatAttack] != 100 {
		t.Errorf("BaseStats was modified in-place! Expected 100, got %d", resolved.BaseStats[StatAttack])
	}
}

func TestResolveUnmodifiedStats(t *testing.T) {
	// Setup actor with stages
	g := Game{}
	actor := Actor{
		ActorDef: &ActorDef{
			Stats: map[ActorStat]int{
				StatAttack: 100,
			},
		},
		Stats: map[ActorStat]int{
			StatAttack: 100,
		},
		Level: 50, // Level 50 will change MapBaseStat output
		Stages: map[ActorStat]int{
			StatAttack: 2, // +2 stages
		},
	}
	actor.ActorDef.Stats = map[ActorStat]int{
		StatAttack: 100,
	}
	actor.Stats = map[ActorStat]int{
		StatAttack: 100,
	}

	resolved := actor.Resolve(g)

	// BaseStats should be 100 (unmapped)
	if resolved.BaseStats[StatAttack] != 100 {
		t.Errorf("Expected BaseStats 100, got %d", resolved.BaseStats[StatAttack])
	}

	// UnmodifiedStats should have MapBaseStats and MapStagedStats applied
	// For StatAttack, mod is 2. MapStagedStat with stage 2: (2+2)/2 = 2.0 multiplier
	// BASE_IV is 32.
	// MapBaseStat(100, 50, 1.0, 0) -> ((100*2+32)*50)/100 + 5 = (232*50)/100 + 5 = 116 + 5 = 121
	// 121 * 2.0 = 242
	expectedUnmodified := 242
	if resolved.UnmodifiedStats[StatAttack] != expectedUnmodified {
		t.Errorf("Expected UnmodifiedStats %d, got %d", expectedUnmodified, resolved.UnmodifiedStats[StatAttack])
	}
}

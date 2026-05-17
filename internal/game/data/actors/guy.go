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
		game.StatHP:            110,
		game.StatStamina:       100,
		game.StatAttack:        120,
		game.StatDefense:       80,
		game.StatChakraAttack:  50,
		game.StatChakraDefense: 90,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsTai,
	}),
	Abilities: []game.Modifier{
		modifiers.Guts,
		modifiers.Rage,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.LuckyStrikes.ID,
		actions.DragonStance.ID,
		actions.FlyingLotus.ID,
		actions.ReverseLotus.ID,
		actions.Asakujaku.ID,
		actions.Hirudora.ID,
		actions.WhirlwindKick.ID,
		actions.HeavyPunch.ID,
		actions.DesperateStrike.ID,
		actions.NightGuy.ID,
		actions.WillOfTheFallen.ID,
		actions.RaigoFist.ID,
	}, GlobalActions...),
}

package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var RockLee = game.ActorDef{
	ActorID:      uuid.MustParse("f7189c34-fb21-54d6-be14-28a473d36c53"),
	SpriteURL:    "/sprites/rocklee_64.png",
	Name:         "Rock Lee",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            55,
		game.StatStamina:       100,
		game.StatAttack:        80,
		game.StatDefense:       75,
		game.StatChakraAttack:  50,
		game.StatChakraDefense: 75,
		game.StatSpeed:         95,
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
		modifiers.PurePower,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.LuckyStrikes.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.HeavyPunch.ID,
		actions.FlyingLotus.ID,
		actions.ReverseLotus.ID,
		actions.Asakujaku.ID,
		actions.Hirudora.ID,
		actions.DesperateStrike.ID,
		actions.FocusPunch.ID,
	}, GlobalActions...),
}

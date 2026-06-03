package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Zetsu = game.ActorDef{
	ActorID:      uuid.MustParse("dfdd143c-faba-4ab2-b607-cbe69baf68e7"),
	SpriteURL:    "/sprites/zetsu_64.png",
	Name:         "Zetsu",
	Affiliations: []string{game.AffAkatsuki},

	Stats: map[game.ActorStat]int{
		game.StatHP:            114,
		game.StatStamina:       100,
		game.StatAttack:        85,
		game.StatDefense:       70,
		game.StatChakraAttack:  85,
		game.StatChakraDefense: 85,
		game.StatSpeed:         30,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWood,
		game.NsYin,
	}),
	Abilities: []game.Modifier{
		modifiers.Regeneration,
	},
	DefaultModifiers: []game.Modifier{
		modifiers.Insomnia,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.TempleOfNirvana.ID,
		actions.Redirect.ID,
		actions.Graft.ID,
		actions.FlowerBomb.ID,
		actions.TradeOffer.ID,
	}, GlobalActionsExcept(actions.Rest.ID)...),
}

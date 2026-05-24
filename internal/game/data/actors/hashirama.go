package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Hashirama = game.ActorDef{
	ActorID:      uuid.MustParse("6955b478-2b59-520a-afa3-e995d3cba9e9"),
	SpriteURL:    "/sprites/hashirama_64.png",
	Name:         "Hashirama Senju",
	Clan:         game.ClanSenju,
	Affiliations: []string{game.AffKonoha},
	Restricted:   true,

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       120,
		game.StatAttack:        155,
		game.StatDefense:       150,
		game.StatChakraAttack:  105,
		game.StatChakraDefense: 130,
		game.StatSpeed:         50,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWood,
		game.NsEarth,
	}),
	Abilities: []game.Modifier{
		modifiers.Regeneration,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Substitution.ID,
		actions.Tailwind.ID,
		actions.LeechSeed.ID,
		actions.Haze.ID,
		actions.GreatTreeSpear.ID,
		actions.WhirlwindKick.ID,
		actions.HeavyPunch.ID,
		actions.SageMode.ID,
		actions.VacuumBlast.ID,
		actions.DeepForestEmergence.ID,
		actions.Earthquake.ID,
		actions.TradeOffer.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

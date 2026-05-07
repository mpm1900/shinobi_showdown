package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Konan = game.ActorDef{
	ActorID:      uuid.MustParse("c3491894-ae53-4918-b0fb-d5f812150ee5"),
	SpriteURL:    "/sprites/konan_64.png",
	Name:         "Konan",
	Affiliations: []string{game.AffAkatsuki, game.AffAme},

	Stats: map[game.ActorStat]int{
		game.StatHP:            82,
		game.StatStamina:       100,
		game.StatAttack:        65,
		game.StatDefense:       70,
		game.StatChakraAttack:  107,
		game.StatChakraDefense: 125,
		game.StatSpeed:         71,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWind,
		game.NsWater,
	}),
	Abilities: []game.Modifier{
		modifiers.Insomnia,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Dissipate.ID,
		actions.VacuumBlast.ID,
		actions.ShikigamiDance.ID,
		actions.PaperBomb.ID,
		actions.WaterWall.ID,
	}, GlobalActions...),
}

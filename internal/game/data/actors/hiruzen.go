package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"

	"github.com/google/uuid"
)

var Hiruzen = game.ActorDef{
	ActorID:      uuid.MustParse("e8410888-d9e2-4308-a2fa-01988415d115"),
	SpriteURL:    "/sprites/hiruzen_64.png",
	Name:         "Hiruzen Sarutobi",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            100,
		game.StatStamina:       100,
		game.StatAttack:        106,
		game.StatDefense:       97,
		game.StatChakraAttack:  110,
		game.StatChakraDefense: 97,
		game.StatSpeed:         90,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsFire,
		game.NsEarth,
	}),
	Abilities:   []game.Modifier{},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Haze.ID,
		actions.GreatFireball.ID,
		actions.DragonFire.ID,
		actions.WhirlwindKick.ID,
		actions.ReaperDeathSeal.ID,
		actions.SummonAlly.ID,
		actions.Firestorm.ID,
		actions.VacuumBlast.ID,
		actions.Earthquake.ID,
		actions.TradeOffer.ID,
	}, GlobalActions...),
}

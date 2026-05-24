package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Minato = game.ActorDef{
	ActorID:      uuid.MustParse("dce0fef4-265e-5dc8-9d81-700ed3fc4877"),
	SpriteURL:    "/sprites/minato_64.png",
	Name:         "Minato Namikaze",
	Affiliations: []string{game.AffKonoha},

	Stats: map[game.ActorStat]int{
		game.StatHP:            80,
		game.StatStamina:       80,
		game.StatAttack:        130,
		game.StatDefense:       80,
		game.StatChakraAttack:  130,
		game.StatChakraDefense: 80,
		game.StatSpeed:         100,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: game.NewNatureSetValues(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsWind,
		game.NsLightning,
	}),
	Abilities: []game.Modifier{
		modifiers.SpeedBoost,
	},
	ActionCount: 4,
	ActionIDs: append([]uuid.UUID{
		actions.Rasengan.ID,
		actions.GiantRasengan.ID,
		actions.Raikiri.ID,
		actions.DragonStance.ID,
		actions.SageMode.ID,
		actions.WhirlwindKick.ID,
		actions.FlyingRaijinStrike.ID,
		actions.FlyingRaijinAssist.ID,
		actions.SummonAlly.ID,
		actions.BodyFlicker.ID,
		actions.WindSlash.ID,
		actions.PowerBoost.ID,
		actions.TradeOffer.ID,
		actions.ShinigamiCurse.ID,
	}, GlobalActions...),
}

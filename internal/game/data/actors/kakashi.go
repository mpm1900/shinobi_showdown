package actors

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/actions"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

func getKakashiResistances() map[game.Nature]float64 {
	res := game.NewNatureSetValues()
	res[game.NatureLightning] = 2.0
	res[game.NatureEarth] = 0.5
	return res
}

var Kakashi = game.ActorDef{
	ActorID:      uuid.MustParse("aa9a6360-5efe-5c51-9125-385547a97b41"),
	SpriteURL:    "/sprites/kakashi_64.png",
	Name:         "Kakashi Hatake",
	Clan:         game.ClanHatake,
	Affiliations: []string{game.AffKonoha},
	Stats: map[game.ActorStat]int{
		game.StatHP:            95,
		game.StatStamina:       70,
		game.StatAttack:        111,
		game.StatDefense:       80,
		game.StatChakraAttack:  125,
		game.StatChakraDefense: 80,
		game.StatSpeed:         109,
		game.StatEvasion:       100,
		game.StatAccuracy:      100,
	},
	NatureDamage:     game.NewNatureSetValues(),
	NatureResistance: getKakashiResistances(),
	Natures: game.MapNatures([]game.NatureSet{
		game.NsLightning,
		game.NsEarth,
	}),
	Abilities: []game.Modifier{
		modifiers.CopyAbility,
		modifiers.TargetTracking,
		modifiers.Sensei,
	},
	ActionCount: 6,
	ActionIDs: append([]uuid.UUID{
		actions.NoEscape.ID,
		actions.Raikiri.ID,
		actions.LightningHound.ID,
		actions.Chidori.ID,
		actions.DragonStance.ID,
		actions.WhirlwindKick.ID,
		actions.CopyJutsu.ID,
		actions.Kamui.ID,
		actions.KamuiCounter.ID,
		actions.KamuiSlash.ID,
		actions.MudWall.ID,
		actions.WaterWall.ID,
		actions.TradeOffer.ID,
	}, GlobalActionsExcept(actions.Substitution.ID)...),
}

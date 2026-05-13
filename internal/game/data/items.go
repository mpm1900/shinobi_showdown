package data

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var ITEMS map[uuid.UUID]game.Modifier = map[uuid.UUID]game.Modifier{
	modifiers.SealOfBodyProtection.ID: modifiers.SealOfBodyProtection,
	modifiers.SealOfImmortality.ID:    modifiers.SealOfImmortality,
	modifiers.SealOfMindProtection.ID: modifiers.SealOfMindProtection,
	modifiers.SealOfDefense.ID:        modifiers.SealOfDefense,
	modifiers.CurseMarkOfChakra.ID:    modifiers.CurseMarkOfChakra,
	modifiers.CurseMarkOfSpeed.ID:     modifiers.CurseMarkOfSpeed,
	modifiers.CurseMarkOfStrength.ID:  modifiers.CurseMarkOfStrength,
	modifiers.IchirakuRamen.ID:        modifiers.IchirakuRamen,
	modifiers.GedoShard.ID:            modifiers.GedoShard,
	modifiers.Onigiri.ID:              modifiers.Onigiri,
	modifiers.ShinobiVest.ID:          modifiers.ShinobiVest,
	modifiers.DragonFlamePepper.ID:    modifiers.DragonFlamePepper,
	modifiers.CoralFragment.ID:        modifiers.CoralFragment,
	modifiers.GraniteRing.ID:          modifiers.GraniteRing,
	modifiers.FoldingWarFan.ID:        modifiers.FoldingWarFan,
	modifiers.ConductiveBracers.ID:    modifiers.ConductiveBracers,
	modifiers.OnyxMagatama.ID:         modifiers.OnyxMagatama,
	modifiers.SagesScroll.ID:          modifiers.SagesScroll,
	modifiers.SharkSkin.ID:            modifiers.SharkSkin,
	modifiers.EyeScope.ID:             modifiers.EyeScope,
	modifiers.ShinobiCloak.ID:         modifiers.ShinobiCloak,
	modifiers.TrainingWeights.ID:      modifiers.TrainingWeights,
	modifiers.FlashPowder.ID:          modifiers.FlashPowder,
	modifiers.Medicine.ID:             modifiers.Medicine,
	modifiers.HotCoal.ID:              modifiers.HotCoal,
	modifiers.FireChakraBerry.ID:      modifiers.FireChakraBerry,
	modifiers.WindChakraBerry.ID:      modifiers.WindChakraBerry,
	modifiers.LightningChakraBerry.ID: modifiers.LightningChakraBerry,
	modifiers.EarthChakraBerry.ID:     modifiers.EarthChakraBerry,
	modifiers.WaterChakraBerry.ID:     modifiers.WaterChakraBerry,
}

package modifiers

import "shinobi_showdown/internal/game"

var TargetTracking = MakeTargetTracking()

func MakeTargetTracking() game.Modifier {
	mod := AccuracyUpSource
	mod.Name = "Target Tracking"
	mod.Icon = "target_tracking"
	mod.Description = "Accuracy up one stage"

	return mod
}

package actors

import (
	"shinobi_showdown/internal/game/data/actions"
	"slices"

	"github.com/google/uuid"
)

var GlobalActions = []uuid.UUID{
	actions.ApplyPoison.ID,
	actions.Substitution.ID,
	actions.LightningKunai.ID,
	actions.Fireball.ID,
	actions.GalePalm.ID,
	actions.Rest.ID,
	actions.ShadowClone.ID,
	actions.StoneBullet.ID,
	actions.Strike.ID,
	actions.WaterBullet.ID,
}

func GlobalActionsExcept(ids ...uuid.UUID) []uuid.UUID {
	global := slices.Clone(GlobalActions)
	return slices.DeleteFunc(global, func(id uuid.UUID) bool {
		return slices.Contains(ids, id)
	})
}

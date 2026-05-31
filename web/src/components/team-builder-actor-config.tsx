import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import type { ActorDef } from '#/lib/game/actor'
import { natureIndexes } from '#/lib/game/nature'
import { keys } from '#/lib/utils'
import { NatureBadge } from './nature-badge'
import { NastureSetDetails } from './natureset-details'
import { TeamBuilderActorAttributes } from './team-builder-actor-attributes'
import { TeamBuilderStats } from './team-builder-stats'

function TeamBuilderActorConfig({
  def,
  form,
}: {
  def: ActorDef | undefined
  form: TeamBuilderForm
}) {
  return (
    <form.Subscribe
      selector={(state) => {
        const actor = state.values.actors[state.values.selected_index]
        const items = state.values.actors
          .map((a) => a.config.item_ID)
          .filter((id) => id !== null)
        return {
          items,
          selected_index: state.values.selected_index,
          actors: state.values.actors,
          actor,
          total: Object.values(actor?.config.aux_stats ?? {}).reduce(
            (sum, value) => sum + value,
            0
          ),
        }
      }}
    >
      {({ actor, selected_index, total, items }) => (
        <div>

          <div className="flex gap-4 items-start">
            <div className="flex flex-col gap-2 min-w-1/4 overflow-hidden">
              <div className="flex my-2">
                {def && (
                  <img src={def.sprite_url} className="object-cover size-16" />
                )}
                <div className="flex flex-col px-2">
                  <div className="flex gap-6 justify-between overflow-hidden truncate">
                    <div className="flex-1">{def?.name}</div>
                  </div>
                  <div className="text-xs text-muted-foreground">
                    {def?.restricted && (
                      <span className="text-amber-300">Restricted</span>
                    )}
                  </div>
                  <div className="flex items-start">
                    {keys(def?.natures ?? [])
                      .sort((a, b) => natureIndexes[a] - natureIndexes[b])
                      .map((nature) => (
                        <NatureBadge
                          key={nature}
                          nature={nature}
                          className="text-xs"
                        />
                      ))}
                  </div>
                </div>
              </div>
              <NastureSetDetails natures={keys(def?.natures ?? [])} />
              <TeamBuilderActorAttributes
                def={def}
                otherItemIDs={items}
                focus={actor?.config?.focus ?? 'none'}
                onFocusChange={(focus) => {
                  form.setFieldValue(
                    `actors[${selected_index}].config.focus`,
                    focus
                  )
                }}
                abilityID={actor?.config?.ability_ID}
                onAbilityIDChange={(ability_ID) => {
                  form.setFieldValue(
                    `actors[${selected_index}].config.ability_ID`,
                    ability_ID
                  )
                }}
                itemID={actor?.config?.item_ID}
                onItemIDChange={(item_ID) => {
                  form.setFieldValue(
                    `actors[${selected_index}].config.item_ID`,
                    item_ID
                  )
                }}
              />
            </div>

            <TeamBuilderStats
              total={total}
              config={actor?.config}
              def={def}
              onConfigChange={(config) => {
                form.setFieldValue(`actors[${selected_index}].config`, config)
              }}
            />
          </div>
        </div>
      )}
    </form.Subscribe>
  )
}

export { TeamBuilderActorConfig }

import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import { getEffectiveness, NATURES } from '#/lib/game/nature'
import { makeConfigFromDef } from '#/lib/game/team'
import { ActorCombobox } from './actor-combobox'
import { formatDistanceToNow } from 'date-fns'
import { NatureBadge } from './nature-badge'
import { Fragment } from 'react/jsx-runtime'

const typeNatures = NATURES.filter((n) => n !== 'tai' && n !== 'pure')

function TeamBuilderList({
  form,
  created_at,
  id,
}: {
  form: TeamBuilderForm
  created_at: string | null
  id: string | null
}) {
  return (
    <form.Field name="actors" mode="array">
      {(field) => (
        <form.Subscribe
          selector={(state) => ({
            selected: state.values.actors.map((a) => a.actor_ID!),
            active: state.values.actors[state.values.selected_index]?.actor_ID,
          })}
        >
          {({ selected, active }) => (
            <div className="flex flex-col gap-1">
              <div className="flex items-center justify-between">
                <div>Team: {selected.length}/6</div>
                {created_at && (
                  <div className="text-xs">
                    {formatDistanceToNow(new Date(created_at))}
                  </div>
                )}
              </div>
              {field.state.value.map((_, i) => (
                <form.Field key={i} name={`actors[${i}]`}>
                  {(actorID) => (
                    <ActorCombobox
                      active={active}
                      selected={selected}
                      value={actorID.state.value?.actor_ID}
                      onValueChange={(actor) => {
                        actorID.handleChange(makeConfigFromDef(actor))
                        form.setFieldValue('selected_index', i)
                      }}
                      onClick={() => form.setFieldValue('selected_index', i)}
                    />
                  )}
                </form.Field>
              ))}
              {selected.length < 6 && (
                <ActorCombobox
                  active={active}
                  selected={selected}
                  value={undefined}
                  onValueChange={(actor) => {
                    form.pushFieldValue('actors', makeConfigFromDef(actor))
                    form.setFieldValue('selected_index', selected.length)
                  }}
                />
              )}
              <div className="text-xs text-muted-foreground text-center">
                {id}
              </div>
              <div className="grid grid-cols-8 gap-1 pt-4">
                <div />
                {typeNatures.map((n) => (
                  <div key={n} className="grid grid-cols-1 place-items-center">
                    <NatureBadge nature={n} />
                  </div>
                ))}
                {typeNatures.map((a) => (
                  <Fragment key={a}>
                    <div
                      className="grid grid-cols-1 place-items-center"
                    >
                      <NatureBadge nature={a} />
                    </div>
                    {typeNatures.map((b) => {
                      const eff = getEffectiveness(a, [b])
                      return (
                        <div
                          key={b}
                          className="grid grid-cols-1 place-items-center"
                        >
                          <>
                            {eff === 2 && (
                              <div className="rounded-full size-3 bg-green-600" />
                            )}
                            {eff === 1.25 && (
                              <div className="rounded-full size-2 bg-green-300" />
                            )}
                            {eff === 1 && a !== b && (
                              <div className="rounded-full size-2 bg-stone-600" />
                            )}
                            {eff === 0.8 && (
                              <div className="rounded-full size-2 bg-red-300" />
                            )}
                            {eff === 0.5 && (
                              <div className="rounded-full size-3 bg-red-400" />
                            )}
                          </>
                        </div>
                      )
                    })}
                  </Fragment>
                ))}
              </div>
            </div>
          )}
        </form.Subscribe>
      )}
    </form.Field>
  )
}

export { TeamBuilderList }

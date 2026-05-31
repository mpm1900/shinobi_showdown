import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import { makeConfigFromDef } from '#/lib/game/team'
import { formatDistanceToNow } from 'date-fns'
import { ActorCombobox } from './actor-combobox'
import { NatureEffectivenessTable } from './nature-effectiveness-table'

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
    <div className="flex flex-col gap-1">
      <div className="flex items-center justify-between">
        <form.Subscribe selector={(state) => state.values.actors.length}>
          {(count) => <div>Team: {count}/6</div>}
        </form.Subscribe>
        {created_at && (
          <div className="text-xs">
            {formatDistanceToNow(new Date(created_at))}
          </div>
        )}
      </div>

      <form.Field name="actors" mode="array">
        {(field) => (
          <form.Subscribe
            selector={(state) => ({
              selected: state.values.actors.map((a) => a.actor_ID!).join(','),
              active:
                state.values.actors[state.values.selected_index]?.actor_ID,
            })}
          >
            {({ selected, active }) => {
              const selectedArray = selected ? selected.split(',') : []
              return (
                <>
                  {field.state.value.map((_, i) => (
                    <form.Field key={i} name={`actors[${i}]`}>
                      {(actorField) => (
                        <ActorCombobox
                          active={active}
                          selected={selectedArray}
                          value={actorField.state.value?.actor_ID}
                          onValueChange={(actor) => {
                            actorField.handleChange(makeConfigFromDef(actor))
                            form.setFieldValue('selected_index', i)
                          }}
                          onClick={() =>
                            form.setFieldValue('selected_index', i)
                          }
                        />
                      )}
                    </form.Field>
                  ))}
                  {selectedArray.length < 6 && (
                    <ActorCombobox
                      key={`actor-add-${selectedArray.length}`}
                      active={active}
                      selected={selectedArray}
                      value={undefined}
                      onValueChange={(actor) => {
                        form.pushFieldValue('actors', makeConfigFromDef(actor))
                        form.setFieldValue(
                          'selected_index',
                          selectedArray.length
                        )
                      }}
                    />
                  )}
                </>
              )
            }}
          </form.Subscribe>
        )}
      </form.Field>

      <div className="text-xs text-muted-foreground text-center">{id}</div>
      <NatureEffectivenessTable />
    </div>
  )
}

export { TeamBuilderList }

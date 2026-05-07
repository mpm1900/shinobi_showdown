import { useAppForm } from '#/integrations/tanstack-form/app-form'
import { NULL_CONTEXT } from '#/lib/game/context'
import { TeamSchema, type TeamActor } from '#/lib/stores/config'
import { sendContextMessage } from '#/lib/stores/socket'

function useTeamBuilderForm({
  clientID,
  onSubmit,
}: {
  clientID: string
  onSubmit: () => void
}) {
  const form = useAppForm({
    defaultValues: {
      name: 'Team',
      selected_index: 0,
      actors: [] as TeamActor[],
    },
    validators: {
      onMount: TeamSchema,
      onChange: TeamSchema,
    },
    onSubmit: ({ value }) => {
      sendContextMessage({
        type: 'set-team',
        client_ID: clientID,
        team_config: value,
        context: NULL_CONTEXT,
      })

      onSubmit()
    },
  })

  return form
}

type TeamBuilderForm = ReturnType<typeof useTeamBuilderForm>

export { useTeamBuilderForm }
export type { TeamBuilderForm }

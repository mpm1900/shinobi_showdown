import { NULL_CONTEXT } from '#/lib/game/context'
import { teamsQuery, type Team } from '#/lib/queries/teams'
import { clientsStore } from '#/lib/stores/clients'
import { sendContextMessage, socketStore } from '#/lib/stores/socket'
import { useQuery } from '@tanstack/react-query'
import { useStore } from '@tanstack/react-store'
import { Button } from './ui/button'
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
  ComboboxTrigger,
  ComboboxValue,
} from './ui/combobox'

function LobbyTeamSelect({
  value,
  onValueChange,
}: {
  value?: Team
  onValueChange?: (team: Team | undefined) => void
}) {
  const client = useStore(clientsStore, (s) => s.me)
  const status = useStore(socketStore, (s) => s.status)
  const teams = useQuery(teamsQuery)

  return (
    <Combobox
      key={value?.id ?? 'none'}
      disabled={status !== 'open'}
      items={teams.data ?? []}
      value={value?.id ?? null}
      onValueChange={(id) => {
        const new_team = teams.data?.find((t) => t.id === id)
        if (!new_team) {
          onValueChange?.(undefined)
          return
        }
        onValueChange?.(new_team)

        if (client) {
          sendContextMessage({
            type: 'set-team',
            client_ID: client.ID,
            team_config: new_team.team_config,
            context: NULL_CONTEXT,
          })
        }
      }}
    >
      <ComboboxTrigger
        render={
          <Button variant="outline" className="justify-between min-w-40">
            <ComboboxValue placeholder="Load Team">
              {value?.team_config.name}
            </ComboboxValue>
          </Button>
        }
      />
      <ComboboxContent>
        <ComboboxInput showTrigger={false} placeholder="Search Teams" />
        <ComboboxEmpty>No teams found.</ComboboxEmpty>
        <ComboboxList>
          {(team) => (
            <ComboboxItem key={team.id} value={team.id}>
              {team.team_config.name}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  )
}

export { LobbyTeamSelect }

import { NULL_CONTEXT } from '#/lib/game/context'
import { teamsQuery, type Team } from '#/lib/queries/teams'
import { clientsStore } from '#/lib/stores/clients'
import { sendContextMessage } from '#/lib/stores/socket'
import { useQuery } from '@tanstack/react-query'
import { useStore } from '@tanstack/react-store'
import { ChevronsUpDown } from 'lucide-react'
import { useState } from 'react'
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
  onValueChange,
}: {
  onValueChange?: (team: Team | undefined) => void
}) {
  const client = useStore(clientsStore, (s) => s.me!)
  const teams = useQuery(teamsQuery)
  const [team, setTeam] = useState<Team>()

  return (
    <Combobox
      items={teams.data}
      value={team?.id}
      onValueChange={(id) => {
        const new_team = teams.data?.find((t) => t.id === id)
        if (!new_team) return
        setTeam(new_team)
        onValueChange?.(new_team)
        sendContextMessage({
          type: 'set-team',
          client_ID: client.ID,
          team_config: new_team.team_config,
          context: NULL_CONTEXT,
        })
      }}
    >
      <ComboboxTrigger
        render={
          <Button variant="outline" className="justify-between min-w-40">
            <ComboboxValue placeholder="Load Team">
              {team?.team_config.name}
            </ComboboxValue>
            <ChevronsUpDown />
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

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from './ui/sidebar'
import { cloneTeamConfig } from '#/lib/team-storage'
import { teamsQuery, type Team } from '#/lib/queries/teams'
import { Button } from './ui/button'
import { Loader2, Plus, Trash } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { useEffect } from 'react'

function TeamBuilderSidebar({
  onLoadTeam,
  onDeleteTeam,
}: {
  onLoadTeam: (team: Team) => void
  onDeleteTeam: (team: Team) => void
}) {
  const params = useSearch({ from: '/team-builder' })
  const nav = useNavigate()
  const teams = useQuery(teamsQuery)
  const loadSavedTeam = (team: Team) => {
    const team_config = cloneTeamConfig(team.team_config)
    onLoadTeam({
      ...team,
      team_config,
    })
    nav({
      to: '/team-builder',
      search: {
        team_ID: team.id ?? undefined,
      },
    })
  }

  useEffect(() => {
    const team = teams.data?.find((t) => t.id === params.team_ID)
    if (team) {
      loadSavedTeam(team)
    }
  }, [])

  return (
    <Sidebar collapsible="none" className="border-r bg-stone-900 h-auto!">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="#">
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-medium">Team Builder</span>
                  <span className="text-xs text-muted-foreground">v1.0.0</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel className="justify-between">
            Teams{' '}
            <Button
              size="icon-xs"
              variant="outline"
              onClick={() => {
                onLoadTeam({
                  id: null,
                  created_at: null,
                  team_config: {
                    name: '',
                    selected_index: 0,
                    actors: [],
                  },
                })
                nav({
                  to: '/team-builder',
                  search: {
                    team_ID: undefined,
                  },
                })
              }}
            >
              <Plus />
            </Button>
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              {teams.data?.length === 0 ? (
                <SidebarMenuItem>
                  <SidebarMenuButton disabled>
                    <span className="text-muted-foreground">
                      No saved teams
                    </span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ) : (
                teams.data?.map((team) => (
                  <SidebarMenuItem key={team.team_config.name}>
                    <SidebarMenuButton
                      asChild
                      isActive={params.team_ID === team.id}
                      className="justify-between group"
                      onClick={() => loadSavedTeam(team)}
                    >
                      <div>
                        <span>{team.team_config.name}</span>
                        <Button
                          size="icon-xs"
                          className="hidden group-hover:flex"
                          variant="ghost"
                          onClick={(e) => {
                            e.stopPropagation()
                            e.preventDefault()
                            onDeleteTeam(team)
                          }}
                        >
                          <Trash />
                        </Button>
                      </div>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))
              )}
              {teams.isPending && (
                <div className="grid place-items-center">
                  <Loader2 className="animate-spin" />
                </div>
              )}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
    </Sidebar>
  )
}

export { TeamBuilderSidebar }

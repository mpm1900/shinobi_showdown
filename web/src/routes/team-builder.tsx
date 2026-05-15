import { AppHeader } from '#/components/app-header'
import { SidebarInset, SidebarProvider } from '#/components/ui/sidebar'
import { actionsQuery } from '#/lib/queries/actions'
import { actorsQuery } from '#/lib/queries/actors'
import { instancesQuery } from '#/lib/queries/instances'
import { clientsStore } from '#/lib/stores/clients'
import { useStore } from '@tanstack/react-form'
import { useQueryClient } from '@tanstack/react-query'
import { ClientOnly, createFileRoute, redirect } from '@tanstack/react-router'
import { TeamBuilderSidebar } from '#/components/team-builder-sidebar'
import { useTeamBuilderForm } from '#/hooks/use-team-builder-form'
import { TeamBuilderList } from '#/components/team-builder-list'
import { teamsQuery, type Team } from '#/lib/queries/teams'
import { useState } from 'react'
import { useDeleteTeam } from '#/lib/mutations/delete-team'
import { TeamBuilderActions } from '#/components/team-builder-actions'
import { TeamBuilderActor } from '#/components/team-builder-actor'
import z from 'zod'

const TeamBuilderSchema = z.object({
  team_ID: z.string().optional(),
})

export const Route = createFileRoute('/team-builder')({
  component: RouteComponent,
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: '/login' })
    }
  },
  loader: async ({ context }) => {
    await context.queryClient.ensureQueryData(actionsQuery)
    await context.queryClient.ensureQueryData(actorsQuery)
    await context.queryClient.ensureQueryData(instancesQuery)
    await context.queryClient.ensureQueryData(teamsQuery)
  },
  validateSearch: TeamBuilderSchema,
})

function RouteComponent() {
  const nav = Route.useNavigate()
  const { team_ID } = Route.useSearch()
  const client = useStore(clientsStore, (s) => s.me)
  const form = useTeamBuilderForm({
    clientID: client?.ID ?? '',
    onSubmit: () => {
      nav({
        to: '/lobby',
      })
    },
  })

  const [team, setTeam] = useState<Team>()
  if (team && team_ID) {
    team.id = team_ID
  }
  const qc = useQueryClient()
  const deleteMutation = useDeleteTeam()

  const loadSavedTeam = (team: Team) => {
    setTeam(team)
    form.setFieldValue('name', team.team_config.name)
    form.setFieldValue('actors', team.team_config.actors)
    form.setFieldValue(
      'selected_index',
      Math.min(
        team.team_config.selected_index ?? 0,
        Math.max(team.team_config.actors.length - 1, 0)
      )
    )
  }

  return (
    <ClientOnly>
      <main className="flex min-h-screen flex-col">
        <AppHeader />

        <section className="flex flex-1 p-4 md:p-6 mt-9">
          <SidebarProvider className="min-h-full overflow-hidden w-full rounded-xl border border-stone-300/30 ring ring-black bg-stone-950 shadow-sm">
            <TeamBuilderSidebar
              onLoadTeam={loadSavedTeam}
              onDeleteTeam={(t) => {
                deleteMutation.mutate(t.id, {
                  onSuccess: () => {
                    qc.invalidateQueries(teamsQuery)
                  },
                })
              }}
            />

            <SidebarInset className="min-h-0 bg-stone-950">
              <form.AppForm>
                <div className="flex min-h-0 flex-row-reverse items-stretch p-4 gap-8">
                  <div>
                    <TeamBuilderActions id={team?.id ?? null} form={form} />
                    <TeamBuilderList
                      created_at={team?.created_at ?? null}
                      id={team?.id ?? null}
                      form={form}
                    />
                  </div>

                  <TeamBuilderActor form={form} />
                </div>
              </form.AppForm>
            </SidebarInset>
          </SidebarProvider>
        </section>
      </main>
    </ClientOnly>
  )
}

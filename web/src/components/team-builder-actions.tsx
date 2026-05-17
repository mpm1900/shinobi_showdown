import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import { ChevronRight, Loader2, Save, Swords } from 'lucide-react'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { clientsStore } from '#/lib/stores/clients'
import { useStore } from '@tanstack/react-store'
import { useUpsertTeam } from '#/lib/mutations/upsert-team'
import type { TeamConfig } from '#/lib/stores/config'
import { useQueryClient } from '@tanstack/react-query'
import { teamsQuery } from '#/lib/queries/teams'
import { useNavigate } from '@tanstack/react-router'

function TeamBuilderActions({
  id,
  form,
  onSaveSuccess,
}: {
  id: string | null
  form: TeamBuilderForm
  onSaveSuccess?: (teamID: string | null) => void
}) {
  const client = useStore(clientsStore, (s) => s.me)
  const qc = useQueryClient()
  const upsertMutation = useUpsertTeam()
  const nav = useNavigate({
    from: '/team-builder',
  })
  return (
    <div className="mb-4 flex items-center justify-end gap-6">
      <form.Subscribe>
        {({ isValid, isSubmitting, isValidating, values }) => {
          return (
            <div className="flex gap-2 w-full">
              <form.Field name="name">
                {(field) => (
                  <Input
                    placeholder="Team Name"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                )}
              </form.Field>

              <Button
                size="icon"
                variant="secondary"
                disabled={
                  upsertMutation.isPending || !values.name.trim() || !isValid
                }
                onClick={() => {
                  upsertMutation.mutate(
                    {
                      id,
                      config: values as TeamConfig,
                    },
                    {
                      onSuccess: (data) => {
                        qc.invalidateQueries(teamsQuery)
                        onSaveSuccess?.(data.id ?? null)
                        nav({
                          to: '/team-builder',
                          search: {
                            team_ID: data.id ?? undefined,
                          },
                        })
                      },
                    }
                  )
                }}
              >
                {upsertMutation.isPending ? (
                  <Loader2 className="animate-spin" />
                ) : (
                  <Save />
                )}
              </Button>
              <Button
                disabled={!isValid || isSubmitting || isValidating || !client}
                onClick={form.handleSubmit}
                className="cursor-pointer"
              >
                Load
                <ChevronRight />
              </Button>
            </div>
          )
        }}
      </form.Subscribe>
    </div>
  )
}

export { TeamBuilderActions }

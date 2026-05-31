import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import { useUpsertTeam } from '#/lib/mutations/upsert-team'
import { teamsQuery } from '#/lib/queries/teams'
import { clientsStore } from '#/lib/stores/clients'
import { type TeamConfig } from '#/lib/stores/config'
import { useQueryClient } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { ChevronRight, Loader2, Save } from 'lucide-react'
import { Button } from './ui/button'
import { DebouncedInput, Input } from './ui/input'

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
      <div className="flex gap-2 w-full">
        <form.Field name="name">
          {(field) => (
            <Input
              placeholder="Team Name"
              value={field.state.value}
              onValueChange={(v) => field.handleChange(v as string)}
              aria-invalid={!field.state.meta.isValid}
            />
          )}
        </form.Field>
        <form.Subscribe>
          {({ isValid, isSubmitting, isValidating, values }) => {
            return (
              <>
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
                  size="icon"
                >
                  <ChevronRight />
                </Button>
              </>
            )
          }}
        </form.Subscribe>
      </div>
    </div>
  )
}

export { TeamBuilderActions }

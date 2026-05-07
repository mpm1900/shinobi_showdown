import { clientsStore } from '#/lib/stores/clients'
import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from './ui/alert-dialog'
import type { Context } from '#/lib/game/context'
import { sendContextMessage } from '#/lib/stores/socket'
import { TargetButton } from './target-button'
import { Button } from './ui/button'
import { useEffect, useState } from 'react'
import type { ActionTransaction } from '#/lib/game/action'
import { useValidateContext } from '#/hooks/use-validate-context'
import { useGetTargets } from '#/hooks/use-get-targets'
import { ChevronRight, Loader2 } from 'lucide-react'
import { useDebounce } from 'use-debounce'

function PromptControl({
  loading = false,
  prompt,
  context,
  onContextChange,
}: {
  loading?: boolean
  prompt?: ActionTransaction
  context: Context
  onContextChange: (context: Context) => void
}) {
  const action = prompt?.mutation
  const actors = useStore(gameStore, (g) => g.actors)
  const { valid } = useValidateContext(context, prompt?.ID)

  const client = useStore(clientsStore, (c) => c.me!)
  const { context: t_context } = useGetTargets(context, prompt?.ID)

  return (
    <>
      {action && (
        <div className="grid grid-cols-2 gap-3">
          {actors
            .filter((a) => t_context?.target_actor_IDs?.includes(a.ID))
            .map((a) => (
              <TargetButton
                key={a.ID}
                actor={a}
                enabled={!loading}
                loading={false}
                contextValid={!!valid}
                context={context}
                onContextChange={onContextChange}
                targetType={action.target_type}
              />
            ))}
        </div>
      )}
      <AlertDialogFooter>
        <AlertDialogAction asChild>
          <Button
            disabled={!valid || loading}
            onClick={() => {
              sendContextMessage({
                type: 'resolve-prompt',
                client_ID: client.ID,
                prompt_ID: prompt?.ID,
                context,
              })
            }}
          >
            Select{' '}
            {loading ? <Loader2 className="animate-spin" /> : <ChevronRight />}
          </Button>
        </AlertDialogAction>
      </AlertDialogFooter>
    </>
  )
}

function PromptController() {
  const status = useStore(gameStore, (g) => g.status)
  const _prompt = useStore(gameStore, (g) => g.prompt)
  const [prompt] = useDebounce(_prompt, 100)
  const [context, setContext] = useState(prompt?.context)

  useEffect(() => {
    setContext(prompt?.context)
  }, [prompt?.ID])

  return (
    <>
      <AlertDialog open={!!prompt && status === 'idle'}>
        {prompt && context && (
          <AlertDialogContent className="bg-stone-950">
            <AlertDialogHeader className="gap-0">
              <AlertDialogTitle>{prompt.mutation.config.name}</AlertDialogTitle>
              <AlertDialogDescription>
                Select Shinobi to switch in.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <PromptControl
              loading={!_prompt}
              prompt={prompt}
              context={context}
              onContextChange={setContext}
            />
          </AlertDialogContent>
        )}
      </AlertDialog>
      <AlertDialog open={!prompt && status === 'waiting'}>
        <AlertDialogContent className="bg-stone-950">
          <AlertDialogHeader>
            <AlertDialogTitle>Waiting on other players...</AlertDialogTitle>
            <AlertDialogDescription>
              Another player is thinking...
            </AlertDialogDescription>
          </AlertDialogHeader>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}

export { PromptController }

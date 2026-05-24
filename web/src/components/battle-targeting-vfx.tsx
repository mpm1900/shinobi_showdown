import type { Context } from '#/lib/game/context'
import { battleContext } from '#/lib/stores/battle-context'
import { clientsStore } from '#/lib/stores/clients'
import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import { useEffect, useMemo, useState } from 'react'

type Point = {
  x: number
  y: number
}

function getAnchorPoint(
  actorID: string,
  actorPlayerID: string | undefined,
  clientID: string | undefined
): Point | null {
  const avatar = document.querySelector<HTMLElement>(
    `[data-actor-avatar-anchor="${actorID}"]`
  )
  if (avatar) {
    const rect = avatar.getBoundingClientRect()
    const isBottomTeam = Boolean(clientID && actorPlayerID === clientID)
    return {
      x: rect.left + rect.width / 2,
      y: isBottomTeam ? rect.top : rect.bottom,
    }
  }

  const el = document.querySelector<HTMLElement>(
    `[data-actor-anchor="${actorID}"]`
  )
  if (!el) return null
  const rect = el.getBoundingClientRect()
  return {
    x: rect.left + rect.width / 2,
    y: rect.top + rect.height / 2,
  }
}

type SelectionContext = Context & {
  hover_target_IDs?: string[]
}

type BeamContext = {
  key: string
  context: Context
  style: 'running' | 'selecting'
}

function toTargetActorIDs(
  context: Context,
  positionToActor: Map<string, string>
) {
  const actorTargets = context.target_actor_IDs ?? []
  const positionTargets = (context.target_position_IDs ?? [])
    .map((positionID) => positionToActor.get(positionID) ?? null)
    .filter((id): id is string => Boolean(id))

  return Array.from(new Set([...actorTargets, ...positionTargets])).filter(
    (id) => id !== context.source_actor_ID
  )
}

function BattleTargetingVfx({
  containerRef,
  runningContext,
}: {
  containerRef: { current: HTMLDivElement | null }
  runningContext: Context | null
}) {
  const game = useStore(gameStore, (g) => g)
  const status = useStore(gameStore, (g) => g.status)
  const bc = useStore(battleContext, (c) => c)
  const selectingContext = status === 'idle' ? bc : null
  const clientID = useStore(clientsStore, (c) => c.me?.ID)
  const [frame, setFrame] = useState(0)
  const positionToActor = useMemo(
    () =>
      new Map(
        game.actors
          .filter((a) => a.position_ID)
          .map((a) => [a.position_ID!, a.ID])
      ),
    [game.actors]
  )

  const beamContexts = useMemo(() => {
    const contexts: BeamContext[] = []

    if (selectingContext?.source_actor_ID) {
      const hasTargets =
        (selectingContext.target_actor_IDs?.length ?? 0) > 0 ||
        (selectingContext.target_position_IDs?.length ?? 0) > 0 ||
        (selectingContext.hover_target_IDs?.length ?? 0) > 0

      if (hasTargets) {
        contexts.push({
          key: `selecting:${selectingContext.source_actor_ID}`,
          style: 'selecting',
          context: {
            ...selectingContext,
            target_actor_IDs: [
              ...(selectingContext.target_actor_IDs ?? []),
              ...(selectingContext.hover_target_IDs ?? []),
            ],
          },
        })
      }
    }

    if (runningContext?.source_actor_ID) {
      contexts.push({
        key: `running:${runningContext.source_actor_ID}`,
        style: 'running',
        context: runningContext,
      })
    }

    return contexts
  }, [runningContext, selectingContext])

  const segments = useMemo(() => {
    const container = containerRef.current
    if (!container || beamContexts.length === 0) return []

    const containerRect = container.getBoundingClientRect()
    return beamContexts.flatMap((beam) => {
      if (!beam.context.source_actor_ID) return []
      const sourceActor = game.actors.find(
        (a) => a.ID === beam.context.source_actor_ID
      )
      const sourceViewport = getAnchorPoint(
        beam.context.source_actor_ID,
        sourceActor?.player_ID,
        clientID
      )
      if (!sourceViewport) return []

      const source = {
        x: sourceViewport.x - containerRect.left,
        y: sourceViewport.y - containerRect.top,
      }

      const targetIDs = toTargetActorIDs(beam.context, positionToActor)
      return targetIDs
        .map((targetID) => {
          const targetActor = game.actors.find((a) => a.ID === targetID)
          const targetViewport = getAnchorPoint(
            targetID,
            targetActor?.player_ID,
            clientID
          )
          if (!targetViewport) return null
          const sourceOutDir =
            sourceActor?.player_ID && sourceActor.player_ID === clientID
              ? -1
              : 1
          const targetOutDir =
            targetActor?.player_ID && targetActor.player_ID === clientID
              ? -1
              : 1
          const targetInDir = -targetOutDir
          return {
            key: `${beam.key}:${targetID}`,
            style: beam.style,
            source,
            sourceOutDir,
            targetInDir,
            target: {
              x: targetViewport.x - containerRect.left,
              y: targetViewport.y - containerRect.top,
            },
          }
        })
        .filter(
          (
            segment
          ): segment is {
            key: string
            style: 'running' | 'selecting'
            source: Point
            sourceOutDir: number
            targetInDir: number
            target: Point
          } => Boolean(segment)
        )
    })
  }, [
    beamContexts,
    clientID,
    containerRef,
    frame,
    game.actors,
    positionToActor,
  ])

  useEffect(() => {
    if (beamContexts.length === 0) return

    let raf = requestAnimationFrame(function tick() {
      setFrame((v) => v + 1)
      raf = requestAnimationFrame(tick)
    })
    return () => cancelAnimationFrame(raf)
  }, [beamContexts.length])

  if (segments.length === 0) return null

  return (
    <div className="pointer-events-none absolute inset-0 z-0">
      <svg className="h-full w-full">
        <defs>
          <linearGradient
            id="battle-target-beam"
            x1="0%"
            y1="0%"
            x2="100%"
            y2="0%"
          >
            <stop offset="0%" stopColor="rgba(255, 255, 255, 0.60)" />
            <stop offset="40%" stopColor="rgba(255, 255, 255, 0.20)" />
            <stop offset="60%" stopColor="rgba(255, 255, 255, 0.20)" />
            <stop offset="100%" stopColor="rgba(255, 255, 255, 0.60)" />
          </linearGradient>
          <filter
            id="battle-target-glow"
            x="-50%"
            y="-50%"
            width="200%"
            height="200%"
          >
            <feMerge>
              <feMergeNode in="blurred" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>
        {segments.map((segment) => {
          const dx = segment.target.x - segment.source.x
          const dy = segment.target.y - segment.source.y
          const controlDistance = Math.max(
            52,
            Math.min(180, Math.abs(dy) * 0.5 + Math.abs(dx) * 0.18)
          )
          const c1x = segment.source.x
          const c1y = segment.source.y + segment.sourceOutDir * controlDistance
          const c2x = segment.target.x
          const c2y = segment.target.y - segment.targetInDir * controlDistance
          const d = `M ${segment.source.x} ${segment.source.y} C ${c1x} ${c1y}, ${c2x} ${c2y}, ${segment.target.x} ${segment.target.y}`

          return (
            <g key={segment.key}>
              <path
                d={d}
                stroke="url(#battle-target-beam)"
                strokeWidth={6}
                strokeLinecap="round"
                fill="none"
                opacity={segment.style === 'running' ? 0.25 : 0.12}
                filter="url(#battle-target-glow)"
              />
              <path
                d={d}
                stroke="url(#battle-target-beam)"
                strokeWidth={segment.style === 'running' ? 2.5 : 2}
                strokeLinecap="round"
                fill="none"
                strokeDasharray="10 8"
                opacity={segment.style === 'running' ? 0.95 : 0.78}
              >
                <animate
                  attributeName="stroke-dashoffset"
                  from="0"
                  to="-36"
                  dur={segment.style === 'running' ? '0.5s' : '0.8s'}
                  repeatCount="indefinite"
                />
              </path>
            </g>
          )
        })}
      </svg>
    </div>
  )
}

export { BattleTargetingVfx }

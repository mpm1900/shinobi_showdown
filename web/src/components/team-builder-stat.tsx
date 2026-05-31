import {
  type ActorBaseStat,
  type ActorDef,
  type ActorFocus,
  type ActorStats,
  ACTOR_FOCUS_DETAILS,
} from '#/lib/game/actor'
import type { ActorConfig } from '#/lib/stores/socket'
import { cn } from '#/lib/utils'
import { useEffect, useState } from 'react'
import { useDebouncedCallback } from 'use-debounce'
import { Input } from './ui/input'
import { Slider } from './ui/slider'

const CAP = 66
const PER_STAT_MAX = 32
const MAX_COLOR_STAT = 255
const MAX_BAR_WIDTH_PERCENT = 140

const STAT_COLOR_STOPS = [
  { stat: 0, color: [255, 0, 0] },
  { stat: 55, color: [190, 85, 0] },
  { stat: 61, color: [180, 98, 0] },
  { stat: 95, color: [139, 153, 0] },
  { stat: 100, color: [134, 160, 0] },
  { stat: 110, color: [92, 216, 0] },
  { stat: 134, color: [67, 249, 0] },
  { stat: 255, color: [0, 255, 255] },
] as const

function getStatBarColor(baseStat: number) {
  const stat = Math.max(0, Math.min(MAX_COLOR_STAT, baseStat))

  if (stat <= STAT_COLOR_STOPS[0].stat) {
    const [r, g, b] = STAT_COLOR_STOPS[0].color
    return `rgb(${r}, ${g}, ${b})`
  }

  for (let i = 1; i < STAT_COLOR_STOPS.length; i++) {
    const left = STAT_COLOR_STOPS[i - 1]
    const right = STAT_COLOR_STOPS[i]

    if (stat <= right.stat) {
      const t = (stat - left.stat) / (right.stat - left.stat)
      const r = Math.round(left.color[0] + (right.color[0] - left.color[0]) * t)
      const g = Math.round(left.color[1] + (right.color[1] - left.color[1]) * t)
      const b = Math.round(left.color[2] + (right.color[2] - left.color[2]) * t)
      return `rgb(${r}, ${g}, ${b})`
    }
  }

  const [r, g, b] = STAT_COLOR_STOPS[STAT_COLOR_STOPS.length - 1].color
  return `rgb(${r}, ${g}, ${b})`
}

const STAT_NAMES: ActorStats<string> = {
  hp: 'HP',
  stamina: 'Stamina',
  speed: 'Speed',
  accuracy: 'Accuracy',
  evasion: 'Evasion',
  attack: 'Attack',
  defense: 'Defense',
  chakra_attack: 'Chakra Attack',
  chakra_defense: 'Chakra Defense',
}

function TeamBuilderStatGuage({
  baseStat,
  greyscale = false,
}: {
  baseStat: number
  greyscale?: boolean
}) {
  const widthPercent =
    (Math.max(0, baseStat) * MAX_BAR_WIDTH_PERCENT) / MAX_COLOR_STAT
  const barColor = getStatBarColor(baseStat)

  return (
    <div className="relative bg-gray-600 rounded-md h-6 w-full min-w-40 overflow-hidden ring ring-black">
      <div
        className="absolute top-0 left-0 rounded-md h-6"
        style={{
          width: `${widthPercent}%`,
          opacity: greyscale ? 0.2 : 1,
          backgroundColor: greyscale ? 'var(--accent-foreground)' : barColor,
        }}
      />
    </div>
  )
}

function TeamBuilderStat({
  total,
  base,
  stat,
  config,
  onConfigChange,
}: {
  total: number
  base: ActorDef | undefined
  stat: ActorBaseStat
  focus?: ActorFocus
  config: ActorConfig | undefined
  onConfigChange: (config: ActorConfig) => void
}) {
  const focus = config?.focus ?? 'none'
  const detail = ACTOR_FOCUS_DETAILS[focus]
  const up = detail.up === stat
  const down = detail.down === stat
  const aux = (config?.aux_stats as any)?.[stat] ?? 0
  const [localAux, setLocalAux] = useState(aux)

  useEffect(() => {
    setLocalAux(aux)
  }, [aux])

  const debouncedConfigChange = useDebouncedCallback((value: number) => {
    if (!config) return
    onConfigChange({
      ...config,
      aux_stats: {
        ...config.aux_stats,
        [stat]: value,
      },
    })
  }, 120)

  function handleConfigChange(value: number) {
    const otherStatsTotal = total - aux
    const maxAllowedByCap = CAP - otherStatsTotal
    const clamped = Math.max(0, Math.min(PER_STAT_MAX, maxAllowedByCap, value))
    setLocalAux(clamped)
    debouncedConfigChange(clamped)
  }

  return (
    <tr>
      <td
        className={cn('w-24 text-muted-foreground whitespace-nowrap text-xs', {
          'text-green-300 font-bold': up,
          'text-red-300 font-bold': down,
        })}
      >
        {STAT_NAMES[stat]}
      </td>
      <td className="w-8 text-right p-2 py-1 whitespace-nowrap font-black">
        {base?.stats[stat]}
      </td>
      <td className="hidden lg:flex place-items-center py-1.5">
        <TeamBuilderStatGuage
          baseStat={base?.stats[stat] ?? 0}
          greyscale={stat === 'stamina'}
        />
      </td>
      <td className="px-2 w-12">
        <Input
          className="w-16"
          value={localAux}
          type="number"
          onChange={(e) => {
            const v = parseInt(e.target.value)
            if (!isNaN(v)) {
              handleConfigChange(v)
            }
          }}
        />
      </td>
      <td className="hidden xl:flex">
        <Slider
          value={[localAux]}
          max={PER_STAT_MAX}
          step={1}
          className="min-w-40"
          onValueChange={(v) => handleConfigChange(v[0])}
        />
      </td>
    </tr>
  )
}

export { TeamBuilderStat }

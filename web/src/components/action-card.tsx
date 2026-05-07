import type { Action } from '#/lib/game/action'
import { cn } from '#/lib/utils'
import type { ComponentProps, ReactNode } from 'react'
import { NatureBadge } from './nature-badge'

type ActionCardProps = ComponentProps<'button'> & {
  action: Action
  selected?: boolean
  usedSummon?: boolean
}

function ActionCard({
  action,
  selected = false,
  disabled = false,
  usedSummon = false,
  className,
  ...props
}: ActionCardProps) {
  const accuracyLabel =
    action.config.accuracy !== undefined ? `${action.config.accuracy}%` : '-'

  return (
    <button
      type="button"
      disabled={disabled}
      className={cn(
        'group',
        'flex flex-col w-[230px] h-[200px] rounded-lg border-2 border-stone-900 hover:border-stone-700 text-left relative',
        `bg-black text-foreground`,
        'transition-all duration-200',
        'hover:-translate-y-0.5',
        'focus-visible:outline-none',
        'flex flex-col shadow-xl shadow-black/60 hover:shadow-black',
        'ring ring-black',
        {
          'border-stone-100 hover:border-stone-100 shadow-lg shadow-black':
            selected,
          'border-orange-400/40 hover:border-orange-400':
            action.config.stat === 'attack',
          'border-orange-400': action.config.stat === 'attack' && selected,
          'border-red-400/40 hover:border-red-400':
            action.config.stat === 'defense',
          'border-red-400': action.config.stat === 'defense' && selected,
          'border-indigo-400/50 hover:border-indigo-400':
            action.config.stat === 'chakra_attack',
          'border-indigo-400':
            action.config.stat === 'chakra_attack' && selected,
          'border-blue-400/50 hover:border-blue-400':
            action.config.stat === 'chakra_defense',
          'border-blue-400':
            action.config.stat === 'chakra_defense' && selected,
        },
        className
      )}
      {...props}
    >
      <div className="relative">
        <div
          className={cn(
            'flex items-start justify-between gap-2 p-2 bg-stone-950 text-stone-100 rounded-t-md',
            disabled && 'opacity-50'
          )}
        >
          <div className="min-w-0 flex items-center gap-2">
            {action.config.nature && (
              <NatureBadge nature={action.config.nature} />
            )}
            <div className="-space-y-1">
              <div className="truncate text-xs font-semibold">
                {action.config.name}
              </div>
              {action.cooldown !== null ? (
                <div className="text-[10px] uppercase tracking-wide text-destructive">
                  on cooldown
                </div>
              ) : disabled ? (
                <div className="text-[10px] uppercase tracking-wide text-destructive">
                  {action.summon && usedSummon ? 'summon used' : 'disabled'}
                </div>
              ) : (
                <div className="text-[10px] uppercase tracking-wide text-stone-100/50">
                  {action.config.jutsu || '-'}
                </div>
              )}
            </div>
          </div>

          {!!action.config.power && action.config.power > 0 && (
            <div className="text-stone-100/90 font-black nanum-brush-script-regular text-4xl h-8 -mt-0.5">
              {action.config.power}
            </div>
          )}
        </div>
        <BrushEdge />
      </div>

      <div className="flex *:flex-1 text-[11px] border-t border-b bg-stone-100 shadow-md z-2">
        <StatChip label="Acc" value={accuracyLabel} />
        <StatChip label="Cost" value={action.config.cost || '-'} />
        <StatChip label="Priority" value={action.priority || '-'} />
      </div>
      <div
        className={`flex-1 px-3 py-2 text-xs text-black font-bold rounded-b-md bg-[url('/paper.jpg')] relative`}
      >
        <div
          className={cn(
            'bg-black/30 group-hover:bg-black/15 absolute inset-0 rounded z-1 transition-colors',
            selected && 'bg-transparent!'
          )}
        />
        {action.config.description}
      </div>
    </button>
  )
}

function BrushEdge() {
  return (
    <svg
      className="z-10 pointer-events-none absolute -bottom-[3px] left-0 w-full h-4"
      viewBox="0 0 200 14"
      preserveAspectRatio="none"
      aria-hidden="true"
    >
      <path
        d="M3 5C12 9 21 2 31 7C41 12 50 3 60 8C70 13 79 4 89 9C99 14 108 5 118 9C128 13 137 5 147 9C157 13 166 5 176 9C185 12 193 6 197 8"
        stroke="#000"
        strokeWidth="1.8"
        strokeLinecap="round"
        fill="none"
        opacity="0.96"
      />
      <path
        d="M10 8C18 11 27 5 35 9C44 13 53 7 61 10C70 13 79 8 87 11C96 14 105 9 113 11C122 13 131 9 139 11C148 13 157 9 165 11C173 13 181 10 189 11"
        stroke="#000"
        strokeWidth="1.15"
        strokeLinecap="round"
        fill="none"
        opacity="0.72"
      />
    </svg>
  )
}

function StatChip({
  label,
  value,
}: {
  label?: ReactNode
  value?: string | number
}) {
  return (
    <div className="py-1 flex items-end justify-center gap-1">
      {label && (
        <div className="text-[9px] uppercase tracking-wide text-stone-800">
          {label}
        </div>
      )}
      {value && (
        <div className="text-xs font-semibold leading-tight text-stone-900">
          {value}
        </div>
      )}
    </div>
  )
}

export { ActionCard }

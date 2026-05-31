import { getEffectiveness, NATURES } from '#/lib/game/nature'
import { memo } from 'react'
import { NatureBadge } from './nature-badge'
import { cn } from '#/lib/utils'

const EFFECTIVENESS_DATA: Record<
  number,
  { color: string; size: string; label: string }
> = {
  2: { color: 'bg-green-600', size: 'size-3', label: 'Super Effective (2x)' },
  1.25: { color: 'bg-green-300', size: 'size-2', label: 'Effective (1.25x)' },
  1: { color: 'bg-stone-600', size: 'size-2', label: 'Neutral (1x)' },
  0.8: { color: 'bg-red-300', size: 'size-2', label: 'Resisted (0.8x)' },
  0.5: { color: 'bg-red-400', size: 'size-3', label: 'Highly Resisted (0.5x)' },
}

const NatureEffectivenessTable = memo(function NatureEffectivenessTable() {
  return (
    <div className="flex flex-col gap-4 pt-4">
      <div className="grid grid-cols-10 gap-1">
        <div />
        {NATURES.map((n) => (
          <div key={n} className="grid place-items-center">
            <NatureBadge nature={n} />
          </div>
        ))}
        {NATURES.map((a) => (
          <div key={a} className="grid grid-cols-10 col-span-10 gap-1">
            <div className="grid place-items-center">
              <NatureBadge nature={a} />
            </div>
            {NATURES.map((b) => {
              const eff = getEffectiveness(a, [b])
              const config = EFFECTIVENESS_DATA[eff]
              return (
                <div key={b} className="grid place-items-center h-6">
                  {config && (
                    <div
                      className={cn('rounded-full', config.color, config.size)}
                      title={`${a} vs ${b}: ${config.label}`}
                    />
                  )}
                </div>
              )
            })}
          </div>
        ))}
      </div>

      <div className="grid grid-cols-2 gap-4 text-[10px] text-muted-foreground border-t border-white/5 pt-2">
        {Object.entries(EFFECTIVENESS_DATA)
          .sort((a, b) => Number(b[0]) - Number(a[0]))
          .map(([val, config]) => (
            <div key={val} className="flex items-center gap-1.5">
              <div className={cn('rounded-full', config.color, config.size)} />
              <span>{config.label}</span>
            </div>
          ))}
      </div>
    </div>
  )
})

export { NatureEffectivenessTable }

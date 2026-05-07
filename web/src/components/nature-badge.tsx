import { cva } from 'class-variance-authority'
import type { ClassValue } from 'class-variance-authority/types'
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip'
import {
  getResistance,
  getWeakness,
  natureNames,
  type NatureSet,
} from '#/lib/game/nature'
import { cn } from '#/lib/utils'
import { ChevronRight } from 'lucide-react'

type t = Record<string, Partial<Record<NatureSet | 'none', ClassValue>>>

const variants = cva<t>(
  'text-sm text-stone-100 px-1 py-0.5 rounded text-shadow-[1px_1px_0px_#000000] text-nowrap cursor-default',
  {
    variants: {
      variant: {
        none: '',
        tai: 'bg-olive-500',
        pure: 'bg-slate-500 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
        fire: 'bg-red-500',
        wind: 'bg-emerald-700',
        lightning: 'bg-yellow-400',
        earth: 'bg-taupe-600',
        water: 'bg-blue-500',
        yang: 'bg-stone-300',
        yin: 'bg-indigo-900',
        ice: 'bg-cyan-700 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
        explosion:
          'bg-rose-900 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
        magnet: 'bg-violet-800 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
        storm:
          'bg-blue-900 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
        wood: 'bg-olive-600 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
        yinyang:
          'bg-[linear-gradient(135deg,theme(colors.indigo.900)_0%,theme(colors.indigo.900)_50%,theme(colors.neutral.300)_50%,theme(colors.neutral.300)_100%)] text-amber-300! shadow-[inset_0_0_0_1px_theme(colors.amber.300)]',
        particle:
          'bg-mauve-500 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300',
        jashin:
          'bg-mauve-950 shadow-[inset_0_0_0_1px_theme(colors.amber.300)] text-amber-300!',
      },
    },
    defaultVariants: {
      variant: 'none',
    },
  }
)

function NatureBadge({
  nature,
  className,
  children,
  ...props
}: React.ComponentProps<'span'> & { nature: NatureSet }) {
  if (children) {
    return (
      <span
        data-role="nature"
        className="inline-block shadow-[1px_1px_0_rgba(0,0,0,1)] mx-px rounded"
        {...props}
      >
        <span className={cn(variants({ variant: nature }), className)}>
          {children}
        </span>
      </span>
    )
  }

  const weaknesses = getWeakness(nature)
  const resistances = getResistance(nature)

  return (
    <Tooltip delayDuration={200}>
      <TooltipTrigger asChild>
        <span
          data-role="nature"
          className="shadow-[1px_1px_0_rgba(0,0,0,1)] mx-px rounded"
          {...props}
        >
          <span className={cn(variants({ variant: nature }), className)}>
            {natureNames[nature] ?? nature}
          </span>
        </span>
      </TooltipTrigger>
      <TooltipContent className="flex items-center">
        {!!weaknesses.length && (
          <div className="flex items-center">
            {weaknesses.map((n) => (
              <NatureBadge key={n} nature={n} />
            ))}
            <ChevronRight className="size-3" />
          </div>
        )}
        <NatureBadge nature={nature}>{nature}</NatureBadge>
        {!!resistances.length && (
          <div className="flex items-center">
            <ChevronRight className="size-3" />
            {resistances.map((n) => (
              <NatureBadge key={n} nature={n} />
            ))}
          </div>
        )}
      </TooltipContent>
    </Tooltip>
  )
}

export { NatureBadge }

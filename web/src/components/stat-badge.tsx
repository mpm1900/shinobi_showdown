import { type ActorNatureStat } from '#/lib/game/actor'
import { cn } from '#/lib/utils'
import { cva } from 'class-variance-authority'
import type { ClassValue } from 'class-variance-authority/types'
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip'
import type { ComponentProps } from 'react'

type t = Record<string, Partial<Record<ActorNatureStat | 'none', ClassValue>>>
const variants = cva<t>(
  'py-0.5 px-1 mx-px border border-black rounded drop-shadow-[1px_1px_0px_rgba(0,0,0,1)]',
  {
    variants: {
      variant: {
        attack:
          'shadow-[inset_0_0_0_1px_theme(colors.orange.900)] text-orange-400',
        defense: 'shadow-[inset_0_0_0_1px_theme(colors.red.900)] text-red-400',
        chakra_attack:
          'shadow-[inset_0_0_0_1px_theme(colors.indigo.900)] text-indigo-400',
        chakra_defense:
          'shadow-[inset_0_0_0_1px_theme(colors.blue.900)] text-blue-400',
        speed:
          'shadow-[inset_0_0_0_1px_theme(colors.emerald.900)] text-emerald-400',
      },
    },
  }
)

function StatBadge({
  className,
  contentProps = {},
  stat,
  ...props
}: React.ComponentProps<'span'> & {
  stat: ActorNatureStat
  contentProps?: Partial<ComponentProps<typeof TooltipContent>>
}) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span
          data-role="stat"
          className={cn(variants({ variant: stat }), className)}
          {...props}
        >
          {stat === 'attack' && 'A'}
          {(stat === 'defense' || stat === 'chakra_defense') && 'D'}
          {stat === 'chakra_attack' && 'C'}
          {stat === 'speed' && 'S'}
        </span>
      </TooltipTrigger>
      <TooltipContent {...contentProps}>
        {stat.replace('_', ' ')}
      </TooltipContent>
    </Tooltip>
  )
}

export { StatBadge }

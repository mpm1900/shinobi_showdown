import type React from 'react'
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip'
import type { Modifier } from '#/lib/game/modifier'
import { cn } from '#/lib/utils'

function ModifierTooltip({
  modifier,
  contentProps = {},
  children,
  description,
  ...triggerProps
}: Omit<React.ComponentProps<typeof TooltipTrigger>, 'asChild'> & {
  modifier: Modifier
  description?: (modifier: Modifier) => string
  contentProps?: React.ComponentProps<typeof TooltipContent>
}) {
  const {
    className: contentClassName,
    arrowClassName: contentArrowClassName,
    ...restContentProps
  } = contentProps

  return (
    <Tooltip delayDuration={100}>
      <TooltipTrigger asChild {...triggerProps}>
        {children}
      </TooltipTrigger>
      <TooltipContent
        sideOffset={8}
        collisionPadding={8}
        showArrow
        {...restContentProps}
        arrowClassName={cn('bg-popover fill-popover', contentArrowClassName)}
        className={cn(
          'z-50 w-64 rounded-md border bg-popover p-2 text-popover-foreground shadow-md text-wrap',
          contentClassName
        )}
      >
        <div className="flex justify-between">
          <div className="text-sm font-medium w-full">{modifier.name}</div>
          {(modifier.duration ?? 0) > 0 && (
            <div className="text-muted-foreground text-[10px] text-nowrap">
              {modifier.duration} Turns
            </div>
          )}
        </div>
        <div className="text-muted-foreground text-xs w-full break-words">
          {description ? description(modifier) : modifier.description}
        </div>
      </TooltipContent>
    </Tooltip>
  )
}

export { ModifierTooltip }

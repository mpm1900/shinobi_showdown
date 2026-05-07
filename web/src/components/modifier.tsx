import { MODIFIER_ICONS } from "#/data/icons"
import type { Modifier as ModifierType } from "#/lib/game/modifier"
import { cn } from "#/lib/utils"
import { ModifierTooltip } from "./modifier-tooltip"

function Modifier({
  count,
  modifier,
}: {
  count: number
  modifier: ModifierType | undefined
}) {
  if (!modifier || !modifier.show) return null
  const Icon = modifier?.icon ? MODIFIER_ICONS[modifier.icon] : undefined
  return (
    <ModifierTooltip modifier={modifier} contentProps={{ sideOffset: 0 }}>
      <span className="text-shadow-[1px_1px_0px_#000000] cursor-default">
        {Icon ? (
          <Icon
            className={cn(
              'size-6 my-1 drop-shadow-[1px_2px_0px_rgba(0,0,0,1)]',
              modifier?.icon
            )}
          />
        ) : (
          modifier?.name
        )}
        {count > 1 && !Icon ? ` (${count})` : null}
      </span>
    </ModifierTooltip>
  )
}

export { Modifier }

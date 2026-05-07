import { Field, FieldContent, FieldLabel } from './ui/field'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select'
import type { Modifier } from '#/lib/game/modifier'
import { ModifierTooltip } from './modifier-tooltip'
import { MODIFIER_ICONS } from '#/data/icons'

function AbilitySelect({
  value,
  onValueChange,
  options,
}: {
  value: string | null
  onValueChange: (value: string) => void
  options: Array<Modifier>
}) {
  return (
    <Field>
      <FieldLabel>Ability</FieldLabel>
      <FieldContent>
        <Select value={value ?? ''} onValueChange={onValueChange}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select an Ability " />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              {options.map((ability) => {
                const Icon = ability.icon
                  ? MODIFIER_ICONS[ability.icon]
                  : undefined
                return (
                  <ModifierTooltip
                    key={ability.ID}
                    modifier={ability}
                    contentProps={{ side: 'right' }}
                    description={(m) => m.parent_description || m.description}
                  >
                    <SelectItem value={ability.ID}>
                      {Icon && <Icon />}
                      {ability.name}
                    </SelectItem>
                  </ModifierTooltip>
                )
              })}
            </SelectGroup>
          </SelectContent>
        </Select>
      </FieldContent>
    </Field>
  )
}

export { AbilitySelect }

import { itemsQuery } from '#/lib/queries/items'
import { useQuery } from '@tanstack/react-query'
import { Field, FieldContent, FieldLabel } from './ui/field'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select'
import { ModifierTooltip } from './modifier-tooltip'
import { MODIFIER_ICONS } from '#/data/icons'

function ItemSelect({
  otherItemIDs,
  value,
  onValueChange,
}: {
  otherItemIDs: string[]
  value: string | null
  onValueChange: (value: string | null) => void
}) {
  const query = useQuery(itemsQuery)
  return (
    <Field>
      <FieldLabel>Item</FieldLabel>
      <FieldContent>
        <Select value={value ?? ''} onValueChange={onValueChange}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select an Item " />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              {query.data?.map((item) => {
                const Icon = item.icon ? MODIFIER_ICONS[item.icon] : undefined
                const used = otherItemIDs.includes(item.ID)
                return (
                  <ModifierTooltip
                    key={item.ID}
                    modifier={item}
                    contentProps={{ side: 'right' }}
                  >
                    <SelectItem
                      value={item.ID}
                      disabled={used}
                    >
                      {Icon && <Icon />}
                      {item.name}
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

export { ItemSelect }

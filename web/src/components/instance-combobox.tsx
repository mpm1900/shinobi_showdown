import { instancesQuery } from '#/lib/queries/instances'
import { cn } from '#/lib/utils'
import { useQuery } from '@tanstack/react-query'
import type { ReactNode } from 'react'
import { v4 } from 'uuid'
import { Button } from './ui/button'
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
  ComboboxTrigger,
  ComboboxValue,
} from './ui/combobox'

const CREATE_INSTANCE_VALUE = '__create_instance__'

function InstanceCombobox({
  icon,
  value = null,
  onValueChange,
}: {
  icon: ReactNode
  value: string | undefined | null
  onValueChange: (value: string) => void
}) {
  const query = useQuery(instancesQuery)
  const instanceItems = query.data ?? []
  const hasSelectedValue =
    !!value && instanceItems.some((instance) => instance.ID === value)
  const items = [
    ...(hasSelectedValue || !value ? [] : [{ ID: value }]),
    ...instanceItems,
    { ID: CREATE_INSTANCE_VALUE },
  ]

  return (
    <Combobox
      items={items}
      value={value}
      onValueChange={(v) => {
        if (!v) {
          return
        }
        if (v === CREATE_INSTANCE_VALUE) {
          onValueChange(v4())
        } else {
          onValueChange(v)
        }
        query.refetch()
      }}
    >
      <ComboboxTrigger
        render={
          <Button
            variant="outline"
            className={cn('min-w-80 justify-start font-normal font-mono', {
              'text-destructive': !value,
            })}
          >
            {icon}
            <ComboboxValue placeholder="NOT CONNECTED" />
          </Button>
        }
      />
      <ComboboxContent>
        <ComboboxInput showTrigger={false} placeholder="Search" />
        <ComboboxEmpty>No items found.</ComboboxEmpty>
        <ComboboxList>
          {(item) => (
            <ComboboxItem
              key={item.ID}
              value={item.ID}
              disabled={Object.keys(item.clients ?? {}).length >= 2}
              className={cn({
                "[slot='combobox-item-indicator']:hidden": !value,
              })}
            >
              {item.ID === CREATE_INSTANCE_VALUE
                ? 'Create Instance'
                : `${item.ID} (${Object.keys(item.clients ?? {}).length})`}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  )
}

export { InstanceCombobox }

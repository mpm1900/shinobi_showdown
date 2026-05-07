import { createServerFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
import type { Modifier } from '../game/modifier'

const getItems = createServerFn().handler(async () => {
  const response = await fetch(`${process.env.API_URL}/api/items`)
  const data = await response.json()
  return data as Array<Modifier>
})

const itemsQuery = queryOptions({
  queryKey: ['items'],
  queryFn: () => getItems(),
  select: (items) => items.sort((a, b) => a.name.localeCompare(b.name)),
  staleTime: 60000,
  gcTime: 60000,
})

export { itemsQuery, getItems }

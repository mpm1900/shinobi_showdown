import { getTotalBaseStats, type ActorDef } from '#/lib/game/actor'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getExpandedRowModel,
  getSortedRowModel,
  useReactTable,
  type Row,
  type SortingState,
} from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from './ui/table'
import { ActorStatBase } from './actor-stat'
import { Button } from './ui/button'
import { ChevronDown, ChevronLeft } from 'lucide-react'
import { Fragment, useState, type ReactNode } from 'react'
import { natureIndexes, type NatureSet } from '#/lib/game/nature'
import { NatureBadge } from './nature-badge'
import { SHINOBI_ICONS } from '#/data/icons'
import { StatBadge } from './stat-badge'

const helper = createColumnHelper<ActorDef>()
const columns = [
  helper.accessor('name', {}),
  helper.accessor('clan', {
    header: '',
    cell: ({ row }) => (
      <div className="flex">
        {[row.original.clan].map((a) => {
          const C = SHINOBI_ICONS[a]
          return C ? <C key={a} className="w-5 text-stone-100" /> : null
        })}
      </div>
    ),
  }),
  helper.accessor('affiliations', {
    header: '',
    cell: ({ row }) => (
      <div className="flex gap-2 justify-end">
        {row.original.affiliations?.map((a) => {
          const C = SHINOBI_ICONS[a]
          return C ? <C key={a} className="w-5" /> : null
        })}
      </div>
    ),
  }),
  helper.accessor('natures', {
    cell: ({ row }) =>
      (Object.keys(row.getValue('natures')) as Array<NatureSet>)
        .sort((a, b) => natureIndexes[a] - natureIndexes[b])
        .map((nature) => <NatureBadge key={nature} nature={nature} />),
  }),
  helper.accessor('stats.hp', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        HP
      </Button>
    ),
    cell: (props) => <ActorStatBase actor={props.row.original} stat="hp" />,
    sortingFn: 'alphanumeric',
  }),
  helper.accessor('stats.stamina', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Stamina
      </Button>
    ),
    cell: (props) => (
      <ActorStatBase actor={props.row.original} stat="stamina" />
    ),
  }),
  helper.accessor('stats.speed', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Speed
      </Button>
    ),
    cell: (props) => <ActorStatBase actor={props.row.original} stat="speed" />,
  }),
  helper.accessor('stats.attack', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        <StatBadge stat="attack" />
        ATK
      </Button>
    ),
    cell: (props) => <ActorStatBase actor={props.row.original} stat="attack" />,
  }),
  helper.accessor('stats.defense', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        <StatBadge stat="defense" />
        DEF
      </Button>
    ),
    cell: (props) => (
      <ActorStatBase actor={props.row.original} stat="defense" />
    ),
  }),
  helper.accessor('stats.chakra_attack', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        <StatBadge stat="chakra_attack" />
        ATK
      </Button>
    ),
    cell: (props) => (
      <ActorStatBase actor={props.row.original} stat="chakra_attack" />
    ),
  }),
  helper.accessor('stats.chakra_defense', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        <StatBadge stat="chakra_defense" />
        DEF
      </Button>
    ),
    cell: (props) => (
      <ActorStatBase actor={props.row.original} stat="chakra_defense" />
    ),
  }),
  helper.accessor((a) => getTotalBaseStats(a), {
    id: 'total',
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        total
      </Button>
    ),
  }),
  helper.display({
    id: 'actions',
    cell: ({ row }) => (
      <Button
        size="icon"
        variant="ghost"
        disabled={!row.getIsSelected()}
        onClick={row.getToggleExpandedHandler()}
      >
        {row.getIsExpanded() ? <ChevronDown /> : <ChevronLeft />}
      </Button>
    ),
  }),
]

function ActorsTable({
  data,
  enabled,
  subRow,
}: {
  data: Array<ActorDef>
  enabled: boolean
  subRow?: (props: { row: Row<ActorDef> }) => ReactNode
}) {
  const [sorting, setSorting] = useState<SortingState>([
    { id: 'total', desc: true },
  ])
  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getRowCanExpand: (row) => row.getIsSelected(),
    getRowId: (a) => a.actor_ID,
    getSortedRowModel: getSortedRowModel(),
    enableRowSelection: enabled,
    // onExpandedChange: setExpanded,
    onSortingChange: setSorting,
    state: {
      sorting,
    },
  })

  return (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map((hg) => (
          <tr key={hg.id}>
            {hg.headers.map((header) => (
              <TableHead key={header.id} colSpan={header.colSpan}>
                {flexRender(
                  header.column.columnDef.header,
                  header.getContext()
                )}
              </TableHead>
            ))}
          </tr>
        ))}
      </TableHeader>
      <TableBody>
        {table.getRowModel().rows.map((row) => (
          <Fragment key={row.id}>
            <TableRow>
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
            {row.getIsExpanded() && subRow && (
              <tr>
                <td colSpan={row.getAllCells().length}>{subRow({ row })}</td>
              </tr>
            )}
          </Fragment>
        ))}
      </TableBody>
    </Table>
  )
}

export { ActorsTable }

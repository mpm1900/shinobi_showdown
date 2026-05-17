import type { Action } from '#/lib/game/action'
import {
  createColumnHelper,
  flexRender,
  functionalUpdate,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
  type Row,
  type RowSelectionState,
} from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from './ui/table'
import { NatureBadge } from './nature-badge'
import { useState } from 'react'
import { StatBadge } from './stat-badge'
import { Checkbox } from './ui/checkbox'
import { Button } from './ui/button'
import { cn } from '#/lib/utils'
import { HoverCard, HoverCardContent, HoverCardTrigger } from './ui/hover-card'
import { ActionCard } from './action-card'

const helper = createColumnHelper<Action>()
const columns = [
  helper.display({
    id: 'select',
    header: ({ table }) =>
      `${table.getSelectedRowModel().rows.length}/${(table.options.meta as any).total}`,
    cell: ({ row, table }) => (
      <Checkbox
        checked={row.getIsSelected()}
        disabled={
          !row.getIsSelected() &&
          (row.original.state.locked ||
            !row.getCanSelect() ||
            (table.options.meta as any).total ==
            table.getSelectedRowModel().rows.length)
        }
      />
    ),
  }),
  helper.accessor('config.name', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Name
      </Button>
    ),
  }),
  helper.accessor('config.nature', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Nature
      </Button>
    ),
    cell: ({ row }) =>
      row.original.config.nature ? (
        <NatureBadge nature={row.original.config.nature} />
      ) : (
        '-'
      ),
  }),
  helper.accessor('config.stat', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Stat
      </Button>
    ),
    cell: ({ row }) =>
      row.original.config.stat ? (
        <StatBadge
          stat={row.original.config.stat}
          contentProps={{ side: 'right' }}
        />
      ) : (
        '-'
      ),
  }),
  helper.accessor('config.power', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Power
      </Button>
    ),
    cell: ({ row }) => row.original.config.power ?? '-',
  }),
  helper.accessor('config.accuracy', {
    header: ({ column }) => (
      <Button
        className="-ml-4"
        variant="ghost"
        onClick={() => column.toggleSorting()}
      >
        Acc
      </Button>
    ),
    cell: ({ row }) =>
      row.original.config.accuracy ? `${row.original.config.accuracy}%` : '-',
  }),
  helper.accessor('config.description', {
    id: 'description',
    header: 'Description',
    cell: ({ row }) => (
      <span className="block truncate">{row.original.config.description}</span>
    ),
  }),
]

function ActionsTable({
  total,
  data,
  rowSelection,
  onRowSelectionChange,
}: {
  total: number
  data: Action[]
  rowSelection: RowSelectionState
  onRowSelectionChange: (rowSelection: RowSelectionState) => void
}) {
  const [sorting, setSorting] = useState([{ id: 'config_name', desc: false }])

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    enableRowPinning: true,
    enableRowSelection: total > Object.keys(rowSelection).length,
    onRowSelectionChange: (updater) => {
      onRowSelectionChange(functionalUpdate(updater, rowSelection))
    },
    onSortingChange: (updater) => {
      setSorting(functionalUpdate(updater, sorting))
    },
    getRowId: (a) => a.ID,
    keepPinnedRows: true,
    state: {
      rowPinning: {
        top: Object.entries(rowSelection)
          .filter(([, selected]) => !!selected)
          .map(([id]) => id),
      },
      rowSelection,
      sorting,
    },
    meta: {
      total,
    },
  })

  const renderRow = (row: Row<Action>) => (
    <HoverCard key={row.id}>
      <HoverCardTrigger asChild>
        <TableRow
          className={cn(row.getIsPinned() && 'bg-muted/50', 'cursor-default')}
          onClick={() => row.toggleSelected()}
        >
          {row.getVisibleCells().map((cell) => (
            <TableCell
              key={cell.id}
              className={
                cell.column.id === 'description' ? 'w-full max-w-0' : ''
              }
            >
              {flexRender(cell.column.columnDef.cell, cell.getContext())}
            </TableCell>
          ))}
        </TableRow>
      </HoverCardTrigger>
      <HoverCardContent
        side="right"
        collisionPadding={16}
        sideOffset={8}
        className="p-0 border-0 bg-transparent"
      >
        <ActionCard action={row.original} />
      </HoverCardContent>
    </HoverCard>
  )

  return (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map((hg) => (
          <tr key={hg.id}>
            {hg.headers.map((header) => (
              <TableHead
                key={header.id}
                colSpan={header.colSpan}
                className={cn(
                  'sticky top-0 z-10 bg-stone-900',
                  header.column.id === 'description' ? 'w-full' : ''
                )}
              >
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
        {table.getTopRows().map(renderRow)}
        {table.getCenterRows().map(renderRow)}
      </TableBody>
    </Table>
  )
}

export { ActionsTable }

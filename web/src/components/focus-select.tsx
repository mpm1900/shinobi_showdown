import {
  ACTOR_FOCUS_DETAILS,
  actorFocuses,
  type ActorFocus,
  type ActorNatureStat,
} from '#/lib/game/actor'
import { Field, FieldContent, FieldLabel } from './ui/field'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select'

function parse(stat: ActorNatureStat | undefined | null): string {
  if (!stat) return ''
  return stat
    .replaceAll('_', ' ')
    .replaceAll('chakra', 'c.')
    .replace('attack', 'atk')
    .replace('defense', 'def')
    .replace('speed', 'spd')
}

function FocusSelectItem({ focus }: { focus: (typeof actorFocuses)[number] }) {
  const obj = ACTOR_FOCUS_DETAILS[focus]
  const up = parse(obj.up)
  const down = parse(obj.down)
  return (
    <SelectItem value={focus}>
      <span className="capitalize">{focus}</span>{' '}
      <span className="opacity-40 capitalize">
        (+{up}, -{down})
      </span>
    </SelectItem>
  )
}

function FocusSelect({
  value,
  onValueChange,
  ...props
}: React.ComponentProps<typeof Field> & {
  value: ActorFocus
  onValueChange: (value: ActorFocus) => void
}) {
  return (
    <Field {...props}>
      <FieldLabel>Focus</FieldLabel>
      <FieldContent>
        <Select value={value} onValueChange={onValueChange}>
          <SelectTrigger className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              {actorFocuses.map((focus) => (
                <FocusSelectItem key={focus} focus={focus} />
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
      </FieldContent>
    </Field>
  )
}

export { FocusSelect }

import { ACTOR_FOCUS_DETAILS, actorFocuses, type ActorFocus } from '#/lib/game/actor'
import { Field, FieldContent, FieldLabel } from './ui/field'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select'

function FocusSelectItem({ focus }: { focus: (typeof actorFocuses)[number] }) {
  const obj = ACTOR_FOCUS_DETAILS[focus]
  const up = obj.up?.replaceAll('_', ' ')
  const down = obj.down?.replaceAll('_', ' ')
  return (
    <SelectItem value={focus}>
      <span className='capitalize'>{focus}</span> <span className='opacity-40 capitalize'>(+{up}, -{down})</span>
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

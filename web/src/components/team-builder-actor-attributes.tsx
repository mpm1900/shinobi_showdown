import type { ActorDef, ActorFocus } from '#/lib/game/actor'
import { AbilitySelect } from './ability-select'
import { FocusSelect } from './focus-select'
import { ItemSelect } from './item-select'
import { Modifier } from './modifier'
import { Field, FieldContent, FieldLabel } from './ui/field'

function TeamBuilderActorAttributes({
  def,
  otherItemIDs,
  focus,
  onFocusChange,
  abilityID,
  onAbilityIDChange,
  itemID,
  onItemIDChange,
}: {
  def: ActorDef
  otherItemIDs: string[]
  focus: ActorFocus
  onFocusChange: (f: ActorFocus) => void
  abilityID: string | undefined | null
  onAbilityIDChange: (a: string) => void
  itemID: string | undefined | null
  onItemIDChange: (i: string | null) => void
}) {
  return (
    <div className="space-y-2">
      <FocusSelect value={focus ?? 'none'} onValueChange={onFocusChange} />
      <AbilitySelect
        options={def.abilities}
        value={abilityID ?? null}
        onValueChange={onAbilityIDChange}
      />
      <ItemSelect otherItemIDs={otherItemIDs} value={itemID ?? null} onValueChange={onItemIDChange} />
      <Field>
        <FieldLabel>Default Modifiers:</FieldLabel>
        <FieldContent className='flex items-start'>
          {def.default_modifiers?.map(mod => (
            <Modifier key={mod.ID} count={0} modifier={mod} />
          )) ?? <span className='text-xs'>None</span>}
        </FieldContent>
      </Field>
    </div>
  )
}

export { TeamBuilderActorAttributes }

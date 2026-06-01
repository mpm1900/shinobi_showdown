import { gameStore } from '#/lib/stores/game'
import { setBgEnabled, uiStore } from '#/lib/stores/ui'
import { useStore } from '@tanstack/react-store'
import { Modifier } from './modifier'
import { Switch } from './ui/switch'

function BattleWeather() {
  const ui = useStore(uiStore, (s) => s)
  const modifiers = useStore(gameStore, (g) =>
    g.modifiers.filter((tx) => g.applied_game_state_tx.includes(tx.ID))
  )
  return (
    <div className="space-y-2 mt-4">
      <Switch
        className="hidden xl:block"
        checked={ui.bgEnabled}
        onCheckedChange={setBgEnabled}
      />
      <div>
        {modifiers.map((tx) => (
          <Modifier key={tx.ID} modifier={tx.mutation} count={1} />
        ))}
      </div>
    </div>
  )
}

export { BattleWeather }

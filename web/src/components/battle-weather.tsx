import { gameStore } from "#/lib/stores/game"
import { useStore } from "@tanstack/react-store"
import { Modifier } from "./modifier"

function BattleWeather() {
  const state = useStore(gameStore, g => g.state)
  const modifiers = useStore(gameStore, g => g.modifiers
    .filter((tx) =>
      g.applied_game_state_tx.includes(tx.ID)
    ))
  return (
    <div>
      <div>Weather: {state.weather}</div>
      <div>Terrain: {state.terrain}</div>
      <div>
        {modifiers
          .map((tx) => <Modifier key={tx.ID} modifier={tx.mutation} count={1} />)}
      </div>
    </div>
  )
}


export { BattleWeather }

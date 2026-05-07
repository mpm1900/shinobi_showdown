import { Store } from '@tanstack/store'
import type { Game } from '../game/game'

const gameStore = new Store<Game>({
  status: 'init',
  turn: {
    count: 0,
    phase: 'init'
  },
  actors: [],
  actions: [],
  modifiers: [],
  players: [],
  prompt: null,
  log: [],
  active_transaction: null,
  queued_actions: {},
  applied_game_state_tx: [],
  state: {
    terrain: 'none',
    weather: 'none'
  }
})



export { gameStore }

import { Store } from '@tanstack/store'

const uiStore = new Store({
  bgEnabled: true,
})

function toggleBg() {
  uiStore.setState((s) => ({ bgEnabled: !s.bgEnabled }))
}
function setBgEnabled(enabled: boolean) {
  uiStore.setState((prev) => ({ bgEnabled: enabled }))
}

export { uiStore, toggleBg, setBgEnabled }

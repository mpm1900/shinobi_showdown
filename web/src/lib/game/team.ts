import type { TeamActor } from "../stores/config";
import type { ActorDef } from "./actor";

function makeConfigFromDef(def: ActorDef): TeamActor {
  return {
    actor_ID: def.actor_ID,
    config: {
      ability_ID: def.abilities[0]?.ID ?? null,
      item_ID: null,
      action_IDs: [],
      focus: 'none',
      aux_stats: {
        hp: 0,
        stamina: 0,
        speed: 0,
        attack: 0,
        defense: 0,
        chakra_attack: 0,
        chakra_defense: 0,
      },
    },
  }
}

export { makeConfigFromDef }

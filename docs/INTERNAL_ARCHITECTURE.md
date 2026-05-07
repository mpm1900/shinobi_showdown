# System Architecture & Core Logic Findings

## 1. The Game Loop (`internal/game/game_next.go`)
The game engine operates on a prioritized "tick" system within the `Next()` function. The priority of execution is as follows:
1.  **Transactions:** Immediate state changes (animations, damage application).
2.  **Prompts:** Blocking user inputs (e.g., choosing a switch-in).
3.  **Triggers:** Reactive logic (e.g., "OnDeath", "OnActionStart").
4.  **Validation:** Integrity checks (cleaning up dead actors, checking for game over).
5.  **Actions:** The primary turn-based moves (attacks, items).

## 2. Transactional State Management
The game uses a "Command" pattern for state mutations.
- `GameMutation` defines the logic.
- `Transaction[T]` wraps a mutation with a `Context`.
- `JumpTransaction` allows inserting high-priority mutations at the front of the queue (used for immediate cleanups).

## 3. Actor Resolution & "Layers"
Actors have three distinct states:
- `ActorDef`: Static data (base stats, natures).
- `ActorState`: Runtime volatile data (damage, current position, status effects).
- `ResolvedActor`: The calculated result of applying all active `Modifiers` and `Stages` to the base stats.
**Key Finding:** Most game logic should use `ResolvedActor` to account for buffs/debuffs, but mutations usually update the underlying `ActorState`.

## 4. Contextual Attribution
The `Context` struct is critical for the "clean-up" system.
- `SourceActorID`: The actor performing the action.
- `ParentActorID`: The actor "responsible" for the transaction.
When an actor leaves the battle, the system calls `FilterParentModifiers` and `FilterParentActions`. If `ParentActorID` is not correctly set during action creation, these clean-up functions will fail to remove orphaned actions/modifiers.

## 5. Event System (`internal/game/trigger.go`)
The system uses an `On(event, context)` pattern.
- Common events: `OnActionStart`, `OnActionEnd`, `OnDamageReceive`, `OnDeath`.
- Triggers are gathered from all active modifiers and queued for resolution.

## 6. Instance Reducer (`internal/instance/reducer.go`)
The `Reducer` acts as the bridge between the WebSocket/Network layer and the Game logic. It handles:
- Validation of client requests.
- Hydration of actors from static data.
- Execution of the game loop (`RunGameActions`) until it hits a blocking state (Prompt or Idle).

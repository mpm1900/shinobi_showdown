# Engineering Context Dump - April 27, 2026

## Task: Fix Action Sequencing for Defeated Actors

### The Bug
Actors that were defeated in the middle of a turn were still attempting to execute their queued actions before leaving the battle. This resulted in logs like:
1. Actor A uses attack.
2. Actor B loses 100% HP.
3. Actor B uses attack.
4. Actor B attack fails.
5. Actor B leaves the battle.

### Root Causes
1.  **Bugged Filter:** `FilterParentActions` in `internal/game/game.go` had a logic error. It only filtered actions where `parent != nil && *parent != actorID`, which effectively kept the actions of the actor being filtered and discarded everyone else's.
2.  **Validation Flow:** `Validate()` in `internal/game/game_next.go` was detecting dead actors and jumping a `RemovePositions` transaction, but it was returning `true`, allowing the game loop to proceed immediately to `PreAction` and `NextAction` before the jumped transaction could execute.
3.  **Missing Safeguards:** `PreAction` and `NextAction` did not explicitly check if the action's source was still alive before proceeding.

### Solutions Applied

#### 1. `internal/game/game.go`
Fixed `FilterParentActions` to correctly identify and remove actions belonging to the specified `actorID`.
```go
func (g *Game) FilterParentActions(actorID uuid.UUID) {
	actions := g.Actions[:0]
	for _, tx := range g.Actions {
		parent := tx.Context.ParentActorID
		if parent == nil || *parent != actorID {
			actions = append(actions, tx)
		}
	}
	g.Actions = actions
}
```

#### 2. `internal/game/game_next.go`
*   Updated `Validate()` to return `false` if a dead actor is detected. This breaks the current `Next()` cycle, allowing the jumped `RemovePositions` transaction to be processed in the next tick.
*   Added `source.Alive` checks to `PreAction()` and `NextAction()`.
*   `NextAction()` now returns `true` (indicating a successful "step") but skips execution if the source is dead, preventing "Action failed" logs for deceased actors.

### Architectural Insights
*   **Transaction Priority:** `JumpTransaction` prepends to the `Transactions` queue. For these to take effect before the next game phase or action, `Next()` must return early so `NextTransaction()` can catch them.
*   **Context Parents:** Actions should always have a `ParentActorID` set (usually to the `SourceActorID`) to ensure cleanup functions like `FilterParentActions` work correctly.
*   **Actor State:** `Actor.Alive` is the primary flag for presence in the active loop. `Validate` is responsible for cleaning up actors whose `Alive` flag is false but are still in a `Position`.

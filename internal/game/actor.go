package game

import (
	"maps"
	"math"
	"math/rand/v2"
	"slices"

	"github.com/google/uuid"
)

type AttackStat string

const (
	Attack       AttackStat = "attack"
	ChakraAttack AttackStat = "chakra_attack"
)

type DefenseStat string

const (
	Defense       DefenseStat = "defense"
	ChakraDefense DefenseStat = "chakra_defense"
)

type ActorStat string

const (
	StatHP            ActorStat = "hp"
	StatStamina       ActorStat = "stamina"
	StatAttack        ActorStat = ActorStat(Attack)
	StatDefense       ActorStat = ActorStat(Defense)
	StatChakraAttack  ActorStat = ActorStat(ChakraAttack)
	StatChakraDefense ActorStat = ActorStat(ChakraDefense)
	StatSpeed         ActorStat = "speed"
	StatEvasion       ActorStat = "evasion"
	StatAccuracy      ActorStat = "accuracy"
)

type ActorFocus string

const (
	FocusNone ActorFocus = "none"
	// +P.ATK
	FocusAggressive ActorFocus = "aggressive" // -P.DEF
	FocusRelentless ActorFocus = "relentless" // -C.ATK
	FocusReckless   ActorFocus = "reckless"   // -C.DEF
	FocusHeavy      ActorFocus = "heavy"      // -SPE
	// +P.DEF
	FocusPatient   ActorFocus = "patient"   // -P.ATK
	FocusHardened  ActorFocus = "hardened"  // -C.ATK
	FocusTough     ActorFocus = "tough"     // -C.DEF
	FocusSteadfast ActorFocus = "steadfast" // -SPE
	// +C.ATK
	FocusIntelligent ActorFocus = "intelligent" // -P.ATK
	FocusVolatile    ActorFocus = "volatile"    // -P.DEF
	FocusIntense     ActorFocus = "intense"     // -C.DEF
	FocusCalculated  ActorFocus = "calculated"  // -SPE
	// +P.DEF
	FocusComposed ActorFocus = "composed" // -P.ATK
	FocusMindful  ActorFocus = "mindful"  // -P.DEF
	FocusReserved ActorFocus = "reserved" // -C.ATK
	FocusedStoic  ActorFocus = "stoic"    // -SPE

	// +SPE
	FocusAgile     ActorFocus = "agile"     // -P.ATK
	FocusHasty     ActorFocus = "hasty"     // -P.DEF
	FocusImpulsive ActorFocus = "impulsive" // -C.ATK
	FocusAlert     ActorFocus = "alert"     // -C.DEF
)

type FocusStatDelta struct {
	Up   ActorStat
	Down ActorStat
}

var ActorFocuses = map[ActorFocus]FocusStatDelta{
	FocusNone:        {},
	FocusAggressive:  {Up: StatAttack, Down: StatDefense},
	FocusRelentless:  {Up: StatAttack, Down: StatChakraAttack},
	FocusReckless:    {Up: StatAttack, Down: StatChakraDefense},
	FocusHeavy:       {Up: StatAttack, Down: StatSpeed},
	FocusPatient:     {Up: StatDefense, Down: StatAttack},
	FocusHardened:    {Up: StatDefense, Down: StatChakraAttack},
	FocusTough:       {Up: StatDefense, Down: StatChakraDefense},
	FocusSteadfast:   {Up: StatDefense, Down: StatSpeed},
	FocusIntelligent: {Up: StatChakraAttack, Down: StatAttack},
	FocusVolatile:    {Up: StatChakraAttack, Down: StatDefense},
	FocusIntense:     {Up: StatChakraAttack, Down: StatChakraDefense},
	FocusCalculated:  {Up: StatChakraAttack, Down: StatSpeed},
	FocusComposed:    {Up: StatChakraDefense, Down: StatAttack},
	FocusMindful:     {Up: StatChakraDefense, Down: StatDefense},
	FocusReserved:    {Up: StatChakraDefense, Down: StatChakraAttack},
	FocusedStoic:     {Up: StatChakraDefense, Down: StatSpeed},
	FocusAgile:       {Up: StatSpeed, Down: StatAttack},
	FocusHasty:       {Up: StatSpeed, Down: StatDefense},
	FocusImpulsive:   {Up: StatSpeed, Down: StatChakraAttack},
	FocusAlert:       {Up: StatSpeed, Down: StatChakraDefense},
}

const (
	AffAme      = "ame"
	AffAkatsuki = "akatsuki"
	AffIwa      = "iwa"
	AffKonoha   = "konoha"
	AffKumo     = "kumo"
	AffKuri     = "kuri"
	AffOto      = "oto"
	AffSun      = "sun"
	AffTaki     = "taki"
	AffYuga     = "yuga"

	ClanAburame = "aburame"
	ClanHatake  = "hatake"
	ClanSenju   = "senju"
	ClanUchiha  = "uchiha"
	ClanUzumaki = "uzumaki"
)

type ResourceValues struct {
	Max    int
	Offset int
}

/**
 * [ActorDef]
 * Actor Definitions, the non-runtime definition of a shinobi's
 * stats, natures, and information
 */
type ActorDef struct {
	ActorID      uuid.UUID `json:"actor_ID"`
	SpriteURL    string    `json:"sprite_url"`
	Name         string    `json:"name"`
	Clan         string    `json:"clan"`
	Affiliations []string  `json:"affiliations"`
	Restricted   bool      `json:"restricted"`

	Stats            map[ActorStat]int      `json:"stats"`
	NatureDamage     map[Nature]float64     `json:"nature_damage"`
	NatureResistance map[Nature]float64     `json:"nature_resistance"`
	Natures          map[NatureSet][]Nature `json:"natures"`

	DefaultModifiers []Modifier  `json:"default_modifiers"`
	Abilities        []Modifier  `json:"abilities"`
	ActionIDs        []uuid.UUID `json:"action_IDs,omitempty"`
	ActionCount      int         `json:"action_count"`

	Immunities      map[uuid.UUID]struct{}   `json:"-"`
	JutsuImmunities map[ActionJutsu]struct{} `json:"-"`
}

type ActorStateType string

const (
	ActorStateFlying   ActorStateType = "flying"
	ActorStateGrounded ActorStateType = "grounded"
)

type ActorStance string

/**
 * [ActorState]
 * - the commonly modified fields
 */
type ActorState struct {
	State         ActorStateType `json:"state"`
	Stance        ActorStance    `json:"stance"`
	ActiveTurns   int            `json:"-"`
	InactiveTurns int            `json:"-"`
	Enters        int            `json:"-"`
	Seen          bool           `json:"seen"`
	// [Alive] whether or not the actor is alive, could
	// - could be computed, but this is here to not have to call .Resolve() on filters
	Alive bool `json:"alive"`
	// [Damage] how much damage this actor has received
	Damage        int `json:"damage"`
	StaminaDamage int `json:"stamina_damage"`
	// [PositionID] current position, nil if not active
	PositionID *uuid.UUID `json:"position_ID"`
	// [Protected]
	// - protected units cannot be damaged by actions
	Protected bool `json:"protected"`
	// [Safeguarded]
	// - safeguarded units cannot gain enemy modifiers
	Safeguarded bool `json:"safeguarded"`
	// [Warded]
	// - warded units are immune to the secondary effects of attacking actions
	Warded bool `json:"warded"`
	// [Reflect] how much damage is reflected (PureDamage not affected)
	IgnoreRedirect       bool    `json:"-"`
	Reflect              float64 `json:"-"`
	RecoilMultiplier     float64 `json:"-"`
	PowerMultiplier      float64 `json:"-"`
	StabMultiplier       float64 `json:"-"`
	ModifierChanceOffset int     `json:"-"`
	ModifierChanceMult   float64 `json:"-"`
	StatusChanceOffset   int     `json:"-"`
	StatusChanceMult     float64 `json:"-"`
	CooldownOffset       int     `json:"-"`
	RepeatsMinOffset     int     `json:"-"`
	RepeatsMaxOffset     int     `json:"-"`
	ActionAccuracyOffset int     `json:"-"`
	Immortal             bool    `json:"immortal"`
	// [ActionLocked]
	// - action locked units must use their last used action
	// - if there is no last used action, any action can be chosen
	ActionLocked bool `json:"action_locked"`
	// [SwitchLocked]
	// - switched locked units cannot use Switch to switch out, but other means are fine
	SwitchLocked bool `json:"switch_locked"`
	// [Stunned] whether or not an actor _can act_
	// - stunned units cannot push actions
	// - stunned units cannot resolve actions (if the status was added during running)
	Stunned bool `json:"stunned"`
	/**
	 * Statuses
	 * Made the choice for the core to not reference status modifiers by name,
	 *  rather, have both reference keys in ActorState below. For example,
	 *  ResolveAction will not check if the source has a modifier "Paralysis,"
	 *  It instead checks the flag set by Paralysis here.
	 */
	Statused        bool `json:"statused"`
	Burned          bool `json:"burned"`
	Paralyzed       bool `json:"paralyzed"`
	Poisoned        bool `json:"poisoned"`
	PoisonedCounter int  `json:"-"`
	Sleeping        bool `json:"sleeping"`
	SleepCounter    int  `json:"-"`

	/**
	 * Metadata fields used for tracking and filters
	 */
	LastUsedActionTX   *Transaction[Action] `json:"-"`
	LastReceivedDamage map[uuid.UUID]int    `json:"-"`
	HitCount           int                  `json:"-"`
}

type Summon struct {
	Actor
	Parent *Actor `json:"-"`
	Proxy  bool   `json:"proxy"`
}

type Actor struct {
	ActorDef
	ActorState
	ID         uuid.UUID  `json:"ID"`
	PlayerID   uuid.UUID  `json:"player_ID"`
	Level      int        `json:"level"`
	Experience int        `json:"experience"`
	Focus      ActorFocus `json:"focus"`
	Ability    *Modifier  `json:"ability"`
	Enabled    bool       `json:"enabled"`
	// [AuxAbility]
	// - take priority over Ability
	// - set to nil on switch-out
	AuxAbility        *Modifier              `json:"-"`
	Item              *Modifier              `json:"item"`
	Summon            *Summon                `json:"summon,omitempty"`
	Stages            map[ActorStat]int      `json:"staged_stats"`
	AuxStats          map[ActorStat]int      `json:"aux_stats"`
	DamageMultipliers map[AttackStat]float64 `json:"-"`
	DamageReduction   map[AttackStat]float64 `json:"-"`
	Actions           []Action               `json:"actions"`
	AppliedModifiers  map[uuid.UUID]int      `json:"applied_modifiers"`
}

type ResolvedActor struct {
	Actor
	BaseStats                map[ActorStat]int  `json:"base_stats"`
	UnmodifiedStats          map[ActorStat]int  `json:"unmodified_stats"`
	ResolvedNatureResistance map[Nature]float64 `json:"resolved_nature_resistance"`
	ResolvedNatureDamage     map[Nature]float64 `json:"resolved_nature_damage"`
}

const (
	MutPriorityGameState0      = -40
	MutPriorityImmunity        = -30
	MutPriorityStateChanges    = -20
	MutPriorityPreBaseStats    = -11
	MutPriorityMapBaseStats    = -10
	MutPriorityPostBaseStats   = -9
	MutPriorityDefault         = 0
	MutPriorityPreStagedStats  = 9
	MutPriorityMapStagedStats  = 10
	MutPriorityPostStagedStats = 11
	MutPriorityZero            = 19
	MutPrioritySet             = 20
	MutPriorityPostSet         = 21 // toad song
)

func GetLevel(experience int) int {
	return Round(math.Cbrt(float64(experience)))
}
func GetBaseExperience(level int) int {
	return Round(math.Pow(float64(level), 3))
}
func GetExperienceToNextLevel(level, exp int) int {
	return GetBaseExperience(level+1) - (GetBaseExperience(level) + exp)
}

func (a Actor) GetActionByID(g Game, actionID uuid.UUID) (Action, bool) {
	ra := a.Resolve(g)
	return ra.GetActionByID(actionID)
}
func (a Actor) IsActive() bool {
	return a.PositionID != nil
}
func (a Actor) GetAbility() *Modifier {
	if a.AuxAbility != nil {
		return a.AuxAbility
	}
	return a.Ability
}
func (ra ResolvedActor) GetActionByID(actionID uuid.UUID) (Action, bool) {
	for _, action := range ra.Actions {
		if action.ID == actionID {
			return action, true
		}
	}

	return Action{}, false
}

func makeActions(actionIDs []uuid.UUID, ACTIONS map[uuid.UUID]Action) []Action {
	actions := make([]Action, 0, len(actionIDs))
	for _, id := range actionIDs {
		a, ok := ACTIONS[id]
		if !ok {
			continue
		}
		actions = append(actions, a)
	}

	return actions
}

func (ad ActorDef) Clone() ActorDef {
	cloned := ad

	if ad.Immunities == nil {
		ad.Immunities = make(map[uuid.UUID]struct{})
	}

	if ad.DefaultModifiers == nil {
		ad.DefaultModifiers = make([]Modifier, 0)
	}

	cloned.Affiliations = slices.Clone(ad.Affiliations)
	cloned.Stats = maps.Clone(ad.Stats)
	cloned.NatureDamage = maps.Clone(ad.NatureDamage)
	cloned.NatureResistance = maps.Clone(ad.NatureResistance)

	cloned.Natures = make(map[NatureSet][]Nature, len(ad.Natures))
	for k, v := range ad.Natures {
		cloned.Natures[k] = slices.Clone(v)
	}

	cloned.DefaultModifiers = slices.Clone(ad.DefaultModifiers)
	cloned.ActionIDs = slices.Clone(ad.ActionIDs)
	cloned.Immunities = maps.Clone(ad.Immunities)
	cloned.JutsuImmunities = maps.Clone(ad.JutsuImmunities)

	return cloned
}

func MakeActorState() ActorState {
	return ActorState{
		State:                ActorStateGrounded,
		ActiveTurns:          0,
		Alive:                true,
		Damage:               0,
		InactiveTurns:        0,
		Enters:               0,
		PositionID:           nil,
		LastUsedActionTX:     nil,
		LastReceivedDamage:   map[uuid.UUID]int{},
		ActionLocked:         false,
		SwitchLocked:         false,
		Protected:            false,
		Safeguarded:          false,
		Warded:               false,
		IgnoreRedirect:       false,
		PowerMultiplier:      1.0,
		StabMultiplier:       1.5,
		Reflect:              0.0,
		RecoilMultiplier:     1.0,
		ModifierChanceOffset: 0,
		ModifierChanceMult:   1.0,
		StatusChanceOffset:   0,
		StatusChanceMult:     1.0,
		CooldownOffset:       0,
		RepeatsMinOffset:     0,
		RepeatsMaxOffset:     0,
		ActionAccuracyOffset: 0,
		Immortal:             false,
		Seen:                 false,
		StaminaDamage:        0,
		Stunned:              false,
		Statused:             false,
		Sleeping:             false,
		SleepCounter:         0,
	}
}

func MakeActor(
	def ActorDef,
	playerID uuid.UUID,
	experience int,
	ability *Modifier,
	item *Modifier,
	actions []Action,
	focus ActorFocus,
	auxStats map[ActorStat]int,
) Actor {
	clonedDef := def.Clone()
	actions = append(actions, Switch)
	return Actor{
		ActorDef:   clonedDef,
		ID:         uuid.New(),
		PlayerID:   playerID,
		Enabled:    false,
		Level:      GetLevel(experience),
		Experience: experience,
		Focus:      focus,
		Item:       item,
		Ability:    ability,
		ActorState: MakeActorState(),
		DamageMultipliers: map[AttackStat]float64{
			Attack:       1.0,
			ChakraAttack: 1.0,
		},
		DamageReduction: map[AttackStat]float64{
			Attack:       1.0,
			ChakraAttack: 1.0,
		},
		Stages: map[ActorStat]int{
			StatHP:            0,
			StatStamina:       0,
			StatAttack:        0,
			StatDefense:       0,
			StatChakraAttack:  0,
			StatChakraDefense: 0,
			StatSpeed:         0,
			StatEvasion:       0,
			StatAccuracy:      0,
		},
		AuxStats:         maps.Clone(auxStats),
		AppliedModifiers: map[uuid.UUID]int{},
		Actions:          actions,
		Summon:           nil,
	}
}

func (a *Actor) Transform(def ActorDef) {
	clone := def.Clone()
	a.ActorDef = clone
	a.Summon = nil
	a.Ability = nil
	if len(a.ActorDef.Abilities) > 0 {
		a.Ability = &a.ActorDef.Abilities[0]
	}
}

func (a *Actor) Reset() {
	a.Enabled = false
	a.ActorState = MakeActorState()
	a.SetSummon(nil)
}

func (a *Actor) SetActions(actionIDs []uuid.UUID, ACTIONS map[uuid.UUID]Action) {
	a.Actions = makeActions(actionIDs, ACTIONS)
}
func (a *Actor) SetPosition(positionID *uuid.UUID) {
	a.PositionID = positionID
	if positionID != nil {
		a.ActiveTurns = 0
		a.Seen = true
		a.Enters++
	} else {
		a.AuxAbility = nil
		a.LastUsedActionTX = nil
		a.InactiveTurns = 0
		a.PoisonedCounter = 0
		a.SetSummon(nil)
	}
}
func (a *Actor) PushImmunities(ids ...uuid.UUID) {
	if a.Immunities == nil {
		a.Immunities = make(map[uuid.UUID]struct{})
	}
	for _, id := range ids {
		a.Immunities[id] = struct{}{}
	}
}
func (a Actor) HasImmunity(id uuid.UUID) bool {
	_, ok := a.Immunities[id]
	return ok
}
func (a Actor) HasJutsuImmunity(jutsu ActionJutsu) bool {
	_, ok := a.JutsuImmunities[jutsu]
	return ok
}
func (a *Actor) SetSummon(summon *Summon) {
	if summon != nil {
		summon.Parent = a
	}
	a.Summon = summon
}
func (a *Actor) SetSummonFromActor(actor *Actor, proxy bool) {
	if actor == nil {
		a.Summon = nil
		return
	}

	summon := Summon{
		Actor:  *actor,
		Parent: a,
		Proxy:  proxy,
	}
	a.SetSummon(&summon)
	a.LastUsedActionTX = nil
}
func (a *Actor) SetActionCooldown(actionID uuid.UUID, cooldown int) {
	for i := range a.Actions {
		if a.Actions[i].ID == actionID {
			a.Actions[i].State.Cooldown = &cooldown
		}
	}
}
func (a *Actor) DecrementCooldowns() {
	for i := range a.Actions {
		if a.Actions[i].State.Cooldown != nil {
			cd := *a.Actions[i].State.Cooldown - 1
			if cd < 0 {
				a.Actions[i].State.Cooldown = nil
			} else {
				a.Actions[i].State.Cooldown = &cd
			}
		}
	}
}
func (a *Actor) RecoverStamina(g Game, ratio float64) {
	resolved := a.ResolveStats(g)
	amount := Round(float64(resolved.Stats[StatStamina]) * ratio)
	a.StaminaDamage = max(a.StaminaDamage-amount, 0)
}

func (a *Actor) IncrementTurns() {
	if a.IsActive() {
		a.ActiveTurns++
	} else {
		a.InactiveTurns++
	}
}

func (a Actor) GetFocusModifier(stat ActorStat) float64 {
	delta, ok := ActorFocuses[a.Focus]
	if !ok {
		return 1.0
	}

	if delta.Up == stat {
		return FOCUS_UP_COEF
	}
	if delta.Down == stat {
		return FOCUS_DOWN_COEF
	}

	return 1.0
}
func (a Actor) GetModifiers() []Modifier {
	modifiers := make([]Modifier, 0)
	ability := a.GetAbility()
	if ability != nil {
		modifiers = append(modifiers, *ability)
	}
	if a.Item != nil {
		modifiers = append(modifiers, *a.Item)
	}

	for _, mod := range a.ActorDef.DefaultModifiers {
		modifiers = append(modifiers, mod)
	}

	return modifiers
}

func MapBaseStat(stat, level int, focus float64, ev int) int {
	base := float64((stat * 2) + BASE_IV)
	ratio := float64((base+(float64(ev*2)))*float64(level)) / 100
	return Round((ratio + 5) * focus)
}

func MapResourceStat(stat, level int, focus float64, ev int) int {
	return MapBaseStat(stat, level, focus, ev) + level + 5
}

func (actor *Actor) MapBase(stat ActorStat) {
	actor.Stats[stat] = MapBaseStat(actor.Stats[stat], actor.Level, actor.GetFocusModifier(stat), actor.AuxStats[stat])
}

func (actor *Actor) MapResource(stat ActorStat) {
	actor.Stats[stat] = MapResourceStat(actor.Stats[stat], actor.Level, 1.0, actor.AuxStats[stat])
}

func MapBaseStats(actor Actor) Actor {
	actor.MapResource(StatHP)
	actor.MapResource(StatStamina)

	actor.MapBase(StatAttack)
	actor.MapBase(StatDefense)
	actor.MapBase(StatChakraAttack)
	actor.MapBase(StatChakraDefense)
	actor.MapBase(StatSpeed)

	return actor
}

func MapStagedStat(stat, stage, mod int) int {
	m := 1.0
	if stage > 0 {
		stage = min(stage, 6)
		m = float64(stage+mod) / float64(mod)
	} else if stage < 0 {
		stage = max(stage, -6)
		m = float64(mod) / float64(-stage+mod)
	}

	return Round(float64(stat) * m)
}

func (actor *Actor) MapStaged(stat ActorStat, mod int) {
	actor.Stats[stat] = MapStagedStat(actor.Stats[stat], actor.Stages[stat], mod)
}

func MapStagedStats(actor Actor) Actor {
	actor.MapStaged(StatAttack, 2)
	actor.MapStaged(StatDefense, 2)
	actor.MapStaged(StatChakraAttack, 2)
	actor.MapStaged(StatChakraDefense, 2)
	actor.MapStaged(StatSpeed, 2)
	actor.MapStaged(StatEvasion, 3)
	actor.MapStaged(StatAccuracy, 3)
	return actor
}

func newActorContext(actor Actor) Context {
	context := MakeContextForActor(actor)
	context.TargetActorIDs = []uuid.UUID{}
	return context
}

func GetActorModifiers(game Game) []Transaction[Modifier] {
	var modifiers []Transaction[Modifier]
	activeActors := game.GetActorsFilters(
		Context{},
		ActiveFilter,
	)

	for _, actor := range activeActors {
		context := newActorContext(actor)
		a_modifiers := actor.GetModifiers()

		for _, mod := range a_modifiers {
			transaction := MakeTransaction(mod, context)
			modifiers = append(modifiers, transaction)
		}
	}

	return modifiers
}

var specialMutations = []ActorMutation{
	MakeActorMutation(
		nil,
		MutPriorityMapBaseStats,
		AllFilter,
		func(g Game, input Actor, context Context) Actor {
			return MapBaseStats(input)
		},
	),
	MakeActorMutation(
		nil,
		MutPriorityMapStagedStats,
		AllFilter,
		func(g Game, input Actor, context Context) Actor {
			return MapStagedStats(input)
		},
	),
}

func (a Actor) Clone() Actor {
	b := a

	b.Stats = maps.Clone(a.Stats)
	b.Stages = maps.Clone(a.Stages)
	b.AuxStats = maps.Clone(a.AuxStats)
	b.NatureDamage = maps.Clone(a.NatureDamage)
	b.NatureResistance = maps.Clone(a.NatureResistance)
	b.DamageMultipliers = maps.Clone(a.DamageMultipliers)
	b.DamageReduction = maps.Clone(a.DamageReduction)

	b.Natures = maps.Clone(a.Natures)
	b.Actions = slices.Clone(a.Actions)

	b.Immunities = maps.Clone(a.Immunities)
	b.JutsuImmunities = maps.Clone(a.JutsuImmunities)

	if a.Summon != nil {
		s := new(Summon)
		*s = *a.Summon
		s.Actor = a.Summon.Actor.Clone()
		s.Parent = &b
		b.Summon = s
	}

	if a.AppliedModifiers != nil {
		b.AppliedModifiers = maps.Clone(a.AppliedModifiers)
	} else {
		b.AppliedModifiers = make(map[uuid.UUID]int)
	}

	return b
}

func getContext(actor Actor, transactions []Transaction[Modifier], mutation ActorMutation) Context {
	context := newActorContext(actor)
	return ResolveModifierTransactionContext(context, transactions, mutation.TransactionID)
}

func resolveActor(actor Actor, g Game, bypassModifiers bool) ResolvedActor {
	handler := newActorResolveHandler(actor, g, bypassModifiers)

	resolved := handler.resolveMutations()
	handler.resolveNatures(&resolved)
	handler.resolveActions(&resolved)
	return resolved
}

func resolveActorStats(actor Actor, g Game, bypassModifiers bool) ResolvedActor {
	handler := newActorResolveHandler(actor, g, bypassModifiers)

	resolved := handler.resolveMutations()
	handler.resolveNatures(&resolved)
	return resolved
}

func (a Actor) getActor() Actor {
	actor := a
	if actor.Summon != nil && !actor.Summon.Proxy {
		form := actor.Summon.Actor
		actor.ActorDef = form.ActorDef
		actor.Actions = form.Actions
	}

	return actor
}

func (a Actor) ResolveStats(g Game) ResolvedActor {
	return resolveActorStats(a.getActor(), g, false)
}
func (a Actor) Resolve(g Game) ResolvedActor {
	actor := a.getActor()
	resolved := resolveActor(actor, g, false)
	unmodified := resolveActorStats(actor, g, true)
	resolved.UnmodifiedStats = maps.Clone(unmodified.Stats)
	resolved.Ability = actor.GetAbility()

	return resolved
}

func (r ResolvedActor) CanAct(game *Game, context Context) bool {
	if r.Stunned {
		log := NewLogContext("$source$ was stunned.", context)
		game.PushLog(log)
		return false
	}

	if r.Paralyzed {
		// check for 1/4 chance
		roll := rand.IntN(100)
		if roll > 75 {
			log := NewLogContext("$source$ could not move.", context)
			game.PushLog(log)
			return false
		}
	}

	if r.Sleeping {
		if r.SleepCounter == 0 {
			game.UpdateActor(r.ID, func(a Actor) Actor {
				a.Statused = false
				a.Sleeping = false
				return a
			})
			game.FilterModifiers(func(modifier Transaction[Modifier]) bool {
				if modifier.Mutation.Status {
					for _, t := range game.GetTargets(modifier.Context) {
						if t.ID == r.ID {
							return false
						}
					}
					return true
				}
				if modifier.Context.SourceActorID == nil {
					return true
				}
				return *modifier.Context.SourceActorID != r.ID
			})
			log := NewLogContext("$source$ woke up.", context)
			game.PushLog(log)
		} else {
			game.UpdateActor(r.ID, func(a Actor) Actor {
				a.SleepCounter--
				return a
			})
			log := NewLogContext("$source$ was sleeping.", context)
			game.PushLog(log)
			return false
		}
	}

	return true
}
func (r ResolvedActor) HasChakra(amount int) bool {
	return (r.Stats[StatStamina] - r.StaminaDamage) >= amount
}
func (r ResolvedActor) IsProtected(mut GameMutation) (GameTransaction, bool) {
	if r.Protected {
		return MakeTransaction(mut, NewContext()), true
	}

	return GameTransaction{}, false
}
func (r ResolvedActor) GetHealthRatio() float64 {
	hp := float64(r.Stats[StatHP])
	dmg := float64(r.Damage)
	return (hp - dmg) / hp
}

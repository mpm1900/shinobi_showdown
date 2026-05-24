import { getVitals, type Actor } from '#/lib/game/actor'
import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import { motion, useMotionValue, useSpring, useTransform } from 'motion/react'
import { useEffect, useMemo, useRef } from 'react'
import { HealthDeltaBurst } from './health-delta-burst'

function HealthBar({
  actor,
  selected = false,
  percentage = true,
}: {
  actor: Actor
  selected?: boolean
  percentage?: boolean
}) {
  const turn = useStore(gameStore, (g) => g.turn.count)
  const prevTurnRef = useRef(turn)
  const vitals = getVitals(actor)

  const hpTarget = useMotionValue(vitals.hp.ratio * 100)
  const staminaTarget = useMotionValue(vitals.stamina.ratio * 100)
  const hpGhostTarget = useMotionValue(vitals.hp.ratio * 100)

  useEffect(() => {
    hpTarget.set(vitals.hp.ratio * 100)
    staminaTarget.set(vitals.stamina.ratio * 100)
  }, [vitals.hp.ratio, vitals.stamina.ratio, hpTarget, staminaTarget])

  useEffect(() => {
    if (prevTurnRef.current !== turn) {
      prevTurnRef.current = turn
      hpGhostTarget.set(vitals.hp.ratio * 100)
    }
  }, [turn, vitals.hp.ratio, hpGhostTarget])

  const hpWidth = useTransform(
    useSpring(hpTarget, { stiffness: 280, damping: 30, mass: 0.8 }),
    (v) => `${v}%`
  )

  const hpGhostWidth = useTransform(
    useSpring(hpGhostTarget, { stiffness: 50, damping: 26, mass: 1.2 }),
    (v) => `${v}%`
  )

  const chakraWidth = useTransform(
    useSpring(staminaTarget, { stiffness: 200, damping: 26, mass: 0.95 }),
    (v) => `${v}%`
  )

  const hpFillStyle = useMemo(() => {
    if (vitals.hp.ratio > 0.8) {
      return 'bg-gradient-to-r from-emerald-500 to-emerald-700'
    }
    if (vitals.hp.ratio > 0.6) {
      return 'bg-gradient-to-r from-green-600 to-lime-400/50'
    }
    if (vitals.hp.ratio > 0.3) {
      return 'bg-gradient-to-r from-yellow-300 via-orange-400 to-yellow-600'
    }
    if (vitals.hp.ratio > 0.2) {
      return 'bg-gradient-to-r from-amber-300 via-orange-400 to-orange-600'
    }
    return 'bg-gradient-to-r from-rose-500 via-rose-600 to-rose-800'
  }, [vitals.hp.ratio])

  return (
    <div className="relative w-full select-none nanum-brush-script-regular">
      {selected && (
        <div className="pointer-events-none absolute -inset-1 rounded-sm bg-gradient-to-r from-cyan-400/15 via-indigo-400/10 to-red-500/15 blur-md" />
      )}

      <div className="relative h-6 overflow-hidden rounded-sm border border-transparent bg-black shadow-[1px_1px_0_rgba(0,0,0,1)]">
        {/* CHAKRA LAYER (behind, border-like energy shell) */}
        <div className="absolute inset-0 bg-slate-800" />
        <motion.div
          className={`absolute inset-y-0 left-0 rounded-sm bg-gradient-to-r from-sky-300 via-sky-200 to-blue-300 ${selected ? 'shadow-lg ring-2 ring-sky-300/40' : ''}`}
          style={{ width: chakraWidth }}
        />

        {/* Inner chamber where HP lives */}
        <div className="absolute inset-px overflow-hidden rounded-sm border border-black/50 bg-black/35 shadow-inner">
          {/* low-contrast track */}
          <div className="absolute inset-0 bg-gradient-to-b from-stone-100/5 to-black/25" />

          {/* delayed damage ghost */}
          <motion.div
            className="absolute inset-y-0 left-0 bg-stone-100/40"
            style={{ width: hpGhostWidth }}
          />

          {/* main HP fill */}
          <motion.div
            className={`absolute inset-y-0 left-0 ${hpFillStyle} ${selected ? 'shadow-md' : ''}`}
            style={{ width: hpWidth }}
          />

          {/* moving highlight for premium sheen */}
          {selected && (
            <motion.div
              className="absolute inset-y-0 left-[-30%] w-[30%] bg-gradient-to-r from-transparent via-stone-100/20 to-transparent mix-blend-screen"
              animate={{ x: ['0%', '460%'] }}
              transition={{ duration: 2.6, ease: 'linear', repeat: Infinity }}
            />
          )}

          {/* subtle grid texture */}
          <div className="absolute inset-0 bg-stone-100/20 opacity-[0.15]" />
        </div>

        {/* text + values */}
        <div className="pointer-events-none absolute inset-0 z-10 flex items-center justify-end px-2">
          <div className="flex flex-row-reverse items-center gap-2">
            <span className="text-2xl -mb-1 tabular-nums text-stone-100 drop-shadow text-shadow-[1px_1px_0px_#000000]">
              {percentage && (
                <span className="font-bold">
                  {Math.round(vitals.hp.ratio * 100)}%
                </span>
              )}
              {!percentage && (
                <span>
                  <span className="font-bold">{vitals.hp.current}</span>
                  <span className="opacity-60">/{vitals.hp.max}</span>
                </span>
              )}
            </span>
          </div>
        </div>
      </div>

      <HealthDeltaBurst
        hpRatio={vitals.hp.ratio}
        direction={percentage ? 'down' : 'up'}
      />
    </div>
  )
}

export { HealthBar }

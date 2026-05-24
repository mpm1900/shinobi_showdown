import { AnimatePresence, motion } from 'motion/react'
import { useEffect, useRef, useState } from 'react'

type FloatingDelta = {
  id: number
  delta: number
  driftX: number
}

const DRIFT_PATTERN = [-16, -9, 0, 9, 16]

function HealthDeltaBurst({
  hpRatio,
  direction,
}: {
  hpRatio: number
  direction: 'up' | 'down'
}) {
  const prevRatioRef = useRef(hpRatio)
  const idRef = useRef(0)
  const [deltas, setDeltas] = useState<Array<FloatingDelta>>([])

  useEffect(() => {
    const prev = prevRatioRef.current
    if (prev === hpRatio) return

    const delta = Math.round((hpRatio - prev) * 100)
    prevRatioRef.current = hpRatio
    if (delta === 0) return

    const id = ++idRef.current
    const driftX = DRIFT_PATTERN[id % DRIFT_PATTERN.length]

    setDeltas((curr) => [...curr.slice(-2), { id, delta, driftX }])
  }, [hpRatio])

  return (
    <div className="pointer-events-none absolute inset-0 z-20 overflow-visible">
      <AnimatePresence initial={false}>
        {deltas.map((entry) => {
          const isDamage = entry.delta < 0

          return (
            <div
              key={entry.id}
              className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
            >
              <motion.span
                className={`relative block whitespace-nowrap text-4xl font-black tabular-nums tracking-wide will-change-transform [text-shadow:0_0_6px_rgba(0,0,0,0.85),0_1px_0_rgba(0,0,0,1),0_0_14px_rgba(255,255,255,0.25)] ${isDamage ? 'text-rose-200 drop-shadow-[0_2px_3px_rgba(120,0,30,0.95)]' : 'text-emerald-100 drop-shadow-[0_2px_3px_rgba(0,95,40,0.95)]'}`}
                initial={{
                  opacity: 0,
                  scale: 0.82,
                  y: 0,
                  x: 0,
                }}
                animate={{
                  opacity: [0, 1, 1, 0],
                  scale: [0.82, 1.2, 1.08, 1.02],
                  y: direction === 'up' ? -36 : 36,
                  x: entry.driftX,
                }}
                exit={{ opacity: 0 }}
                transition={{ duration: 2.8, ease: [0.2, 0.65, 0.22, 0.98] }}
                onAnimationComplete={() => {
                  setDeltas((curr) => curr.filter((d) => d.id !== entry.id))
                }}
              >
                {entry.delta > 0 ? `+${entry.delta}%` : `${entry.delta}%`}
              </motion.span>
            </div>
          )
        })}
      </AnimatePresence>
    </div>
  )
}

export { HealthDeltaBurst }

import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import { useEffect, useRef } from 'react'

declare global {
  interface Window {
    VANTA: any
  }
}

export const VantaBackground = () => {
  const vantaRef = useRef<HTMLDivElement>(null)
  const vantaEffectRef = useRef<any>(null)
  const g_state = useStore(gameStore, g => g.state)
  const weather = g_state.weather

  useEffect(() => {
    const initVanta = () => {
      if (
        window.VANTA &&
        window.VANTA.FOG &&
        vantaRef.current &&
        !vantaEffectRef.current
      ) {
        try {
          let colors = {
            highlightColor: 0x8089,
            midtoneColor: 0xffa600,
            lowlightColor: 0xffffff,
            baseColor: 0x0,
          }

          if (weather === 'rain') {
            colors = {
              highlightColor: 0x38608e,
              midtoneColor: 0xb3d2de,
              lowlightColor: 0x59ff,
              baseColor: 0xf20,
            }
          }

          if (weather === 'sandstorm') {
            colors = {
              highlightColor: 0x827052,
              midtoneColor: 0x725d34,
              lowlightColor: 0x848465,
              baseColor: 0x201400,
            }
          }

          vantaEffectRef.current = window.VANTA.FOG({
            el: vantaRef.current,
            mouseControls: false,
            touchControls: false,
            gyroControls: false,
            minHeight: 200.0,
            minWidth: 200.0,
            blurFactor: 0.35,
            speed: 0.3,
            zoom: 0.4,
            ...colors
          })
        } catch (err) {
          console.error('Vanta initialization failed:', err)
        }
      }
    }

    // Check every 100ms if VANTA is available, up to 50 times (5 seconds)
    let attempts = 0
    const interval = setInterval(() => {
      if (window.VANTA && window.VANTA.FOG) {
        initVanta()
        clearInterval(interval)
      }
      attempts++
      if (attempts > 50) {
        clearInterval(interval)
      }
    }, 100)

    return () => {
      clearInterval(interval)
      if (vantaEffectRef.current) {
        vantaEffectRef.current.destroy()
        vantaEffectRef.current = null
      }
    }
  }, [weather])

  return (
    <>
      <div
        ref={vantaRef}
        style={{
          position: 'fixed',
          inset: 0,
          zIndex: -2,
          pointerEvents: 'none',
        }}
      />
      <div
        style={{
          position: 'fixed',
          inset: 0,
          zIndex: -1,
          pointerEvents: 'none',
          boxShadow: 'inset 0 0 500px rgba(0,0,0,1)',
        }}
      />
      <div
        style={{
          position: 'fixed',
          inset: 0,
          zIndex: -1,
          pointerEvents: 'none',
          boxShadow: 'inset 0 0 200px rgba(0,0,0,1)',
        }}
      />
    </>
  )
}

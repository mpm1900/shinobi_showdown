import {
  getResistance,
  getWeakness,
  natureIndexes,
  type NatureSet,
} from '#/lib/game/nature'
import { NatureBadge } from './nature-badge'

function NastureSetDetails({ natures }: { natures: Array<NatureSet> }) {
  const weaknesses = getWeakness(...natures)
  const resistances = getResistance(...natures)
  if (weaknesses.length === 0 && resistances.length === 0) return null
  return (
    <div className="flex justify-between items-center gap-4">
      <div className="flex items-center gap-1">
        <div className="text-muted-foreground text-xs flex-1 text-center text-nowrap">
          is weak to
        </div>
        {weaknesses
          .filter((w) => !resistances.includes(w))
          .sort((a, b) => natureIndexes[a] - natureIndexes[b])
          .map((nature) => (
            <NatureBadge key={nature} nature={nature} />
          ))}
      </div>
      <div className="flex items-center gap-1">
        <div className="text-muted-foreground text-xs flex-1 text-center text-nowrap">
          resists
        </div>
        {resistances
          .filter((r) => !weaknesses.includes(r))
          .sort((a, b) => natureIndexes[a] - natureIndexes[b])
          .map((nature) => (
            <NatureBadge key={nature} nature={nature} />
          ))}
      </div>
    </div>
  )
}

export { NastureSetDetails }

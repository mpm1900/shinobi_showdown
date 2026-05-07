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
  return (
    <div>
      <div className="flex justify-between items-center">
        {natures
          .sort((a, b) => natureIndexes[a] - natureIndexes[b])
          .map((nature) => (
            <NatureBadge key={nature} nature={nature} />
          ))}
        <div className="text-muted-foreground text-xs flex-1 text-center">
          {' '}
          is weak to{' '}
        </div>
        {weaknesses
          .filter((w) => !resistances.includes(w))
          .sort((a, b) => natureIndexes[a] - natureIndexes[b])
          .map((nature) => (
            <NatureBadge key={nature} nature={nature} />
          ))}
      </div>
      <div className="flex justify-between items-center">
        {natures
          .sort((a, b) => natureIndexes[a] - natureIndexes[b])
          .map((nature) => (
            <NatureBadge key={nature} nature={nature} />
          ))}
        <div className="text-muted-foreground text-xs flex-1 text-center">
          {' '}
          resists{' '}
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

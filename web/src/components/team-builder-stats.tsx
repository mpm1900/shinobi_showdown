import { getTotalBaseStats, type ActorDef } from '#/lib/game/actor'
import type { ActorConfig } from '#/lib/stores/config'
import { cn } from '#/lib/utils'
import { TeamBuilderStat } from './team-builder-stat'

function TeamBuilderStats({
  config,
  def,
  total,
  onConfigChange,
}: {
  config: ActorConfig
  def: ActorDef
  total: number
  onConfigChange: (config: ActorConfig) => void
}) {
  return (
    <table className="flex-1">
      <tbody>
        <tr>
          <td colSpan={1}>Stats</td>
          <td
            colSpan={1}
            className="text-end w-8 p-2 whitespace-nowrap font-black"
          >
            {getTotalBaseStats(def)}
          </td>
          <td
            colSpan={3}
            className={cn('text-end', total > 66 ? 'text-destructive' : '')}
          >
            {total}
            /66
          </td>
        </tr>
        <TeamBuilderStat
          total={total}
          base={def}
          stat="stamina"
          config={config}
          onConfigChange={onConfigChange}
        />
        <TeamBuilderStat
          total={total}
          base={def}
          stat="hp"
          config={config}
          onConfigChange={onConfigChange}
        />


        <TeamBuilderStat
          total={total}
          base={def}
          stat="attack"
          config={config}
          onConfigChange={onConfigChange}
        />
        <TeamBuilderStat
          total={total}
          base={def}
          stat="defense"
          config={config}
          onConfigChange={onConfigChange}
        />
        <TeamBuilderStat
          total={total}
          base={def}
          stat="chakra_attack"
          config={config}
          onConfigChange={onConfigChange}
        />
        <TeamBuilderStat
          total={total}
          base={def}
          stat="chakra_defense"
          config={config}
          onConfigChange={onConfigChange}
        />
        <TeamBuilderStat
          total={total}
          base={def}
          stat="speed"
          config={config}
          onConfigChange={onConfigChange}
        />
      </tbody>
    </table>
  )
}

export { TeamBuilderStats }

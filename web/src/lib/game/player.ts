type PlayerPosition = {
  ID: string
  actor_ID: string | null
}

type PlayerUser = {
  id: string
  username: string
  email: string
}

type Player = {
  ID: string
  user: PlayerUser
  positions_capacity: number
  positions: Array<PlayerPosition>
  team_capacity: number
  used_summon: boolean
  ready: boolean
}

export type { Player, PlayerUser }

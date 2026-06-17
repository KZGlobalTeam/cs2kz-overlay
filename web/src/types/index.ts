export type Mode = 'classic' | 'vanilla'
export type Tier =
  | 'very-easy'
  | 'easy'
  | 'medium'
  | 'advanced'
  | 'hard'
  | 'very-hard'
  | 'extreme'
  | 'death'
  | 'unfeasible'
  | 'impossible'

export interface GameState {
  mapName: string | null
  courseName: string | null
  mode: Mode | null
  steamId: string | null
}

export interface Record {
  id: string
  player: {
    id: string
    name: string
  }
  server: {
    id: number
    name: string
  }
  map: {
    id: number
    name: string
  }
  course: {
    id: number
    name: string
    nub_tier:
      | 'very-easy'
      | 'easy'
      | 'medium'
      | 'advanced'
      | 'hard'
      | 'very-hard'
      | 'extreme'
      | 'death'
      | 'unfeasible'
      | 'impossible'
    pro_tier:
      | 'very-easy'
      | 'easy'
      | 'medium'
      | 'advanced'
      | 'hard'
      | 'very-hard'
      | 'extreme'
      | 'death'
      | 'unfeasible'
      | 'impossible'
    state: 'unranked' | 'pending' | 'ranked'
  }
  mode: 'vanilla' | 'classic'
  styles: 'auto-bhop'[]
  teleports: number
  time: number
  nub_rank: number | null
  nub_points: number | null
  pro_rank: number | null
  pro_points: number | null
  replay_available: boolean
}

export interface RecordResponse {
  total: number
  values: Record[]
}

export interface PlayerProfile {
  id: string
  name: string
  profile_url: string
  avatar_url: string
}

export interface Map {
  id: number
  workshop_id: number
  name: string
  description?: string
  state: 'invalid' | 'in-testing' | 'approved'
  vpk_checksum: string
  mappers: {
    id: string
    name: string
  }[]
  courses: {
    name: string
    description?: string
    mappers: {
      id: string
      name: string
    }[]
    filters: {
      vanilla: {
        nub_tier: Tier
        pro_tier: Tier
        state: 'unranked' | 'pending' | 'ranked'
        notes?: string
      }
      classic: {
        nub_tier: Tier
        pro_tier: Tier
        state: 'unranked' | 'pending' | 'ranked'
        notes?: string
      }
    }
  }[]
  approved_at: string
}

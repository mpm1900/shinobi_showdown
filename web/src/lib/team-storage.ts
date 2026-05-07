import z from 'zod'
import { TeamSchema, type TeamConfig } from '#/lib/stores/config'

export const SAVED_TEAMS_KEY = 'team-builder:saved-teams'

const SavedTeamsSchema = z.array(TeamSchema).default([])

export const parseSavedTeams = (value: unknown): TeamConfig[] => {
  const parsed = SavedTeamsSchema.safeParse(value)
  if (!parsed.success) return []
  return parsed.data
}

export const parseSavedTeamsJSON = (raw: string | null): TeamConfig[] => {
  if (!raw) return []

  try {
    return parseSavedTeams(JSON.parse(raw))
  } catch {
    return []
  }
}

export const serializeSavedTeams = (teams: TeamConfig[]): string =>
  JSON.stringify(teams)

export const upsertSavedTeam = (
  existing: TeamConfig[],
  candidate: TeamConfig
): TeamConfig[] => {
  const parsed = TeamSchema.safeParse(candidate)
  if (!parsed.success) return existing

  return [
    parsed.data,
    ...existing.filter((team) => team.name !== parsed.data.name),
  ]
}

export const cloneTeamConfig = (team: TeamConfig): TeamConfig =>
  structuredClone(team)

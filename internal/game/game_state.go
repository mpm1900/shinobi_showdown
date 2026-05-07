package game

type GameWeather string
type GameTerrain string

const (
	GameWeatherNone      GameWeather = "none"
	GameWeatherRain      GameWeather = "rain"
	GameWeatherSandstorm GameWeather = "sandstorm"
	GameWeatherSunlight  GameWeather = "sunlight"
)

const (
	GameTerrainNone        GameTerrain = "none"
	GameTerrainFlamable    GameTerrain = "flamable"
	GameTerrainElectrified GameTerrain = "electrified"
	GameTerrainFlooded     GameTerrain = "flooded"
	GameTerrainRocky       GameTerrain = "rocky"
)

/**
 * [GameState]
 * This object exists to bypass a complicated system of
 * checks to see if more than one modifier of a given
 * type can exist at once.
 *
 * Instead, just modifity this object so only one value can exist
 * examples can be weather, terrain,
 *
 * the design space does overlap a bit with modifiers if not for the singular nature
 */
type GameState struct {
	Weather GameWeather `json:"weather"`
	Terrain GameTerrain `json:"terrain"`
}

func NewGameState() GameState {
	return GameState{
		Terrain: GameTerrainNone,
		Weather: GameWeatherNone,
	}
}

func GS_TrueFilter(g Game, gs GameState, c Context) bool {
	return true
}
func GS_SourceIsActiveFilter(g Game, gs GameState, context Context) bool {
	source, ok := g.GetSource(context)
	if !ok {
		return false
	}

	return source.IsActive()
}

package mcpserver

// RouteKey identifies a standard MCP REST mirror route.
type RouteKey string

const (
	RouteRadiation    RouteKey = "radiation"
	RouteArea         RouteKey = "area"
	RouteTracks       RouteKey = "tracks"
	RouteTrackByID    RouteKey = "track_by_id"
	RouteDevice       RouteKey = "device"
	RouteSensors      RouteKey = "sensors"
	RouteSensorByID   RouteKey = "sensor_by_id"
	RouteSpectra      RouteKey = "spectra"
	RouteSpectrumByID RouteKey = "spectrum_by_id"
	RouteStats        RouteKey = "stats"
	RouteExtreme      RouteKey = "extreme"
	RouteInfo         RouteKey = "info"
	RouteGPTRadiation RouteKey = "gpt_radiation"
	RouteGPTArea      RouteKey = "gpt_area"
	RouteGPTStats     RouteKey = "gpt_stats"
)

var defaultRESTRouteOrder = []RouteKey{
	RouteRadiation,
	RouteArea,
	RouteTracks,
	RouteTrackByID,
	RouteDevice,
	RouteSensors,
	RouteSensorByID,
	RouteSpectra,
	RouteSpectrumByID,
	RouteStats,
	RouteExtreme,
	RouteInfo,
	RouteGPTRadiation,
	RouteGPTArea,
	RouteGPTStats,
}

// RegisterRESTRoutes iterates the canonical REST mirror route set.
func RegisterRESTRoutes(register func(RouteKey)) {
	if register == nil {
		return
	}
	for _, key := range defaultRESTRouteOrder {
		register(key)
	}
}

package mcpserver

// ToolKey identifies a standard MCP tool registration slot.
type ToolKey string

const (
	ToolPing                   ToolKey = "ping"
	ToolQueryRadiation         ToolKey = "query_radiation"
	ToolSearchArea             ToolKey = "search_area"
	ToolListTracks             ToolKey = "list_tracks"
	ToolGetTrack               ToolKey = "get_track"
	ToolDeviceHistory          ToolKey = "device_history"
	ToolGetSpectrum            ToolKey = "get_spectrum"
	ToolListSpectra            ToolKey = "list_spectra"
	ToolRadiationInfo          ToolKey = "radiation_info"
	ToolDBInfo                 ToolKey = "db_info"
	ToolListSensors            ToolKey = "list_sensors"
	ToolSensorCurrent          ToolKey = "sensor_current"
	ToolSensorHistory          ToolKey = "sensor_history"
	ToolQueryAnalytics         ToolKey = "query_analytics"
	ToolRadiationStats         ToolKey = "radiation_stats"
	ToolQueryDuckDBLogs        ToolKey = "query_duckdb_logs"
	ToolQueryExtremeReadings   ToolKey = "query_extreme_readings"
	ToolTopUploaders           ToolKey = "top_uploaders"
	ToolSearchTracksByLocation ToolKey = "search_tracks_by_location"
)

var defaultToolOrder = []ToolKey{
	ToolPing,
	ToolQueryRadiation,
	ToolSearchArea,
	ToolListTracks,
	ToolGetTrack,
	ToolDeviceHistory,
	ToolGetSpectrum,
	ToolListSpectra,
	ToolRadiationInfo,
	ToolDBInfo,
	ToolListSensors,
	ToolSensorCurrent,
	ToolSensorHistory,
	ToolQueryAnalytics,
	ToolRadiationStats,
	ToolQueryDuckDBLogs,
	ToolQueryExtremeReadings,
	ToolTopUploaders,
	ToolSearchTracksByLocation,
}

// RegisterTools iterates the canonical MCP tool set in stable order.
func RegisterTools(register func(ToolKey)) {
	if register == nil {
		return
	}
	for _, key := range defaultToolOrder {
		register(key)
	}
}

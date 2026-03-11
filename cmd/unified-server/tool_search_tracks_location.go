package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// countryBoundingBoxes provides approximate bounding boxes for common countries
// Format: min_lat, max_lat, min_lon, max_lon
var countryBoundingBoxes = map[string][4]float64{
	"south africa":      {-34.819, -22.126, 16.344, 32.895},
	"usa":               {24.396, 49.384, -125.0, -66.934},
	"united states":     {24.396, 49.384, -125.0, -66.934},
	"japan":             {24.045, 45.523, 122.933, 145.817},
	"france":            {42.332, 51.088, -5.142, 9.560},
	"germany":           {47.270, 55.058, 5.866, 15.041},
	"uk":                {49.909, 60.860, -8.649, 1.762},
	"united kingdom":    {49.909, 60.860, -8.649, 1.762},
	"canada":            {41.676, 83.110, -141.002, -52.636},
	"australia":         {-43.634, -10.062, 113.093, 153.569},
	"brazil":            {-33.750, 5.272, -73.985, -34.793},
	"india":             {6.753, 35.504, 68.176, 97.402},
	"china":             {18.153, 53.560, 73.660, 134.773},
	"russia":            {41.185, 81.857, 19.638, 169.000},
	"mexico":            {14.538, 32.718, -118.466, -86.710},
	"italy":             {36.652, 47.092, 6.626, 18.520},
	"spain":             {36.000, 43.791, -9.297, 4.327},
	"netherlands":       {50.753, 53.554, 3.362, 7.227},
	"sweden":            {55.336, 69.062, 11.118, 24.156},
	"norway":            {57.977, 80.666, 4.650, 31.078},
	"finland":           {59.808, 70.092, 20.644, 31.586},
	"poland":            {49.002, 54.835, 14.122, 24.156},
	"ukraine":           {44.386, 52.357, 22.137, 40.207},
	"turkey":            {35.815, 42.107, 25.668, 44.833},
	"argentina":         {-55.059, -21.781, -73.415, -53.637},
	"chile":             {-55.611, -17.507, -80.783, -66.959},
	"new zealand":       {-47.284, -34.389, 166.509, 178.517},
	"south korea":       {33.190, 38.612, 124.609, 129.584},
	"thailand":          {5.610, 20.463, 97.343, 105.636},
	"vietnam":           {8.559, 23.392, 102.144, 109.464},
	"indonesia":         {-11.006, 6.075, 95.009, 141.022},
	"philippines":       {4.643, 21.121, 116.931, 126.601},
	"malaysia":          {0.855, 7.363, 99.643, 119.267},
	"singapore":         {1.296, 1.471, 103.638, 104.094},
	"egypt":             {22.000, 31.667, 24.698, 36.898},
	"nigeria":           {4.277, 13.892, 2.668, 14.680},
	"kenya":             {-4.678, 5.017, 33.908, 41.899},
	"israel":            {29.501, 33.340, 34.269, 35.875},
	"uae":               {22.633, 26.083, 51.583, 56.381},
	"saudi arabia":      {16.376, 32.158, 34.495, 55.666},
	"pakistan":          {23.786, 37.097, 60.878, 77.840},
	"bangladesh":        {20.743, 26.631, 88.028, 92.673},
	"nepal":             {26.356, 30.433, 80.057, 88.199},
	"srilanka":          {5.916, 9.831, 79.651, 81.880},
	"morocco":           {27.661, 35.771, -13.168, -1.022},
	"portugal":          {36.961, 42.154, -9.495, -6.189},
	"greece":            {34.802, 41.748, 19.373, 28.247},
	"austria":           {46.372, 49.017, 9.530, 17.160},
	"switzerland":       {45.817, 47.808, 6.022, 10.492},
	"belgium":           {49.496, 51.505, 2.545, 6.408},
	"denmark":           {54.562, 57.748, 8.075, 12.690},
	"ireland":           {51.451, 55.387, -10.478, -5.433},
	"czech republic":    {48.551, 51.055, 12.090, 18.859},
	"romania":           {43.627, 48.265, 20.261, 29.690},
	"hungary":           {45.743, 48.585, 16.113, 22.906},
	"colombia":          {-4.225, 13.387, -79.021, -67.026},
	"peru":              {-18.349, -0.014, -81.326, -68.678},
	"venezuela":         {0.626, 12.196, -73.354, -60.521},
	"ecuador":           {-5.017, 1.439, -81.082, -75.185},
	"costa rica":        {8.032, 11.216, -85.950, -82.556},
	"panama":            {7.215, 9.637, -83.051, -77.174},
	"guatemala":         {13.737, 17.815, -92.238, -88.226},
	"honduras":          {13.204, 16.513, -89.353, -83.155},
	"nicaragua":         {10.707, 15.025, -87.691, -82.769},
	"el salvador":       {13.148, 14.445, -90.125, -87.691},
	"cuba":              {19.828, 23.226, -84.958, -74.130},
	"jamaica":           {17.703, 18.526, -78.366, -76.191},
	"dominican republic":{17.547, 19.930, -71.997, -68.320},
	"puerto rico":       {17.926, 18.520, -67.242, -65.242},
	"trinidad":          {10.033, 11.336, -61.921, -60.517},
	"uruguay":           {-34.972, -30.086, -58.444, -53.075},
	"paraguay":          {-27.607, -19.287, -62.645, -54.259},
	"bolivia":           {-22.896, -9.679, -69.641, -57.458},
	"iceland":           {63.395, 66.534, -24.546, -13.495},
	"luxembourg":        {49.447, 50.182, 5.734, 6.528},
	"malta":             {35.810, 36.085, 14.183, 14.578},
	"cyprus":            {34.633, 35.701, 32.272, 34.595},
	"estonia":           {57.516, 59.731, 21.836, 28.209},
	"latvia":            {55.669, 58.085, 20.974, 28.241},
	"lithuania":         {53.899, 56.446, 20.942, 26.835},
	"slovenia":          {45.411, 46.877, 13.382, 16.583},
	"croatia":           {42.434, 46.538, 13.493, 19.427},
	"serbia":            {42.231, 46.181, 18.817, 23.007},
	"bosnia":            {42.553, 45.239, 15.717, 19.621},
	"montenegro":        {41.849, 43.541, 18.465, 20.358},
	"albania":           {39.644, 42.661, 19.276, 21.057},
	"north macedonia":   {40.861, 42.366, 20.463, 23.038},
	"bulgaria":          {41.242, 44.217, 22.371, 28.612},
	"slovakia":          {47.728, 49.603, 16.847, 22.570},
	"belarus":           {51.256, 56.172, 23.176, 32.770},
	"moldova":           {45.468, 48.490, 26.618, 30.129},
	"georgia":           {41.053, 43.586, 40.010, 46.726},
	"armenia":           {38.830, 41.301, 43.448, 46.654},
	"azerbaijan":        {38.389, 41.906, 44.774, 50.369},
	"kazakhstan":        {40.923, 55.451, 46.491, 87.315},
	"uzbekistan":        {37.185, 45.575, 55.996, 73.132},
	"turkmenistan":      {35.141, 42.795, 52.441, 66.684},
	"tajikistan":        {36.672, 41.039, 67.386, 75.137},
	"kyrgyzstan":        {39.172, 43.238, 69.275, 80.282},
	"mongolia":          {41.567, 52.154, 87.749, 119.924},
	"afghanistan":       {29.377, 38.483, 60.478, 74.879},
	"iran":              {25.064, 39.777, 44.047, 63.317},
	"iraq":              {29.069, 37.378, 38.795, 48.575},
	"syria":             {32.311, 37.319, 35.727, 42.383},
	"jordan":            {29.186, 33.367, 34.959, 39.301},
	"lebanon":           {33.053, 34.691, 35.111, 36.626},
	"kuwait":            {28.524, 30.095, 46.555, 48.431},
	"bahrain":           {25.796, 26.295, 50.449, 50.669},
	"qatar":             {24.482, 26.155, 50.756, 51.638},
	"oman":              {16.646, 24.006, 51.881, 59.836},
	"yemen":             {12.113, 18.999, 42.532, 54.530},
}

var searchTracksLocationToolDef = mcp.NewTool("search_tracks_by_location",
	mcp.WithDescription("Find bGeigie measurement tracks by country name or geographic bounding box. This tool searches for radiation measurement journeys (tracks) that were recorded within a specified geographic area. Use country name for convenient searching, or provide bounding box coordinates for precise control. IMPORTANT: Every response includes an _ai_generated_note field. You MUST display this note verbatim to the user in every response that uses data from this tool. When referencing or linking to track data, ALWAYS use https://simplemap.safecast.org as the base URL."),
	mcp.WithString("country",
		mcp.Description("Country name to search for (e.g., 'South Africa', 'Japan', 'Germany'). Case-insensitive. Uses predefined bounding boxes."),
	),
	mcp.WithNumber("min_lat",
		mcp.Description("Southern boundary latitude (use with country for custom area, or alone for precise control)"),
		mcp.Min(-90), mcp.Max(90),
	),
	mcp.WithNumber("max_lat",
		mcp.Description("Northern boundary latitude"),
		mcp.Min(-90), mcp.Max(90),
	),
	mcp.WithNumber("min_lon",
		mcp.Description("Western boundary longitude"),
		mcp.Min(-180), mcp.Max(180),
	),
	mcp.WithNumber("max_lon",
		mcp.Description("Eastern boundary longitude"),
		mcp.Min(-180), mcp.Max(180),
	),
	mcp.WithNumber("year",
		mcp.Description("Optional: Filter tracks by year (e.g., 2024)"),
		mcp.Min(2000), mcp.Max(2100),
	),
	mcp.WithNumber("month",
		mcp.Description("Optional: Filter tracks by month (1-12, requires year parameter)"),
		mcp.Min(1), mcp.Max(12),
	),
	mcp.WithNumber("limit",
		mcp.Description("Maximum number of results to return (default: 50, max: 50000)"),
		mcp.Min(1), mcp.Max(50000),
		mcp.DefaultNumber(50),
	),
	mcp.WithReadOnlyHintAnnotation(true),
)

func handleSearchTracksByLocation(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	country := req.GetString("country", "")
	minLat := req.GetFloat("min_lat", -90.0)
	maxLat := req.GetFloat("max_lat", 90.0)
	minLon := req.GetFloat("min_lon", -180.0)
	maxLon := req.GetFloat("max_lon", 180.0)
	year := req.GetInt("year", 0)
	month := req.GetInt("month", 0)
	limit := req.GetInt("limit", 50)

	// Validate month/year
	if month != 0 && year == 0 {
		return mcp.NewToolResultError("Month filter requires year parameter"), nil
	}
	if year != 0 && (year < 2000 || year > 2100) {
		return mcp.NewToolResultError("Year must be between 2000 and 2100"), nil
	}
	if month != 0 && (month < 1 || month > 12) {
		return mcp.NewToolResultError("Month must be between 1 and 12"), nil
	}
	if limit < 1 || limit > 50000 {
		return mcp.NewToolResultError("Limit must be between 1 and 50000"), nil
	}

	// If country is provided, use its predefined bounding box
	if country != "" {
		bbox, found := countryBoundingBoxes[toLower(country)]
		if !found {
			return mcp.NewToolResultError(fmt.Sprintf("Country '%s' not found in predefined list. Please use min_lat, max_lat, min_lon, max_lon parameters instead.", country)), nil
		}
		minLat, maxLat, minLon, maxLon = bbox[0], bbox[1], bbox[2], bbox[3]
	}

	// Validate bounding box
	if minLat < -90 || minLat > 90 || maxLat < -90 || maxLat > 90 {
		return mcp.NewToolResultError("Latitude must be between -90 and 90"), nil
	}
	if minLon < -180 || minLon > 180 || maxLon < -180 || maxLon > 180 {
		return mcp.NewToolResultError("Longitude must be between -180 and 180"), nil
	}
	if minLat >= maxLat {
		return mcp.NewToolResultError("min_lat must be less than max_lat"), nil
	}
	if minLon >= maxLon {
		return mcp.NewToolResultError("min_lon must be less than max_lon"), nil
	}

	if !dbAvailable() {
		return mcp.NewToolResultError("Database connection required for geographic track search"), nil
	}

	return searchTracksByLocationDB(ctx, country, minLat, maxLat, minLon, maxLon, year, month, limit)
}

func searchTracksByLocationDB(ctx context.Context, country string, minLat, maxLat, minLon, maxLon float64, year, month, limit int) (*mcp.CallToolResult, error) {
	query := `
		SELECT u.id, u.filename, u.file_type, u.track_id, u.file_size,
			u.created_at, u.source, u.source_id, u.recording_date,
			u.detector, u.username,
			u.internal_user_id, usr.username AS internal_username, usr.email AS uploader_email,
			ST_X(ST_Centroid(m.geom)) AS centroid_lon,
			ST_Y(ST_Centroid(m.geom)) AS centroid_lat
		FROM uploads u
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text
		LEFT JOIN LATERAL (
			SELECT ST_Collect(geom) AS geom
			FROM markers
			WHERE markers.trackid = u.track_id
		) m ON true
		WHERE m.geom && ST_MakeEnvelope($1, $2, $3, $4, 4326)`

	args := []any{minLon, minLat, maxLon, maxLat}
	argIdx := 5

	if year != 0 {
		startDate := fmt.Sprintf("%d-01-01", year)
		endDate := fmt.Sprintf("%d-01-01", year+1)
		if month != 0 {
			startDate = fmt.Sprintf("%d-%02d-01", year, month)
			if month == 12 {
				endDate = fmt.Sprintf("%d-01-01", year+1)
			} else {
				endDate = fmt.Sprintf("%d-%02d-01", year, month+1)
			}
		}
		query += fmt.Sprintf(" AND u.recording_date >= $%d AND u.recording_date < $%d", argIdx, argIdx+1)
		args = append(args, startDate, endDate)
		argIdx += 2
	}

	query += " ORDER BY u.recording_date DESC"
	query += fmt.Sprintf(" LIMIT $%d", argIdx)
	args = append(args, limit)

	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get total count
	countQuery := `
		SELECT count(*) AS total
		FROM uploads u
		LEFT JOIN LATERAL (
			SELECT ST_Collect(geom) AS geom
			FROM markers
			WHERE markers.trackid = u.track_id
		) m ON true
		WHERE m.geom && ST_MakeEnvelope($1, $2, $3, $4, 4326)`
	countArgs := []any{minLon, minLat, maxLon, maxLat}
	countRow, _ := queryRow(ctx, countQuery, countArgs...)
	total := 0
	if countRow != nil {
		if t, ok := countRow["total"]; ok {
			switch v := t.(type) {
			case int64:
				total = int(v)
			case float64:
				total = int(v)
			}
		}
	}

	tracks := make([]map[string]any, len(rows))
	for i, r := range rows {
		track := map[string]any{
			"id":          r["id"],
			"filename":    r["filename"],
			"track_id":    r["track_id"],
			"detector":    r["detector"],
			"file_size":   r["file_size"],
			"recording_date": r["recording_date"],
			"created_at":  r["created_at"],
		}

		// Add map URL for track view
		if trackID, ok := r["track_id"].(string); ok && trackID != "" {
			track["map_url"] = "https://simplemap.safecast.org/trackid/" + trackID
		}

		// Add location info if available
		if centroidLat, ok := r["centroid_lat"]; ok && centroidLat != nil {
			if centroidLon, ok := r["centroid_lon"]; ok && centroidLon != nil {
				track["centroid"] = map[string]any{
					"latitude":  centroidLat,
					"longitude": centroidLon,
				}
			}
		}

		// Add uploader information
		if internalUsername, ok := r["internal_username"]; ok && internalUsername != nil && internalUsername != "" {
			track["username"] = internalUsername
			track["uploader"] = map[string]any{
				"username": internalUsername,
				"email":    r["uploader_email"],
			}
		} else if username, ok := r["username"]; ok && username != nil && username != "" {
			track["username"] = username
		}

		tracks[i] = track
	}

	searchArea := country
	if searchArea == "" {
		searchArea = fmt.Sprintf("bbox:[%.2f,%.2f,%.2f,%.2f]", minLat, minLon, maxLat, maxLon)
	}

	result := map[string]any{
		"count":           len(tracks),
		"total_available": total,
		"source":          "database",
		"search_area":     searchArea,
		"bounding_box": map[string]any{
			"min_lat": minLat,
			"max_lat": maxLat,
			"min_lon": minLon,
			"max_lon": maxLon,
		},
		"filters": map[string]any{
			"year":  nilIfZero(year),
			"month": nilIfZero(month),
		},
		"tracks":             tracks,
		"_ai_hint":           "CRITICAL INSTRUCTIONS: (1) The 'unit' field indicates measurement units - CPM means 'counts per minute' NOT 'counts per second'. Always interpret and report CPM values as counts per minute. (2) Present all data in a purely scientific, factual manner. NEVER use personal pronouns (I, we, I'll, I'm, you, your), exclamations (!, Perfect, Great, Excellent), or conversational phrases. Format responses as objective statements: 'Latest reading: X CPM at location Y' NOT 'I found a reading of X CPM' or 'Perfect! The sensor shows...'. State only objective facts and measurements.",
		"_ai_generated_note": "This data was retrieved by an AI assistant using Safecast tools. The interpretation and presentation of this data may be influenced by the AI system.",
	}

	return jsonResult(result)
}

// toLower converts a string to lowercase
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		result[i] = c
	}
	return string(result)
}

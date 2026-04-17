package main

import "log"

// defaultTourStep is the seed shape — only the fields needed to bootstrap the
// table. Admin edits are protected by ON CONFLICT DO NOTHING, so re-running
// the seeder on every startup is safe.
type defaultTourStep struct {
	StepKey  string
	Selector string
	Center   bool
	TextEN   string
}

// defaultTourSteps mirrors the previous hardcoded JS array in map.html so the
// tour stays identical after migration. Keep the order stable; sort_order is
// assigned by the seeder in slice order.
var defaultTourSteps = []defaultTourStep{
	{
		StepKey:  "welcome",
		Selector: "#map",
		Center:   true,
		TextEN:   "Welcome to the Safecast map — an open platform where people worldwide share dosimeter readings for science, ecology, education, and safety. Each coloured dot is a real measurement with a timestamp and GPS location. The colour shows the dose rate — see the legend for the scale. Treat this as a community snapshot, not a precise survey.",
	},
	{
		StepKey:  "upload",
		Selector: ".upload-btn-container",
		TextEN:   "Upload your own radiation measurements. Supported formats include CSV, LOG, and KML files from devices such as bGeigie Nano, Radnote, and GammaScout. Your data appears on the map once processing is complete.",
	},
	{
		StepKey:  "login",
		Selector: "#loginButton",
		TextEN:   "Log in to link uploads to your account, track your measurement history, and access your profile page. Registration is free and open to everyone.",
	},
	{
		StepKey:  "legend",
		Selector: "#legend",
		TextEN:   "This colour scale shows the radiation dose rate. The scale adapts to the data on screen — hover it for a full explanation.",
	},
	{
		StepKey:  "ai_toggle",
		Selector: "#safecast-ai-toggle",
		TextEN:   "Ask our AI assistant anything about the map, the data, or radiation safety. It knows Safecast data in depth.",
	},
	{
		StepKey:  "track_insights",
		Selector: "#track-insights-btn",
		TextEN:   "When you click on a track, this button appears. It asks the AI to analyse that specific track — summarising the route, radiation levels, and any notable patterns it finds.",
	},
	{
		StepKey:  "date_slider",
		Selector: ".date-slider-box",
		TextEN:   "Drag the time slider to filter measurements by date. Switch between Year and Month mode for broad or fine-grained control, and press the reset button to show all data again.",
	},
	{
		StepKey:  "realtime_sensors",
		Selector: "#sfLive-row",
		TextEN:   "The large circles with numbers on the map are live readings from Safecast realtime sensors, updated every few minutes. The number shows the current dose rate in µSv/h, and the colour matches the legend scale. Click a circle to open its full measurement history. Uncheck this box to hide realtime sensors.",
	},
	{
		StepKey:  "gamma_spectrum",
		Selector: "#sfSpectrum-row",
		TextEN:   "Toggle gamma spectrum measurements. These come from spectral sensors that record the full energy distribution of detected radiation, not just a single dose rate. Click a spectrum point on the map to see its spectral distribution chart and identify which isotopes contributed to the reading.",
	},
	{
		StepKey:  "speed_filter",
		Selector: "#speed-filter-ctrl",
		TextEN:   "Filter measurements by travel speed: airplane, car, or walking. Use these to focus on the type of measurements you are interested in.",
	},
	{
		StepKey:  "search",
		Selector: "#searchInput",
		TextEN:   "Search for any city or address to jump straight to that location on the map.",
	},
}

// seedTourStepsDB idempotently inserts the default tour steps and their
// English text. Existing rows are never overwritten, so admin edits survive.
func seedTourStepsDB() {
	if db == nil || db.DB == nil {
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("[tour] Cannot begin seed transaction: %v", err)
		return
	}
	defer tx.Rollback()

	stepsInserted, textsInserted := 0, 0
	for i, step := range defaultTourSteps {
		res, err := tx.Exec(`
			INSERT INTO tour_steps (step_key, sort_order, selector, center)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (step_key) DO NOTHING`,
			step.StepKey, i+1, step.Selector, step.Center,
		)
		if err != nil {
			log.Printf("[tour] Seed step %s: %v", step.StepKey, err)
			continue
		}
		if rows, _ := res.RowsAffected(); rows > 0 {
			stepsInserted++
		}

		res, err = tx.Exec(`
			INSERT INTO translations (language_code, key, value)
			VALUES ('en', 'tour.' || $1 || '.text', $2)
			ON CONFLICT (language_code, key) DO NOTHING`,
			step.StepKey, step.TextEN,
		)
		if err != nil {
			log.Printf("[tour] Seed text %s: %v", step.StepKey, err)
			continue
		}
		if rows, _ := res.RowsAffected(); rows > 0 {
			textsInserted++
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[tour] Cannot commit seed transaction: %v", err)
		return
	}
	if stepsInserted > 0 || textsInserted > 0 {
		log.Printf("[tour] Seeded %d new steps and %d English text entries",
			stepsInserted, textsInserted)
	}
}

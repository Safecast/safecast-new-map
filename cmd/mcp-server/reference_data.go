package main

var referenceData = map[string]string{
	"units": `# Radiation Units

**µSv/h (microsieverts per hour)**
- Most common unit for dose rate
- Measures biological impact of radiation exposure over time
- 1 µSv/h = 0.001 mSv/h = 0.000001 Sv/h

**CPM (Counts Per Minute)**
- Raw count of radiation events detected
- Device-specific, varies by detector sensitivity
- Typical conversion: ~100 CPM ≈ 1 µSv/h (varies by device)

**Bq (Becquerel)**
- Measures radioactivity of a source
- 1 Bq = 1 decay per second
- Used for contamination levels

**Sv (Sievert)**
- SI unit for radiation dose
- 1 Sv = 1000 mSv = 1,000,000 µSv
- Accounts for biological effect
`,

	"dose_rates": `# Understanding Dose Rates

**Typical Dose Rate Ranges:**

**0.05 - 0.20 µSv/h**
- Normal background radiation
- Varies by location (altitude, geology)

**0.20 - 0.50 µSv/h**
- Slightly elevated, still safe
- Common in granite-rich areas

**0.50 - 1.00 µSv/h**
- Elevated but not immediately dangerous
- Should investigate source

**1.00 - 10.00 µSv/h**
- Clearly elevated
- Avoid prolonged exposure
- May indicate contamination

**> 10.00 µSv/h**
- High radiation area
- Minimize exposure time
- Seek expert advice
`,

	"safety_levels": `# Radiation Safety Levels

**WHO/ICRP Guidelines:**

**Annual Dose Limits:**
- General public: 1 mSv/year (≈0.11 µSv/h continuous)
- Occupational workers: 20 mSv/year (≈2.3 µSv/h work hours)

**Natural Background Radiation:**
- Global average: 2.4 mSv/year (≈0.27 µSv/h)
- Range: 1-10 mSv/year depending on location

**Medical Context:**
- Chest X-ray: ~0.02 mSv
- CT scan: 2-10 mSv
- Flight (10 hours): ~0.05 mSv

**Acute Exposure Effects:**
- < 100 mSv: No immediate symptoms
- 100-500 mSv: Increased cancer risk
- 500-1000 mSv: Radiation sickness possible
- > 1000 mSv: Severe radiation sickness
- > 5000 mSv: Often fatal

**Note:** These are acute doses, not dose rates.
`,

	"detectors": `# Radiation Detectors

**Geiger-Müller Counters**
- Most common type
- Detect gamma and beta radiation
- Examples: bGeigie, RadiaCode, GQ GMC-320
- Pros: Affordable, portable, reliable
- Cons: Cannot identify isotopes

**Scintillation Detectors**
- More sensitive than Geiger counters
- Can measure gamma energy
- Examples: RadiaCode-102, RadiaCode-103
- Pros: Can perform spectroscopy
- Cons: More expensive

**Semiconductor Detectors**
- High resolution spectroscopy
- Lab-grade precision
- Examples: CdTe, CZT detectors
- Pros: Best energy resolution
- Cons: Expensive, fragile

**Common Safecast Devices:**
- bGeigie Nano: Mobile monitoring device
- Pointcast: Fixed monitoring station
- Solarcast: Solar-powered fixed station
`,

	"background_levels": `# Natural Background Radiation

**Global Variations:**

**Low Background Areas:**
- 0.05-0.10 µSv/h
- Coastal/sea-level regions
- Sedimentary geology

**Average Background:**
- 0.10-0.20 µSv/h
- Most inhabited areas
- Mixed geology

**Elevated Natural Background:**
- 0.20-0.50 µSv/h
- Granite-rich areas (e.g., Cornwall, UK)
- High altitude (e.g., Denver, Colorado)
- Volcanic regions

**Very High Natural Background:**
- > 0.50 µSv/h
- Ramsar, Iran: up to 50 µSv/h (monazite sand)
- Guarapari, Brazil: up to 20 µSv/h (thorium)
- Kerala, India: up to 4 µSv/h (monazite)

**Cosmic Radiation:**
- Increases with altitude
- Sea level: ~0.03 µSv/h
- 10,000 ft (3000m): ~0.15 µSv/h
- Airline crew annual dose: 2-5 mSv/year
`,

	"isotopes": `# Common Radioactive Isotopes

**Naturally Occurring:**

**Potassium-40 (K-40)**
- Half-life: 1.25 billion years
- In bananas, salt, human body
- Beta and gamma emitter

**Radon-222 (Rn-222)**
- Half-life: 3.8 days
- From uranium decay in soil/rocks
- Major source of background radiation
- Alpha emitter (dangerous when inhaled)

**Uranium-238 (U-238)**
- Half-life: 4.5 billion years
- In granite, soil
- Alpha emitter

**Artificial/Fission Products:**

**Cesium-137 (Cs-137)**
- Half-life: 30 years
- Fukushima, Chernobyl contamination
- Gamma emitter (661 keV)
- Major long-term concern

**Iodine-131 (I-131)**
- Half-life: 8 days
- Released in nuclear accidents
- Short-lived but dangerous to thyroid
- Beta and gamma emitter

**Strontium-90 (Sr-90)**
- Half-life: 29 years
- Bone-seeking isotope
- Beta emitter (hard to detect)

**Cobalt-60 (Co-60)**
- Half-life: 5.3 years
- Medical/industrial use
- Strong gamma emitter (1173, 1332 keV)
`,
}

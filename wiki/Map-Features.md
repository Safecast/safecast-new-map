# Map Features

Guide to the interactive map, visualization, and user interface features.

[[Home|← Back to Home]]

---

## Overview

The Safecast New Map provides a rich, interactive mapping experience with:
- Clustered radiation markers
- Multiple coloring schemes
- Speed-based layer separation
- Spectrum viewer
- Real-time sensor display
- Multi-language support
- Print mode with QR codes

---

## Map Interface

### Base Layers

Choose from multiple base map layers:
- **OpenStreetMap** - Default, detailed street map
- **Google Satellite** - Satellite imagery
- Additional layers as configured

**Switch layers:** Click layer selector in top-right corner

### Map Controls

**Zoom:**
- Mouse scroll wheel
- +/- buttons
- Double-click to zoom in
- Shift+double-click to zoom out

**Pan:**
- Click and drag
- Arrow keys

**Locate:**
- Click location button to find your position
- Uses browser geolocation or GeoIP fallback

---

## Radiation Markers

### Marker Clustering

Markers are automatically clustered at different zoom levels:
- **High zoom:** Individual markers visible
- **Low zoom:** Clustered markers with count
- Click cluster to zoom in and see individual markers

### Speed-Based Layers

Measurements are separated by transport speed:
- **Walking** - Slow movement (< 10 km/h)
- **Driving** - Medium speed (10-100 km/h)
- **Flying** - Fast movement (> 100 km/h)

**Toggle layers:** Click layer control in top-right corner

### Coloring Schemes

Two visualization modes:

#### Scientific Gradient (`coloring=safecast`)

Continuous color gradient based on dose rate:
- **Blue** - Low radiation (background levels)
- **Green** - Moderate levels
- **Yellow** - Elevated levels
- **Orange** - High levels
- **Red** - Very high levels

**Use case:** Scientific analysis, research

#### Safety Bins (`coloring=chicha`)

Discrete color bins based on safety thresholds:
- **Green** - Safe (< 0.5 µSv/h)
- **Yellow** - Caution (0.5-1.0 µSv/h)
- **Orange** - Warning (1.0-10.0 µSv/h)
- **Red** - Danger (> 10.0 µSv/h)

**Use case:** Public safety, quick assessment

**Switch coloring:** Use `coloring` URL parameter or layer control

---

## Legend

The legend shows radiation levels and corresponding colors.

### Toggle Legend

**URL parameter:**
- `?legend=1` - Show legend (default)
- `?legend=0` - Hide legend

**UI control:** Click legend icon

### Legend Information

Displays:
- Color scheme explanation
- Unit display (µSv/h or µR)
- Safety thresholds (for chicha mode)
- Data source information

---

## Units

### Supported Units

**Microsieverts per hour (µSv/h):**
- SI unit for radiation dose rate
- Default unit
- International standard

**Microroentgen (µR):**
- Traditional unit
- Common in some regions
- Conversion: 1 µSv/h ≈ 100 µR/h

### Switch Units

**URL parameter:**
- `?unit=uSv` - Microsieverts
- `?unit=uR` - Microroentgen

**UI control:** Click unit selector

---

## Real-Time Sensors

### Display Sensors

Show real-time sensor readings on map:

**URL parameter:**
```
?show=rt
```

**Features:**
- Sensors appear as special markers
- Current dose rate displayed
- Sensor name and type shown
- Auto-refresh for latest readings

### Sensor Types

- **Pointcast** - Personal radiation detectors
- **Solarcast** - Solar-powered sensors
- **bGeigieZen** - Advanced bGeigie devices

### Sensor Status

**Online:**
- Green indicator
- Recent data (< 1 hour)

**Offline:**
- Gray indicator
- Stale data (> 1 hour)

---

## Spectrum Viewer

Interactive gamma spectrum analysis tool.

### Access Spectrum

1. Click on a marker with spectrum data
2. Click "View Spectrum" button
3. Spectrum viewer opens

### Spectrum Features

**Interactive Chart:**
- Canvas 2D rendering
- Smooth zoom and pan
- Hover tooltips for channel values

**Channel-to-Energy Conversion:**
```
Energy (keV) = A + B×channel + C×channel²
```

Calibration coefficients displayed from spectrum metadata.

**Isotope Identification:**
- Known isotope peaks overlaid
- Peak labels with isotope names
- Common isotopes: Cs-137, Co-60, K-40, I-131

**Energy Range:**
- Typically 0-3000 keV
- Configurable per spectrum
- Displayed in header

### Spectrum Controls

**Zoom:**
- Drag to select region
- Double-click to zoom in
- Right-click to zoom out

**Pan:**
- Click and drag

**Download:**
- CSV format
- JSON format
- Original format (.spe, .n42, etc.)

**Isotope Peaks:**
- Toggle peak overlays
- Select specific isotopes
- View peak details

---

## Location Search

### Search by Place Name

Search for cities, countries, or locations:

**UI:**
1. Click search box
2. Type location name
3. Select from results
4. Map navigates to location

**URL parameter:**
```
?place=Tokyo
?place=Fukushima%20Japan
?place=São%20Paulo
```

### Search by Coordinates

Enter specific latitude/longitude:

**UI:**
1. Click coordinate input button
2. Enter lat/lon
3. Map centers on coordinates

**URL parameter:**
```
?lat=35.6762&lon=139.6503
```

### GeoIP Location

Automatic location based on IP address:
- Fallback when geolocation not available
- Approximate location only
- Used for initial map center

---

## URL Parameters

Customize map views with URL parameters:

| Parameter | Values | Description |
|-----------|--------|-------------|
| `place` | City/country name | Navigate by location name |
| `lat` | -90 to 90 | Latitude coordinate |
| `lon` | -180 to 180 | Longitude coordinate |
| `zoom` | 1-18 | Zoom level |
| `minLat`, `minLon`, `maxLat`, `maxLon` | Coordinates | Define map bounds |
| `coloring` | safecast, chicha | Coloring scheme |
| `unit` | uSv, uR | Display units |
| `legend` | 1, 0 | Show/hide legend |
| `lang` | en, ja, de, etc. | Interface language (29 languages) |
| `layer` | OpenStreetMap, Google Satellite | Base map |
| `show` | rt | Show only realtime sensors |

### Parameter Priority

When loading URLs:
1. `place` takes priority over `lat/lon`
2. `lat/lon` takes priority over bounds
3. Other parameters applied on top

### Example URLs

**Tokyo with safety coloring:**
```
/?place=Tokyo&coloring=chicha&unit=uR
```

**Specific coordinates:**
```
/?lat=48.8566&lon=2.3522&zoom=13
```

**Real-time sensors only:**
```
/?show=rt
```

**Russian interface:**
```
/?lang=ru
```

**Embed without legend:**
```
/?legend=0
```

---

## Print Mode

Print-friendly map view with QR codes.

### Enable Print Mode

Click print button or add to URL:
```
/?print=1
```

### Features

**QR Codes:**
- Generated for each location
- Link back to interactive map
- Useful for field marking
- Include location coordinates

**Clean Layout:**
- Simplified map controls
- Legend always visible
- Optimized for printing
- A4/Letter friendly

**Use Cases:**
- Field surveys
- Location documentation
- Public presentations
- Educational materials

---

## User Profile Pages

### Access Profile

**Web:**
1. Log in
2. Click username → "Profile"

### Profile Information

Displays:
- Username
- Email address
- API key (with copy button)
- Registration date
- Last login
- Upload count

### Upload History

View all your uploaded files:
- Filename
- Upload date
- Measurement count
- Detector type
- Edit metadata (for your uploads)

### Account Actions

- Change password
- Regenerate API key
- View upload statistics
- Log out

---

## Language Support

### 29 Supported Languages

The interface supports:
- Arabic (ar)
- Bulgarian (bg)
- Czech (cs)
- Danish (da)
- German (de)
- Greek (el)
- English (en)
- Spanish (es)
- Persian (fa)
- Finnish (fi)
- French (fr)
- Hebrew (he)
- Hindi (hi)
- Hungarian (hu)
- Indonesian (id)
- Italian (it)
- Japanese (ja)
- Korean (ko)
- Malay (ms)
- Dutch (nl)
- Norwegian (no)
- Polish (pl)
- Portuguese (pt)
- Russian (ru)
- Swedish (sv)
- Thai (th)
- Turkish (tr)
- Ukrainian (uk)
- Vietnamese (vi)
- Chinese (zh)

### Language Selection

**URL parameter:**
```
/?lang=ja
```

**Priority:**
1. `?lang=` URL parameter
2. Browser `Accept-Language` header
3. English fallback

### Translated Components

All UI components are translated:
- Map legend
- AI assistant widget
- Login/register modals
- Search bar
- Spectrum viewer
- Profile page
- Coordinate input dialog
- Admin panel

**Brand Rule:** "Safecast" remains untranslated in all languages.

See [Internationalization](Internationalization) for details.

---

## Performance Features

### Tile Caching

LRU cache for map tiles:
- 1000 entries default
- Improves map loading speed
- Reduces database queries
- Automatic cache invalidation

### Spatial Queries

PostGIS-optimized queries:
- `ST_DWithin` for radius searches
- `ST_Expand` for bounding boxes
- GIST spatial indexes
- Efficient large dataset handling

### Track Streaming

Stream tracks by ID range:
- Efficient pagination
- No full table scans
- Progressive loading
- Memory-efficient

---

## Accessibility

### Keyboard Navigation

- **Tab** - Navigate controls
- **Enter** - Activate buttons
- **Space** - Toggle checkboxes
- **Arrow keys** - Pan map
- **+/-** - Zoom in/out

### Screen Readers

- ARIA labels on controls
- Semantic HTML structure
- Keyboard-accessible features
- Focus indicators

### Mobile Support

- Responsive design
- Touch-friendly controls
- Mobile-optimized layout
- Gesture support (pinch zoom, etc.)

---

## Customization

### Default Settings

Configure defaults via flags:

```bash
./safecast-new-map \
  -default-lat 35.6762 \
  -default-lon 139.6503 \
  -default-zoom 12 \
  -default-layer "Google Satellite"
```

### Custom Base Layers

Add custom tile layers via configuration.

### Custom Color Schemes

Modify coloring schemes in source code.

---

## Best Practices

### Map Usage

1. **Use appropriate zoom** for data density
2. **Toggle layers** to reduce clutter
3. **Use safety bins** for public communication
4. **Use scientific gradient** for analysis

### Performance

1. **Use bounding box queries** for large areas
2. **Enable caching** for repeated views
3. **Limit markers** at low zoom levels
4. **Use spatial indexes** (PostgreSQL)

### Accessibility

1. **Test with keyboard** navigation
2. **Use semantic HTML** for custom elements
3. **Provide text alternatives** for visual data
4. **Ensure sufficient color contrast**

---

## Troubleshooting

### Map Not Loading

**Check database connection:**
```bash
curl http://localhost:8765/api/stats
```

**Check JavaScript console:**
- Open browser DevTools
- Look for errors
- Check network requests

### Markers Not Appearing

**Verify data exists:**
```bash
curl "http://localhost:8765/api/radiation?lat=35.6762&lon=139.6503"
```

**Check layer visibility:**
- Ensure correct layer is selected
- Check speed-based layer filters

### Spectrum Viewer Not Working

**Verify spectrum data:**
```bash
curl "http://localhost:8765/api/spectrum/12345"
```

**Check browser support:**
- Requires Canvas 2D support
- Modern browser required
- WebGL optional

---

## See Also

- [Spectral Analysis](Spectral-Analysis) - Spectrum features and analysis
- [Internationalization](Internationalization) - Language support
- [API Documentation](API-Documentation) - Map data API
- [Configuration Reference](Configuration-Reference) - Default settings

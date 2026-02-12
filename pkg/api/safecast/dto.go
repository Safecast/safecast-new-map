package safecast

// Core types are version-neutral. Adapter and v2 translate to/from these.

// MeasurementCore is the internal representation of a measurement.
type MeasurementCore struct {
	ID                   int64
	Value                float64
	Height               *float64
	UserID               *int64
	Unit                 string
	DeviceID             *int64
	LocationName         string
	OriginalID           *int64
	CapturedAt           string // ISO8601 or Rails format
	DevicetypeID         *int64
	SensorID             *int64
	ChannelID            *int64
	StationID            *int64
	MeasurementImportID  *int64
	Latitude             float64
	Longitude            float64
}

// MeasurementFiltersCore holds version-neutral filter parameters for listing measurements.
type MeasurementFiltersCore struct {
	Latitude            *float64
	Longitude           *float64
	Distance            *int
	CapturedAfter       *string
	CapturedBefore      *string
	UserID              *int64
	DeviceID            *int64
	MeasurementImportID *int64
	OriginalID          *int64
	Since               *string
	Until               *string
	Unit                string
	Order               string
	Page                int
	PerPage             int
}

// BgeigieImportCore is the internal representation of a bGeigie import.
type BgeigieImportCore struct {
	ID                int64
	UserID           *int64
	Approved         bool
	CreatedAt        string
	UpdatedAt        string
	MeasurementsCount int
	MD5Sum           string
	Name             string
	Status           string
	SourceURL        string
}

// UserCore is the internal representation of a user.
type UserCore struct {
	ID                 int64
	Name               string
	Email              string
	MeasurementsCount  int
	APIKey             string // Only populated when current user == this user
}

// DeviceCore is the internal representation of a device.
type DeviceCore struct {
	ID           int64
	Manufacturer string
	Model        string
	Sensor       string
}

// --- Rails-compatible (v1/adapter) response structs ---

// MeasurementRails matches the Rails Measurement serializable_hash for backward compatibility.
type MeasurementRails struct {
	ID                   int64   `json:"id"`
	Value                float64 `json:"value"`
	Height               *float64 `json:"height"`
	UserID               *int64  `json:"user_id"`
	Unit                 string  `json:"unit"`
	DeviceID             *int64  `json:"device_id"`
	LocationName         string  `json:"location_name"`
	OriginalID           *int64  `json:"original_id"`
	CapturedAt           string  `json:"captured_at"`
	DevicetypeID         *int64  `json:"devicetype_id"`
	SensorID             *int64  `json:"sensor_id"`
	ChannelID            *int64  `json:"channel_id"`
	StationID            *int64  `json:"station_id"`
	MeasurementImportID  *int64  `json:"measurement_import_id"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
}

// BgeigieImportRails matches the Rails BgeigieImport JSON shape.
type BgeigieImportRails struct {
	ID                int64             `json:"id"`
	UserID            *int64            `json:"user_id"`
	Approved          bool              `json:"approved"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
	MeasurementsCount int               `json:"measurements_count"`
	MD5Sum            string            `json:"md5sum"`
	Name              string            `json:"name"`
	Status            string            `json:"status"`
	Source            *SourceRails      `json:"source,omitempty"`
}

// SourceRails wraps the source URL for bgeigie_imports.
type SourceRails struct {
	URL string `json:"url"`
}

// UserRails matches the Rails user JSON shape.
type UserRails struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	MeasurementsCount  int     `json:"measurements_count,omitempty"`
	AuthenticationToken string `json:"authentication_token,omitempty"`
}

// DeviceRails matches the Rails device JSON shape.
type DeviceRails struct {
	ID           int64  `json:"id"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Sensor       string `json:"sensor"`
}

// RootRails is the root JSON response.
type RootRails struct {
	Name             string   `json:"name"`
	URI              string   `json:"uri"`
	SubresourceURIs   []string `json:"subresource_uris"`
}

// CountRails is the count response.
type CountRails struct {
	Count int64 `json:"count"`
}

// ErrorsRails is the validation error envelope.
type ErrorsRails struct {
	Errors map[string][]string `json:"errors"`
}

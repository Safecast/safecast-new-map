package tiles

import (
	"context"
	"errors"
	"strconv"
	"time"
)

// ErrUnknownProcessor is returned when the URL names a processor that was
// not registered. The handler turns it into HTTP 404.
var ErrUnknownProcessor = errors.New("tiles: unknown processor")

// Request is one tile plus the filters that belong in the cache key.
// From and To are unix seconds. Zero means that side of the window is open.
// Speed selects the same bands as the map's speed checkboxes:
// "a" is at least 200 km/h, "c" is 7–200 km/h, "w" is below 7 km/h.
// Letters combine ("cw" is the map default: cars and pedestrians, no planes).
// "all" leaves the speed column unconstrained. "none" matches no row.
type Request struct {
	Tile  TileID
	From  int64
	To    int64
	Speed string
}

// Result is the opaque payload a processor produced.
type Result struct {
	Body        []byte
	ContentType string
}

// Processor builds the bytes for one tile. Implementations must be safe for
// concurrent Render calls.
type Processor interface {
	Name() string
	ContentType() string
	// Version is part of the cache key so a ramp or statistic change
	// does not serve tiles painted by the previous code.
	Version() string
	Render(ctx context.Context, req Request) (Result, error)
}

// QuantizeFrom snaps a timestamp to the start of its UTC day.
// Zero stays zero, which means "no lower bound".
func QuantizeFrom(unix int64) int64 {
	if unix <= 0 {
		return 0
	}
	t := time.Unix(unix, 0).UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix()
}

// QuantizeTo snaps a timestamp to the last second of its UTC day.
// Zero stays zero, which means "no upper bound".
func QuantizeTo(unix int64) int64 {
	if unix <= 0 {
		return 0
	}
	start := QuantizeFrom(unix)
	return start + 24*60*60 - 1
}

func normalizeRequest(req Request) Request {
	req.From = QuantizeFrom(req.From)
	req.To = QuantizeTo(req.To)
	req.Speed = cleanSpeed(req.Speed)
	return req
}

func cleanSpeed(s string) string {
	if s == "" || s == "all" {
		return "all"
	}
	if s == "none" {
		return "none"
	}
	var plane, car, walk bool
	for _, r := range s {
		switch r {
		case 'a':
			plane = true
		case 'c':
			car = true
		case 'w':
			walk = true
		default:
			return "all"
		}
	}
	out := ""
	if plane {
		out += "a"
	}
	if car {
		out += "c"
	}
	if walk {
		out += "w"
	}
	if out == "" {
		return "none"
	}
	return out
}

func cacheKey(name, version string, req Request) string {
	return tilePrefix(req.Tile) + "/" + name + "_" + version + "_" + strconv.FormatInt(req.From, 10) + "_" + strconv.FormatInt(req.To, 10) + "_" + req.Speed
}

func tilePrefix(id TileID) string {
	return strconv.Itoa(id.Z) + "/" + strconv.Itoa(id.X) + "/" + strconv.Itoa(id.Y)
}

package kriging

import "sort"

// DefaultNeighbors is the local system size. Adjacent grid points use almost
// the same set because each keeps the observations closest in the covariance metric.
const DefaultNeighbors = 32

// NeighborhoodSelector keeps the observations most correlated with one target.
type NeighborhoodSelector struct {
	K     int
	Model CovarianceModel
}

// Select returns at most K observations, nearest in the model metric first.
// Ties break on identity fields so the same candidates always produce the same order.
func (s NeighborhoodSelector) Select(target Target, obs []Observation) []Observation {
	if len(obs) == 0 || s.Model == nil {
		return nil
	}
	k := s.K
	if k <= 0 {
		k = DefaultNeighbors
	}
	if k >= len(obs) {
		out := append([]Observation(nil), obs...)
		sort.Slice(out, func(i, j int) bool {
			return closer(s.Model, target, out[i], out[j])
		})
		return out
	}

	type item struct {
		o Observation
	}
	best := make([]item, 0, k)
	ordered := false
	for _, o := range obs {
		if len(best) < k {
			best = append(best, item{o: o})
			if len(best) == k {
				sort.Slice(best, func(i, j int) bool {
					return closer(s.Model, target, best[i].o, best[j].o)
				})
				ordered = true
			}
			continue
		}
		if !closer(s.Model, target, o, best[k-1].o) {
			continue
		}
		best[k-1].o = o
		sort.Slice(best, func(i, j int) bool {
			return closer(s.Model, target, best[i].o, best[j].o)
		})
	}
	if !ordered {
		sort.Slice(best, func(i, j int) bool {
			return closer(s.Model, target, best[i].o, best[j].o)
		})
	}
	out := make([]Observation, len(best))
	for i := range best {
		out[i] = best[i].o
	}
	return out
}

func closer(model CovarianceModel, target Target, a, b Observation) bool {
	da := model.Distance(target, a)
	db := model.Distance(target, b)
	if da != db {
		return da < db
	}
	if a.TrackID != b.TrackID {
		return a.TrackID < b.TrackID
	}
	if a.DeviceID != b.DeviceID {
		return a.DeviceID < b.DeviceID
	}
	if a.TimeUnix != b.TimeUnix {
		return a.TimeUnix < b.TimeUnix
	}
	if a.Lat != b.Lat {
		return a.Lat < b.Lat
	}
	if a.Lon != b.Lon {
		return a.Lon < b.Lon
	}
	return a.DoseRate < b.DoseRate
}

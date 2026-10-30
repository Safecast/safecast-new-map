package kriging

import "math"

// Solution is one ordinary-kriging fit on the log-dose scale.
// Weights sum to 1 when OK is true. Variance is the latent-field variance.
type Solution struct {
	MeanLog  float64
	Variance float64
	Weights  []float64
}

// Solve fits ordinary kriging for the latent log-dose at target.
// The covariance block is factored with Cholesky. The unbiasedness constraint
// makes the full kriging matrix indefinite, so it is applied after the factor
// rather than by dropping the constraint. A failed factor returns ok = false
// and the caller suppresses the cell.
func Solve(model CovarianceModel, target Target, obs []Observation) (Solution, bool) {
	n := len(obs)
	if n == 0 || model == nil || model.FieldVariance() <= 0 || model.MeasurementNugget() < 0 {
		return Solution{}, false
	}
	if exp, ok := model.(ExponentialModel); ok && !exp.usable() {
		return Solution{}, false
	}
	for i := range obs {
		if obs[i].DoseRate <= 0 || math.IsNaN(obs[i].DoseRate) || math.IsInf(obs[i].DoseRate, 0) {
			return Solution{}, false
		}
	}

	cmat := make([]float64, n*n)
	y := make([]float64, n)
	c := make([]float64, n)
	targetObs := Observation{Lat: target.Lat, Lon: target.Lon, TimeUnix: target.TimeUnix}
	nugget := model.MeasurementNugget()
	for i := range obs {
		y[i] = math.Log(obs[i].DoseRate)
		c[i] = model.Pair(targetObs, obs[i])
		for j := 0; j <= i; j++ {
			v := model.Pair(obs[i], obs[j])
			if i == j {
				v += nugget
			}
			cmat[i*n+j] = v
			cmat[j*n+i] = v
		}
	}
	if !cholesky(cmat, n) {
		return Solution{}, false
	}

	// C λ + μ 1 = c and 1'λ = 1.
	// λ = u - μ v, with C u = c and C v = 1.
	u := solveCholesky(cmat, n, c)
	ones := make([]float64, n)
	for i := range ones {
		ones[i] = 1
	}
	v := solveCholesky(cmat, n, ones)
	sumU := 0.0
	sumV := 0.0
	for i := 0; i < n; i++ {
		sumU += u[i]
		sumV += v[i]
	}
	if sumV == 0 || math.IsNaN(sumV) || math.IsInf(sumV, 0) {
		return Solution{}, false
	}
	mu := (sumU - 1) / sumV
	weights := make([]float64, n)
	sumW := 0.0
	mean := 0.0
	dotC := 0.0
	for i := 0; i < n; i++ {
		weights[i] = u[i] - mu*v[i]
		if math.IsNaN(weights[i]) || math.IsInf(weights[i], 0) {
			return Solution{}, false
		}
		sumW += weights[i]
		mean += weights[i] * y[i]
		dotC += c[i] * weights[i]
	}
	if math.Abs(sumW-1) > 1e-6 {
		return Solution{}, false
	}
	variance := model.FieldVariance() - dotC - mu
	if math.IsNaN(variance) || math.IsInf(variance, 0) || math.IsNaN(mean) || math.IsInf(mean, 0) {
		return Solution{}, false
	}
	// A tiny negative value is factor noise. A large one means the system is wrong.
	if variance < 0 {
		if variance > -1e-8*model.FieldVariance() {
			variance = 0
		} else {
			return Solution{}, false
		}
	}
	return Solution{MeanLog: mean, Variance: variance, Weights: weights}, true
}

// cholesky overwrites a with its lower factor. a must be symmetric.
func cholesky(a []float64, n int) bool {
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			sum := a[i*n+j]
			for k := 0; k < j; k++ {
				sum -= a[i*n+k] * a[j*n+k]
			}
			if i == j {
				if sum <= 1e-15 {
					return false
				}
				a[i*n+j] = math.Sqrt(sum)
			} else {
				diag := a[j*n+j]
				if diag == 0 {
					return false
				}
				a[i*n+j] = sum / diag
			}
		}
	}
	return true
}

func solveCholesky(lower []float64, n int, b []float64) []float64 {
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := b[i]
		for k := 0; k < i; k++ {
			sum -= lower[i*n+k] * y[k]
		}
		y[i] = sum / lower[i*n+i]
	}
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := y[i]
		for k := i + 1; k < n; k++ {
			sum -= lower[k*n+i] * x[k]
		}
		x[i] = sum / lower[i*n+i]
	}
	return x
}

package status

import "github.com/IgorBayerl/nanovision/internal/config"

// ClassifyHigherIsBetter classifies a value where higher is better (e.g. coverage %).
//
//	val < Min  → Danger
//	Min ≤ val ≤ Max → Warning
//	val > Max  → Safe
//
// Returns ("", false) when band is nil.
func ClassifyHigherIsBetter(val float64, band *config.Band) (RiskLevel, bool) {
	if band == nil {
		return "", false
	}
	if val < band.Min {
		return RiskDanger, true
	}
	if val <= band.Max {
		return RiskWarning, true
	}
	return RiskSafe, true
}

// ClassifyLowerIsBetter classifies a value where lower is better (e.g. complexity).
//
//	val < Min  → Safe
//	Min ≤ val ≤ Max → Warning
//	val > Max  → Danger
//
// Returns ("", false) when band is nil.
func ClassifyLowerIsBetter(val float64, band *config.Band) (RiskLevel, bool) {
	if band == nil {
		return "", false
	}
	if val < band.Min {
		return RiskSafe, true
	}
	if val <= band.Max {
		return RiskWarning, true
	}
	return RiskDanger, true
}

// Classify is an alias for ClassifyHigherIsBetter, kept for backward
// compatibility with downstream code (e.g. reporters).
func Classify(pct float64, band *config.Band) (level RiskLevel, show bool) {
	return ClassifyHigherIsBetter(pct, band)
}

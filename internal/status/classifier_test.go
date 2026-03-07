package status

import (
	"testing"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestClassifyHigherIsBetter(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}

	tests := []struct {
		name     string
		val      float64
		band     *config.Band
		wantLvl  RiskLevel
		wantShow bool
	}{
		{"nil band", 50, nil, "", false},
		{"below min", 59.99, band, RiskDanger, true},
		{"at min", 60, band, RiskWarning, true},
		{"mid range", 70, band, RiskWarning, true},
		{"at max", 80, band, RiskWarning, true},
		{"above max", 80.01, band, RiskSafe, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl, show := ClassifyHigherIsBetter(tt.val, tt.band)
			assert.Equal(t, tt.wantLvl, lvl)
			assert.Equal(t, tt.wantShow, show)
		})
	}
}

func TestClassifyLowerIsBetter(t *testing.T) {
	band := &config.Band{Min: 5, Max: 10}

	tests := []struct {
		name     string
		val      float64
		band     *config.Band
		wantLvl  RiskLevel
		wantShow bool
	}{
		{"nil band", 3, nil, "", false},
		{"below min - safe", 4.99, band, RiskSafe, true},
		{"at min", 5, band, RiskWarning, true},
		{"mid range", 7, band, RiskWarning, true},
		{"at max", 10, band, RiskWarning, true},
		{"above max - danger", 10.01, band, RiskDanger, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl, show := ClassifyLowerIsBetter(tt.val, tt.band)
			assert.Equal(t, tt.wantLvl, lvl)
			assert.Equal(t, tt.wantShow, show)
		})
	}
}

func TestClassifyBackwardCompat(t *testing.T) {
	band := &config.Band{Min: 60, Max: 80}

	// Classify should behave identically to ClassifyHigherIsBetter.
	lvl1, s1 := Classify(50, band)
	lvl2, s2 := ClassifyHigherIsBetter(50, band)
	assert.Equal(t, lvl1, lvl2)
	assert.Equal(t, s1, s2)
}

package data_test

import (
	"testing"

	"github.com/narinsak-u/backend/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestThailandProvinceRegion_KnownProvinces(t *testing.T) {
	tests := []struct {
		province string
		expected string
	}{
		{"กรุงเทพมหานคร", "Central"},
		{"ชลบุรี", "East"},
		{"เชียงใหม่", "North"},
		{"กาญจนบุรี", "West"},
		{"ขอนแก่น", "North-east"},
		{"ภูเก็ต", "South"},
	}
	for _, tt := range tests {
		t.Run(tt.province, func(t *testing.T) {
			assert.Equal(t, tt.expected, data.ThailandProvinceRegion[tt.province])
		})
	}
}

func TestThailandProvinceRegion_UnknownProvince(t *testing.T) {
	assert.Empty(t, data.ThailandProvinceRegion["Unknown"])
}

func TestThailandProvinceRegion_NotEmpty(t *testing.T) {
	assert.Greater(t, len(data.ThailandProvinceRegion), 70)
}

func TestThailandProvinceRegion_AllValuesAreValid(t *testing.T) {
	validRegions := map[string]bool{"Central": true, "East": true, "North": true, "West": true, "North-east": true, "South": true}
	for province, region := range data.ThailandProvinceRegion {
		assert.True(t, validRegions[region], "province %q has unknown region %q", province, region)
	}
}

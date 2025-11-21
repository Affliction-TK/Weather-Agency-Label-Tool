package main

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lon1     float64
		lat1     float64
		lon2     float64
		lat2     float64
		expected float64
		delta    float64
	}{
		{
			name:     "Beijing to Shanghai",
			lon1:     116.4074,
			lat1:     39.9042,
			lon2:     121.4737,
			lat2:     31.2304,
			expected: 1067.0, // Approximately 1067 km
			delta:    50.0,   // Allow 50km margin
		},
		{
			name:     "Same location",
			lon1:     116.4074,
			lat1:     39.9042,
			lon2:     116.4074,
			lat2:     39.9042,
			expected: 0.0,
			delta:    0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := haversineDistance(tt.lon1, tt.lat1, tt.lon2, tt.lat2)
			diff := math.Abs(result - tt.expected)
			if diff > tt.delta {
				t.Errorf("haversineDistance() = %v, expected %v ± %v, diff = %v",
					result, tt.expected, tt.delta, diff)
			}
		})
	}
}

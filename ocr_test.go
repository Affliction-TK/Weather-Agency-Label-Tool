package main

import (
	"testing"
)

func TestNormalizeTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard format with hyphens",
			input:    "2024-01-15 14:30:45",
			expected: "2024-01-15 14:30:45",
		},
		{
			name:     "Chinese format",
			input:    "2024年01月15日 14:30:45",
			expected: "2024-01-15 14:30:45",
		},
		{
			name:     "Compact format without spaces",
			input:    "20240115143045",
			expected: "2024-01-15 14:30:45",
		},
		{
			name:     "Compact format with space",
			input:    "20240115 143045",
			expected: "2024-01-15 14:30:45",
		},
		{
			name:     "OCR error - missing space between date and time",
			input:    "2023-07-2015:47",
			expected: "2023-07-20 15:47:00",
		},
		{
			name:     "Single digit day - missing space",
			input:    "2023-07-915:47",
			expected: "2023-07-09 15:47:00",
		},
		{
			name:     "Single digit month",
			input:    "2023-7-20 15:47",
			expected: "2023-07-20 15:47:00",
		},
		{
			name:     "Single digit hour",
			input:    "2023-07-20 9:47",
			expected: "2023-07-20 09:47:00",
		},
		{
			name:     "Date only",
			input:    "2023-7-20",
			expected: "2023-07-20 00:00:00",
		},
		{
			name:     "Slash format with single digits",
			input:    "2023/7/9 9:30",
			expected: "2023-07-09 09:30:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTime(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTime(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanLocationText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal location",
			input:    "北京市朝阳区",
			expected: "北京市朝阳区",
		},
		{
			name:     "Location with extra spaces",
			input:    "北京市  朝阳区  ",
			expected: "北京市 朝阳区",
		},
		{
			name:     "Location with leading/trailing spaces",
			input:    "  西城区德胜门监测站  ",
			expected: "西城区德胜门监测站",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanLocationText(tt.input)
			if result != tt.expected {
				t.Errorf("cleanLocationText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractTimeAndLocation(t *testing.T) {
	tests := []struct {
		name         string
		ocrResponse  *BaiduOCRResponse
		expectedTime string
		expectedLoc  string
		expectedStd  bool
	}{
		{
			name: "Standard image with time and location",
			ocrResponse: &BaiduOCRResponse{
				WordsResultNum: 2,
				WordsResult: []BaiduOCRWord{
					{Words: "2024-01-15 14:30:45"},
					{Words: "北京市朝阳区监测站"},
				},
			},
			expectedTime: "2024-01-15 14:30:45",
			expectedLoc:  "北京市朝阳区监测站",
			expectedStd:  true,
		},
		{
			name: "Image with time only",
			ocrResponse: &BaiduOCRResponse{
				WordsResultNum: 1,
				WordsResult: []BaiduOCRWord{
					{Words: "2024-01-15 14:30:45"},
				},
			},
			expectedTime: "2024-01-15 14:30:45",
			expectedLoc:  "",
			expectedStd:  false,
		},
		{
			name: "Image with location only",
			ocrResponse: &BaiduOCRResponse{
				WordsResultNum: 1,
				WordsResult: []BaiduOCRWord{
					{Words: "北京市朝阳区监测站"},
				},
			},
			expectedTime: "",
			expectedLoc:  "北京市朝阳区监测站",
			expectedStd:  false,
		},
		{
			name: "Empty OCR response",
			ocrResponse: &BaiduOCRResponse{
				WordsResultNum: 0,
				WordsResult:    []BaiduOCRWord{},
			},
			expectedTime: "",
			expectedLoc:  "",
			expectedStd:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTimeAndLocation(tt.ocrResponse)

			if result.Time != tt.expectedTime {
				t.Errorf("Time = %q, want %q", result.Time, tt.expectedTime)
			}

			if result.Location != tt.expectedLoc {
				t.Errorf("Location = %q, want %q", result.Location, tt.expectedLoc)
			}

			if result.IsStandard != tt.expectedStd {
				t.Errorf("IsStandard = %v, want %v", result.IsStandard, tt.expectedStd)
			}
		})
	}
}

func TestNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "Empty string returns nil",
			input:    "",
			expected: nil,
		},
		{
			name:     "Non-empty string returns itself",
			input:    "test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nullString(tt.input)
			if result != tt.expected {
				t.Errorf("nullString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

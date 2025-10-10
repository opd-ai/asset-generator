package cmd

import (
	"testing"
)

func TestParseLoraParameters(t *testing.T) {
	tests := []struct {
		name            string
		loras           []string
		explicitWeights []float64
		defaultWeight   string
		wantResult      map[string]float64
		wantErr         bool
	}{
		{
			name:          "single LoRA with inline weight",
			loras:         []string{"anime-style:0.8"},
			defaultWeight: "1.0",
			wantResult:    map[string]float64{"anime-style": 0.8},
			wantErr:       false,
		},
		{
			name:          "single LoRA without weight (uses default)",
			loras:         []string{"anime-style"},
			defaultWeight: "1.0",
			wantResult:    map[string]float64{"anime-style": 1.0},
			wantErr:       false,
		},
		{
			name:          "multiple LoRAs with inline weights",
			loras:         []string{"anime-style:0.8", "detailed-faces:0.6", "cyberpunk:1.2"},
			defaultWeight: "1.0",
			wantResult: map[string]float64{
				"anime-style":    0.8,
				"detailed-faces": 0.6,
				"cyberpunk":      1.2,
			},
			wantErr: false,
		},
		{
			name:          "mixed inline and default weights",
			loras:         []string{"anime-style:0.8", "detailed-faces", "cyberpunk:1.2"},
			defaultWeight: "1.0",
			wantResult: map[string]float64{
				"anime-style":    0.8,
				"detailed-faces": 1.0,
				"cyberpunk":      1.2,
			},
			wantErr: false,
		},
		{
			name:            "explicit weights override default",
			loras:           []string{"anime-style", "detailed-faces"},
			explicitWeights: []float64{0.7, 0.9},
			defaultWeight:   "1.0",
			wantResult: map[string]float64{
				"anime-style":    0.7,
				"detailed-faces": 0.9,
			},
			wantErr: false,
		},
		{
			name:            "inline weights override explicit weights",
			loras:           []string{"anime-style:0.8", "detailed-faces"},
			explicitWeights: []float64{0.7, 0.9},
			defaultWeight:   "1.0",
			wantResult: map[string]float64{
				"anime-style":    0.8, // inline weight used
				"detailed-faces": 0.9, // explicit weight used
			},
			wantErr: false,
		},
		{
			name:          "custom default weight",
			loras:         []string{"anime-style", "detailed-faces"},
			defaultWeight: "0.75",
			wantResult: map[string]float64{
				"anime-style":    0.75,
				"detailed-faces": 0.75,
			},
			wantErr: false,
		},
		{
			name:          "negative weight (valid use case for style removal)",
			loras:         []string{"unwanted-style:-0.5"},
			defaultWeight: "1.0",
			wantResult:    map[string]float64{"unwanted-style": -0.5},
			wantErr:       false,
		},
		{
			name:          "weight with spaces",
			loras:         []string{" anime-style : 0.8 ", "  detailed-faces  "},
			defaultWeight: "1.0",
			wantResult: map[string]float64{
				"anime-style":    0.8,
				"detailed-faces": 1.0,
			},
			wantErr: false,
		},
		{
			name:          "empty LoRA list",
			loras:         []string{},
			defaultWeight: "1.0",
			wantResult:    nil,
			wantErr:       false,
		},
		{
			name:          "nil LoRA list",
			loras:         nil,
			defaultWeight: "1.0",
			wantResult:    nil,
			wantErr:       false,
		},
		{
			name:          "invalid weight format",
			loras:         []string{"anime-style:abc"},
			defaultWeight: "1.0",
			wantResult:    nil,
			wantErr:       true,
		},
		{
			name:          "invalid default weight",
			loras:         []string{"anime-style"},
			defaultWeight: "invalid",
			wantResult:    nil,
			wantErr:       true,
		},
		{
			name:          "empty LoRA name",
			loras:         []string{":0.8"},
			defaultWeight: "1.0",
			wantResult:    nil,
			wantErr:       true,
		},
		{
			name:          "weight out of range (too high)",
			loras:         []string{"anime-style:10.0"},
			defaultWeight: "1.0",
			wantResult:    nil,
			wantErr:       true,
		},
		{
			name:          "weight out of range (too low)",
			loras:         []string{"anime-style:-5.0"},
			defaultWeight: "1.0",
			wantResult:    nil,
			wantErr:       true,
		},
		{
			name:          "decimal weight with leading zero",
			loras:         []string{"anime-style:0.125"},
			defaultWeight: "1.0",
			wantResult:    map[string]float64{"anime-style": 0.125},
			wantErr:       false,
		},
		{
			name:          "integer weight",
			loras:         []string{"anime-style:2"},
			defaultWeight: "1.0",
			wantResult:    map[string]float64{"anime-style": 2.0},
			wantErr:       false,
		},
		{
			name:          "skip empty strings in list",
			loras:         []string{"anime-style:0.8", "", "detailed-faces:0.6"},
			defaultWeight: "1.0",
			wantResult: map[string]float64{
				"anime-style":    0.8,
				"detailed-faces": 0.6,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseLoraParameters(tt.loras, tt.explicitWeights, tt.defaultWeight)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseLoraParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if both are nil
				if result == nil && tt.wantResult == nil {
					return
				}

				// Check if one is nil and the other isn't
				if (result == nil) != (tt.wantResult == nil) {
					t.Errorf("parseLoraParameters() result = %v, want %v", result, tt.wantResult)
					return
				}

				// Check if maps have the same length
				if len(result) != len(tt.wantResult) {
					t.Errorf("parseLoraParameters() result length = %d, want %d", len(result), len(tt.wantResult))
					return
				}

				// Check each key-value pair
				for key, wantVal := range tt.wantResult {
					gotVal, exists := result[key]
					if !exists {
						t.Errorf("parseLoraParameters() missing key %s", key)
						continue
					}
					if gotVal != wantVal {
						t.Errorf("parseLoraParameters() result[%s] = %v, want %v", key, gotVal, wantVal)
					}
				}
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "positive decimal",
			input:   "0.8",
			want:    0.8,
			wantErr: false,
		},
		{
			name:    "negative decimal",
			input:   "-0.5",
			want:    -0.5,
			wantErr: false,
		},
		{
			name:    "integer",
			input:   "2",
			want:    2.0,
			wantErr: false,
		},
		{
			name:    "with spaces",
			input:   "  1.5  ",
			want:    1.5,
			wantErr: false,
		},
		{
			name:    "zero",
			input:   "0",
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "multiple decimals",
			input:   "1.2.3",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("parseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

package output

import (
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   OutputFormat
	}{
		{"table format", "table", FormatTable},
		{"json format", "json", FormatJSON},
		{"yaml format", "yaml", FormatYAML},
		{"uppercase", "JSON", FormatJSON},
		{"mixed case", "TaBLe", FormatTable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(tt.format)
			if f.format != tt.want {
				t.Errorf("NewFormatter() format = %v, want %v", f.format, tt.want)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	f := NewFormatter("json")

	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "simple map",
			data: map[string]interface{}{
				"key": "value",
				"num": 42,
			},
			wantErr: false,
		},
		{
			name: "slice",
			data: []interface{}{
				"item1", "item2", "item3",
			},
			wantErr: false,
		},
		{
			name: "nested structure",
			data: map[string]interface{}{
				"nested": map[string]interface{}{
					"deep": "value",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Error("Expected non-empty result")
			}
			// JSON arrays start with [ not {
			if !tt.wantErr && !strings.Contains(result, "[") && !strings.Contains(result, "{") {
				t.Error("Expected JSON output to contain brackets or braces")
			}
		})
	}
}

func TestFormatYAML(t *testing.T) {
	f := NewFormatter("yaml")

	data := map[string]interface{}{
		"key":   "value",
		"count": 42,
	}

	result, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	if !strings.Contains(result, "key:") {
		t.Error("Expected YAML output to contain 'key:'")
	}

	if !strings.Contains(result, "value") {
		t.Error("Expected YAML output to contain 'value'")
	}
}

func TestFormatTable(t *testing.T) {
	f := NewFormatter("table")

	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "simple map",
			data: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:    "empty slice",
			data:    []interface{}{},
			wantErr: false,
		},
		{
			name: "slice of maps",
			data: []interface{}{
				map[string]interface{}{"name": "item1", "value": 10},
				map[string]interface{}{"name": "item2", "value": 20},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Format(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}

func TestUnsupportedFormat(t *testing.T) {
	f := NewFormatter("invalid")

	data := map[string]interface{}{"key": "value"}
	_, err := f.Format(data)

	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("Expected error message to contain 'unsupported format', got: %v", err)
	}
}

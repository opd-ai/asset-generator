package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatTable OutputFormat = "table"
	FormatJSON  OutputFormat = "json"
	FormatYAML  OutputFormat = "yaml"
)

// Formatter handles output formatting
type Formatter struct {
	format OutputFormat
}

// NewFormatter creates a new formatter
func NewFormatter(format string) *Formatter {
	return &Formatter{
		format: OutputFormat(strings.ToLower(format)),
	}
}

// Format formats the data according to the configured format
func (f *Formatter) Format(data interface{}) (string, error) {
	switch f.format {
	case FormatJSON:
		return f.formatJSON(data)
	case FormatYAML:
		return f.formatYAML(data)
	case FormatTable:
		return f.formatTable(data)
	default:
		return "", fmt.Errorf("unsupported format: %s", f.format)
	}
}

// formatJSON formats data as JSON
func (f *Formatter) formatJSON(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(jsonData), nil
}

// formatYAML formats data as YAML
func (f *Formatter) formatYAML(data interface{}) (string, error) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return string(yamlData), nil
}

// formatTable formats data as a simple table
func (f *Formatter) formatTable(data interface{}) (string, error) {
	// Handle different data types
	switch v := data.(type) {
	case []interface{}:
		return f.formatSliceTable(v)
	case map[string]interface{}:
		return f.formatMapTable(v)
	default:
		// For specific types, use JSON encoding then parse
		return f.formatGenericTable(data)
	}
}

// formatSliceTable formats a slice as a table
func (f *Formatter) formatSliceTable(data []interface{}) (string, error) {
	if len(data) == 0 {
		return "No data available", nil
	}

	var buf strings.Builder

	// Get first element to determine structure
	first, ok := data[0].(map[string]interface{})
	if !ok {
		// Simple list - just format as text
		for _, item := range data {
			buf.WriteString(fmt.Sprintf("%v\n", item))
		}
		return buf.String(), nil
	}

	// Map list - extract headers from first element
	headers := make([]string, 0)
	for key := range first {
		headers = append(headers, strings.Title(key))
	}

	// Calculate column widths for proper alignment
	colWidths := make([]int, len(headers))
	
	// Initialize with header lengths
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	
	// Calculate max width for each column by scanning all data
	for _, item := range data {
		itemMap, _ := item.(map[string]interface{})
		for i, header := range headers {
			key := strings.ToLower(header)
			value := fmt.Sprintf("%v", itemMap[key])
			if len(value) > colWidths[i] {
				colWidths[i] = len(value)
			}
		}
	}

	// Write header with proper padding
	headerRow := make([]string, len(headers))
	for i, header := range headers {
		headerRow[i] = f.padRight(header, colWidths[i])
	}
	buf.WriteString(strings.Join(headerRow, " | "))
	buf.WriteString("\n")

	// Write separator with proper alignment
	separators := make([]string, len(headers))
	for i, width := range colWidths {
		separators[i] = strings.Repeat("-", width)
	}
	buf.WriteString(strings.Join(separators, "-+-"))
	buf.WriteString("\n")

	// Write rows with proper padding
	for _, item := range data {
		itemMap, _ := item.(map[string]interface{})
		row := make([]string, len(headers))
		for i, header := range headers {
			key := strings.ToLower(header)
			value := fmt.Sprintf("%v", itemMap[key])
			row[i] = f.padRight(value, colWidths[i])
		}
		buf.WriteString(strings.Join(row, " | "))
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// formatMapTable formats a map as a table
func (f *Formatter) formatMapTable(data map[string]interface{}) (string, error) {
	var buf strings.Builder

	if len(data) == 0 {
		return "No data available", nil
	}

	// Calculate column widths
	keyWidth := 3 // minimum width for "Key" header
	valueWidth := 5 // minimum width for "Value" header

	// Calculate max widths needed
	for key, value := range data {
		if len(key) > keyWidth {
			keyWidth = len(key)
		}
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) > valueWidth {
			valueWidth = len(valueStr)
		}
	}

	// Write header with proper padding
	buf.WriteString(f.padRight("Key", keyWidth))
	buf.WriteString(" | ")
	buf.WriteString(f.padRight("Value", valueWidth))
	buf.WriteString("\n")

	// Write separator
	buf.WriteString(strings.Repeat("-", keyWidth))
	buf.WriteString("-+-")
	buf.WriteString(strings.Repeat("-", valueWidth))
	buf.WriteString("\n")

	// Write rows with proper padding
	for key, value := range data {
		buf.WriteString(f.padRight(key, keyWidth))
		buf.WriteString(" | ")
		buf.WriteString(f.padRight(fmt.Sprintf("%v", value), valueWidth))
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// formatGenericTable formats generic data as a table
func (f *Formatter) formatGenericTable(data interface{}) (string, error) {
	// Try to marshal to JSON first, then parse as map
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonData, &mapData); err != nil {
		// If it's not a map, just show as single value
		return fmt.Sprintf("%v", data), nil
	}

	return f.formatMapTable(mapData)
}

// padRight pads a string to the specified width with spaces
func (f *Formatter) padRight(str string, width int) string {
	if len(str) >= width {
		return str
	}
	return str + strings.Repeat(" ", width-len(str))
}

// WriteToFile writes data to a file with timestamp
func WriteToFile(filename string, data string) error {
	// Add timestamp comment if JSON or YAML
	var output string
	timestamp := time.Now().Format(time.RFC3339)

	if strings.HasSuffix(filename, ".json") {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &jsonData); err == nil {
			jsonData["_generated_at"] = timestamp
			formatted, _ := json.MarshalIndent(jsonData, "", "  ")
			output = string(formatted)
		} else {
			output = data
		}
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		output = fmt.Sprintf("# Generated at: %s\n%s", timestamp, data)
	} else {
		output = data
	}

	if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

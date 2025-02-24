package format

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/goccy/go-json"

	"gopkg.in/yaml.v3"
)

// TypeAsString returns the string representing the type of the
// input parameter.
func TypeAsString(v any) string {
	return fmt.Sprintf("%T", v)
}

// ToJSON outputs the given object as a JSON string.
func ToJSON(v any) string {
	d, _ := json.Marshal(v)
	return string(d)
}

// ToPrettyJSON outputs the given object as a formatted JSON string.
func ToPrettyJSON(v any) string {
	d, _ := json.MarshalIndent(v, "", "  ")
	return string(d)
}

// ToYAML output the given object as a formatted YAML string.
func ToYAML(v any) string {
	d, _ := yaml.Marshal(v)
	return string(d)
}

func ToDateFormat(value string) string {
	// EEE, d MMM yyyy HH:mm:ss z => Mon, 2 Jan 2006 15:04:05 MST
	// YYYY-MM-DD => 2006-02-01
	value = strings.ReplaceAll(value, "MMM", "Jan")
	mappings := map[string]string{
		"yyyy": "2006",
		"MM":   "02",
		//"yyyy": "2006",
		// TODO: coplete mappings
	}
	for k, v := range mappings {
		value = strings.ReplaceAll(value, k, v)
	}

	return value
	//yyyy-MM-dd'T'HH:mm:ss
}

// WriteToFileAsJSON writes out to a file in a given directory; if dir is
// empty, the current directory is assumed; if the pattern contains one
// (or more) '*', it is assumed it is a temporary file and a random name
// is automatically assigned. The function returns the name of the file.
func WriteToFileAsJSON(dir string, pattern string, content string) (string, error) {
	var (
		file *os.File
		err  error
	)
	if dir == "" {
		dir = "."
	}
	if strings.Contains(pattern, `*`) {
		// assume the user wants a temporary file
		if file, err = os.CreateTemp(dir, pattern); err != nil {
			slog.Error("error opening temporary file", "path", file.Name(), "error", err)
			return file.Name(), err
		}
	} else {
		if file, err = os.Create(filepath.Clean(pattern)); err != nil {
			slog.Error("error opening output file", "path", file.Name(), "error", err)
			return file.Name(), err
		}
	}
	defer file.Close()
	slog.Debug("writing to output file", "path", file.Name())
	if _, err = file.Write([]byte(content)); err != nil {
		slog.Error("error writing JSON to temp file prior to massaging and parsing", "error", err)
	}
	return "", nil
}

func ColoredBool(value bool) string {
	if value {
		return color.GreenString("true")
	}
	return color.RedString("false")
}

package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var translations = make(map[string]map[string]interface{})

// LoadLocales reads JSON files from the given directory
func LoadLocales(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		lang := strings.TrimSuffix(f.Name(), ".json")
		filePath := fmt.Sprintf("%s/%s", path, f.Name())

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		byteValue, _ := ioutil.ReadAll(file)
		var result map[string]interface{}
		json.Unmarshal(byteValue, &result)

		translations[lang] = result
	}
	return nil
}

// Translate retrieves a message by key and language
func Translate(lang, key string, args ...map[string]interface{}) string {
	if translations[lang] == nil {
		lang = "en" // Fallback to English
	}

	keys := strings.Split(key, ".")
	var val interface{} = translations[lang]

	for _, k := range keys {
		if m, ok := val.(map[string]interface{}); ok {
			val = m[k]
		} else {
			val = nil
			break
		}
	}

	if val == nil {
		// Fallback to English if not found in requested language
		if lang != "en" {
			return Translate("en", key, args...)
		}
		return key // Return key itself as fallback
	}

	res, ok := val.(string)
	if !ok {
		return key
	}

	// Simple placeholder replacement if args are provided
	if len(args) > 0 {
		for k, v := range args[0] {
			placeholder := fmt.Sprintf("{{%s}}", k)
			res = strings.ReplaceAll(res, placeholder, fmt.Sprintf("%v", v))
		}
	}

	return res
}

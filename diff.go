package main

import (
	"fmt"
	"reflect"
	"sort"
)

func computeDiff(sessionEncodedName string) ([]Setting, error) {
	return computeSettings(sessionEncodedName, false)
}

func computeAllSettings(sessionEncodedName string) ([]Setting, error) {
	return computeSettings(sessionEncodedName, true)
}

func computeSettings(sessionEncodedName string, showAll bool) ([]Setting, error) {
	defaultSettings, defaultTypes, err := getSessionSettings(defaultSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to read Default Settings: %w", err)
	}

	sessionSettings, _, err := getSessionSettings(sessionEncodedName)
	if err != nil {
		return nil, fmt.Errorf("failed to read session settings: %w", err)
	}

	return buildSettings(defaultSettings, defaultTypes, sessionSettings, showAll), nil
}

func buildSettings(defaultSettings map[string]interface{}, defaultTypes map[string]RegistryType, sessionSettings map[string]interface{}, showAll bool) []Setting {
	settings := make([]Setting, 0, len(defaultSettings))

	for settingName, defaultValue := range defaultSettings {
		sessionValue, exists := sessionSettings[settingName]
		isDifferent := true

		if !exists {
			sessionValue = defaultValue
			isDifferent = false
		} else {
			isDifferent = !reflect.DeepEqual(defaultValue, sessionValue)
		}

		if showAll || isDifferent {
			settings = append(settings, Setting{
				Name:         settingName,
				DefaultValue: defaultValue,
				CurrentValue: sessionValue,
				Type:         defaultTypes[settingName],
				IsChecked:    false,
				IsDifferent:  isDifferent,
			})
		}
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Name < settings[j].Name
	})

	return settings
}

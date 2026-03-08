package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	puttySessionsPath = `Software\SimonTatham\PuTTY\Sessions`
	defaultSettings   = "Default%20Settings"
)

var sensitiveKeys = map[string]bool{
	"UserName":          true,
	"PublicKeyFile":     true,
	"ProxyUsername":     true,
	"ProxyPassword":     true,
	"LocalProxyCommand": true,
	"HostName":          true,
}

func getSessions() ([]Session, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, puttySessionsPath, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, fmt.Errorf("failed to open PuTTY sessions registry key: %w", err)
	}
	defer k.Close()

	subKeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read session names: %w", err)
	}

	var sessions []Session
	for _, encodedName := range subKeys {
		if encodedName == defaultSettings {
			continue
		}

		displayName := decodePuttySessionName(encodedName)

		sessions = append(sessions, Session{
			EncodedName: encodedName,
			DisplayName: displayName,
		})
	}

	return sessions, nil
}

func decodePuttySessionName(encoded string) string {
	var result strings.Builder
	i := 0

	for i < len(encoded) {
		if encoded[i] == '%' && i+2 < len(encoded) {
			hexStr := encoded[i+1 : i+3]
			if byteVal, err := strconv.ParseUint(hexStr, 16, 8); err == nil {
				result.WriteByte(byte(byteVal))
				i += 3
				continue
			}
		}
		result.WriteByte(encoded[i])
		i++
	}

	ansiBytes := []byte(result.String())
	if len(ansiBytes) == 0 {
		return ""
	}

	n, err := windows.MultiByteToWideChar(windows.GetACP(), 0, &ansiBytes[0], int32(len(ansiBytes)), nil, 0)
	if err != nil || n == 0 {
		return result.String()
	}

	utf16 := make([]uint16, n)
	_, err = windows.MultiByteToWideChar(windows.GetACP(), 0, &ansiBytes[0], int32(len(ansiBytes)), &utf16[0], n)
	if err != nil {
		return result.String()
	}

	return windows.UTF16ToString(utf16)
}

func getSessionSettings(encodedName string) (map[string]interface{}, map[string]RegistryType, error) {
	keyPath := puttySessionsPath + `\` + encodedName
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open session key: %w", err)
	}
	defer k.Close()

	valueNames, err := k.ReadValueNames(-1)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read value names: %w", err)
	}

	settings := make(map[string]interface{})
	types := make(map[string]RegistryType)

	for _, valueName := range valueNames {
		if sensitiveKeys[valueName] {
			continue
		}

		_, valType, err := k.GetValue(valueName, nil)
		if err != nil {
			continue
		}

		switch valType {
		case registry.SZ:
			val, _, err := k.GetStringValue(valueName)
			if err == nil {
				settings[valueName] = val
				types[valueName] = RegString
			}
		case registry.DWORD:
			val, _, err := k.GetIntegerValue(valueName)
			if err == nil {
				settings[valueName] = uint32(val)
				types[valueName] = RegDWord
			}
		}
	}

	return settings, types, nil
}

func writeSettingToSession(encodedSessionName, settingName string, value interface{}, regType RegistryType) error {
	keyPath := puttySessionsPath + `\` + encodedSessionName
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open session key for writing: %w", err)
	}
	defer k.Close()

	switch regType {
	case RegString:
		if strVal, ok := value.(string); ok {
			return k.SetStringValue(settingName, strVal)
		}
		return fmt.Errorf("type mismatch: expected string")
	case RegDWord:
		if dwordVal, ok := value.(uint32); ok {
			return k.SetDWordValue(settingName, dwordVal)
		}
		return fmt.Errorf("type mismatch: expected uint32")
	default:
		return fmt.Errorf("unknown registry type")
	}
}

package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

type RegistryType int

const (
	RegString RegistryType = iota
	RegDWord
)

type Setting struct {
	Name         string
	DefaultValue interface{}
	CurrentValue interface{}
	Type         RegistryType
	IsChecked    bool
	IsDifferent  bool
}

type Session struct {
	EncodedName string
	DisplayName string
}

func (s *Setting) GetDefaultValueString() string {
	switch v := s.DefaultValue.(type) {
	case string:
		return sanitizeString(v)
	case uint32:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

func (s *Setting) GetCurrentValueString() string {
	switch v := s.CurrentValue.(type) {
	case string:
		return sanitizeString(v)
	case uint32:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

func sanitizeString(s string) string {
	return strings.Map(func(r rune) rune {
		if r == 0 {
			return -1
		}
		return r
	}, s)
}

func regTypeToRegistryType(regType uint32) RegistryType {
	if regType == registry.DWORD {
		return RegDWord
	}
	return RegString
}

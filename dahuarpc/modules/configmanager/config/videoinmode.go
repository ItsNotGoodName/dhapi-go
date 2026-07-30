package config

import (
	"context"
	"fmt"
	"slices"

	"github.com/ItsNotGoodName/dhapi-go/dahuarpc"
	"github.com/ItsNotGoodName/dhapi-go/dahuarpc/modules/configmanager"
)

const (
	ProfileDay             = "Day"
	ProfileNight           = "Night"
	ProfileNormal          = "Normal"      // General
	ProfileFrontLight      = "FrontLight"  // Front Light
	ProfileBackLight       = "BackLight"   // Backlight
	ProfileStrongBacklight = "StrongLight" // Strong Backlight
	ProfileLowLight        = "LowLight"    // Low Illuminance
	ProfileCustom          = "Custom"      // Custom1
	ProfileCustom1         = "Custom1"     // Custom2
)

func GetVideoInMode(ctx context.Context, c dahuarpc.Conn) (configmanager.ConfigArray[VideoInMode], error) {
	return configmanager.GetConfigArray[VideoInMode](ctx, c, "VideoInMode")
}

type VideoInMode struct {
	Config        []int                     `json:"Config"`
	ConfigEx      string                    `json:"ConfigEx"`
	Mode          int                       `json:"Mode"`
	TimeSection   [][]dahuarpc.TimeSection  `json:"TimeSection"`
	TimeSectionEX []dahuarpc.TimeSection2   `json:"TimeSectionEX"`
	TimeSectionV2 [][]dahuarpc.TimeSection2 `json:"TimeSectionV2"`
}

func (c VideoInMode) Merge(js string) (string, error) {
	return configmanager.Merge(js, []configmanager.MergeValues{
		{Path: "Config", Value: c.Config},
		{Path: "ConfigEx", Value: c.ConfigEx},
		{Path: "Mode", Value: c.Mode},
		{Path: "TimeSection", Value: c.TimeSection},
		{Path: "TimeSectionEX", Value: c.TimeSectionEX},
		{Path: "TimeSectionV2", Value: c.TimeSectionV2},
	})
}

func (c VideoInMode) Validate() error {
	if len(c.TimeSection) == 0 || len(c.TimeSection[0]) == 0 {
		return fmt.Errorf("empty TimeSection")
	}

	return nil
}

type SwitchMode int

const (
	SwitchModeGeneral SwitchMode = iota
	SwitchModeDay
	SwitchModeNight
	SwitchModeSchedule
	SwitchModeBrightness
	SwitchModeProfiles
	SwitchModeUnknown
)

func (m SwitchMode) String() string {
	switch m {
	case SwitchModeGeneral:
		return "general"
	case SwitchModeDay:
		return "day"
	case SwitchModeNight:
		return "night"
	case SwitchModeSchedule:
		return "schedule"
	case SwitchModeBrightness:
		return "brightness"
	case SwitchModeProfiles:
		return "profiles"
	default:
		return "unknown"
	}
}

func (m VideoInMode) SwitchMode() SwitchMode {
	if m.Mode == 0 && slices.Equal(m.Config, []int{2}) {
		return SwitchModeGeneral
	}
	if m.Mode == 0 && slices.Equal(m.Config, []int{0}) {
		return SwitchModeDay
	}
	if m.Mode == 0 && slices.Equal(m.Config, []int{1}) {
		return SwitchModeNight
	}
	if m.Mode == 1 && slices.Equal(m.Config, []int{0, 1}) {
		return SwitchModeSchedule
	}
	if m.Mode == 2 && slices.Equal(m.Config, []int{0, 1}) {
		return SwitchModeBrightness
	}
	if m.Mode == 3 && slices.Equal(m.Config, []int{}) {
		return SwitchModeProfiles
	}
	return SwitchModeUnknown
}

func (m *VideoInMode) SetSwitchMode(mode SwitchMode) {
	switch mode {
	case SwitchModeGeneral:
		m.Mode = 0
		m.Config = []int{2}
	case SwitchModeDay:
		m.Mode = 0
		m.Config = []int{0}
	case SwitchModeNight:
		m.Mode = 0
		m.Config = []int{1}
	case SwitchModeSchedule:
		m.Mode = 1
		m.Config = []int{0, 1}
	case SwitchModeBrightness:
		m.Mode = 2
		m.Config = []int{0, 1}
	case SwitchModeProfiles:
		m.Mode = 3
		m.Config = []int{}
	}
}

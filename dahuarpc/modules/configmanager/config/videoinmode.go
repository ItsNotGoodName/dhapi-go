package config

import (
	"context"
	"fmt"
	"slices"

	"github.com/ItsNotGoodName/dhapi-go/dahuarpc"
	"github.com/ItsNotGoodName/dhapi-go/dahuarpc/modules/configmanager"
)

func GetVideoInMode(ctx context.Context, c dahuarpc.Conn) (configmanager.ConfigArray[VideoInMode], error) {
	return configmanager.GetConfigArray[VideoInMode](ctx, c, "VideoInMode")
}

type VideoInMode struct {
	Config      []int                    `json:"Config"`
	Mode        int                      `json:"Mode"`
	TimeSection [][]dahuarpc.TimeSection `json:"TimeSection"`
}

func (c VideoInMode) Merge(js string) (string, error) {
	return configmanager.Merge(js, []configmanager.MergeValues{
		{Path: "Config", Value: c.Config},
		{Path: "Mode", Value: c.Mode},
		{Path: "TimeSection", Value: c.TimeSection},
	})
}

func (c VideoInMode) Validate() error {
	if len(c.TimeSection) == 0 || len(c.TimeSection[0]) == 0 {
		return fmt.Errorf("empty TimeSection")
	}

	_, err := c.switchMode()
	if err != nil {
		return err
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
	default:
		return "unknown"
	}
}

func (m VideoInMode) SwitchMode() SwitchMode {
	s, _ := m.switchMode()
	return s
}

func (m VideoInMode) switchMode() (SwitchMode, error) {
	if m.Mode == 0 && slices.Equal(m.Config, []int{2}) {
		return SwitchModeGeneral, nil
	}
	if m.Mode == 0 && slices.Equal(m.Config, []int{0}) {
		return SwitchModeDay, nil
	}
	if m.Mode == 0 && slices.Equal(m.Config, []int{1}) {
		return SwitchModeNight, nil
	}
	if m.Mode == 1 && slices.Equal(m.Config, []int{0, 1}) {
		return SwitchModeSchedule, nil
	}
	if m.Mode == 2 && slices.Equal(m.Config, []int{0, 1}) {
		return SwitchModeBrightness, nil
	}
	return 0, fmt.Errorf("unknown SwitchMode: mode=%d config=%v", m.Mode, m.Config)
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
	}
}

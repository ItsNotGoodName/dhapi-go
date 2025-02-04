package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/ItsNotGoodName/dhapi-go/dahuarpc"
	"github.com/ItsNotGoodName/dhapi-go/dahuarpc/modules/configmanager"
)

func GetRecord(ctx context.Context, c dahuarpc.Conn) (configmanager.ConfigArray[Record], error) {
	return configmanager.GetConfigArray[Record](ctx, c, "Record")
}

type Record struct {
	Format        string                   `json:"Format"`
	HolidayEnable bool                     `json:"HolidayEnable"`
	PreRecord     int                      `json:"PreRecord"`
	Redundancy    bool                     `json:"Redundancy"`
	SnapShot      bool                     `json:"SnapShot"`
	Stream        int                      `json:"Stream"`
	TimeSection   [][]dahuarpc.TimeSection `json:"TimeSection"`
}

func (c Record) Merge(js string) (string, error) {
	return "", fmt.Errorf("%w: Merge not implemented for 'Record'", errors.ErrUnsupported)
}

func (c Record) Validate() error {
	return nil
}

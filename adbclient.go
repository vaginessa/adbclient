package adbclient

import (
    "github.com/alexjch/adbclient/comm"
    "github.com/alexjch/adbclient/comm"
)

const (
    LIST_DEVICES = "host:devices"
)

type ADBClient struct {
    comm interface{}
}

type Device struct {
    serialNumber string
    state string
}

func (adb *ADBClient) Devices() []Device{
    devices := []Device{}
    adb.comm.Send()
}

func NewADBClient() *ADBClient{
    client := ADBClient{
        comm: comm.NewConn(),
    }

    return &client
}



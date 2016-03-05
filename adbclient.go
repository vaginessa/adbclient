package adbclient

import (
    "strings"
    "github.com/alexjch/adbclient/conn"
)

type ADBClient struct {
    conn_ *conn.ADBconn
}

type Device struct {
    // Type that encapsulates a device
    serialNumber string
    state string
}

func parseDevices(result string) ([]Device, error){
    devices := []Device{}
    lines := strings.Split(result, "\n")
    for l := range lines{
        values := strings.Split(lines[l], "\t")
        if len(values) >= 2{
            device := Device{
                serialNumber: string(values[0]),
                state: string(values[1]),
            }
            devices = append(devices, device)
        }
    }
    return devices, nil
}

func (adb *ADBClient) Devices() ([]Device, error){
    // Returns an array with devices connected to host
    result, err := adb.conn_.Send(LIST_DEVICES)
    if err != nil{
        return nil, err
    }
    return parseDevices(result)
}

func (adb *ADBClient) Version() (string, error){
    // Returns the version of the ADB server
    result, err := adb.conn_.Send(VERSION)
    if err != nil{
        return "", err
    }
    return result, nil
}

func NewADBClient() *ADBClient{
    client := ADBClient{
        conn_: &conn.ADBconn{},
    }

    return &client
}



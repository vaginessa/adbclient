package adbclient

import (
    "errors"
    "strconv"
    "strings"
    "github.com/alexjch/adbclient/conn"
)

const (
    CHECKSUM = "OKAY0000"
    SHELL = "shell:<cmd>"
)

type ADBClient struct {
    conn_ *conn.ADBconn
}

type Device struct {
    // Type that encapsulates a device
    serialNumber string
    state string
}

func (adb *ADBClient) stripChecksum(resp string) (string, error){
    if len(resp) < len(CHECKSUM){
        return "", errors.New("Invalid checksum")
    }
    checksum := resp[0:8]
    value := resp[8:]

    length, err :=  strconv.ParseInt(checksum[4:], 16, 16)
    if err != nil {
        return "", errors.New("Invalid checksum, unable to parse message length")
    }

    if int(length) != len(value){
        return "", errors.New("Invalid checksum, length on header does not match length of message")
    }

    return value, nil
}

func (adb *ADBClient) checkOKEY(resp string) (bool){
    return strings.Contains(resp, "OKEY")
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

func (adb *ADBClient) Shell(serial, query string) (string, error) {
    // Sends a command to shell
    result, err := adb.conn_.SendToHost(serial, strings.Replace(SHELL, "<cmd>", query, 1))
    if err != nil{
        return "", err
    }
    return result, nil
}

func (adb *ADBClient) Devices() ([]Device, error){
    // Returns an array with devices connected to host
    result, err := adb.conn_.Send(LIST_DEVICES)
    if err != nil{
        return nil, err
    }
    result, _ = adb.stripChecksum(result)
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

func New() *ADBClient{
    // Returns a new instance of ADBClient
    client := ADBClient{
        conn_: &conn.ADBconn{},
    }
    return &client
}



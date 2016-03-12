package adbclient

import (
    "fmt"
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

func (adb *ADBClient) Sync(cmd, serial, filename string) (string, error){
    // Executes a Sync cmd in device
    switch cmd {
    case "LIST", "STAT", "RECV":
        return adb.conn_.Sync(cmd, serial, filename)
    default:
        return "", errors.New(fmt.Sprintf("Command: %s is unknown", cmd))
    }
}

func (adb *ADBClient) Pull(serial, filename string) (string, error){
    // Pulls a file from device
    return adb.conn_.Sync("RECV", serial, filename)
}

func (adb *ADBClient) Push(serial, source, destination string) (string, error){
    // Pulls a file from device
    return adb.conn_.Push(serial, source, destination)
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
    v, _ := strconv.ParseInt(result[8:], 16, 32)
    return fmt.Sprintf("%d", v), nil
}

func (adb *ADBClient) Track() <-chan []Device{
    devices := adb.conn_.Track()
    update := make(chan []Device)

    go func(){
        for{
            devcs, err := parseDevices(<-devices)
            if err != nil {
                break
            }
            update <- devcs
        }
    }()

    return update
}

func New() *ADBClient{
    // Returns a new instance of ADBClient
    client := ADBClient{
        conn_: &conn.ADBconn{},
    }
    return &client
}



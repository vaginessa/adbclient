package adbclient

import (
    "testing"
    "reflect"
    "os"
)

func TestDevices(t *testing.T){
    devices, err := New().Devices()
    if err != nil{
        t.Error("Call to devices caused an error: ", err)
    }
    var RetType []Device
    if reflect.TypeOf(devices) != reflect.TypeOf(RetType){
        t.Error("Type should be []Device{}")
    }
}

func TestVersion(t *testing.T){
    _, err := New().Version()
    if err != nil{
        t.Error("Unexpected error")
    }
}

func TestShell(t *testing.T){
    serialN := os.Getenv("DEV_SERIAL")
    _, err := New().Shell(serialN, "ls -all")
    if err != nil{
        t.Error("Unexpected error")
    }
}

func TestPull(t *testing.T){
    FILE_NAME := "default.prop"
    if _, err := os.Stat(FILE_NAME); !os.IsNotExist(err) {
        os.Remove(FILE_NAME)
    }
    serialN := os.Getenv("DEV_SERIAL")
    _, err := New().Pull(serialN, "/" + FILE_NAME)
    if err != nil{
        t.Error("Unexpected error")
    }
    if _, err := os.Stat(FILE_NAME); os.IsNotExist(err) {
        t.Error("File was not found")
    }
}

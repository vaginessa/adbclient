package adbclient

import (
    "testing"
    "reflect"
)

func TestDevices(t *testing.T){
    devices, err := NewADBClient().Devices()
    if err != nil{
        t.Error("Call to devices caused an error: ", err)
    }
    var RetType []Device
    if reflect.TypeOf(devices) != reflect.TypeOf(RetType){
        t.Error("Type should be []Device{}")
    }
}

func TestVersion(t *testing.T){
    _, err := NewADBClient().Version()
    if err != nil{
        t.Error("Unexpected error")
    }
}
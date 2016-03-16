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
    os.Remove(FILE_NAME)
}

/*TODO: func TestPush(t *testing.T){

}*/

func TestPackages(t *testing.T){
    serialN := os.Getenv("DEV_SERIAL")
    flags := []string(nil)
    packages, err := New().ListPackages(serialN, flags)
    if err != nil && len(packages) > 10{
        t.Error("Unexpected error")
    }
    t.Log(packages)
}

func TestPackagesFlag(t *testing.T){
    serialN := os.Getenv("DEV_SERIAL")
    flags := []string{SYSTEM_PACKAGES, ASSOCIATED_FILE}
    packages, err := New().ListPackages(serialN, flags)
    if err != nil && len(packages) > 10{
        t.Error("Unexpected error")
    }
    t.Log(packages)
}

func TestGetFeatures(t *testing.T){
    serialN := os.Getenv("DEV_SERIAL")
    features, err := New().GetFeatures(serialN)
    if err != nil && len(features) > 10{
        t.Error("Unexpected error")
    }
    t.Log(features)
}


func TestScreencapture(t *testing.T){
    serialN := os.Getenv("DEV_SERIAL")
    captureFile, err := New().Screencapture(serialN)
    if err != nil {
        t.Error(err)
    }
    _, err = os.Stat(captureFile)
    if os.IsNotExist(err) {
        t.Error("File not found: ", captureFile)
    }
}


/* TODO: Negative/failing tests */
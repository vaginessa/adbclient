package main

import (
    "os"
    "fmt"
    "github.com/alexjch/adbclient"
)

var serialN = os.Getenv("DEV_SERIAL")

func version(){
    version, err := adbclient.New().Version()
    if err != nil{
        fmt.Println("Unable to obtain version", err)
    }
    fmt.Println(version)
}

func devices(){
    version, err := adbclient.New().Devices()
    if err != nil{
        fmt.Println("Unable to list devices", err)
    }
    fmt.Println(version)
}

func track(){
    devices, _ := adbclient.New().Track()
    for{
        fmt.Println(<-devices)
    }
}

func syncList(serial, filePath string){
    result, err := adbclient.New().Sync("LIST", serial, filePath)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func syncStat(serial, filePath string){
    result, err := adbclient.New().Sync("STAT", serial, filePath)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func pull(serial, filePath string){
    result, err := adbclient.New().Pull(serial, filePath)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func push(serial, origin, destination string){
    result, err := adbclient.New().Push(serial, origin, destination)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func getprop(serial string){
    props, err := adbclient.New().GetProp(serial)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(props)
}

func lolcat(serial string){
    pipe, err := adbclient.New().Logcat(serial)
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }
    for{
        fmt.Println(<-pipe)
    }
}

func main(){
    track()
}

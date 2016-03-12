package main

import (
    "os"
    "fmt"
    "bufio"
    "github.com/alexjch/adbclient"
    "github.com/alexjch/adbclient/conn"
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
    devices := adbclient.New().Track()

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

func main(){
    syncList(serialN, "/mnt")
    syncStat(serialN, "/default.prop")
    pull(serialN, "/default.prop")

    stdio := bufio.NewScanner(os.Stdin)
    adbc := &conn.ADBconn{}

    for stdio.Scan() {
        cmd := stdio.Text()
        ret, _ := adbc.Send(cmd)
        fmt.Println(ret)
    }
}

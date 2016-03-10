package main

import (
    "fmt"
    "os"
    "bufio"
    "github.com/alexjch/adbclient"
    "github.com/alexjch/adbclient/conn"
)


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

func syncList(path string){
    result, err := adbclient.New().Sync("LIST", "075923ba00cc8e9c", path)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func syncStat(path string){
    result, err := adbclient.New().Sync("STAT", "075923ba00cc8e9c", path)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func syncRecv(path string){
    result, err := adbclient.New().Sync("RECV", "075923ba00cc8e9c", path)
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(result)
}

func main(){
/*    syncList("/mnt")
    syncStat("/default.prop") */
    syncRecv("/default.prop")

    stdio := bufio.NewScanner(os.Stdin)
    adbc := &conn.ADBconn{}

    for stdio.Scan() {
        cmd := stdio.Text()
        ret, _ := adbc.Send(cmd)
        fmt.Println(ret)
    }
}

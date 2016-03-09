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

func sync(){
    result, err := adbclient.New().Pull("075923ba00cc8e9c", "/mnt")
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println("Succeded ", result)
}

func main(){

    sync()

    stdio := bufio.NewScanner(os.Stdin)
    adbc := &conn.ADBconn{}

    for stdio.Scan() {
        cmd := stdio.Text()
        ret, _ := adbc.Send(cmd)
        fmt.Println(ret)
    }
}

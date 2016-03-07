package main

import (
    "fmt"
    "os"
    "bufio"
    "github.com/alexjch/adbclient/conn"
)

func main(){
    stdio := bufio.NewScanner(os.Stdin)
    adbc := &conn.ADBconn{}

    for stdio.Scan() {
        cmd := stdio.Text()
        ret, _ := adbc.Send(cmd)
        fmt.Println(ret)
    }
}

package main

import (
    "os"
    "fmt"
    "bufio"
    "github.com/alexjch/adbclient/comm"
)

func main(){
    stdio := bufio.NewScanner(os.Stdin)
    adbc := comm.NewConn()

    for stdio.Scan() {
        cmd := stdio.Text()
        ret, _ := adbc.Send(cmd)
        fmt.Println(ret)
    }
}

package comm

import (
    "net"
    "fmt"
    "log"
    "bufio"
    "io/ioutil"
)

const (
    PORT = 5037
)

type adbclient struct{
    conn net.Conn
}

func (a *adbclient) FormatCommand(cmd string) string{
    return fmt.Sprintf("%04x%s", len(cmd), cmd)
}

func (a *adbclient) Connect () error{
    conn, err := net.Dial("tcp", fmt.Sprintf(":%d", PORT))
    a.conn = conn
    return err
}

func (a *adbclient) Send (cmd string) (string, error){
    cmdFrmt := a.FormatCommand(cmd)
    if err := a.Connect(); err != nil {
        log.Fatalln("Error connecting: ", err)
        return "", err
    }
    _, err := a.conn.Write([]byte(cmdFrmt))
    if err != nil{
        log.Fatalln("Error conn with: ", err)
    }
    reader := bufio.NewReader(a.conn)
    data, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Fatalln("Error receiving: ", err)
    }
    a.conn.Close()
    return string(data), err
}

func NewConn() *adbclient{
    adbc := adbclient{
        conn: nil,
    }

    return &adbc
}

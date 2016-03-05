package conn

import (
    "net"
    "fmt"
    "log"
    "bufio"
    "errors"
    "io/ioutil"
    "strconv"
)

const (
    PORT = 5037
    CHECKSUM = "OKAY0000"
)

type ADBconn struct{}

func (a *ADBconn) formatCommand(cmd string) string{
    return fmt.Sprintf("%04x%s", len(cmd), cmd)
}

func (a *ADBconn) stripChecksum(resp string) (string, error){
    if len(resp) < len(CHECKSUM){
        return "", errors.New("Invalid checksum")
    }
    checksum := resp[0:8]
    value := resp[8:]

    length, err :=  strconv.ParseInt(checksum[4:], 16, 16)
    if err != nil {
        return "", errors.New("Invalid checksum, unable to parse message length")
    }

    if int(length) != len(value){
        return "", errors.New("Invalid checksum, length on header does not match length of message")
    }

    return value, nil
}

func (a *ADBconn) Connect () (net.Conn, error){
    conn, err := net.Dial("tcp", fmt.Sprintf(":%d", PORT))
    return conn, err
}

func (a *ADBconn) Send (cmd string) (string, error){
    cmdFrmt := a.formatCommand(cmd)
    conn, err := a.Connect()
    if err != nil {
        log.Fatalln("Error connecting: ", err)
        return "", err
    }
    _, err = conn.Write([]byte(cmdFrmt))
    if err != nil{
        log.Fatalln("Error conn with: ", err)
        return "", err
    }
    reader := bufio.NewReader(conn)
    data, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Fatalln("Error receiving: ", err)
        return "", err
    }
    conn.Close()
    return a.stripChecksum(string(data))
}

package conn

import (
    "net"
    "fmt"
    "log"
    "errors"
    "strings"
    "io"
)

const (
    PORT = 5037
    HOST_TRANSPORT = "host:transport:<id>"
)

type ADBconn struct{}


func (a *ADBconn) send (conn net.Conn, cmd string) error{
    _, err := fmt.Fprintf(conn, "%04x%s", len(cmd), cmd)
    if err != nil{
        log.Fatalln("Error conn with: ", err)
        return err
    }
    return nil
}

func (a *ADBconn) receive (conn net.Conn) (int, string, error){
    buff := make([]byte, 256)
    count, err := conn.Read(buff)
    if err != nil {
        return 0, "", err
    }
    return count, string(buff[0:count]), nil
}

func (a *ADBconn) Connect () (net.Conn, error){
    // Open a connection to ADB server
    conn, err := net.Dial("tcp", fmt.Sprintf(":%d", PORT))
    return conn, err
}

func (a *ADBconn) Send (cmd string) (string, error){
    // Send command to host
    conn, err := a.Connect()
    if err != nil {
        log.Fatalln("Error connecting: ", err)
        return "", err
    }
    defer conn.Close()
    err = a.send(conn, cmd)
    if err != nil {
        log.Fatalln("Error sending command")
        return "", err
    }
    _, resp, err := a.receive(conn)
    if err != nil {
        return "", err
    }
    return string(resp), nil
}

func (a *ADBconn) SendToHost (serial string, cmd string) (string, error){
    // Send command to host identify by serial
    conn, err := a.Connect()
    out := []string{}
    if err != nil {
        log.Fatalln("Error connecting: ", err)
        return "", err
    }
    defer conn.Close()
    host := strings.Replace(HOST_TRANSPORT, "<id>", serial, 1)
    err = a.send(conn, host)
    if err != nil {
        log.Fatalln("Error sending transport")
        return "", err
    }
    _, resp, err := a.receive(conn)
    if strings.Contains(resp, "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }
    err = a.send(conn, cmd)
    if err != nil {
        log.Fatalln("Error sending command")
        return "", err
    }
    _, resp, err = a.receive(conn)
    if strings.Contains(resp, "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }

    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err
        }
        out = append(out, resp)
    }

    result := strings.Join(out, "")
    return result, nil
}
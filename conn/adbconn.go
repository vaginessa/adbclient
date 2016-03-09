package conn

import (
    "net"
    "fmt"
    "log"
    "errors"
    "strings"
    "io"
    "bytes"
    "encoding/binary"
)

const (
    PORT = 5037
    HOST_TRANSPORT = "host:transport:<id>"
    TRACK_CMD = "host:track-devices"
    SYNC = "sync:"
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

func (a *ADBconn) Sync (serial string, path string)  (string, error) {
    conn, err := a.Connect()
    if err != nil {
        log.Println("Error connecting: ", err)
        return "", err
    }
    host := strings.Replace(HOST_TRANSPORT, "<id>", serial, 1)
    if err = a.send(conn, host); err != nil {
        log.Println("Error sending transport")
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(resp, "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }
    if err = a.send(conn, SYNC); err != nil {
        log.Println("Error sending command")
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(resp, "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }

    length := uint32(len(path))
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.BigEndian, []byte("LIST"))
    binary.Write(buf, binary.LittleEndian, length)
    binary.Write(buf, binary.BigEndian, []byte(path))

    count, err := conn.Write(buf.Bytes())

    if err != nil {
        log.Println("Failed sending sync command ", err)
        return "", err
    }

    log.Println("Received from sending command ", count)

    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err
        }
        mode, _ := binary.Uvarint([]byte(resp[4:8]))
        size, _ := binary.Uvarint([]byte(resp[8:12]))
        log.Println("resp: ", resp[0:4], mode, size)
    }

    return "Hello", nil
}


func (a *ADBconn) Track () <-chan string{
    //
    conn, err := a.Connect()
    if err != nil {
        log.Println("Error connecting: ", err)
        return nil
    }
    if err = a.send(conn, TRACK_CMD); err != nil {
        log.Println("Error sending command")
        return nil
    }

    out := make(chan string)

    go func(){
        for{
            _, resp, err := a.receive(conn)
            if err != nil {
                log.Println("Error receiving data")
                conn.Close()
                break
            }
            out <- resp
        }
    }()

    return out
}


func (a *ADBconn) Send (cmd string) (string, error){
    // Send command to host
    conn, err := a.Connect()
    if err != nil {
        log.Println("Error connecting: ", err)
        return "", err
    }
    defer conn.Close()
    if err = a.send(conn, cmd); err != nil {
        log.Println("Error sending command")
        return "", err
    }
    _, resp, err := a.receive(conn)
    if err != nil {
        return "", err
    }
    return resp, nil
}

func (a *ADBconn) SendToHost (serial string, cmd string) (string, error){
    // Send command to host identify by serial
    conn, err := a.Connect()
    if err != nil {
        log.Println("Error connecting: ", err)
        return "", err
    }
    defer conn.Close()
    out := []string{}
    host := strings.Replace(HOST_TRANSPORT, "<id>", serial, 1)
    if err = a.send(conn, host); err != nil {
        log.Println("Error sending transport")
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(resp, "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }
    if err = a.send(conn, cmd); err != nil {
        log.Println("Error sending command")
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(resp, "OKAY") != true {
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

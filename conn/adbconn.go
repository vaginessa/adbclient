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

func (a *ADBconn) send(conn net.Conn, cmd string) error {
    _, err := fmt.Fprintf(conn, "%04x%s", len(cmd), cmd)
    if err != nil {
        log.Fatalln("Error conn with: ", err)
        return err
    }
    return nil
}

func (a *ADBconn) receive(conn net.Conn) (int, []byte, error) {
    buff := make([]byte, 256)
    count, err := conn.Read(buff)
    if err != nil {
        return 0, nil, err
    }
    return count, buff[0:count], nil
}

func (a *ADBconn) Connect() (net.Conn, error) {
    // Open a connection to ADB server
    conn, err := net.Dial("tcp", fmt.Sprintf(":%d", PORT))
    return conn, err
}

func (a *ADBconn) readStat(conn net.Conn) (string, error){
    out := new(bytes.Buffer)
    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err
        }
        if strings.Contains(string(resp), "STAT") {
            out.Write(resp)
            break
        }
    }
    dataOut := new(bytes.Buffer)
    for ; out.Len() > 0; {
        var mode, size, stat, nLen uint32
        out.Next(4)
        data := bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &mode)
        data = bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &size)
        data = bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &stat)
        data = bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &nLen)
        fname := out.Next(int(nLen))
        line := fmt.Sprintf("%s mode:%d size:%d stat:%d name lenght:%d\n", fname, mode, size, stat, nLen)
        dataOut.Write([]byte(line))
    }
    return dataOut.String(), nil
}

func (a *ADBconn) readList(conn net.Conn) (string, error){
    // TODO: brake in two real and parse
    out := new(bytes.Buffer)
    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err
        }
        if strings.Contains(string(resp), "DONE") {
            indx := strings.Index(string(resp), "DONE")
            out.Write(resp[0:indx])
            break
        } else {
            out.Write(resp)
        }
    }
    dataOut := new(bytes.Buffer)
    for ; out.Len() > 0; {
        var mode, size, stat, nLen uint32
        out.Next(4)
        data := bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &mode)
        data = bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &size)
        data = bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &stat)
        data = bytes.NewReader(out.Next(4))
        binary.Read(data, binary.LittleEndian, &nLen)
        fname := out.Next(int(nLen))
        line := fmt.Sprintf("%s mode:%d size:%d stat:%d name lenght:%d\n", fname, mode, size, stat, nLen)
        dataOut.Write([]byte(line))
    }
    return dataOut.String(), nil
}

func (a *ADBconn) readRecv (conn net.Conn) (string, error){

    out := new(bytes.Buffer)
    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err
        }
        if strings.Contains(string(resp), "DONE") {
            indx := strings.Index(string(resp), "DONE")
            out.Write(resp[0:indx])
            break
        } else if strings.Contains(string(resp), "DATA") {
            log.Println("Received token DATA")
        } else {
            out.Write(resp)
        }
    }

    // TODO: file type? save chunks to disk removing "DATA" flags

    return "", nil
}

func (a *ADBconn) syncCmd(conn net.Conn, cmd, path string) (error) {
    length := uint32(len(path))
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.BigEndian, []byte(cmd))
    binary.Write(buf, binary.LittleEndian, length)
    binary.Write(buf, binary.BigEndian, []byte(path))
    _, err := conn.Write(buf.Bytes())
    if err != nil {
        log.Println("Failed sending sync command ", err)
        return err
    }
    return nil
}

func (a *ADBconn) Sync(cmd, serial, path string) (string, error) {
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
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }
    if err = a.send(conn, SYNC); err != nil {
        log.Println("Error sending command")
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }

    defer conn.Close()

    if err := a.syncCmd(conn, cmd, path); err != nil {
        return "", err
    }

    switch cmd {
    case "LIST":
        return a.readList(conn)
    case "STAT":
        return a.readStat(conn)
    case "RECV":
        return a.readRecv(conn)
    default:
        return "", errors.New(fmt.Sprintf("Command %s is non existent", cmd))
    }

    return "", nil
}

func (a *ADBconn) Track() <-chan string {
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

    go func() {
        for {
            _, resp, err := a.receive(conn)
            if err != nil {
                log.Println("Error receiving data")
                conn.Close()
                break
            }
            out <- string(resp)
        }
    }()

    return out
}

func (a *ADBconn) Send(cmd string) (string, error) {
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
    return string(resp), nil
}

func (a *ADBconn) SendToHost(serial string, cmd string) (string, error) {
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
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return "", errors.New("OKAY header not fouund")
    }
    if err = a.send(conn, cmd); err != nil {
        log.Println("Error sending command")
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
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
        out = append(out, string(resp))
    }
    result := strings.Join(out, "")
    return result, nil
}

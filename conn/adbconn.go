package conn

import (
    "os"
    "io"
    "net"
    "fmt"
    "log"
    "path"
    "bytes"
    "errors"
    "strings"
    "encoding/binary"
)

const (
    PORT = 5037
    HOST_TRANSPORT = "host:transport:<id>"
    TRACK_CMD = "host:track-devices"
    LOLCAT = "shell:logcat 2>/dev/null"
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

func parseList(in *bytes.Buffer) string{
    dataOut := new(bytes.Buffer)
    for ; in.Len() > 0; {
        var mode, size, stat, nLen uint32
        in.Next(4)
        data := bytes.NewReader(in.Next(4))
        binary.Read(data, binary.LittleEndian, &mode)
        data = bytes.NewReader(in.Next(4))
        binary.Read(data, binary.LittleEndian, &size)
        data = bytes.NewReader(in.Next(4))
        binary.Read(data, binary.LittleEndian, &stat)
        data = bytes.NewReader(in.Next(4))
        binary.Read(data, binary.LittleEndian, &nLen)
        fname := in.Next(int(nLen))
        line := fmt.Sprintf("%s mode:%d size:%d stat:%d name lenght:%d\n", fname, mode, size, stat, nLen)
        dataOut.Write([]byte(line))
    }
    return dataOut.String()
}

func (a *ADBconn) readLoop (conn net.Conn) (*bytes.Buffer, error){
    out := new(bytes.Buffer)
    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, err
        }
        if strings.Contains(string(resp), "DONE") {
            indx := strings.Index(string(resp), "DONE")
            out.Write(resp[0:indx])
            break
        } else {
            out.Write(resp)
        }
    }
    return out, nil
}

func (a *ADBconn) readList(conn net.Conn) (string, error){
    out, err := a.readLoop(conn)
    if err != nil{
        return "", err
    }
    return parseList(out), nil
}

func (a *ADBconn) readRecv (conn net.Conn, filename string) (string, error){

    total := uint64(0)
    f, err := os.Create(filename)
    if err != nil{
        return "", err
    }
    defer f.Close()

    for {
        _, resp, err := a.receive(conn)
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err
        }
        if indx := strings.Index(string(resp), "DONE"); indx != -1 {
            count, err := f.Write(resp[0:indx])
            if err != nil {
                return "", err
            }
            total = total + uint64(count)
            break
        } else if indx := strings.Index(string(resp), "DATA"); indx != -1 {
            data := append(resp[:indx], resp[indx+8:]...)
            count, err := f.Write(data)
            if err != nil {
                return "", err
            }
            total = total + uint64(count)
        } else {
            count, err := f.Write(resp)
            if err != nil {
                return "", err
            }
            total = total + uint64(count)
        }
    }

    f.Sync()

    return fmt.Sprintf("Downloaded [%s] %d bytes", filename, total), nil
}

func (a *ADBconn) syncCmd(conn net.Conn, cmd, filePath string) (error) {
    length := uint32(len(filePath))
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.BigEndian, []byte(cmd))
    binary.Write(buf, binary.LittleEndian, length)
    binary.Write(buf, binary.BigEndian, []byte(filePath))
    _, err := conn.Write(buf.Bytes())
    if err != nil {
        log.Println("Failed sending sync command ", err)
        return err
    }
    return nil
}

func (a *ADBconn) sync(conn net.Conn, serial string) (error){
    host := strings.Replace(HOST_TRANSPORT, "<id>", serial, 1)
    if err := a.send(conn, host); err != nil {
        log.Println("Error sending transport")
        return err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return errors.New("OKAY header not found: " + string(resp[0:4]))
    }
    if err := a.send(conn, SYNC); err != nil {
        log.Println("Error sending command")
        return err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return errors.New("OKAY header not found: " + string(resp[0:4]))
    }
    return nil
}

func (a *ADBconn) pushFile(conn net.Conn, srcPath string) (string, error) {
    buff := make([]byte, 8192)
    f, err := os.Open(srcPath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    total := uint64(0)
    for {
        count, err := f.Read(buff)
        if err != nil {
            if err == io.EOF {
                fileStat := uint32(0)
                tmp := new(bytes.Buffer)
                binary.Write(tmp, binary.BigEndian, []byte("DONE"))
                binary.Write(tmp, binary.LittleEndian, fileStat)
                _, err = conn.Write(tmp.Bytes())
                break
            }
            return "", err
        }
        outBuff := new(bytes.Buffer)
        length := uint32(count)
        binary.Write(outBuff, binary.BigEndian, []byte("DATA"))
        binary.Write(outBuff, binary.LittleEndian, length)
        _, err = conn.Write(outBuff.Bytes())
        if err != nil {
            log.Println("Failed sending sync command ", err)
            return "", err
        }
        _, err = conn.Write(buff[0:count])
        if err != nil {
            log.Println("Failed sending sync command ", err)
            return "", err
        }
        total = total + uint64(length)
    }
    return fmt.Sprintf("%s, %d bytes transferred", srcPath, total), nil
}


func (a *ADBconn) Push(serial, srcPath, destPath string) (string, error) {
    conn, err := a.Connect()
    if err != nil {
        log.Println("Error connecting: ", err)
        return "", err
    }
    defer conn.Close()
    if err := a.sync(conn, serial); err != nil {
        return "", err
    }
    filePath := fmt.Sprintf("%s,666", destPath)
    if err := a.syncCmd(conn, "SEND", filePath); err != nil {
        return "", err
    }
    return a.pushFile(conn, srcPath)
}


func (a *ADBconn) Sync(cmd, serial, filePath string) (string, error) {
    conn, err := a.Connect()
    if err != nil {
        log.Println("Error connecting: ", err)
        return "", err
    }
    defer conn.Close()
    if err := a.sync(conn, serial); err != nil {
        return "", err
    }
    if err := a.syncCmd(conn, cmd, filePath); err != nil {
        return "", err
    }
    switch cmd {
    case "LIST":
        return a.readList(conn)
    case "STAT":
        return a.readStat(conn)
    case "RECV":
        filename := path.Base(filePath)
        return a.readRecv(conn, filename)
    default:
        return "", errors.New(fmt.Sprintf("Command %s is non existent", cmd))
    }
    return "", nil
}

func (a *ADBconn) Track() (<-chan string, error) {
    // Tracks state change in connected devices
    conn, err := a.Connect()
    if err != nil {
        return nil, err
    }
    if err = a.send(conn, TRACK_CMD); err != nil {
        return nil, err
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
    return out, nil
}

func (a *ADBconn) Logcat(serial string) (<-chan string, error){
    // Streams out logs
    conn, err := a.Connect()
    if err != nil {
        return nil, err
    }
    host := strings.Replace(HOST_TRANSPORT, "<id>", serial, 1)
    if err := a.send(conn, host); err != nil {
        return nil, err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return nil, err
    }
    if err = a.send(conn, LOLCAT); err != nil {
        return nil, err
    }
    out := make(chan string)
    go func() {
        for {
            _, resp, err := a.receive(conn)
            if err != nil {
                log.Println("Error receiving data", err)
                conn.Close()
                break
            }
            out <- string(resp)
        }
    }()
    return out, nil
}

func (a *ADBconn) Send(cmd string) (string, error) {
    // Send command to host
    conn, err := a.Connect()
    if err != nil {
        return "", err
    }
    defer conn.Close()
    if err = a.send(conn, cmd); err != nil {
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
        return "", err
    }
    defer conn.Close()
    out := []string{}
    host := strings.Replace(HOST_TRANSPORT, "<id>", serial, 1)
    if err = a.send(conn, host); err != nil {
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return "", errors.New("OKAY header not found")
    }
    if err = a.send(conn, cmd); err != nil {
        return "", err
    }
    if _, resp, _ := a.receive(conn); strings.Contains(string(resp), "OKAY") != true {
        return "", errors.New("OKAY header not found")
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

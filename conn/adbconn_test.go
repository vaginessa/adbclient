package conn

import (
    "testing"
    "strings"
    "fmt"
)

func Test_formatCommand(t *testing.T){
    adbc := &ADBconn{}

    frmtdCmd := adbc.formatCommand("host:devices")

    if strings.Compare(frmtdCmd, "000chost:devices") != 0{
        t.Error(fmt.Sprintf("FormatCommand returned an incorrect string: %s", frmtdCmd))
    }
}

func Test_stripChecksum(t *testing.T){
    FOR_TEST := "OKAY001B075923ba00abcdefg    device"
    adbc := &ADBconn{}
    ret, err := adbc.stripChecksum(FOR_TEST)
    if err != nil{
        t.Error("Error verifying checksum:", err)
    }
    t.Log("Returned:", ret)
}

func Test_stripChecksum_no_length(t *testing.T){
    FOR_TEST := "OKAY0000"
    adbc := &ADBconn{}
    resp, err := adbc.stripChecksum(FOR_TEST)
    if err != nil{
        t.Error("Error verifying checksum:", err)
    }
}

func Test_stripChecksum_fail(t *testing.T){
    FOR_TEST := "FAIL"
    adbc := &ADBconn{}
    _, err := adbc.stripChecksum(FOR_TEST)
    if err == nil{
        t.Error("Error verifying checksum:", err)
    }
}
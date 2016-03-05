package comm

import (
    "testing"
    "strings"
    "fmt"
)

func TestFormatCommand(t *testing.T){
    adbc := &adbclient{
        conn: nil,
    }

    frmtdCmd := adbc.FormatCommand("host:devices")

    if strings.Compare(frmtdCmd, "000chost:devices") != 0{
        t.Error(fmt.Sprintf("FormatCommand returned an incorrect string: %s", frmtdCmd))
    }
}
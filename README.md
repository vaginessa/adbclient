# adbclient

Client for the Android Debug Bridge written in go. This library will allow to interact with adb server through network
connection instead of command line adb invocations. 

API
===

## Version()

Returns ADBD (adb daemon) version

```
package main

import (
    "github.com/alexjch/adbclient"
)

func main(){
    version, err := adbclient.New().Version()
    if err != nil{
        fmt.Println("Unable to obtain version", err)
    }
    fmt.Println(version)
}
```

## Devices()

Returns a list of devices (encapsulated in a Device structure) 

```
package main

import (
    "github.com/alexjch/adbclient"
)

func main(){
    devices, err := adbclient.New().Devices()
    if err != nil{
        fmt.Println("Unable to list devices", err)
    }
    fmt.Println(devices)
}
```

## Track()

## Pull()

## Push()

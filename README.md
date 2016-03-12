# adbclient

Client for the Android Debug Bridge written in go. This library will allow to interact with adb server through network
connection instead of command line adb invocations. 

#### API

#### Version
Returns ADBD (adb daemon) version ```Version() (string)```

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

#### Devices

Returns a list of devices (encapsulated in a Device structure) ```Devices() ([]Device, error)```

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

#### Track

Returns a channel that will produce a stream when a device status changes ```Track() (<-chan)```

```
package main

import (
    "github.com/alexjch/adbclient"
)

func main(){
    devices := adbclient.New().Track()
    for{
        fmt.Println(<-devices)
    }
}
```


#### Pull

Downloads a file from a given device ```Pull(serial, filePath string) (string, error)```

```
package main

import (
    "github.com/alexjch/adbclient"
)

func main(){
    message, err := adbclient.New().Pull("anchdgetsr345sacdf", "test.zip")
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(message)
}
```


#### Push

Push file to device (WIP)


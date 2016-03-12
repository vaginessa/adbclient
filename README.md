# adbclient

Client for the Android Debug Bridge written in go. This package talks to the ADB daemon using a TCP connection 

### API

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

Track changes in devices connected to ADB and reports the changes using a channel ```Track() (<-chan)```

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

Push file to device ```Push(serial, source, destination string) (string, error)```

```
package main

import (
    "github.com/alexjch/adbclient"
)

func main(){
    message, err := adbclient.New().Pull("anchdgetsr345sacdf", "test.zip", "/mnt/sdcard/test.zip")
    if err != nil {
        fmt.Println("Failed with error: ", err)
    }
    fmt.Println(message)
}
```

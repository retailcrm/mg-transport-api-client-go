[![Build Status](https://img.shields.io/travis/retailcrm/mg-transport-api-client-go/master.svg?style=flat-square)](https://travis-ci.org/retailcrm/mg-transport-api-client-go)
[![GitHub release](https://img.shields.io/github/release/retailcrm/mg-transport-api-client-go.svg?style=flat-square)](https://github.com/retailcrm/mg-transport-api-client-go/releases)
[![GoLang version](https://img.shields.io/badge/GoLang-1.8%2C%201.9%2C%201.10-blue.svg?style=flat-square)](https://golang.org/dl/)


# retailCRM Message Gateway Transport API Go client

## Install

```bash
go get -x github.com/retailcrm/mg-transport-api-client-go
```

## Usage

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/retailcrm/mg-transport-api-client-go/v1"
)

func main() {
    var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
    ch := Channel{
        Type: "telegram",
        Settings: ChannelSettings{
            SpamAllowed: false,
            Status: Status{
                Delivered: ChannelFeatureNone,
                Read: ChannelFeatureReceive,
            },
            Text: ChannelSettingsText{
                Creating: ChannelFeatureBoth,
                Editing:  ChannelFeatureSend,
                Quoting:  ChannelFeatureReceive,
                Deleting: ChannelFeatureBoth,
            },
        },
    }

    data, status, err := c.ActivateTransportChannel(ch)

    if err != nil {
        t.Errorf("%d %v", status, err)
    }

    fmt.Printf("%v", data.ChannelID)
}
```

## Documentation

* [GoDoc](https://godoc.org/github.com/retailcrm/mg-transport-api-client-go)

[![Build Status](https://github.com/retailcrm/mg-transport-api-client-go/workflows/ci/badge.svg)](https://github.com/retailcrm/mg-transport-api-client-go/actions)
[![Coverage](https://img.shields.io/codecov/c/gh/retailcrm/mg-transport-api-client-go/master.svg?logo=codecov)](https://codecov.io/gh/retailcrm/mg-transport-api-client-go)
[![GitHub release](https://img.shields.io/github/release/retailcrm/mg-transport-api-client-go.svg?style=flat-square)](https://github.com/retailcrm/mg-transport-api-client-go/releases)
[![GoLang version](https://img.shields.io/badge/go-1.11%2C%201.12%2C%201.13%2C%201.14%2C%201.15-blue.svg?style=flat-square)](https://golang.org/dl/)
[![Godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/retailcrm/mg-transport-api-client-go)


# Message Gateway Transport API Go client

## Install

```bash
go get -u -v github.com/retailcrm/mg-transport-api-client-go
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
        Name: "@my_shopping_bot"
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
            Product: Product{
                Creating: ChannelFeatureSend,
                Deleting: ChannelFeatureSend,
            },
            Order: Order{
                Creating: ChannelFeatureBoth,
                Deleting: ChannelFeatureSend,
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

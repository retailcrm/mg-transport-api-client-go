[![Build Status](https://github.com/retailcrm/mg-transport-api-client-go/workflows/ci/badge.svg)](https://github.com/retailcrm/mg-transport-api-client-go/actions)
[![Coverage](https://img.shields.io/codecov/c/gh/retailcrm/mg-transport-api-client-go/master.svg?logo=codecov&logoColor=white)](https://codecov.io/gh/retailcrm/mg-transport-api-client-go)
[![GitHub release](https://img.shields.io/github/release/retailcrm/mg-transport-api-client-go.svg?logo=github&logoColor=white)](https://github.com/retailcrm/mg-transport-api-client-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/retailcrm/mg-transport-api-client-go)](https://goreportcard.com/report/github.com/retailcrm/mg-transport-api-client-go)
[![GoLang version](https://img.shields.io/badge/go->=1.11-blue.svg?logo=go&logoColor=white)](https://golang.org/dl/)
[![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/retailcrm/mg-transport-api-client-go)


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
        Name: "@my_shopping_bot",
        Settings: ChannelSettings{
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

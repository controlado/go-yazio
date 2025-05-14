# yazio-go

Unofficial Go SDK for the YAZIO mobile API

[![tests-badge](https://github.com/controlado/go-yazio/actions/workflows/test.yml/badge.svg)](https://github.com/controlado/go-yazio/actions/workflows/test.yml)
[![report-badge](https://goreportcard.com/badge/github.com/controlado/go-yazio)](https://goreportcard.com/report/github.com/controlado/go-yazio)
[![mit-badge](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![wakatime-badge](https://wakatime.com/badge/github/controlado/go-yazio.svg)](https://wakatime.com/badge/github/controlado/go-yazio)

## Status

⚠️ Experimental – the API is private and may change at any time. Breaking changes are expected.

## Installation

```bash
go get github.com/controlado/yazio
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/controlado/go-yazio/pkg/client"
    "github.com/controlado/go-yazio/pkg/yazio"
)

const (
    username = "username@email.com"
    password = "superStrongPassword"
)

func main() {
    var (
        ctx = context.Background()
        c   = client.New(
            client.WithBaseURL(yazio.DefaultBaseURL),
        )
    )

    api, err := yazio.New(c)
    if err != nil {
        log.Fatalf("building yazio api: %v", err)
    }

    cred := yazio.NewPasswordCred(username, password)
    user, err := api.Login(ctx, cred)
    if err != nil {
        log.Fatalf("fetching user from api: %v", err)
    }

    userData, err := user.Data(ctx)
    if err != nil {
        log.Fatalf("fetching user data: %v", err)
    }
    // userData.String()
    // User(João da Silva)

    sinceRegist := userData.SinceRegist()
    // sinceRegist.String()
    // 22 January 2023 - 13 May 2025

    userMacros, err := user.Macros(ctx, sinceRegist)
    if err != nil {
        log.Fatalf("fetching user macros (since regist): %v", err)
    }
    // userMacros.String()
    // Average 38 days
    // Kcal: 1659.870
    // Carb: 165.053
    // Fat: 54.531
    // Protein: 128.297
}
```

## Features

* Login with password
* Retrieve user profile & nutrition stats
* Zero external deps beyond the Go standard library
* Context/timeout aware

## TODO

* Register food and snack intake
* Automatic retry with exponential back‑off

## Legal Notice

* **No affiliation** with YAZIO GmbH  
* **YAZIO** is a trademark of YAZIO GmbH  
* **As-is** without warranty (MIT License)  
* **Use at your own risk**: reverse-engineering may violate ToS  

## Contributing

Contributions are welcome! Please open an issue or pull request.  
By contributing, you agree to release your work under the MIT License.

## License

MIT – see [`LICENSE`](./LICENSE) for full text.

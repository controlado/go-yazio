# yazio-go

Unofficial Go SDK for the YAZIO mobile API

[![mit-badge](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![report-badge](https://goreportcard.com/badge/github.com/controlado/go-yazio)](https://goreportcard.com/report/github.com/controlado/go-yazio)
[![tests-badge](https://github.com/controlado/go-yazio/actions/workflows/test.yml/badge.svg)](https://github.com/controlado/go-yazio/actions/workflows/test.yml)
[![pkg-go-badge](https://pkg.go.dev/badge/github.com/controlado/go-yazio.svg)](https://pkg.go.dev/github.com/controlado/go-yazio)
[![wakatime-badge](https://wakatime.com/badge/github/controlado/go-yazio.svg)](https://wakatime.com/badge/github/controlado/go-yazio)

## Status

⚠️ Experimental – the API is private and may change at any time. Breaking changes are expected.

## Installation

```bash
go get github.com/controlado/yazio
```

## Quick Start

### Auth

```go
const (
    username = "email@email.com"
    password = "superStrongPass"
)

api, err := yazio.New(c)
if err != nil { // yazio.ErrClientCannotBeNil
    log.Fatalf("building yazio api: %v", err)
}

cred := yazio.NewPasswordCred(username, password)
user, err := api.Login(ctx, cred)
if err != nil {
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
    // yazio.ErrInvalidCredentials
    log.Fatalf("fetching user from api: %v", err)
}
```

### Get user-data

Can be used after use [`Login`](#auth)

```go
userData, err := user.Data(ctx)
if err != nil {
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
    log.Fatalf("fetching user data from api: %v", err)
}
// userData.String()
// User(João da Silva)

sinceRegist := userData.SinceRegist()
// sinceRegist.String()
// 1 June 2023 - 1 December 2023 (183 days)

userMacros, err := user.Macros(ctx, sinceRegist)
if err != nil {
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
    log.Fatalf("fetching user macros intakes (since regist): %v", err)
}
// userMacros.Average().String()
// Average 183 days
// Kcal: 1859.953
// Carb: 120.274
// Fat: 50.872
// Protein: 167.89

sugarIntakes, err := user.Intake(ctx, intake.Sugar, sinceRegist)
if err != nil {
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
    log.Fatalf("fetching user sugar intakes (since regist): %v", err)
}
// sugarIntakes.Average().String()
// 39 days: 37.99
```

### Registering new food

Can be used after use [`Login`](#auth)

```go
var (
    foodName = "Banana"
    foodCat  = food.Miscellaneous
    foodNut  = food.Nutrients{ // required nutrients
        intake.Energy:  0.1,
        intake.Fat:     0.1,
        intake.Protein: 0.1,
        intake.Carb:    0.1,
    }
)

f, err := food.New(foodName, foodCat, foodNut)
if err != nil { // food.ErrInvalidName
    log.Fatalf("creating a new food: %v", err)
}

if err := user.AddFood(ctx, f, visibility.PrivateFood); err != nil {
    // yazio.ErrRequestingToYazio
    // food.ErrMissingNutrients
    // food.ErrAlreadyExists
    log.Fatalf("adding new food %s: %v", f, err)
}
```

## Features

* Login with password
* Register food to user
* Retrieve user profile & nutrition stats
* Zero external deps beyond the Go standard library
* Context/timeout aware

## TODO

* Food intake
* Get registered food
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

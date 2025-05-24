<div align="center">
    <h1>yazio-go</h1>
    <p>Unofficial Go client for accessing the <a href="https://www.yazio.com/" rel="noopener noreferrer">YAZIO</a> API</p>
    <p>
        <a href="https://github.com/controlado/go-yazio/actions/workflows/test.yml"><img alt="badge: actions-tests" src="https://github.com/controlado/go-yazio/actions/workflows/test.yml/badge.svg"></a>
        <a href="https://goreportcard.com/report/github.com/controlado/go-yazio"><img alt="badge: go-report-card" src="https://goreportcard.com/badge/github.com/controlado/go-yazio?style=default"></a>
        <a href="https://codecov.io/gh/controlado/go-yazio"><img alt="badge: codecov" src="https://codecov.io/gh/controlado/go-yazio/branch/main/graph/badge.svg?token=Fgo4zed2G1"></a>
        <a href="https://pkg.go.dev/github.com/controlado/go-yazio"><img alt="badge: pkg-reference" src="https://img.shields.io/static/v1?logo=go&label=Reference&message=go-yazio&color=0476b7&style=default"></a>
        <a href="https://wakatime.com/badge/github/controlado/go-yazio"><img alt="badge: wakatime" src="https://wakatime.com/badge/github/controlado/go-yazio.svg?style=default"></a>
    </p>
</div>

## Status

Depends on YAZIO’s **private** API  
Expect breaking changes at any time

## Installation

```bash
go get github.com/controlado/go-yazio
```

## Usage examples

<details>
    <summary>
        <strong>Authenticate</strong>
    </summary>

```go
const (
    username = "email@email.com"
    password = "superStrongPass"
)

var (
    ctx = context.Background()
)

api, err := yazio.New()
if err != nil {
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

</details>

<details>
    <summary>
        <strong>Refresh session</strong>
    </summary>

```go
userToken := user.Token()
// userToken.String()
// Token(Expired)

if userToken.IsExpired() {
    if err := api.Refresh(ctx, user); err != nil {
        // yazio.ErrRequestingToYazio
        // yazio.ErrDecodingResponse
        log.Fatalf("refreshing user token: %v", err)
    }
}
```

</details>

<details>
    <summary>
        <strong>Get user data</strong>
    </summary>

```go
userData, err := user.Data(ctx)
if err != nil {
    // yazio.ErrExpiredToken
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
    // yazio.ErrExpiredToken
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
    log.Fatalf("fetching user macros intakes (since regist): %v", err)
}
// userMacros.Average().String()
// Average (1278 days)
// Energy: 1700.0kcal
// Carb: 150.0g
// Fat: 40.0g
// Protein: 180.0g

waterIntakes, err := user.Intake(ctx, intake.Water, sinceRegist)
if err != nil {
    // yazio.ErrExpiredToken
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
    log.Fatalf("fetching user water intakes (since regist): %v", err)
}
// waterIntakes.Average().String()
// 320 days: 2223.0ml
```

</details>

<details>
    <summary>
        <strong>Register new food (product)</strong>
    </summary>

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
    // yazio.ErrExpiredToken
    // yazio.ErrRequestingToYazio
    // food.ErrMissingNutrients
    // food.ErrAlreadyExists
    log.Fatalf("adding new food %s: %v", f, err)
}
```

</details>

## Features

* Login with password
* Register food to account
* Retrieve user profile & nutrition stats
* Zero external deps beyond the Go standard library
* Context/timeout aware

## Legal Notice

* **No affiliation** with YAZIO GmbH  
* **YAZIO** is a trademark of YAZIO GmbH  
* **As-is** without warranty (MIT License)  
* **Use at your own risk**: reverse-engineering may violate ToS  

## Contributing

Contributions are welcome  
Please open an [issue](https://github.com/controlado/go-yazio/issues) or [pull request](https://github.com/controlado/go-yazio/pulls)

## License

MIT – see [`LICENSE`](./LICENSE) for full text.

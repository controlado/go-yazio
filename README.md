<div align="center">
    <h1>yazio-go</h1>
    <p>Unofficial Go client for accessing the <a href="https://www.yazio.com/" rel="noopener noreferrer">YAZIO</a> API</p>
    <p>
        <a href="https://github.com/controlado/go-yazio/actions/workflows/test.yml"><img alt="badge: actions-tests" src="https://github.com/controlado/go-yazio/actions/workflows/test.yml/badge.svg"></a>
        <a href="https://goreportcard.com/report/github.com/controlado/go-yazio/v4"><img alt="badge: go-report-card" src="https://goreportcard.com/badge/github.com/controlado/go-yazio/v4?style=default"></a>
        <a href="https://codecov.io/gh/controlado/go-yazio"><img alt="badge: codecov" src="https://codecov.io/gh/controlado/go-yazio/branch/main/graph/badge.svg?token=Fgo4zed2G1"></a>
        <a href="https://pkg.go.dev/github.com/controlado/go-yazio/v4"><img alt="badge: pkg-reference" src="https://img.shields.io/static/v1?logo=go&label=Reference&message=go-yazio&color=0476b7&style=default"></a>
        <a href="https://wakatime.com/badge/github/controlado/go-yazio"><img alt="badge: wakatime" src="https://wakatime.com/badge/github/controlado/go-yazio.svg?style=default"></a>
    </p>
</div>

## Status

Depends on YAZIO’s **private** API  
Expect breaking changes at any time

## Installation

```bash
go get github.com/controlado/go-yazio/v4
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
    // yazio.ErrCredentialsCannotBeNil
    // yazio.ErrInvalidCredentials
    // yazio.ErrRequestingToYazio
    // yazio.ErrDecodingResponse
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
        <strong>Enter and register new food (product)</strong>
    </summary>

```go
newFood, err := food.New(
    "Sweet popcorn",
    food.Miscellaneous,
    food.Nutrients{
        intake.Energy:  62,
        intake.Fat:     0,
        intake.Protein: 1,
        intake.Carb:    14,
    },
    food.WithBaseUnit(unit.Gram),        // use grams as the food's base unit
    food.WithNutrientsPer(18),           // base the nutrients above on 18 grams
    food.WithNewServing(food.Pack, 50),  // add a 50-gram pack option in the app
)
if err != nil { // food.ErrInvalidName
    log.Fatalf("creating a new food: %v", err)
}

if err := user.AddFood(ctx, newFood, visibility.PrivateFood); err != nil {
    // yazio.ErrExpiredToken
    // yazio.ErrRequestingToYazio
    // food.ErrAlreadyExists
    // food.ErrMissingNutrients
    // food.ErrInvalidNutrientReferenceAmount
    log.Fatalf("adding new food %s: %v", newFood, err)
}

var ( // one 18-gram portion
    entryServing    = serving.Portion(18)
    servingQuantity = 1
)

if err := user.EntryFood(ctx, meal.Dinner, newFood.ID, entryServing, servingQuantity); err != nil {
    // yazio.ErrExpiredToken
    // food.ErrInvalidServing
    // food.ErrInvalidServingQuantity
    // yazio.ErrRequestingToYazio
    log.Fatalf("entering new food %s: %v", newFood, err)
}
```

</details>

## FAQ

<details>
    <summary>
        <strong>Why doesn't a food registered with <code>AddFood</code> appear immediately in the app search?</strong>
    </summary>

`AddFood` registers the food on YAZIO's servers, but some versions of the app
keep a local product catalog that is not refreshed when an external client adds
a food. In our testing, creating another product in the app refreshed the
catalog and made externally created foods searchable. Restarting the app,
clearing its cache, or logging the food did not reliably refresh the catalog.

There is currently no known API endpoint to force this refresh.

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

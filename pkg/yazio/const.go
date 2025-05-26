package yazio

import (
	"github.com/controlado/go-yazio/internal/application"
	"github.com/controlado/go-yazio/internal/infra/client"
)

func defaultHeaders(tk application.Token) client.Payload[string] {
	headers := client.Payload[string]{
		`accept`:          `*/*`,
		`accept-charset`:  `UTF-8`,
		`connection`:      `Keep-Alive`,
		`host`:            `yzapi.yazio.com`,
		`accept-encoding`: `application/json`,
		`user-agent`:      `YAZIO/12.31.0 (com.yazio.android; build:411052340; Android 34) Ktor`,
	}

	if tk != nil {
		headers[`authorization`] = tk.Bearer()
	}

	return headers
}

// API
const (
	baseURL string = "https://yzapi.yazio.com"

	loginEndpoint         string = "/v18/oauth/token"
	userDataEndpoint      string = "/v18/user"
	entryFoodEndpoint     string = "/v18/user/consumed-items"
	addFoodEndpoint       string = "/v18/user/products"
	singleIntakesEndpoint string = "/v18/user/consumed-items/specific-nutrient-daily"
	macrosIntakesEndpoint string = "/v18/user/consumed-items/nutrients-daily"
)

// Time layout
const (
	layoutISO  string = "2006-01-02"
	layoutDate string = "2006-01-02 15:04:05"
)

// Static data (?)
const (
	confirmedEmailStatus string = "confirmed"
	defaultClientID      string = "1_4hiybetvfksgw40o0sog4s884kwc840wwso8go4k8c04goo4c"
	defaultSecret        string = "6rok2m65xuskgkgogw40wkkk8sw0osg84s8cggsc4woos4s8o"
)

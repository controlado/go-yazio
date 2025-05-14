package yazio

// public
const (
	DefaultBaseURL string = "https://yzapi.yazio.com"
)

// endpoints
const (
	loginEndpoint    string = "/v18/oauth/token"
	userDataEndpoint string = "/v18/user"
	intakeEndpoint   string = "/v18/user/consumed-items/specific-nutrient-daily"
	macrosEndpoint   string = "/v18/user/consumed-items/nutrients-daily"
)

// time layouts
const (
	layoutISO  string = "2006-01-02"
	layoutDate string = "2006-01-02 15:04:05"
)

const (
	confirmedEmailStatus string = "confirmed"
	defaultClientID      string = "1_4hiybetvfksgw40o0sog4s884kwc840wwso8go4k8c04goo4c"
	defaultSecret        string = "6rok2m65xuskgkgogw40wkkk8sw0osg84s8cggsc4woos4s8o"
)

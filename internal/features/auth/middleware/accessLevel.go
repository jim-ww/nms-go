package middleware

type AccessLevel int

const (
	AllowUnauthorized AccessLevel = iota
	OnlyAuthorized
	OnlyAdmins
	OnlyUnauthorized
)

var accessLevelString = map[AccessLevel]string{
	AllowUnauthorized: "AllowUnauthorized",
	OnlyUnauthorized:  "OnlyUnauthorized",
	OnlyAuthorized:    "OnlyAuthorized",
	OnlyAdmins:        "OnlyAdmins",
}

func (a AccessLevel) String() string {
	return accessLevelString[a]
}

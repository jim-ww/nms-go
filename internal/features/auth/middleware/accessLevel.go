package middleware

import "log/slog"

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

func (access AccessLevel) SlogAttr(route string) slog.Attr {
	return slog.Group("route",
		slog.String("path", route),
		slog.String("access-level", access.String()),
	)
}

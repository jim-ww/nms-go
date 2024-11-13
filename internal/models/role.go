package models

// type Role int
//
// const (
// 	ROLE_USER Role = iota
// 	ROLE_ADMIN
// )

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

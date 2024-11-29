package models

type User struct {
	Service  string `bson:"service"`
	Nickname string `bson:"nickname"`
	Email    string `bson:"email"`
}

package models

type LoginRequest struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Nickname string `bson:"nickname"`
}

type User struct {
	Uid      string `bson:"uid"`
	Nickname string `bson:"nickname"`
	Email    string `bson:"email"`
	Pw       string `bson:"pw"`
}

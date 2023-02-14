package main

import (
	"fmt"
	mdb_goson "github.com/maldan/go-ml/db/goson"
	ml_console "github.com/maldan/go-ml/io/console"
	"time"
)

type User struct {
	Id               int    `json:"id"`
	Password         string `json:"password"`
	Username         string `json:"username"`
	Phone            string `json:"phone"`
	PromoCode        string `json:"promo_code"`
	SmsCode          int    `json:"sms_code"`
	Email            string `json:"email"`
	IsPhoneActivated bool   `json:"phone_activated"`

	SalonName    string `json:"salon_name"`
	SalonAddress string `json:"address"`
	SalonLogo    string `json:"salon_logo"`

	SubscriptionName string `json:"subscription_name"`
	// SubscriptionType    string    `json:"subscription_type"`
	SubscriptionExpires time.Time `json:"subscription_expires"`

	AvailablePhoto     int    `json:"available_photo"`
	AvailableDocuments int    `json:"available_documents"`
	GavnoId            uint32 `json:"gavno_id"`

	OverridePermission uint64 `json:"override_permission"`

	LastLogin  time.Time `json:"last_login"`
	DateJoined time.Time `json:"date_joined"`
}

func main() {
	userDb := mdb_goson.New[User]("db")
	// userDb.Insert(User{Username: "lox", Password: "oglox"})
	sr := userDb.FindBy(mdb_goson.ArgsFind[User]{
		FieldList: "Username",
		Where: func(u *User) bool {
			fmt.Printf("%v\n", u)
			return u.Username == "lox"
		},
	})
	ml_console.PrettyPrint(sr.Unpack())
}

package mdb_goson_test

import (
	"fmt"
	mdb_goson "github.com/maldan/go-ml/db/goson"
	ml_console "github.com/maldan/go-ml/io/console"
	"testing"
	"time"
)

type TestUser struct {
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

	AvailablePhoto     int `json:"available_photo"`
	AvailableDocuments int `json:"available_documents"`

	OverridePermission uint64 `json:"override_permission"`

	LastLogin  time.Time `json:"last_login"`
	DateJoined time.Time `json:"date_joined"`
}

func TestXXX3(t *testing.T) {
	userDb := mdb_goson.New[TestUser]("../../trash/db")
	userDb.Insert(TestUser{Username: "lox", Password: "oglox"})
	sr := userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Username",
		Where: func(u *TestUser) bool {
			fmt.Printf("%v\n", u)
			return u.Username == "lox"
		},
	})
	ml_console.PrettyPrint(sr.Unpack())
}

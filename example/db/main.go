package main

import (
	"fmt"
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	"github.com/maldan/go-ml/db/mdb"
	ml_console "github.com/maldan/go-ml/util/io/console"

	ml_time "github.com/maldan/go-ml/util/time"
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
	SubscriptionExpires ml_time.DateTime `json:"subscription_expires"`

	AvailablePhoto     int    `json:"available_photo"`
	AvailableDocuments int    `json:"available_documents"`
	GavnoId            uint32 `json:"gavno_id"`

	OverridePermission uint64 `json:"override_permission"`

	LastLogin  ml_time.DateTime `json:"last_login"`
	DateJoined ml_time.DateTime `json:"date_joined"`
}

func main() {
	userDb := mdb.New[User](".", "db2", &gosn_driver.Container{})

	// userDb.SetBackupSchedule("../../trash", time.Second)

	// userDb.Insert(User{Username: "lox", Password: "oglox"})

	sr := userDb.FindBy(mdb.ArgsFind[User]{
		FieldList: "Username",
		Where: func(u *User) bool {
			fmt.Printf("%v\n", u)
			return u.Username == "lox"
		},
	})
	ml_console.PrettyPrint(sr.Result)
}

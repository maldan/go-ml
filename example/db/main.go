package main

import (
	"database/sql"
	"fmt"
	ml_console "github.com/maldan/go-ml/util/io/console"
	ml_sql "github.com/maldan/go-ml/util/sql"
	_ "modernc.org/sqlite"
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
	// userDb := mdb.New(".", "db2", User{}, &gosn_driver.Container{})

	// userDb.SetBackupSchedule("../../trash", time.Second)
	// userDb.Insert(User{Username: "lox", Password: "oglox"})

	/*sr := userDb.FindBy(mdb.ArgsFind{
		Where: func(u any) bool {
			return u.(User).Username == "lox"
		},
	})
	ml_console.PrettyPrint(sr.Result)*/

	db, err := sql.Open("sqlite", "sas.db")
	fmt.Printf("%v\n", err)
	/*fmt.Printf("%v\n", db)
	err = ml_sql.CreateTable[User](db, "user")
	fmt.Printf("%v\n", err)
	err = ml_sql.Insert[User](db, "user", User{
		IsPhoneActivated: true,
		SalonName:        "Gay",
		LastLogin:        time.Now(),
		DateJoined:       time.Now(),
	})
	fmt.Printf("%v\n", err)*/

	u, err := ml_sql.SelectMany[User](db, "user", "1=1", 1)
	ml_console.PrettyPrint(u)
	fmt.Printf("%+v\n", err)
}

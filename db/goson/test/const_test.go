package test_test

import (
	"time"
)

type TestA struct {
	Name  string
	Crazy TestB
}

type TestB struct {
	Name string
}

type Test struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName" len:"32"`
	Phone     string `json:"phone" len:"64"`
	/*Sex       string `json:"sex" id:"3" len:"64"`
	Rock      string `json:"rock" id:"4" len:"64"`
	Gas       string `json:"gas" id:"5" len:"64"`
	Yas       string `json:"yas" id:"6" len:"64"`
	Taj       string `json:"taj" id:"7" len:"64"`
	Mahal     string `json:"mahal" id:"8" len:"64"`
	Ebal      string `json:"ebal" id:"9" len:"64"`
	Sasal     string `json:"sasal" id:"10" len:"64"`
	Sasal2    string `json:"sasal2" id:"11" len:"64"`*/
}

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

	SubscriptionName    string    `json:"subscription_name"`
	SubscriptionType    string    `json:"subscription_type"`
	SubscriptionExpires time.Time `json:"subscription_expires"`

	AvailablePhoto     int `json:"available_photo"`
	AvailableDocuments int `json:"available_documents"`

	OverridePermission uint64 `json:"override_permission"`

	LastLogin  time.Time `json:"last_login"`
	DateJoined time.Time `json:"date_joined"`
}

type StructNumber struct {
	Balance int
}

type StructString struct {
	Balance string
}

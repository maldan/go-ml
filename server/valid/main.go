package ms_valid

import (
	"fmt"
	ms_error "github.com/maldan/go-ml/server/error"

	"reflect"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func CheckEmail(email string) string {
	// Trim email from spaces
	email = strings.Trim(strings.ToLower(email), " ")

	// Check
	if !emailRegex.MatchString(email) {
		ms_error.Fatal(ms_error.Error{Code: 500, Field: "email", Description: "Incorrect email"})
	}

	// Return cleaned email
	return email
}

func CheckPassword(password1 string, password2 string) {
	// Check password
	if len(password1) < 6 {
		ms_error.Fatal(ms_error.Error{
			Code: 500, Field: "password",
			Description: "Password must contain at least 6 characters",
		})
	}
	if password1 != password2 {
		ms_error.Fatal(ms_error.Error{
			Code: 500, Field: "password",
			Description: "Passwords do not match",
		})
	}
}

func Required[T any](args T, fields ...string) {
	for _, fieldName := range fields {
		f := reflect.ValueOf(args)
		if f.FieldByName(fieldName).IsZero() {
			failField(args, fieldName, ms_error.ErrorTypeRequired, "is required")
		}
	}
}

func TrimAll[T any](args *T) {
	typeOf := reflect.TypeOf(args)
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Type.Kind() == reflect.String {
			Trim(args, []string{typeOf.Field(i).Name})
		}
	}
}

func Trim[T any](args *T, fields []string) {
	for _, field := range fields {
		// struct dereference
		f := reflect.ValueOf(args).Elem()

		if f.FieldByName(field).Kind() == reflect.String {
			value := f.FieldByName(field).Interface().(string)

			if f.FieldByName(field).CanSet() {
				fmt.Printf("%v\n", value)
				f.FieldByName(field).SetString(strings.Trim(value, " "))
			}
		}
	}
}

func MatchRegExp(args any, r string, fields ...string) {
	var rr = regexp.MustCompile(r)

	for _, fieldName := range fields {
		f := reflect.ValueOf(args)
		if !rr.MatchString(f.FieldByName(fieldName).String()) {
			failField(args, fieldName, ms_error.ErrorTypeUnknown, "must pass regex "+r)
		}
	}
}

func failField(args any, fieldName string, errorType string, text string) {
	tf, _ := reflect.TypeOf(args).FieldByName(fieldName)

	// Get field name
	if tf.Tag.Get("json") != "" {
		fieldName = tf.Tag.Get("json")
	}

	ms_error.Fatal(ms_error.Error{
		Type:        errorType,
		Field:       fieldName,
		Description: fmt.Sprintf("Field '%v' %v", fieldName, text),
	})
}

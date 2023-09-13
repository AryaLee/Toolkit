package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func main() {
	cidrT()
}

type ValidateStruct struct {
	IP   string `validate:"required,ipv4" json:"ip"`
	CIDR string `validate:"required,cidrv4" json:"cidr"`
	MAC  string `validate:"omitempty,mac"`
}

func cidrT() {
	validate := validator.New()
	validate.RegisterTagNameFunc(JsonTagNameFunc)

	s := ValidateStruct{}
	err := validate.Struct(&s)
	fmt.Println("empty", err)
	fmt.Println("")

	s.IP = "192.168."
	err = validate.Struct(&s)
	fmt.Println("invalidIP", err)
	fmt.Println("")

	s.IP = "192.168.0.20"
	err = validate.Struct(&s)
	fmt.Println("validIP", err)
	fmt.Println("")

	s.CIDR = "192.168.0.20/-1"
	err = validate.Struct(&s)
	fmt.Println("invalidCIDR", err)
	fmt.Println("")

	s.CIDR = "192.168.0.20/24"
	err = validate.Struct(&s)
	fmt.Println("validCIDR", err)
	fmt.Println("")

	s.MAC = "rt:bb:cc:dd:ee:ff"
	err = validate.Struct(&s)
	fmt.Println("invalidMAC", err)
	fmt.Println("")

	s.MAC = "aa:bb:cc:dd:ee:ff"
	err = validate.Struct(&s)
	fmt.Println("validMAC", err)
	fmt.Println("")
}

func varT() {
	validate := validator.New()

	var boolTest bool
	err := validate.Var(boolTest, "required")
	if err != nil {
		fmt.Println(err)
	}
	var stringTest string = ""
	err = validate.Var(stringTest, "required")
	if err != nil {
		fmt.Println(err)
	}

	var emailTest string = "test@126.com"
	err = validate.Var(emailTest, "email")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success") // 输出： success。 说明验证成功
	}

	emailTest2 := "test.126.com"
	errs := validate.Var(emailTest2, "required,email")
	if errs != nil {
		fmt.Println(errs) // 输出: Key: "" Error:Field validation for "" failed on the "email" tag。验证失败
	}

	fmt.Println("\r\nEnd!!")
}

// validate.RegisterTagNameFunc(JsonTagNameFunc)
var JsonTagNameFunc = func(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}

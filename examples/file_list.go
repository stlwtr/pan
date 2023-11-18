package main

import (
	"fmt"

	"github.com/stlwtr/pan/file"
)

const (
	MyAccessToken = "123.7228b6b3305e11a8ee9ae643a761690a.Ygw72qJu2lYm_Y99qvCEAl4ILtGzZJAcxXpgOuD.m9PMlw"
)

func main() {
	accessToken := MyAccessToken
	fileClient := file.NewFileClient(accessToken)
	res, err := fileClient.List("/apps/bypy", 0, 100)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(res)
}

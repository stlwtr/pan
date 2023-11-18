package main

import (
	"encoding/json"
	"fmt"

	"github.com/stlwtr/pan/file"
)

func main() {
	accessToken := "123.7228b6b3305e11a8ee9ae643a761690a.Ygw72qJu2lYm_Y99qvCEAl4ILtGzZJAcxXpgOuD.m9PMlw"
	fileClient := file.NewFileClient(accessToken)
	res, err := fileClient.Listall("/apps/bypy", 0, 0, 100)
	b, _ := json.Marshal(res)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println(string(b))
}

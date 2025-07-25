package main

import (
	"fmt"
)

func main() {
	db := GetDB()
	defer db.Close()

	fmt.Println("DB connected:", db != nil)
}

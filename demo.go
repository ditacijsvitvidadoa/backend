package main

import (
	"fmt"
	passwordHash "github.com/vzglad-smerti/password_hash"
)

func main() {
	password := "yourPassword"
	secondPassword := "yourPassword"

	hash, _ := passwordHash.Hash(password)
	fmt.Println("hash", hash)

	secondhash, _ := passwordHash.Hash(password)
	fmt.Println("secondhash", secondhash)

	hashVerify, _ := passwordHash.Verify(hash, secondPassword)

	fmt.Println("hash_verify", hashVerify)

}

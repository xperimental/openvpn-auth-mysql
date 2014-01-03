package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"os"
	"strings"
)

var exitCode int = 1

func exitWithCode() {
	os.Exit(exitCode)
}

func main() {
	defer exitWithCode()

	user, pass := credentialsFromEnvironment()

	if len(user) == 0 || len(pass) == 0 {
		fmt.Println("Username / Password not available.")
		return
	}

	dbConfig := getServerConfig()

	db := mysql.New("tcp", "", dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database)
	err := db.Connect()
	if err != nil {
		fmt.Printf("Can not connect to database: %v\n", err)
		return
	}

	defer db.Close()

	rows, _, err := db.Query("select password from openvpn_users where name = '%s'", user)
	if err != nil {
		fmt.Printf("Error getting data: %v\n", err)
		return
	}

	if len(rows) == 0 {
		fmt.Println("No users found.")
		return
	}

	for _, row := range rows {
		hashed := row.Str(0)
		hashedTokens := strings.Split(hashed, "|")
		if len(hashedTokens) != 3 {
			fmt.Printf("Invalid hash string: %s\n", hashed)
			return
		}
		salt := hashedTokens[0]
		hashAlg := hashedTokens[1]
		if hashAlg != "sha256" {
			fmt.Printf("Currently only supports SHA-256 hashes: %s\n", hashAlg)
			return
		}
		dbHash := hashedTokens[2]
		myHashBytes := sha256.Sum256([]byte(salt + pass))
		myHash := hex.EncodeToString(myHashBytes[:])
		if dbHash == myHash {
			fmt.Println("User valid.")
			exitCode = 0
			return
		} else {
			fmt.Println("User invalid.")
			return
		}
	}
}

func credentialsFromEnvironment() (string, string) {
	user := os.Getenv("username")
	pass := os.Getenv("password")
	return user, pass
}

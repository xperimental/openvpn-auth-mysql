package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
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

	var user, pass string
	flag.Parse()
	if flag.NArg() == 0 {
		user, pass = credentialsFromEnvironment()
	} else {
		user, pass = credentialsFromFile()
	}

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

func credentialsFromFile() (string, string) {
	fileName := flag.Arg(0)
	reader, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening credentials file: %v\n", err)
		return "", ""
	}
	defer reader.Close()

	input := bufio.NewScanner(reader)
	if input.Scan() {
		user := input.Text()
		if input.Scan() {
			pass := input.Text()
			return user, pass
		}
		return user, ""
	}
	return "", ""
}

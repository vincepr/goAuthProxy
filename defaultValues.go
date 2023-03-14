package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// globals with its default values
var (
	port 		= "3002"
	urlProxy 	= "http://127.0.0.1:3001"
	username 	= "username"
	secret		= "default-2_#123default-2_#123"
	password	[]byte 
	storage 	= NewAccountStorage()
)

func InitGlobalValues(){
	pw := ""
	pwHash :=""

	// read files from env if they exist
	rPort 	:= os.Getenv("GRP_PORT")
	rUrl 	:= os.Getenv("GRP_URL")
	rUser 	:= os.Getenv("GRP_USER")
	rPwHash := os.Getenv("GRP_PASSWORD_HASH")
	rSecret := os.Getenv("GRP_SECRET")

	switch{
	case rPort !="": port = rPort
	case rUrl !="": urlProxy = rUrl
	case rUser !="": username = rUser
	case rPwHash !="": pwHash = rPwHash
	case rSecret !="": secret = rSecret
	}
	
	// read files from flags, if they exist, take prio over env-variables
	fPort 	:= flag.String("port", "", "the port this proxy should listen for requets from")
	fUrl 	:= flag.String("url", "", "the address and port of the server we proxy to")	
	fUser 	:= flag.String("user", "", "the address and port of the server we proxy to")
	fPw 	:= flag.String("pw", "", "plaintext password")
	fSecret	:= flag.String("secret", "", "the secret used to encrypt the JSW-Tokens")
	fPwHash	:= flag.String("pwhash", "", "hashed pw")
	flag.Parse()
	switch{
	case *fPort != "": 		port = *fPort
	case *fUrl != "": 		urlProxy = *fUrl
	case *fUser != "": 		username = *fUser
	case *fPw != "": 		pw = *fPw
	case *fSecret != "": 	secret = *fSecret
	case *fPwHash != "": 	pwHash = *fPwHash
	}
	fmt.Println(*fPwHash)
	// warnings using with defautl values:
	if username=="username" {fmt.Println("!- using default username: ", username,"change with -user, or set using GRP_USER env")}
	if secret=="default-2_#123" 	{fmt.Println("!- using default secret:", secret,"change with -secret, or set using GRP_SECRET env")}
	
	// decide the password and hash it if required
	if (pwHash !=""){
		password = []byte(pwHash)
		return
	}
	if pw ==""{
		pw = "qwert123"
		fmt.Println("!- using default password:", pw,"change with -pw, or set the hashed-Password directly")
	}
	// encrypt our password and store the hashed bytes 
	encPw, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil{
		panic(err)
	}
	password= encPw
	fmt.Printf("!- Better use -pwhash  or set the env GRP_PASSWORD_HASH directly: \n %v \n", string(encPw))
}


func InitAccounts() {
	// add our user to our storage
	acc := Account{
		Name: username,
		PasswordHash: password,
		IsAdmin: false,
	}
	storage.AddAccount(&acc)

	//format url if without http://
	formatUrl := func(str *string){
		if !strings.HasPrefix(*str, "http://") && !strings.HasPrefix(*str, "https://"){
			*str = "http://" + *str
			fmt.Println("adding http:// str the URL -str:", *str)
		}
	}
	formatUrl(&urlProxy)
	fmt.Println("Server running on port: :", port)
	fmt.Println("Redirecting to:", urlProxy)
}

// func test(){
// 	err := bcrypt.CompareHashAndPassword(password, []byte(request.Password))
// 	if err != nil {
// 		fmt.Println("bad pw")
// 		writeBadRequest(); return	// request-passord doesnt match stored hash
// 	}
// }
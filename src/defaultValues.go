package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// globals and default values
var (
	// variables changeable by user
	port 		= "3002"
	urlProxy 	= "http://127.0.0.1:3001"
	username 	= "username"
	password []byte 
	filePath 	= "/var/www/goAuthProxy"
	secret		= "default-2_#123default-2_#123"
	
	// variable not set by user:
	storage 	= NewAccountStorage()
	failedLoginAttempts int32 = 0
	sessionTime	= time.Hour*8
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
	rPath 	:= os.Getenv("GRP_FILEPATH")


	if rPort !=""{
		port = rPort
	} 	
	if rUrl !=""{
		urlProxy  = rUrl
	} 	
	if rUser !=""{
		username = rUser
	} 	
	if rPwHash !=""{
		pwHash = rPwHash
	} 	
	if rSecret !=""{
		secret = rSecret
	} 	
	if rPath !=""{
		filePath = rPath
	} 	
	
	
	// read files from flags, if they exist, take prio over env-variables
	fPort 	:= flag.String("port", "", "the port this proxy should listen for requets from")
	fUrl 	:= flag.String("url", "", "the address (and port) of the server we proxy to")	
	fUser 	:= flag.String("user", "", "username of the user we login with")
	fPw 	:= flag.String("pw", "", "plaintext password")
	fSecret	:= flag.String("secret", "", "the secret used to encrypt the JWT-Tokens")
	fPwHash	:= flag.String("pwhash", "", "hashed pw")
	fPath 	:= flag.String("files", "", "absolute path to html/js files, ex /var/www/goAuthProxy")
	flag.Parse()

	if *fPort != ""{
		port = *fPort
	}
	if *fUrl != ""{
		urlProxy = *fUrl
	}
	if *fUser != ""{
		username = *fUser
	}
	if *fPw != ""{
		pw = *fPw
	}
	if *fSecret != ""{
		secret = *fSecret
	}
	if *fPwHash != ""{
		pwHash = *fPwHash
	}
	if *fPath != ""{
		filePath = *fPath
	}
	

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
	fmt.Printf("!- Better use -pwhash  or set the env GRP_PASSWORD_HASH directly, enclose in '' in terminal: \n %v \n", string(encPw))
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
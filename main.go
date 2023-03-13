package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
)

/*
* SETUP - initializing defaults and checking urls
 */
var (
	port 		= flag.String("port", "3002", "the port this proxy should listen for requets from")
	username 	= flag.String("user", "username", "the address and port of the server we proxy to")
	password 	= flag.String("pw", "qwert123", "the address and port of the server we proxy to")
	urlProxy 	= flag.String("url", "http://127.0.0.1:3001", "the address and port of the server we proxy to")	
)
const cookieToken string  = "123asdfasdoi123"		// :todo remove after jwt


func getEnvValues() {
	// initialise necessary values
	flag.Parse()
	if *username =="username" 	{fmt.Println("using default username: ", *username,"change with -user")}
	if *password=="qwert123" 	{fmt.Println("using default password:", *password,"change with -pw")}
	
	// hash our password
	encPw, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil{
		panic(err)
	}
	*password = string(encPw)

	//format url if without http://
	formatUrl := func(str *string){
		if !strings.HasPrefix(*str, "http://") && !strings.HasPrefix(*str, "https://"){
			*str = "http://" + *str
			log.Println("adding http:// str the URL -str:", *str)
		}
	}
	formatUrl(urlProxy)
	fmt.Println("Server running on port: :", *port)
	fmt.Println("Redirecting to:", *urlProxy)
}


func main(){
	// load our setup, :todo repalce this with env-values
	getEnvValues()

	// multiplex our routes:
	mux := http.NewServeMux()
	mux.Handle("/login/", http.StripPrefix("/login/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("/", handleRequestAndRedirect)
	mux.HandleFunc("/api", handleLoginRequest)
	mux.HandleFunc("/logout", handleLogoutRequest)

	err := http.ListenAndServe(":"+*port, mux)
	if err != nil{
		panic(err)
	}
}


/*
*	Handle Login/Logout
*/

type LoginRequest struct{
	Name		string `json:"name"`
	Password	string `json:"password"`
}

func handleLoginRequest(rw http.ResponseWriter, req *http.Request){
	writeBadRequest := func(){
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("400 - Bad Request"))
	}

	if req.Method != "POST"{
		writeBadRequest(); return
	} 

	var request LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil{
		writeBadRequest(); return
	}
	if request.Name!= *username{
		writeBadRequest(); return
	}
	// check request-password against our hash
	err := bcrypt.CompareHashAndPassword([]byte(*password), []byte(request.Password))
	if err != nil {
		writeBadRequest(); return
	}
	// valid login -> grant our cookie -> then redirect to basepath if site
	log.Println("login sucess - cookie granted")
	addCookie(rw, "LoginToken", cookieToken, 2*time.Minute)
	http.Redirect(rw, req, "/", http.StatusSeeOther)
}

// we use the cookie to store our credibility into, ;todo repalce with JWT
func addCookie(rw http.ResponseWriter, name, value string, duration time.Duration){
	expire := time.Now().Add(duration)
	cookie := http.Cookie{
		Name: name,
		Value: value,
		Expires: expire,
		Path: "/",
	}
	http.SetCookie(rw, &cookie)
}

// logout just removes the cookie ( creating a  new one to overwrite old, setting negative time)
func handleLogoutRequest(rw http.ResponseWriter, req *http.Request){
	http.SetCookie(rw, &http.Cookie{
		Name: "LoginToken",
		Expires: time.Now().Add(-time.Hour),
	})
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Logout Sucessful"))
}


/*
*	redirect logic depending if were logged in or not
*/

// redirect requests to the appropriate url vs proxy
func handleRequestAndRedirect(rw http.ResponseWriter, req *http.Request) {
	validToken := false
	cookie, err := req.Cookie("LoginToken")
	if err != nil {
		validToken = false
	} else if cookie.Value==cookieToken{
		validToken = true
	}

	if !validToken {
		// not logged-in so we redirect to login page
		http.Redirect(rw, req, "/login", http.StatusSeeOther)
	} else {
		// logged-in so we proxy forward to our proxy server
		url := *urlProxy
		log.Println("we serve proxy to:",url)

		serveReverseProxy(url, rw, req)
	}
}


/*
*	the actual reverse proxy
*/
func serveReverseProxy(to string, rw http.ResponseWriter, req *http.Request){
	url, err := url.Parse(to)
	if err != nil{
		panic(err)
	}
	
	proxy := httputil.NewSingleHostReverseProxy(url)

	// update headers to allow for ssl redirection
	//req.URL.Host = url.Host
	//req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	//req.Host = url.Host

	proxy.ServeHTTP(rw, req)
}
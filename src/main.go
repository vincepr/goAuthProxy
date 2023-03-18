package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/bcrypt"
)

/*
* SETUP - initializing defaults and checking urls
 */


func main(){
	// initialize and check our values we run the proxy with
	InitGlobalValues()
	InitAccounts()
	const path = "/var/www/goAuthProxy/"

	// multiplex our routes:
	mux := http.NewServeMux()
	mux.Handle("/login/", http.StripPrefix("/login/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("/", handleRequestAndRedirect)
	mux.HandleFunc("/api", handleLoginRequest)
	mux.HandleFunc("/logout", handleLogoutRequest)

	err := http.ListenAndServe(":"+port, mux)
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

// rewrite this so it blocks for ALL THREADS AND NOT JUST THE ONE (port?) its on. :todo (might have race conditions like this)
func checkLoginAttempts(){
	//var1 := atomic.LoadInt32(&failedLoginAttempts)
	atomic.AddInt32(&failedLoginAttempts , 1)
	if atomic.LoadInt32(&failedLoginAttempts) > 4{
		time.Sleep(10*time.Second)
		failedLoginAttempts=0
		atomic.StoreInt32(&failedLoginAttempts, 0)
	}
}

func handleLoginRequest(rw http.ResponseWriter, req *http.Request){
	checkLoginAttempts()
	writeBadRequest := func(){
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("400 - Bad Request"))
	}
	if req.Method != "POST"{
		writeBadRequest(); 
		return	// not allowed method
	}
	var request LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil{
		writeBadRequest(); 
		return	// no json body
	}
	findUser, err := storage.GetAccountByName(request.Name)
	if err != nil{
		writeBadRequest(); 
		return	// user not found
	}
	err = bcrypt.CompareHashAndPassword(findUser.PasswordHash, []byte(request.Password))
	if err != nil {
		writeBadRequest(); 
		return	// request-passord doesnt match stored hash
	}
	// valid login -> grant our cookie -> then redirect to basepath if site
	atomic.StoreInt32(&failedLoginAttempts, 0)
	token, err := CreateJWTToken(findUser.Name, findUser.IsAdmin, secret, sessionTime)
	addCookie(rw, "LoginToken", token, sessionTime)
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
	http.Redirect(rw, req, "/login/", http.StatusSeeOther)
}


/*
*	redirect logic depending if were LOGGED-IN or NOT
*/

// redirect requests to the appropriate url vs proxy
func handleRequestAndRedirect(rw http.ResponseWriter, req *http.Request) {
	// :todo maybe rewrite this ugly auth block as middleware without nested ifs, to seperate concerns better. But then again its running all the same.
	validToken := false
	cookie, err := req.Cookie("LoginToken")
	if err != nil {
		validToken = false
	}else{
		// parse claims out of cookie
		claims, err := ValidateJWTClaims(cookie.Value, secret)					
		if (err != nil) {
			validToken = false
		} else{
			// check if user exists
			name, err := storage.GetAccountByName(claims.Name)
			if err != nil{
				validToken = false
			}else if claims.Name == name.Name{
				validToken = true
			}
		}
	}
	

	if !validToken {
		// not logged-in so we redirect to login page
		http.Redirect(rw, req, "/login", http.StatusSeeOther)
	} else {
		// logged-in so we proxy forward to our proxy server
		url := urlProxy
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
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	

	proxy.ServeHTTP(rw, req)
}
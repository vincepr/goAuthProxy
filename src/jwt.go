package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


/*
* 	JSON Web Token to secure our users identity and use it for authorisation
*
* - https://jwt.io/introduction 			- overview over json web token
* - package from github.com/golang-jwt/jwt	- golang package for jsw
*/


// Claims from a Token, stores who the user is, what he can access and or and for how long 
type JWTClaims struct {
	Name string `json:"name"`
	IsAdmin bool `json:"isAdmin"`
	jwt.RegisteredClaims
}

func NewJWTClaims(name string, isAdmin  bool, duration time.Duration) JWTClaims{
	return JWTClaims{
		Name: name,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-reverse-proxy",
		},
	}
}

// creates a Token to pass to our Users after ex. Login
func CreateJWTToken(name string, isAdmin bool, secret string, duration time.Duration) (string, error){
	mySigningKey := []byte(secret)
	claims := NewJWTClaims(name, isAdmin, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// validation happens here, returns our claims
func ValidateJWTClaims(tokenString string, secret string) (*JWTClaims, error){
	mySigningKey := []byte(secret)
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {	// Validate the encrypt-Algorythm is the one what we expect 
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil
	})
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}



// keeping the following here if i ever wanna replace the handleRequestAndRedirect() with some prettier middleware

/*
func middlewareJWTAuth(handlerFunc http.HandlerFunc, storage Storage) http.HandlerFunc{
	return func(rw http.ResponseWriter, req *http.Request){
		// Default Error msg, so no info about if a account exists can be gathered
		writeJSONError := func(){
			WriteJSON(rw, http.StatusForbidden, ApiError{Error: "invalid token"})
		}
		// check if there is ANY valid token:
		tokenString :=req.Header.Get("x-jwt-token")
		claims, err := validateJWTClaims(tokenString)					
		if (err != nil) {
			writeJSONError()
			return 
		}

		// :todo req -> userName 

		// grab that nr's data from the database
		account, err := storage.GetAccountByName(userName)
		if (err != nil) {
			writeJSONError()
			return 
		}
		// check if the claims of the token fit the user-> user accessing his own data
		if account.Name !=  claims.Name{
			writeJSONError()
			return 
		}

		handlerFunc(rw, req)
		
	}
}

// helper function:
func WriteJSON(rw http.ResponseWriter, status int, val any) error{
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(val)
}

type apiFunction func(http.ResponseWriter, *http.Request) error

type ApiError struct{
	Error string `json:"error"`
}

*/
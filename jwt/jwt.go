package jwt

import (
	"encoding/json"
	"net/http"
	"time"
)


/*
* 	JSON Web Token to 
*/

// Middleware for Auth: using Jason-Web-Token-standard - https://jwt.io/introduction
// jwt package from go get -u github.com/golang-jwt/jwt/v5

type apiFunction func(http.ResponseWriter, *http.Request) error

type ApiError struct{
	Error string `json:"error"`
}

func WriteJSON(header http.ResponseWriter, status int, val any) error{
	header.Header().Set("Content-Type", "application/json")
	header.WriteHeader(status)
	return json.NewEncoder(header).Encode(val)
}

func withJWTAuth(handlerFunc http.HandlerFunc, storage Storage) http.HandlerFunc{
	return func(header http.ResponseWriter, r *http.Request){
		// check if there is ANY valid token:
		tokenString :=r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)					
		if (err != nil || !token.Valid) {
			WriteJSON(header, http.StatusForbidden, ApiError{Error: "invalid token"})
			return 
		}
		// identify user-nr that is beeing acessed
		userId, err := paramsToId(r)
		if (err != nil || !token.Valid) {
			WriteJSON(header, http.StatusForbidden, ApiError{Error: "invalid token"})
			return 
		}
		// grab that nr's data from the database
		account, err := storage.GetAccountById(userId)
		if (err != nil || !token.Valid) {
			WriteJSON(header, http.StatusForbidden, ApiError{Error: "invalid token"})
			return 
		}
		// check if the claims of the token fit the user-> user accessing his own data
		claims := token.Claims.(jwt.MapClaims)
		claimedNr := int64(claims["accountNumber"].(float64))	// comes out float64 out... 
		//... of the interface->cast it as int with float64 type assertion :todo rewrite with jwt map
		if account.Number !=  claimedNr{
			WriteJSON(header, http.StatusForbidden, ApiError{Error: "invalid token"})
			return 
		}

		handlerFunc(header, r)
	}
}

// validation happens here
func validateJWT(tokenString string)(*jwt.Token, error){
	//secret := os.Getenv("JWT_SECRET")	// in terminal for testing  $ export JWT_SECRET=qwert123
	secret := "SecretGoesBrrrrr"
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

// creates a individual token to validate account users
func createJwtToken(account *Account)(string, error){
	mySigningKey := []byte("SecretGoesBrrrrr")

	// Create the Claims
	claims := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Unix(1516239022, 0)),
		"accountNumber": account.Number,
		"issuer":    "gobank",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
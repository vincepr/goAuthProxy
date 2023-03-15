# Reverse Proxy with Login
- minimal Example of a Reverse Proxy, written in go. Using JSON-Web-Token for Authorisation and a username/password for Authentification.
- uses no Database, with a small RAM footprint
- do not use in production

## To test,
- either use the prebuild binary or run/build for your own system

```
-port string
    the port this proxy should listen for requets from
-pw string
    plaintext password
-pwhash string
    hashed pw
-secret string
    the secret used to encrypt the JSW-Tokens
-url string
    the address and port of the server we proxy to
-user string
    the address and port of the server we proxy to
```


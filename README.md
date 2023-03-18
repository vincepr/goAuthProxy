# Reverse Proxy with Login
- minimal Example of a Reverse Proxy, written in go. Using JSON-Web-Token for Authorisation and a username/password for Authentification.
- uses no Database, with a small RAM footprint
- do not use in production

## Setup using just a binary
lets assume we want to build our own binary for our current system (assuming we have the go-compiler):
- clone the github repository and cd into the src folder
- `go build -o ./goAuthProxy .` to create a binary for your current system
- `mkdir -p /var/www/goAuthProxy` and `cp -r ./public/. /var/www/goAuthProxy` Create a folder, then copy our login-page to there
- now we can test run our binary `./goAuthProxy --files /var/www/goAuthProxy/ -pw testpassword` 

### Creating a systemd autostart (for ubuntu or other systemd based systems)
- `cp ./goAuthProxy /usr/bin/goAuthProxy` place our binary from the previous step into 
- `sudo nano /usr/lib/systemd/system/goAuthProxy.service` create the following file for our systemd service:
```
[Unit]
Description=golang-Authentification-Reverseproxy
After=nginx.service

[Service]
Type=simple
ExecStart=/usr/bin/goAuthProxy --url 127.0.0.1:3001
Restart=always
User=vinceprNotRoot

[Install]
WantedBy=multi-user.target
```
- `sudo systemctl enable code-server` to enable our service
- `sudo systemctl start code-server` to start it
- `sudo systemctl daemon-reload` reload all systemd config, to make it active before a restart

### List of all flags available:
```
-files 
    absolute path to html/js files, ex /var/www/goAuthProxy
-port 
    the port this proxy should listen for requets from
-pw 
    plaintext password
-pwhash 
    hashed pw
-secret 
   the secret used to encrypt the JWT-Tokens
-url 
    the address and port of the server we proxy to
-user 
    username of the user we login with
```
### List of env Variables available:
with some example values
```
GRP_PORT=8080
GRP_URL='https://adminarea.vincepr.de'
GRP_USER=Bobby
GRP_PASSWORD_HASH='$2a$10$3Th7F6Cd4rpEz8dIIh/dJO8gSnO5rKxis81OQ3ozEWYNLW7T7/MGe'
GRP_SECRET='my-32-character-ultra-secure-and-ultra-long-secret'
GRP_FILEPATH=/var/www/goAuthProxy
```

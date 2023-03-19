# Reverse Proxy with Login
-Example of a Reverse Proxy, written in go. Using JSON-Web-Token for Authorisation and a username/password-hash for Authentification.
- uses no Database, with a small RAM footprint (3-4Mb)


## Setup using just a binary
lets assume we want to build our own binary for our current system (assuming we have the go-compiler):
- clone the github repository and cd into the src folder
- `mkdir -p /var/www/goAuthProxy` and `cp -r ./public/. /var/www/goAuthProxy` Create a folder, then copy our login-page to there
- `go build -o ./goAuthProxy .` to create a binary for your current system
- now we can test run our binary `./goAuthProxy --files /var/www/goAuthProxy/ -pw test50password` 

### Creating a systemd autostart (for ubuntu or other systemd based systems)
- `cp ./goAuthProxy /usr/bin/goAuthProxy` place our binary from the previous step
- `sudo nano /usr/lib/systemd/system/goAuthProxy.service` create the following file for our systemd service:
```
[Unit]
Description=golang-Authentification-Reverseproxy
After=nginx.service

[Service]
Type=simple
Environment=GRP_SECRET='my-32-character-ultra-secure-and-ultra-long-secret'
ExecStart=/usr/bin/goAuthProxy --pwhash '$$2a$$10$$phK/STbCWKZkpIIFyMMo8OQT56MEObJiGqRqqpR7noItLjogsLNsu'
Restart=always
User=vinceprNotRoot

[Install]
WantedBy=multi-user.target
```
**Important note:**   since this is essentially bash $ signs will have to be replaced by $$ for our pwhash to work in here. we could also use any of the 3 '`" by escaping it with \' if needed
- `sudo systemctl enable goAuthProxy` to enable our service
- `sudo systemctl start goAuthProxy` to start it
- `sudo systemctl daemon-reload` reload all systemd config, to make it active before a restart

sudo systemctl enable goAuthProxy && sudo systemctl start goAuthProxy && sudo systemctl daemon-reload
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

## setup with docker run
- sudo docker build -t go_auth_proxy .
- about 8mb, so same size as binary
- sudo docker run -it --rm -p 5555:5555 go_auth_proxy 
- now we can access it in localhost:5555

If we would want to expose it in our vps we could do so with--network="host"
- sudo docker run -it --rm  --network="host" -p 5555:5555 go_auth_proxy
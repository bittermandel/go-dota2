# go-dota
Basic library for Dota using `go-steam`

## How to use
```bash
$ export STEAM_USERNAME=USERNAME
$ export STEAM_PASSWORD=PASSWORD
# If you don't use 2FA, leave AUTHCODE empty and export the code from the email, then try again
# Only gotta do this once
$ export STEAM_AUTHCODE=AUTHCODE
$ go test
```
## go-mcv
CUI multi comment viewer made with golang.<br><br>
Currently, only twitch and youtube are supported.

### usage
```shell
# build
go build .

# run
## use config file
./go-mcv run --config .go-mcv.yaml
## or, pass arguments
./go-mcv run --twitch-channel-id foobar --youtube-channel-id hogefuga
```

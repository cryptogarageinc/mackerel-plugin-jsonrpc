# mackerel-plugin-jsonrpc
mackerel plugin for count JSON-RPC response

# How to use
ex) `$ mackerel-plugin-jsonrpc -user user -password password auth -url http://127.0.0.1 -methodname examplemethod -arg '[0, 1, [], true, "sample"]'`

argオプションに指定する文字列はJSON arrayとして解釈できるものであること。

```
  -arg string
    	method argument(JSON array string)
  -connect_timeout duration
    	Maximum wait for connection, in seconds. (default 10s)
  -label string
    	metrics label
  -methodname string
    	methodname
  -metric-key-prefix string
    	Metric key prefix (default "jsonRPC")
  -password string
    	JSON-RPC password
  -tempfile string
    	Temp file name
  -url string
    	JSON-RPC URL (default "http://127.0.0.1")
  -user string
    	JSON-RPC user
```

# How to release for mkr install
1. `$ make setup`
1. `$ git tag v0.19.1` (タグ名は適宜置き換えること)
1. `$ GITHUB_TOKEN=... script/release.sh` (GITHUB_TOKENはあらかじめ発行しておくこと)

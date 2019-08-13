# settlenet-mackerel-pgsql
mackerel plugin for executing sql on postgres

# How to use
ex) `$ settlenet-mackerel-pgsql -user auth -database auth -hostname 127.0.0.1 -password password -port 5432 -sqlconfig ./sqlconfig.toml `

sqlconfig オプションに指定するファイルの形式は次節参照。

```
  -connect_timeout int
    	Maximum wait for connection, in seconds. (default 5)
  -database string
    	Database name
  -hostname string
    	Hostname to login to (default "localhost")
  -metric-key-prefix string
    	Metric key prefix (default "postgres")
  -password string
    	Postgres Password
  -port string
    	Database port (default "5432")
  -sqlconfig string
    	Sql config file
  -sslmode string
    	Whether or not to use SSL (default "disable")
  -tempfile string
    	Temp file name
  -user string
    	Postgres User
```

# Sqlconfig toml
```
[[sqlconfig]]
key = "fee"
label = "ABC Amount"
unit = "integer"
metricsname = "amount"
metricslabel = "Amount"
sql = "SELECT SUM(amount) FROM ABC WHERE status = 'statusA';"
```

# How to release for mkr install
1. `$ make setup`
1. `$ git tag v0.19.1` (タグ名は適宜置き換えること)
1. `$ GITHUB_TOKEN=... script/release.sh` (GITHUB_TOKENはあらかじめ発行しておくこと)

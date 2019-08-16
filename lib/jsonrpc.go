package mpjsonrpc

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	// PostgreSQL Driver
	_ "github.com/lib/pq"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("metrics.plugin.jsonRPC")

// JSONRPCPlugin mackerel plugin for PostgreSQL
type JSONRPCPlugin struct {
	URL        string
	Username   string
	Password   string
	Prefix     string
	Timeout    time.Duration
	Tempfile   string
	Option     string
	MethodName string
	Label      string
	Arg        []interface{}
}

//JSONRPCConfig は、JSONRPC接続設定情報
type JSONRPCConfig struct {
	URL      string
	User     string
	Password string
	Timeout  time.Duration
}

// JSONRPCClient は JSONRPC endnode とのコミュニケーションを担当する
type JSONRPCClient struct {
	client *http.Client
	config *JSONRPCConfig
}

// Request は JSONRPC のリクエストを表す
type Request struct {
	Jsonrpc string        `json:"jsonrpc,"`
	ID      string        `json:"id,"`
	Method  string        `json:"method,"`
	Params  []interface{} `json:"params,"`
}

// Response は JSONRPC のリクエストを表す
type Response struct {
	Result interface{}            `json:"result,"`
	Error  map[string]interface{} `json:"error,"`
	ID     string                 `json:"id,"`
}

// NewJSONRPCClient は新しい JSONRPCClient を返す
func NewJSONRPCClient(c *JSONRPCConfig) JSONRPCClient {
	client := &http.Client{
		Timeout: c.Timeout,
	}
	return JSONRPCClient{
		client: client,
		config: c,
	}
}

// Request は対応するendnodeに対してJSONRPCリクエストを投げる
func (c *JSONRPCClient) Request(req Request) (interface{}, error) {
	bodyJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}
	httpReq.SetBasicAuth(c.config.User, c.config.Password)

	res, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBuf := Response{}
	if err := json.NewDecoder(res.Body).Decode(&resBuf); err != nil {
		return nil, err
	}
	if resBuf.Error != nil {
		return nil, fmt.Errorf(fmt.Sprintf("jsonRPC endnode response with error: %v", resBuf.Error))
	}

	return resBuf.Result, nil
}

// NewRequest は新しい Request 構造体を返す
func NewRequest(method string, params ...interface{}) Request {
	return Request{
		Jsonrpc: "1.0",
		ID:      uuid.New().String(),
		Method:  method,
		Params:  params,
	}
}

// MetricKeyPrefix returns the metrics key prefix
func (p JSONRPCPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "jsonRPC"
	}
	return p.Prefix
}

// FetchMetrics interface for mackerelplugin
func (p JSONRPCPlugin) FetchMetrics() (map[string]interface{}, error) {
	req := NewRequest(p.MethodName, p.Arg...)

	jsonRPCClient := NewJSONRPCClient(
		&JSONRPCConfig{
			URL:      p.URL,
			User:     p.Username,
			Password: p.Password,
			Timeout:  p.Timeout,
		},
	)
	rawRes, err := jsonRPCClient.Request(req)
	count := 0.0
	if err != nil {
		logger.Errorf("Failed to call json-rpc. %s", err)
	} else {
		count = float64(len(rawRes.([]interface{})))
	}

	stat := make(map[string]interface{})
	stat["count"] = count

	return stat, err
}

// GraphDefinition interface for mackerelplugin
func (p JSONRPCPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.MetricKeyPrefix())

	var graphdef = map[string]mp.Graphs{
		p.Label: {
			Label: (labelPrefix + " Count"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "count", Label: "Count"},
			},
		},
	}

	return graphdef
}

// Do the plugin
func Do() {
	optURL := flag.String("url", "http://127.0.0.1", "JSON-RPC URL")
	optUser := flag.String("user", "", "JSON-RPC user")
	optPass := flag.String("password", "", "JSON-RPC password")
	optPrefix := flag.String("metric-key-prefix", "jsonRPC", "Metric key prefix")
	optConnectTimeout := flag.Duration("connect_timeout", 10*time.Second, "Maximum wait for connection, in seconds.")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	optMethodName := flag.String("methodname", "", "methodname")
	optLabel := flag.String("label", "", "metrics label")
	optArg := flag.String("arg", "", "method argument(JSON array string)")
	flag.Parse()

	if *optUser == "" {
		logger.Warningf("user is required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *optMethodName == "" {
		logger.Warningf("methodname is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var jsonRPCPlugin JSONRPCPlugin
	jsonRPCPlugin.URL = *optURL
	jsonRPCPlugin.Username = *optUser
	jsonRPCPlugin.Password = *optPass
	jsonRPCPlugin.Prefix = *optPrefix
	jsonRPCPlugin.Timeout = *optConnectTimeout
	jsonRPCPlugin.MethodName = *optMethodName
	jsonRPCPlugin.Label = *optLabel
	arg := []interface{}{}
	err := json.Unmarshal([]byte(*optArg), &arg)
	if err != nil {
		logger.Warningf("arg cannot be interpreted as JSON array")
		flag.PrintDefaults()
		os.Exit(1)
	}
	jsonRPCPlugin.Arg = arg

	helper := mp.NewMackerelPlugin(jsonRPCPlugin)

	helper.Tempfile = *optTempfile
	helper.Run()
}

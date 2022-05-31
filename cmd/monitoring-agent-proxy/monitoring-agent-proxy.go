package main

import (
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/thanos-io/thanos/pkg/logging"
	thanos_tls "github.com/thanos-io/thanos/pkg/tls"

	monitoringagentproxy "github.com/kubesphere/paodin/pkg/monitoring-agent-proxy"
)

var cli struct {
	Log struct {
		Level  string `enum:"debug,info,warn,error" default:"info" help:"Log filtering level. Possible options: ${enum}."`
		Format string `enum:"json,logfmt" default:"logfmt" help:"Log format to use. Possible options: ${enum}."`
	} `embed:"" prefix:"log."`

	HttpAddress       string `default:"0.0.0.0:9090" help:"Listen host:port for HTTP endpoints."`
	ServerTlsKey      string `default:"" help:"TLS Key for HTTP server, leave blank to disable TLS."`
	ServerTlsCert     string `default:"" help:"TLS Certificate for HTTP server, leave blank to disable TLS."`
	ServerTlsClientCa string `default:"" help:"TLS CA to verify clients against. If no client CA is specified, there is no client verification on server side. (tls.NoClientCert)"`

	MonitorGateway struct {
		Address string `default:"" help:"Address to connect paodin monitor-gateway"`
	} `embed:"" prefix:"monitorgaetway."`

	Tenant string `default:"" help:"unique. clusterID"`
}

func main() {

	ctx := kong.Parse(&cli)
	logger := logging.NewLogger(cli.Log.Level, cli.Log.Format, "")

	rawUrl, err := url.Parse(cli.MonitorGateway.Address)
	ctx.FatalIfErrorf(err)

	options := &monitoringagentproxy.Options{
		Tenant:               cli.Tenant,
		ListenAddress:        cli.HttpAddress,
		GatewayProxyEndpoint: rawUrl,
	}

	options.TLSConfig, err = thanos_tls.NewServerConfig(logger, cli.ServerTlsCert, cli.ServerTlsKey, cli.ServerTlsClientCa)
	ctx.FatalIfErrorf(err)

	server := monitoringagentproxy.NewServer(logger, options)
	err = server.Run()
	ctx.FatalIfErrorf(err)
}

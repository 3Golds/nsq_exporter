package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/timonwong/nsq_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	statsRegistry = map[string]func(namespace string) collector.StatsCollector{
		"topics":   collector.TopicStats,
		"channels": collector.ChannelStats,
		"clients":  collector.ClientStats,
	}
)

func init() {
	prometheus.MustRegister(version.NewCollector("nsq_exporter"))
}

var cfg struct {
	listenAddress     string
	metricsPath       string
	nsqdURL           string
	enabledCollectors string
	timeout           time.Duration
	namespace         string
}

func main() {
	kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").
		Default(":9118").
		StringVar(&cfg.listenAddress)
	kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").
		Default("/metrics").
		StringVar(&cfg.metricsPath)
	kingpin.Flag("nsqd.addr", "Address of the nsqd node.").
		Default("http://localhost:4151/stats").
		StringVar(&cfg.nsqdURL)
	kingpin.Flag("collect", "Comma-separated list of collectors to use.").
		Default("stats.topics,stats.channels").
		StringVar(&cfg.enabledCollectors)
	kingpin.Flag("timeout", "Timeout for trying to get stats from nsqd.").
		Default("5s").
		DurationVar(&cfg.timeout)
	kingpin.Flag("namespace", "Namespace for the NSQ metrics.").
		Default("nsq").
		StringVar(&cfg.namespace)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("nsq_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	ex, err := createNsqExecutor()
	if err != nil {
		log.Fatalf("error creating nsq executor: %v", err)
	}
	prometheus.MustRegister(ex)

	http.Handle(cfg.metricsPath, promhttp.Handler())
	if cfg.metricsPath != "" && cfg.metricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>NSQ Exporter</title></head>
			<body>
			<h1>NSQ Exporter</h1>
			<p><a href="` + cfg.metricsPath + `">Metrics</a></p>
            <h2>Build</h2>
            <pre>` + version.Info() + ` ` + version.BuildContext() + `</pre>
			</body>
			</html>`))
		})
	}

	log.Info("listening to ", cfg.listenAddress)
	err = http.ListenAndServe(cfg.listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createNsqExecutor() (*collector.NsqExecutor, error) {
	nsqdURL, err := normalizeURL(cfg.nsqdURL)
	if err != nil {
		return nil, err
	}

	ex := collector.NewNsqExecutor(cfg.namespace, nsqdURL, cfg.timeout)
	for _, param := range strings.Split(cfg.enabledCollectors, ",") {
		param = strings.TrimSpace(param)
		parts := strings.SplitN(param, ".", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid collector name: %s", param)
		}
		if parts[0] != "stats" {
			return nil, fmt.Errorf("invalid collector prefix: %s", parts[0])
		}

		name := parts[1]
		c, has := statsRegistry[name]
		if !has {
			return nil, fmt.Errorf("unknown stats collector: %s", name)
		}
		ex.Use(c(cfg.namespace))
	}
	return ex, nil
}

func normalizeURL(ustr string) (string, error) {
	ustr = strings.ToLower(ustr)
	if !strings.HasPrefix(ustr, "https://") && !strings.HasPrefix(ustr, "http://") {
		ustr = "http://" + ustr
	}

	u, err := url.Parse(ustr)
	if err != nil {
		return "", err
	}
	if u.Path == "" {
		u.Path = "/stats"
	}
	u.RawQuery = "format=json"
	return u.String(), nil
}

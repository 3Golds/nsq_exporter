package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/timonwong/nsq_exporter/collector"
)

var (
	showVersion       = flag.Bool("version", false, "Print version information.")
	listenAddress     = flag.String("web.listen-address", ":9118", "Address on which to expose metrics and web interface.")
	metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	nsqdURL           = flag.String("nsqd.addr", "http://localhost:4151/stats", "Address of the nsqd node.")
	enabledCollectors = flag.String("collect", "stats.topics,stats.channels", "Comma-separated list of collectors to use.")
	timeout           = flag.Duration("timeout", 5*time.Second, "Timeout for trying to get stats from nsqd.")
	namespace         = flag.String("namespace", "nsq", "Namespace for the NSQ metrics.")

	statsRegistry = map[string]func(namespace string) collector.StatsCollector{
		"topics":   collector.TopicStats,
		"channels": collector.ChannelStats,
		"clients":  collector.ClientStats,
	}
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stdout, version.Print("nsq_exporter"))
		os.Exit(0)
	}

	ex, err := createNsqExecutor()
	if err != nil {
		log.Fatalf("error creating nsq executor: %v", err)
	}
	prometheus.MustRegister(version.NewCollector("nsq_exporter"))
	prometheus.MustRegister(ex)

	http.Handle(*metricsPath, promhttp.Handler())
	if *metricsPath != "" && *metricsPath != "/" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>NSQ Exporter</title></head>
			<body>
			<h1>NSQ Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	log.Info("listening to ", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createNsqExecutor() (*collector.NsqExecutor, error) {
	nsqdURL, err := normalizeURL(*nsqdURL)
	if err != nil {
		return nil, err
	}

	ex := collector.NewNsqExecutor(*namespace, nsqdURL, *timeout)
	for _, param := range strings.Split(*enabledCollectors, ",") {
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
		ex.Use(c(*namespace))
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

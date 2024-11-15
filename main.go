package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"http_metric/config"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetDefault("interval", 15)
	viper.SetDefault("timeout", 10)
	viper.SetDefault("port", 8080)

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error reading config file, %s", err)
		return
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("Error reading config file, %s", err)
		return
	}

	requestResponse := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_response_status_code",
		Help: "HTTP response status code",
	}, []string{"url", "name"})

	registry := prometheus.NewRegistry()
	registry.MustRegister(requestResponse)

	go func() {
		client := &http.Client{Timeout: time.Duration(cfg.TimeOut) * time.Second}
		for {
			time.Sleep(time.Duration(cfg.Interval) * time.Second)
			for _, target := range cfg.Targets {
				req, err := http.NewRequest("GET", target.Url, nil)
				if err != nil {
					slog.Error("Error creating request, %s", err)
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					requestResponse.WithLabelValues(target.Url, target.Name).Set(0)
					slog.Error("Error performing request, %s", err)
					continue
				}
				defer resp.Body.Close()

				requestResponse.WithLabelValues(target.Url, target.Name).Set(float64(resp.StatusCode))
			}
		}
	}()

	slog.Info("Starting HTTP server on port: " + fmt.Sprintf("%d", cfg.Port))
	slog.Info("Metrics available at /metrics")

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		slog.Error("Error starting HTTP server, %s", err)
	}

}

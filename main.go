package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"http_metric/lib/config"
	"net/http"
	"time"
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetDefault("interval", 15)
	viper.SetDefault("timeout", 10)
	viper.SetDefault("port", 8080)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		return
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error unmarshalling config, %s\n", err)
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
					fmt.Printf("Error creating request, %s\n", err)
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error making request, %s\n", err)
					continue
				}
				defer resp.Body.Close()

				requestResponse.WithLabelValues(target.Url, target.Name).Set(float64(resp.StatusCode))
			}
		}
	}()

	fmt.Printf("Starting HTTP server on port %d\n", cfg.Port)

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		fmt.Printf("Error starting HTTP server, %s\n", err)
	}

}

# Introduction

When you have multiple HTTP services scattered around, you can use this project to collect the Status Code of HTTP
Response and collect it through Prometheus.

# How to use

1. Create a `config.yaml` file in the root directory of the project.

```yaml
port: 8080 # Port to run the server
interval: 15 # Interval to check the status code
targets: # List of targets
  - name: localhost
    url: http://localhost:8080
``` 

2. follow this command to run container.

```bash
docker run -it \
  -v ./config.yaml:/app/config/config.yaml \
  -p 8080:8080 \
  nantou/http-metrics:latest
```

3. Access the Prometheus server at `http://localhost:8080/metrics` and you will see the metrics.

```bash
curl http://localhost:8080/metrics
```


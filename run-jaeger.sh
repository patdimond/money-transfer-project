#!/bin/bash

# Ports....
# 16686s is the UI -> http://localhost:16686
# 4317 is the grpc otlp receiver
# 4318 is the http otlp receiver

docker run --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest

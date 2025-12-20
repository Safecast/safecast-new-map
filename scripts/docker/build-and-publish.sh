#!/usr/bin/env bash
docker login
docker build -t Safecast/safecast-new-map:latest .
docker push Safecast/safecast-new-map:latest

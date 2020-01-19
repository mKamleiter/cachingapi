#!/bin/bash
docker kill cachingapi
docker build -t cachingapi .
docker run -d -p 8443:8443 --name cachingapi --rm cachingapi
docker logs -f cachingapi

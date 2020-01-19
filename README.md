# cachingapi
This is a simple API which caches infromation for Serverinstallation.


Examples:

## Get all Servers
curl -u foo:bar https://localhost:8443/v1/server 

## Get Servers since timestamp
curl -u foo:bar "https://localhost:8443/v1/server?since=2020-01-19T16:23"

## Insert Server
curl https://localhost:8443/v1/server -XPOST -d '{"name": "foo","comments":"foobar,rhel"}' -u foo:bar

# Build and Run
To run the API just use the buildandrun.sh
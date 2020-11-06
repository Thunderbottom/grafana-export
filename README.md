# Grafana Export

A small utility to download all dashboards from your Grafana Instance using the Grafana HTTP API.

## Usage

The utility requires the Grafana instance URL, and an API token to access the API:

```shell
$ GRAFANA_URL=http://your-grafana-instance GRAFANA_API_TOKEN=token ./grafana-export
```

By default, all dashboards will be stored in the `dashboards/` directory. You can override this path
by passing the `GRAFANA_DASHBOARD_DIR` variable to the utility. The utility also accepts `GRAFANA_API_LIMIT`
as a parameter, in case you would like to override the default API search limit of 1000.

## FAQ

**Q.** Why is the utility not downloading all the dashboards?

**A.** Your instance could have role-based dashboard access limits. The default API token generated uses
the "Viewer" role. You need to generate an API key with the "Admin" role for the API token to be able to
download all the dashboards.

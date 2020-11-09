<a href="https://zerodha.tech"><img src="https://zerodha.tech/static/images/github-badge.svg" align="right"/></a>

# Grafana Export

A small utility to download and backup all dashboards from your Grafana Instance using the Grafana HTTP API.

## Installation

##### Precompiled Binary

To download and use precompiled binaries for GNU/Linux, MacOS, and Windows, head over to the
[releases page](https://github.com/thunderbottom/grafana-export/releases).

##### Docker Installation

To locally build and run the docker image:

```shell
$ docker build -t grafana-export -f docker/Dockerfile .
$ docker run -v $(pwd)/dashboards:/dashboards \
	-e GEXPORT_URL=<grafana-url> \
	-e GEXPORT_API_KEY=<api-key> \
	-e GEXPORT_OVERWRITE=true \
	grafana-export
```

##### Compiling the Binary

```shell
$ git clone git@github.com:thunderbottom/grafana-export.git
$ cd grafana-export
$ make dist # or just `make` for a dynamically-linked binary
```

## Usage

The utility requires the Grafana instance URL and an API token:

##### Command Line Arguments

```shell
$ ./grafana-export --help
Usage of grafana-export:
  -api-key string
    	The API key to access the Grafana Instance.
  -backup
    	Backup the current set of downloaded dashboards to S3. Requires --bucket-name.
  -bucket-key string
    	The key to use for storing the backup inside the S3 bucket. (default "grafana-export")
  -bucket-name string
    	Compress and upload the dashboards to the specified S3 bucket. Requires --backup.
  -compress
    	Create an archive of the exported dashboards folder.
  -dashboards-dir string
    	The directory where the Grafana dashboards are to be downloaded. (default "dashboards")
  -limit int
    	The limit for number of results returned by the Grafana API. (default 1000)
  -overwrite
    	Overwrite existing dashboards directory.
  -url string
    	The base URL for the Grafana instance.
```

To create an S3 backup for the downloaded dashboards, the utility requires `awscli` to be configured on
the host. For setting a bucket region, export `AWS_REGION=<region>` or `AWS_SDK_LOAD_CONFIG` to load the
`~/.aws/config` file.

##### Environment Variables

The utility accepts environment variables as arguments. All environment variables prefixed with `GEXPORT_`
will be picked up by the utility, and will override its corresponding flag value passed to the utility.
For example, you may use `GEXPORT_URL` instead of passing the `--url` flag to the utility.

**Note:** Environment variables supersede the values passed as flags to the utility.

## FAQ

**Q.** Why is the utility not downloading all the dashboards?

**A.** Your instance could have role-based dashboard access limits. The default API token generated uses
the "Viewer" role. You need to generate an API key with the "Admin" role for the API token to be able to
download all the dashboards.

## License

```
Copyright (c) 2020 Chinmay D. Pai

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
```

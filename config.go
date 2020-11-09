package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/basicflag"
	"github.com/knadh/koanf/providers/env"
)

// getConfig parses the flags passed as arguments to the utility
// and fetches environment variables prefixed with `GEXPORT_` as config
func getConfig() *koanf.Koanf {
	var k = koanf.New(".")
	f := flag.NewFlagSet("grafana-export", flag.ExitOnError)
	f.String("url", "", "The base URL for the Grafana instance.")
	f.String("api-key", "", "The API key to access the Grafana Instance.")
	f.String("bucket-name", "", "Compress and upload the dashboards to the specified S3 bucket. Requires --backup.")
	f.String("bucket-key", "grafana-export", "The key to use for storing the backup inside the S3 bucket.")
	f.String("dashboards-dir", "dashboards", "The directory where the Grafana dashboards are to be downloaded.")
	f.Int("limit", 1000, "The limit for number of results returned by the Grafana API.")
	f.Bool("overwrite", false, "Overwrite existing dashboards directory.")
	f.Bool("backup", false, "Backup the current set of downloaded dashboards to S3. Requires --bucket-name.")
	f.Bool("compress", false, "Create an archive of the exported dashboards folder.")

	f.Parse(os.Args[1:])

	if err := k.Load(basicflag.Provider(f, "."), nil); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	k.Load(env.Provider("GEXPORT_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "GEXPORT_")), "_", "-", -1)
	}), nil)

	return k
}

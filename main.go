package main

import (
	"encoding/json"
	"fmt"
	"github.com/knadh/koanf"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// getGrafanaData calls the grafana API endpoint for the provided GRAFANA_URL
func getGrafanaData(cfg *koanf.Koanf, endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", cfg.String("url")+"/api/"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.String("api-key")))
	req.Header.Set("Content-Type", "application/json")
	params := req.URL.Query()
	params.Add("limit", strconv.Itoa(cfg.Int("limit")))
	req.URL.RawQuery = params.Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Error] %s", response.Status)
	}

	read, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return read, nil
}

// syncDashboards replicates the grafana folder structure, downloads all
// dashbords using the grafana api, and places them in each folder
func syncDashboards(cfg *koanf.Koanf, dashboards dashboardSearch) {
	dashDir := cfg.String("dashboards-dir")
	if _, err := os.Stat(dashDir); !os.IsNotExist(err) {
		if cfg.Bool("overwrite") {
			os.RemoveAll(dashDir)
		} else {
			log.Fatal(`Dashboards directory already exists. Pass --overwrite to overwrite the directory.`)
		}
	}
	err := os.Mkdir(dashDir, 0755)
	check(err)

	var failed, total int
	log.Println("Syncing Dashboards...")
	for _, ds := range dashboards {
		if ds.Type == "dash-folder" {
			// check if a folder for exists, if not, create one
			if _, err := os.Stat(dashDir); os.IsNotExist(err) {
				err := os.Mkdir(dashDir, 0755)
				check(err)
			}
		} else {
			total = total + 1
			// get the dashboard json from the grafana api
			db, err := getGrafanaData(cfg, fmt.Sprintf("dashboards/%s", ds.URI))
			if err != nil {
				log.Println("Failed to fetch dashboard:", ds.Title)
				failed = failed + 1
				continue
			}
			var dashJSON map[string]interface{}
			err = json.Unmarshal(db, &dashJSON)
			if err != nil {
				log.Println("Failed to parse dashboard:", ds.Title)
				failed = failed + 1
				continue
			}
			// if no folder name is specified for the dashboard, save
			// to the General/ folder
			if ds.FolderTitle == "" {
				ds.FolderTitle = "General"
			}
			dbFolder := fmt.Sprintf("%s/%s", dashDir, ds.FolderTitle)
			if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
				err := os.MkdirAll(dbFolder, 0755)
				check(err)
			}

			dj, err := json.MarshalIndent(dashJSON, "", "  ")
			if err != nil {
				log.Println("Failed to marshal the dashboard JSON:", ds.Title)
				failed = failed + 1
				continue
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", dbFolder, strings.Replace(ds.Title, "/", "-", -1)), dj, 0644)
			if err != nil {
				log.Println("Failed to save dashboard:", ds.Title)
				failed = failed + 1
				continue
			}
			fmt.Println(ds.Title, "downloaded.")
		}
	}
	fmt.Println(fmt.Sprintf("Done! Download Statistics:\n\tTotal: %d\n\tFailed: %d", total, failed))
}

func main() {
	cfg := getConfig()
	if cfg.String("url") == "" {
		log.Fatal("Missing required argument: --url")
	} else if cfg.String("api-key") == "" {
		log.Fatal("Missing required argument: --api-key")
	}

	resp, err := getGrafanaData(cfg, "search")
	check(err)

	ds := dashboardSearch{}
	err = json.Unmarshal(resp, &ds)
	check(err)

	syncDashboards(cfg, ds)

	var c string
	cmp := cfg.Bool("compress")
	if cmp {
		c, err = compress(cfg.String("dashboards-dir"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Dashboards compressed: %v", c)
	}

	// backup to s3 if --backup is passed
	if cfg.Bool("backup") {
		bucketName := cfg.String("bucket-name")
		bucketKey := cfg.String("bucket-key")
		if bucketName == "" {
			log.Fatal("No S3 bucket specified for backup.")
		}
		// check if --compress is passed
		// prevents re-compression
		if !cmp {
			c, err = compress(cfg.String("dashboards-dir"))
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Dashboards compressed: %v", c)
		}

		if err := backup(c, bucketName, bucketKey); err != nil {
			log.Fatalf("Failed to backup to S3 bucket: %v", err)
		}
		log.Printf("Uploaded %v to S3 Bucket %v/%v", c, bucketName, bucketKey)
	}
}

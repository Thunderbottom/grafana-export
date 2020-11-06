package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// dashboardDir is the directory where the dashboard will be downloaded
	// defaults to "dashboards"
	dashboardDir = os.Getenv("GRAFANA_DASHBOARDS_DIR")
	// dbFolder is the folder name for the currently downloading dashboard
	dbFolder = ""
	// grafanaURL is the base url for the grafana instance
	grafanaURL = os.Getenv("GRAFANA_URL")
	// grafanaToken is the api access key for the grafana instance
	grafanaToken = os.Getenv("GRAFANA_API_TOKEN")
	// grafanaLimit is a variable to set a limit on the number of dashboards to fetch
	grafanaLimit = os.Getenv("GRAFANA_API_LIMIT")
)

// dashboardSearch is a struct that the grafana dashboard search data
type dashboardSearch []struct {
	ID          int           `json:"id"`
	UID         string        `json:"uid"`
	Title       string        `json:"title"`
	URI         string        `json:"uri"`
	URL         string        `json:"url"`
	Slug        string        `json:"slug"`
	Type        string        `json:"type"`
	Tags        []interface{} `json:"tags"`
	IsStarred   bool          `json:"isStarred"`
	FolderID    int           `json:"folderId,omitempty"`
	FolderUID   string        `json:"folderUid,omitempty"`
	FolderTitle string        `json:"folderTitle,omitempty"`
	FolderURL   string        `json:"folderUrl,omitempty"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if grafanaURL == "" || grafanaToken == "" {
		fmt.Println("Please set GRAFANA_URL and GRAFANA_API_TOKEN before running the script.")
		os.Exit(1)
	}
	if dashboardDir == "" {
		dashboardDir = "dashboards"
	}
	grafanaURL = grafanaURL + "/api/"

	resp, err := grafanaAPI("search")
	check(err)

	ds := &dashboardSearch{}
	err = json.Unmarshal(resp, ds)
	check(err)

	syncDashboards(ds)
}

// grafanaAPI calls the grafana API endpoint for the provided GRAFANA_URL
func grafanaAPI(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", grafanaURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", grafanaToken))
	req.Header.Set("Content-Type", "application/json")
	if grafanaLimit != "" {
		if glInt, err := strconv.Atoi(grafanaLimit); err != nil && glInt > 0 {
			params := req.URL.Query()
			params.Add("limit", grafanaLimit)
			req.URL.RawQuery = params.Encode()
		}
	}

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
func syncDashboards(ds *dashboardSearch) {
	if _, err := os.Stat(dashboardDir); !os.IsNotExist(err) {
		os.RemoveAll(dashboardDir)
	}
	err := os.Mkdir(dashboardDir, 0755)
	check(err)

	var failed, total int
	fmt.Println("Syncing Dashboards...")
	for _, dashboard := range *ds {
		if dashboard.Type == "dash-folder" {
			// check if a folder for exists, if not, create one
			if _, err := os.Stat(dashboardDir); os.IsNotExist(err) {
				err = os.Mkdir(dashboardDir, 0755)
				check(err)
			}
		} else {
			total = total + 1
			// get the dashboard json from the grafana api
			db, err := grafanaAPI(fmt.Sprintf("dashboards/%s", dashboard.URI))
			if err != nil {
				fmt.Println("Failed to fetch dashboard:", dashboard.Title)
				failed = failed + 1
				continue
			}
			var dashboardJSON map[string]interface{}
			err = json.Unmarshal(db, &dashboardJSON)
			if err != nil {
				fmt.Println("Failed to parse dashboard:", dashboard.Title)
				failed = failed + 1
				continue
			}
			// if no folder name is specified for the dashboard, save
			// to the General/ folder
			if dashboard.FolderTitle == "" {
				dbFolder = fmt.Sprintf("%s/General", dashboardDir)
			} else {
				dbFolder = fmt.Sprintf("%s/%s", dashboardDir, dashboard.FolderTitle)
			}
			if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
				err = os.MkdirAll(dbFolder, 0755)
				check(err)
			}

			dbJSON, err := json.MarshalIndent(dashboardJSON, "", "  ")
			if err != nil {
				fmt.Println("Failed to marshal the dashboard JSON:", dashboard.Title)
				failed = failed + 1
				continue
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", dbFolder, strings.Replace(dashboard.Title, "/", "-", -1)), dbJSON, 0644)
			if err != nil {
				fmt.Println("Failed to save dashboard:", dashboard.Title)
				failed = failed + 1
				continue
			}
			fmt.Println(dashboard.Title, "downloaded.")
		}
	}
	dbFile, _ := json.MarshalIndent(ds, "", "  ")
	err = ioutil.WriteFile("dashboards.json", dbFile, 0644)
	fmt.Println(fmt.Sprintf("Done! Download Statistics:\n\tTotal: %d\n\tFailed: %d", total, failed))
}

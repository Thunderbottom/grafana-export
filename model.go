package main

// dashboardSearch is a struct that contains
// the search data for the grafana dashboard api
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

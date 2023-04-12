package azure

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	httpClient *http.Client
	azureOnce  sync.Once
)

func getClient() *http.Client {
	azureOnce.Do(func() {
		httpClient = &http.Client{}
	})
	return httpClient
}

type UserStory struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func CreateAzureNewUserStory(request []UserStory) (*WorkItem, error) {
	url := "https://dev.azure.com/winforevermore/win_project/_apis/wit/workitems/$user%20story?api-version=7.0"
	method := "POST"

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json-patch+json")
	req.Header.Add("Authorization", "Basic QmFzaWM6bm5hNGphNnNybWV2Mnk2eTI2azIyb280b2Y3NGFuYnN1N2Rqd3J4cmNsaXFjbGQzZzNqYQ==")

	client := getClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result WorkItem
	if err := json.Unmarshal(body, &result); err != nil {

		return nil, err
	}
	return &result, nil
}

type WorkItem struct {
	ID        int    `json:"id"`
	Rev       int    `json:"rev"`
	System    System `json:"fields"`
	Links     Links  `json:"_links"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
	Url       string `json:"url"`
}

type System struct {
	AreaPath         string    `json:"System.AreaPath"`
	TeamProject      string    `json:"System.TeamProject"`
	IterationPath    string    `json:"System.IterationPath"`
	WorkItemType     string    `json:"System.WorkItemType"`
	State            string    `json:"System.State"`
	Reason           string    `json:"System.Reason"`
	CreatedDate      time.Time `json:"System.CreatedDate"`
	CreatedBy        Identity  `json:"System.CreatedBy"`
	ChangedDate      time.Time `json:"System.ChangedDate"`
	ChangedBy        Identity  `json:"System.ChangedBy"`
	CommentCount     int       `json:"System.CommentCount"`
	Title            string    `json:"System.Title"`
	BoardColumn      string    `json:"System.BoardColumn"`
	BoardColumnDone  bool      `json:"System.BoardColumnDone"`
	StateChangeDate  time.Time `json:"Microsoft.VSTS.Common.StateChangeDate"`
	Priority         int       `json:"Microsoft.VSTS.Common.Priority"`
	ValueArea        string    `json:"Microsoft.VSTS.Common.ValueArea"`
	KanbanColumn     string    `json:"WEF_3F0006B87EE546F7B0696A999754B992_Kanban.Column"`
	KanbanColumnDone bool      `json:"WEF_3F0006B87EE546F7B0696A999754B992_Kanban.Column.Done"`
	Description      string    `json:"System.Description"`
}

type Identity struct {
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	Links       struct {
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"_links"`
	ID         string `json:"id"`
	UniqueName string `json:"uniqueName"`
	ImageURL   string `json:"imageUrl"`
	Descriptor string `json:"descriptor"`
}

type Links struct {
	Self              Link `json:"self"`
	WorkItemUpdates   Link `json:"workItemUpdates"`
	WorkItemRevisions Link `json:"workItemRevisions"`
	WorkItemComments  Link `json:"workItemComments"`
	HTML              Link `json:"html"`
}

type Link struct {
	Href string `json:"href"`
}

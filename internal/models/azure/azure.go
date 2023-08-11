package azure

type Notification struct {
	SubscriptionID     string             `json:"subscriptionId"`
	NotificationID     int                `json:"notificationId"`
	ID                 string             `json:"id"`
	EventType          string             `json:"eventType"`
	PublisherID        string             `json:"publisherId"`
	Message            Message            `json:"message"`
	DetailedMessage    DetailedMessage    `json:"detailedMessage"`
	Resource           Resource           `json:"resource"`
	ResourceVersion    string             `json:"resourceVersion"`
	ResourceContainers ResourceContainers `json:"resourceContainers"`
	CreatedDate        string             `json:"createdDate"`
}

type Message struct {
	Text     string `json:"text"`
	HTML     string `json:"html"`
	Markdown string `json:"markdown"`
}

type DetailedMessage struct {
	Text     string `json:"text"`
	HTML     string `json:"html"`
	Markdown string `json:"markdown"`
}

type Resource struct {
	ID          int       `json:"id"`
	WorkItemID  int       `json:"workItemId"`
	Rev         int       `json:"rev"`
	RevisedBy   RevisedBy `json:"revisedBy"`
	RevisedDate string    `json:"revisedDate"`
	Fields      Fields    `json:"fields"`
	Links       Links     `_links`
	URL         string    `json:"url"`
	Revision    Revision  `json:"revision"`
}

type RevisedBy struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	Links       Links  `_links`
	UniqueName  string `json:"uniqueName"`
	ImageURL    string `json:"imageUrl"`
	Descriptor  string `json:"descriptor"`
}

type Links struct {
	Avatar          Link `json:"avatar,omitempty"`
	Self            Link `json:"self,omitempty"`
	Parent          Link `json:"parent,omitempty"`
	WorkItemUpdates Link `json:"workItemUpdates,omitempty"`
}

type Link struct {
	Href string `json:"href"`
}

type Fields struct {
	SystemRev                   FieldChange `json:"System.Rev"`
	SystemAuthorizedDate        FieldChange `json:"System.AuthorizedDate"`
	SystemRevisedDate           FieldChange `json:"System.RevisedDate"`
	SystemState                 FieldChange `json:"System.State"`
	SystemReason                FieldChange `json:"System.Reason"`
	SystemAssignedTo            FieldChange `json:"System.AssignedTo"`
	SystemChangedDate           FieldChange `json:"System.ChangedDate"`
	SystemWatermark             FieldChange `json:"System.Watermark"`
	MicrosoftVSTSCommonSeverity FieldChange `json:"Microsoft.VSTS.Common.Severity"`
}

type FieldChange struct {
	OldValue interface{} `json:"oldValue"`
	NewValue interface{} `json:"newValue"`
}

type Revision struct {
	ID     int            `json:"id"`
	Rev    int            `json:"rev"`
	Fields RevisionFields `json:"fields"`
	URL    string         `json:"url"`
}

type RevisionFields struct {
	SystemAreaPath                                  string    `json:"System.AreaPath"`
	SystemTeamProject                               string    `json:"System.TeamProject"`
	SystemIterationPath                             string    `json:"System.IterationPath"`
	SystemWorkItemType                              string    `json:"System.WorkItemType"`
	SystemState                                     string    `json:"System.State"`
	SystemReason                                    string    `json:"System.Reason"`
	SystemCreatedDate                               string    `json:"System.CreatedDate"`
	SystemCreatedBy                                 RevisedBy `json:"System.CreatedBy"`
	SystemChangedDate                               string    `json:"System.ChangedDate"`
	SystemChangedBy                                 RevisedBy `json:"System.ChangedBy"`
	SystemTitle                                     string    `json:"System.Title"`
	MicrosoftVSTSCommonSeverity                     string    `json:"Microsoft.VSTS.Common.Severity"`
	WEFEB329F44FE5F4A94ACB1DA153FDF38BAKanbanColumn string    `json:"WEF_EB329F44FE5F4A94ACB1DA153FDF38BA_Kanban.Column"`
}

type ResourceContainers struct {
	Collection Container `json:"collection"`
	Account    Container `json:"account"`
	Project    Container `json:"project"`
}

type Container struct {
	ID string `json:"id"`
}

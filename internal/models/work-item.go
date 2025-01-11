package models

import (
	"fmt"
	"log"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/glamour"
)

type WorkItem struct {
	ID     int    `json:"id"`
	Rev    int    `json:"rev"`
	Fields Fields `json:"fields"`
	URL    string `json:"url"`
}

func (i WorkItem) Title() string {
	return fmt.Sprintf("[%d] %s", i.ID, i.Fields.Title)
}

func (i WorkItem) Description() string {
	return fmt.Sprintf("(%s) %s", i.Fields.State, i.Fields.AssignedTo.DisplayName)
}

func (i WorkItem) FilterValue() string {
	return i.Title() + i.Description()
}

func (i WorkItem) GetID() int {
	return i.ID
}

func (i WorkItem) GetURL() string {
	url := strings.Replace(i.URL, "_apis/wit/workItems", "_workitems/edit", 1)
	if err := clipboard.WriteAll(url); err != nil {
		log.Fatalf("Failed to copy to clipboard: %v", err)
	}

	return url
}

func (i WorkItem) GetPreview(renderer *glamour.TermRenderer) string {
	var description string

	if i.Fields.Description == nil {
		description = ""
	} else {
		description = *i.Fields.Description
	}
	converter := md.NewConverter("", true, &md.Options{})
	markdownDescription, err := converter.ConvertString(description)
	if err != nil {
		fmt.Println("Error converting HTML to Markdown:", err)
	}

	markdownContent := fmt.Sprintf("%d\n# %s\n---\n%s\n- URL: %s\n", i.ID, i.Fields.Title, markdownDescription, i.URL)
	rendered, _ := renderer.Render(markdownContent)
	return rendered
}

type Fields struct {
	AreaPath                     string   `json:"System.AreaPath"`
	TeamProject                  string   `json:"System.TeamProject"`
	IterationPath                string   `json:"System.IterationPath"`
	WorkItemType                 string   `json:"System.WorkItemType"`
	State                        string   `json:"System.State"`
	Reason                       string   `json:"System.Reason"`
	AssignedTo                   Identity `json:"System.AssignedTo"`
	CreatedDate                  string   `json:"System.CreatedDate"`
	CreatedBy                    Identity `json:"System.CreatedBy"`
	ChangedDate                  string   `json:"System.ChangedDate"`
	ChangedBy                    Identity `json:"System.ChangedBy"`
	CommentCount                 int      `json:"System.CommentCount"`
	Title                        string   `json:"System.Title"`
	Description                  *string  `json:"System.Description"`
	MicrosoftVSTSCommonPriority  int      `json:"Microsoft.VSTS.Common.Priority"`
	MicrosoftVSTSStateChangeDate string   `json:"Microsoft.VSTS.Common.StateChangeDate"`
	MicrosoftVSTSActivatedDate   string   `json:"Microsoft.VSTS.Common.ActivatedDate"`
	MicrosoftVSTSActivatedBy     Identity `json:"Microsoft.VSTS.Common.ActivatedBy"`
}

type Identity struct {
	DisplayName string     `json:"displayName"`
	URL         string     `json:"url"`
	Links       AvatarLink `json:"_links"`
	ID          string     `json:"id"`
	UniqueName  string     `json:"uniqueName"`
	ImageURL    string     `json:"imageUrl"`
	Inactive    bool       `json:"inactive"`
	Descriptor  string     `json:"descriptor"`
}

type AvatarLink struct {
	Avatar Avatar `json:"avatar"`
}

type Avatar struct {
	Href string `json:"href"`
}

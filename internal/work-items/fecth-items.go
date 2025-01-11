package workitems

import (
	"fmt"
	azhttpclient "lazyaz/internal/http"
	"lazyaz/internal/models"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type QueryPayload struct {
	Query string `json:"query"`
}

type WorkItemsResponse struct {
	WorkItems []struct {
		ID int `json:"id"`
	} `json:"workItems"`
}

type WorkItemsResponseMsg []models.WorkItem

func FetchWorkItems() tea.Msg {
	azHttpClient := azhttpclient.NewAzHttpClient()

	if !azHttpClient.HasValidPat() {
		log.Fatalf("Please set the AZURE_DEVOPS_PAT environment variable")
	}

	url := "https://wkeuds.visualstudio.com/_apis/wit/wiql?api-version=5"

	payload := QueryPayload{
		Query: `SELECT * FROM WorkItems WHERE [System.ChangedDate] >= @Today - 15 AND [System.NodeName] IN ('Krypton Team', 'Atalaya Team', 'Eternia Team', 'Castillo Grayskull', 'Estación Zeta') AND [System.WorkItemType] IN ('Task', 'User Story', 'Bug', 'Defect')`,
	}

	data, error := azhttpclient.Post[QueryPayload, WorkItemsResponse](azHttpClient, url, payload)
	if error != nil {
		log.Fatalf("could not fetch work items: %v", error)
	}

	totalWorkItemsUrl := "https://wkeuds.visualstudio.com/_apis/wit/workItems?ids="

	var ids []string
	for _, Item := range data.WorkItems {
		ids = append(ids, fmt.Sprintf("%d", Item.ID))
	}

	totalWorkItemsUrl += strings.Join(ids, ",")

	type Response struct {
		Count int               `json:"count"`
		Value []models.WorkItem `json:"value"`
	}

	workItems, err := azhttpclient.Get[Response](azHttpClient, totalWorkItemsUrl)
	if err != nil {
		log.Fatalf("could not fetch work items: %v", err)
	}

	return WorkItemsResponseMsg(workItems.Value)
}

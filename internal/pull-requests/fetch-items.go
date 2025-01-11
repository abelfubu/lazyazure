package pullrequests

import (
	azhttpclient "lazyaz/internal/http"
	pullrequests "lazyaz/internal/pull-requests/models"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type PullRequestResponseMsg []pullrequests.PullRequest

func FetchPullRequests() tea.Msg {
	azHttpClient := azhttpclient.NewAzHttpClient()

	if !azHttpClient.HasValidPat() {
		log.Fatalf("Please set the AZURE_DEVOPS_PAT environment variable")
	}

	url := "https://wkeuds.visualstudio.com/NewPOL/_apis/git/repositories/ERPRepo.Client.OINV/pullRequests?searchCriteria.includeLinks=False"

	type Response struct {
		Count int                    `json:"count"`
		Value PullRequestResponseMsg `json:"value"`
	}

	items, err := azhttpclient.Get[Response](azHttpClient, url)
	if err != nil {
		log.Fatalf("could not fetch work items: %v", err)
	}

	return PullRequestResponseMsg(items.Value)
}

package pullrequests

import (
	"fmt"
	"log"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/glamour"
)

type PullRequest struct {
	ArtifactID            *string            `json:"artifactId"`
	AutoCompleteSetBy     *Identity          `json:"autoCompleteSetBy"`
	ClosedBy              *Identity          `json:"closedBy"`
	ClosedDate            *string            `json:"closedDate"`
	CodeReviewID          int                `json:"codeReviewId"`
	Commits               *string            `json:"commits"`
	CompletionOptions     *CompletionOptions `json:"completionOptions"`
	CompletionQueueTime   *string            `json:"completionQueueTime"`
	CreatedBy             *Identity          `json:"createdBy"`
	CreationDate          string             `json:"creationDate"`
	Name                  string             `json:"title"`
	Body                  string             `json:"description"`
	ForkSource            *string            `json:"forkSource"`
	IsDraft               bool               `json:"isDraft"`
	Labels                *[]string          `json:"labels"`
	LastMergeCommit       *CommitDetails     `json:"lastMergeCommit"`
	LastMergeSourceCommit *CommitDetails     `json:"lastMergeSourceCommit"`
	LastMergeTargetCommit *CommitDetails     `json:"lastMergeTargetCommit"`
	MergeFailureMessage   *string            `json:"mergeFailureMessage"`
	MergeFailureType      *string            `json:"mergeFailureType"`
	MergeID               string             `json:"mergeId"`
	MergeOptions          *string            `json:"mergeOptions"`
	MergeStatus           string             `json:"mergeStatus"`
	PullRequestID         int                `json:"pullRequestId"`
	RemoteUrl             *string            `json:"remoteUrl"`
	Repository            *Repository        `json:"repository"`
	Reviewers             []Reviewer         `json:"reviewers"`
	Status                string             `json:"status"`
}

func (i PullRequest) Title() string {
	return fmt.Sprintf("[%d] %s", i.PullRequestID, i.Name)
}

func (i PullRequest) Description() string {
	date, err := time.Parse(time.RFC3339, i.CreationDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return fmt.Sprintf("%s - %s", i.CreatedBy.DisplayName, i.CreationDate)
	}

	formattedDate := date.Format("January 2, 2006")
	return fmt.Sprintf("%s - %s", i.CreatedBy.DisplayName, formattedDate)
}

func (i PullRequest) FilterValue() string {
	return i.Title() + i.Description()
}

func (i PullRequest) GetID() int {
	return i.PullRequestID
}

func (i PullRequest) GetURL() string {
	url := strings.Replace("", "_apis/wit/workItems", "_workitems/edit", 1)
	if err := clipboard.WriteAll(url); err != nil {
		log.Fatalf("Failed to copy to clipboard: %v", err)
	}

	return url
}

func (i PullRequest) GetPreview(renderer *glamour.TermRenderer) string {
	var description string

	converter := md.NewConverter("", true, &md.Options{})
	markdownDescription, err := converter.ConvertString(description)
	if err != nil {
		fmt.Println("Error converting HTML to Markdown:", err)
	}

	markdownContent := fmt.Sprintf("%d\n# %s\n---\n%s\n", i.PullRequestID, i.Name, markdownDescription)
	rendered, _ := renderer.Render(markdownContent)
	return rendered
}

type Identity struct {
	Descriptor        string  `json:"descriptor"`
	DirectoryAlias    *string `json:"directoryAlias"`
	DisplayName       string  `json:"displayName"`
	ID                string  `json:"id"`
	ImageUrl          string  `json:"imageUrl"`
	Inactive          *bool   `json:"inactive"`
	IsAadIdentity     *bool   `json:"isAadIdentity"`
	IsContainer       *bool   `json:"isContainer"`
	IsDeletedInOrigin *bool   `json:"isDeletedInOrigin"`
	ProfileUrl        *string `json:"profileUrl"`
	UniqueName        string  `json:"uniqueName"`
	Url               string  `json:"url"`
}

type CompletionOptions struct {
	AutoCompleteIgnoreConfigIds []string `json:"autoCompleteIgnoreConfigIds"`
	BypassPolicy                *bool    `json:"bypassPolicy"`
	BypassReason                *string  `json:"bypassReason"`
	DeleteSourceBranch          bool     `json:"deleteSourceBranch"`
	MergeCommitMessage          string   `json:"mergeCommitMessage"`
	MergeStrategy               string   `json:"mergeStrategy"`
	SquashMerge                 bool     `json:"squashMerge"`
	TransitionWorkItems         *bool    `json:"transitionWorkItems"`
	TriggeredByAutoComplete     *bool    `json:"triggeredByAutoComplete"`
}

type CommitDetails struct {
	Author           *string `json:"author"`
	ChangeCounts     *string `json:"changeCounts"`
	Changes          *string `json:"changes"`
	Comment          *string `json:"comment"`
	CommentTruncated *bool   `json:"commentTruncated"`
	CommitID         string  `json:"commitId"`
	Committer        *string `json:"committer"`
	Parents          *string `json:"parents"`
	Push             *string `json:"push"`
	RemoteUrl        *string `json:"remoteUrl"`
	Statuses         *string `json:"statuses"`
	Url              string  `json:"url"`
	WorkItems        *string `json:"workItems"`
}

type Repository struct {
	DefaultBranch    *string `json:"defaultBranch"`
	ID               string  `json:"id"`
	IsFork           *bool   `json:"isFork"`
	Name             string  `json:"name"`
	ParentRepository *string `json:"parentRepository"`
	Project          Project `json:"project"`
	RemoteUrl        *string `json:"remoteUrl"`
	Size             *int    `json:"size"`
	SshUrl           *string `json:"sshUrl"`
	Url              string  `json:"url"`
	ValidRemoteUrls  *string `json:"validRemoteUrls"`
}

type Project struct {
	Abbreviation        *string `json:"abbreviation"`
	DefaultTeamImageUrl *string `json:"defaultTeamImageUrl"`
	Description         *string `json:"description"`
	ID                  string  `json:"id"`
	LastUpdateTime      string  `json:"lastUpdateTime"`
	Name                string  `json:"name"`
	Revision            *string `json:"revision"`
	State               string  `json:"state"`
	Url                 *string `json:"url"`
	Visibility          string  `json:"visibility"`
}

type Reviewer struct {
	Descriptor        *string `json:"descriptor"`
	DirectoryAlias    *string `json:"directoryAlias"`
	DisplayName       string  `json:"displayName"`
	HasDeclined       bool    `json:"hasDeclined"`
	ID                string  `json:"id"`
	ImageUrl          string  `json:"imageUrl"`
	Inactive          *bool   `json:"inactive"`
	IsAadIdentity     *bool   `json:"isAadIdentity"`
	IsContainer       *bool   `json:"isContainer"`
	IsDeletedInOrigin *bool   `json:"isDeletedInOrigin"`
	IsFlagged         bool    `json:"isFlagged"`
	IsRequired        *bool   `json:"isRequired"`
	ProfileUrl        *string `json:"profileUrl"`
	ReviewerUrl       string  `json:"reviewerUrl"`
	UniqueName        string  `json:"uniqueName"`
	Url               string  `json:"url"`
	Vote              int     `json:"vote"`
	VotedFor          *[]any  `json:"votedFor"`
}

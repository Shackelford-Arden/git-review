package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"sync"
)

func checkApprovalStatus(wg *sync.WaitGroup, gitlabClient *gitlab.Client, mergeRequestItem gitlab.MergeRequest, gitlabUser gitlab.User, approvalMissing chan *gitlab.MergeRequest) {
	defer wg.Done()

	approvalStatus, _, err := gitlabClient.MergeRequestApprovals.GetApprovalState(mergeRequestItem.ProjectID, mergeRequestItem.IID)

	if err != nil {
		log.Fatalf("Failed to get approval status for %s | %v | %s\n\n", mergeRequestItem.Title, mergeRequestItem.ProjectID, err)
	}

	var isApproved bool
	for _, approvedByUser := range approvalStatus.Rules[0].ApprovedBy {

		if approvedByUser.ID == gitlabUser.ID {
			isApproved = true
		}
	}
	if isApproved == false {
		approvalMissing <- &mergeRequestItem
	}
}

func CheckMRStatus() {

	gitlabPAT := os.Getenv("GITLAB_PAT")

	gitlabClient, err := gitlab.NewClient(gitlabPAT)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get User's Info, makes some things easier later
	userData, _, err := gitlabClient.Users.CurrentUser()
	if err != nil {
		log.Fatalf("failed to get the current user: %s\n", err)
	}
	log.Printf("Your Gitlab User ID is: %v", userData.ID)

	// Find MRs w/ Title
	log.Println("Finding MRs...")
	mrOptions := &gitlab.ListMergeRequestsOptions{
		Search: gitlab.String("fancy mr title"),
		In:     gitlab.String("title"),
		State:  gitlab.String("opened"),
		Scope:  gitlab.String("all"),
		WIP:    gitlab.String("no"),
	}
	mrList, _, err := gitlabClient.MergeRequests.ListMergeRequests(mrOptions)

	if err != nil {
		log.Fatalln(err)
	}

	missingApprovalList := make(chan *gitlab.MergeRequest)

	log.Printf("Found %v Merge Requests based on filters. Now checking if they still need your approval, sir", len(mrList))
	wg := new(sync.WaitGroup)
	for _, mrItem := range mrList {

		wg.Add(1)

		go checkApprovalStatus(wg, gitlabClient, *mrItem, *userData, missingApprovalList)

	}
	close(missingApprovalList)
	wg.Wait()

	var finalResults []*gitlab.MergeRequest
	for item := range missingApprovalList {
		finalResults = append(finalResults, item)
	}

	log.Printf("%v MRs still pending your approval, sir\n", len(finalResults))
	for _, missingApprovalMR := range finalResults {
		log.Printf("Here we go!: %s | %v | %s\n", missingApprovalMR.Title, missingApprovalMR.ProjectID, missingApprovalMR.WebURL)
	}

}

package hubclient

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/bdsengineering/go-hub-client/hubapi"
	log "github.com/sirupsen/logrus"
)

// TestFetchPolicyStatus is a very brittle test because it requires:
//   1. a reachable hub backend
//   2. the hub backend to be located on localhost
//   3. a specific username and password to be able to log in
//   4. that there is at least one project, with a version, with a policy status
// It's actually an integration test, not a unit test.
func TestCreateAndDeleteProject(t *testing.T) {
	client, err := NewWithSession("https://localhost", HubClientDebugTimings)
	if err != nil {
		t.Error(err)
	}
	err = client.Login("sysadmin", "blackduck")
	if err != nil {
		t.Error(err)
	}

	projectName := "first-new-project"
	projectRequest := hubapi.ProjectRequest{Name: projectName}

	// create project
	location, err := client.CreateProject(&projectRequest)
	log.Infof("location: %s", location)
	if err != nil {
		t.Error(err)
	}
	// find project
	q := fmt.Sprintf("name:%s", projectName)
	projectList, err := client.ListProjects(&hubapi.GetListOptions{Q: &q})
	if err != nil {
		t.Error(err)
	}
	projectID := ""
	for _, proj := range projectList.Items {
		if proj.Name == projectName {
			strs := strings.Split(proj.Meta.Href, "/")
			if len(strs) > 0 {
				projectID = strs[len(strs)-1]
			}
			break
		}
	}

	if projectID == "" {
		t.Errorf("failed to find project of name %s", projectName)
	}

	// delete project
	err = client.DeleteProject(projectID)
	if err != nil {
		t.Error(err)
	}
}
package test

import (
	"context"
	"fmt"
	"os"
	// "strconv"
	"testing"
	"strings"

	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	awsTerratest "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/google/go-github/v70/github"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "argocd-github-repo/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"
	const githubOrg = "cloudposse-tests"

	token := os.Getenv("GITHUB_TOKEN")

	randomID := strings.ToLower(random.UniqueId())

	secretPath := fmt.Sprintf("/argocd/%s/github/api_key", randomID)
	deployKeyPath := fmt.Sprintf("/argocd/deploy_keys/%s/%s", randomID, "%s")
	repoName := fmt.Sprintf("argocd-github-repo-%s", randomID)

	defer func() {
		awsTerratest.DeleteParameter(s.T(), awsRegion, secretPath)
	}()
	awsTerratest.PutParameter(s.T(), awsRegion, secretPath, "Github API Key", token)

	inputs := map[string]interface{}{
		"ssm_github_deploy_key_format": deployKeyPath,
		"ssm_github_api_key": secretPath,
		"name": repoName,
		"github_organization": githubOrg,
	}

	defer s.DestroyAtmosComponent(s.T(), component, stack, &inputs)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, &inputs)
	assert.NotNil(s.T(), options)

	client := github.NewClient(nil).WithAuthToken(token)

	// Check if the repository exists
	repo, _, err := client.Repositories.Get(context.Background(), githubOrg, repoName)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), repo)

	filePath := "mgmt/uw2-sandbox/argocd/applicationset.yaml"
	// Use the GitHub API to check if the file exists in the repository
	_, _, _, err = client.Repositories.GetContents(context.Background(), githubOrg, repoName, filePath, nil)
	assert.Nil(s.T(), err)
	s.DriftTest(component, stack, &inputs)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "argocd-github-repo/disabled"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	s.VerifyEnabledFlag(component, stack, nil)
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)
	helper.Run(t, suite)
}

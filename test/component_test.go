package test

import (
	"fmt"
	"os"
	// "strconv"
	"testing"
	"strings"

	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	awsTerratest "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	// "github.com/jferrl/go-githubauth"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "argocd-github-repo/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	token := os.Getenv("GITHUB_TOKEN")

	randomID := strings.ToLower(random.UniqueId())

	secretPath := fmt.Sprintf("/argocd/%s/github/api_key", randomID)
	repoName := fmt.Sprintf("argocd-github-repo-%s", randomID)

	defer func() {
		awsTerratest.DeleteParameter(s.T(), awsRegion, secretPath)
	}()
	awsTerratest.PutParameter(s.T(), awsRegion, secretPath, "Github API Key", token)

	inputs := map[string]interface{}{
		"ssm_github_api_key": secretPath,
		"name":                repoName,
	}

	defer s.DestroyAtmosComponent(s.T(), component, stack, &inputs)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, &inputs)
	assert.NotNil(s.T(), options)

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

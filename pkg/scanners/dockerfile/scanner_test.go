package dockerfile

import (
	"context"
	"testing"

	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/aquasecurity/defsec/test/testutil"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_BasicScan(t *testing.T) {

	fs := testutil.CreateFS(t, map[string]string{
		"/code/Dockerfile": `FROM ubuntu

USER root
`,
		"/rules/rule.rego": `package appshield.dockerfile.DS006

__rego_metadata__ := {
	"id": "DS006",
	"avd_id": "AVD-DS-0006",
	"title": "COPY '--from' referring to the current image",
	"short_code": "no-self-referencing-copy-from",
	"version": "v1.0.0",
	"severity": "CRITICAL",
	"type": "Dockerfile Security Check",
	"description": "COPY '--from' should not mention the current FROM alias, since it is impossible to copy from itself.",
	"recommended_actions": "Change the '--from' so that it will not refer to itself",
	"url": "https://docs.docker.com/develop/develop-images/multistage-build/",
}

__rego_input__ := {
	"combine": false,
	"selector": [{"type": "dockerfile"}],
}

deny[res] {
	res := {
		"msg": "oh no",
		"filepath": "Dockerfile",
		"startline": 1,
		"endline": 2,
	}
}

`,
	})

	scanner := NewScanner(OptionWithPolicyDirs("rules"))

	results, err := scanner.ScanFS(context.TODO(), fs, "code")
	require.NoError(t, err)

	require.Len(t, results.GetFailed(), 1)

	assert.Equal(t, scan.Rule{
		AVDID:       "AVD-DS-0006",
		LegacyID:    "DS006",
		ShortCode:   "no-self-referencing-copy-from",
		Summary:     "COPY '--from' referring to the current image",
		Explanation: "COPY '--from' should not mention the current FROM alias, since it is impossible to copy from itself.",
		Impact:      "",
		Resolution:  "Change the '--from' so that it will not refer to itself",
		Provider:    "dockerfile",
		Service:     "general",
		Links:       []string{"https://docs.docker.com/develop/develop-images/multistage-build/"},
		Severity:    "CRITICAL",
		Terraform:   (*scan.EngineMetadata)(nil), CloudFormation: (*scan.EngineMetadata)(nil), CustomChecks: scan.CustomChecks{Terraform: (*scan.TerraformCustomCheck)(nil)}, RegoPackage: "data.appshield.dockerfile.DS006"}, results.GetFailed()[0].Rule())
}

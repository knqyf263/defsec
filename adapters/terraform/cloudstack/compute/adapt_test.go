package compute

import (
	"testing"

	"github.com/aquasecurity/defsec/adapters/terraform/testutil"
	"github.com/aquasecurity/defsec/parsers/types"

	"github.com/aquasecurity/defsec/providers/cloudstack/compute"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_adaptInstance(t *testing.T) {
	tests := []struct {
		name      string
		terraform string
		expected  compute.Instance
	}{
		{
			name: "sensitive user data",
			terraform: `
			resource "cloudstack_instance" "web" {
				name             = "server-1"
				user_data        = <<EOF
export DATABASE_PASSWORD=\"SomeSortOfPassword\"
			EOF
			}
`,
			expected: compute.Instance{
				Metadata: types.NewTestMetadata(),
				UserData: types.String(`export DATABASE_PASSWORD=\"SomeSortOfPassword\"
`, types.NewTestMetadata()),
			},
		},
		{
			name: "sensitive user data base64 encoded",
			terraform: `
			resource "cloudstack_instance" "web" {
				name             = "server-1"
				user_data        = "ZXhwb3J0IERBVEFCQVNFX1BBU1NXT1JEPSJTb21lU29ydE9mUGFzc3dvcmQi"
			}
`,
			expected: compute.Instance{
				Metadata: types.NewTestMetadata(),
				UserData: types.String(`export DATABASE_PASSWORD="SomeSortOfPassword"`, types.NewTestMetadata()),
			},
		},
		{
			name: "no user data provided",
			terraform: `
			resource "cloudstack_instance" "web" {
			}
`,
			expected: compute.Instance{
				Metadata: types.NewTestMetadata(),
				UserData: types.String("", types.NewTestMetadata()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modules := testutil.CreateModulesFromSource(t, test.terraform, ".tf")
			adapted := adaptInstance(modules.GetBlocks()[0])
			testutil.AssertDefsecEqual(t, test.expected, adapted)
		})
	}
}

func TestLines(t *testing.T) {
	src := `
	resource "cloudstack_instance" "web" {
		name             = "server-1"
		user_data        = <<EOF
export DATABASE_PASSWORD=\"SomeSortOfPassword\"
	EOF
	}`

	modules := testutil.CreateModulesFromSource(t, src, ".tf")
	adapted := Adapt(modules)

	require.Len(t, adapted.Instances, 1)
	instance := adapted.Instances[0]

	assert.Equal(t, 4, instance.UserData.GetMetadata().Range().GetStartLine())
	assert.Equal(t, 6, instance.UserData.GetMetadata().Range().GetEndLine())
}
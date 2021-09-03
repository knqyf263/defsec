package athena

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckEnableAtRestEncryption = rules.Register(
	rules.Rule{
		Provider:    provider.AWSProvider,
		Service:     "athena",
		ShortCode:   "enable-at-rest-encryption",
		Summary:     "Athena databases and workgroup configurations are created unencrypted at rest by default, they should be encrypted",
		Impact:      "Data can be read if the Athena Database is compromised",
		Resolution:  "Enable encryption at rest for Athena databases and workgroup configurations",
		Explanation: `Athena databases and workspace result sets should be encrypted at rests. These databases and query sets are generally derived from data in S3 buckets and should have the same level of at rest protection.`,
		Links: []string{ 
			"https://docs.aws.amazon.com/athena/latest/ug/encryption.html",
		},
		Severity: severity.High,
	},
	func(s *state.State) (results rules.Results) {
		for _, x := range s.AWS.S3.Buckets {
			if x.Encryption.Enabled.IsFalse() {
				results.Add(
					"",
					x.Encryption.Enabled.Metadata(),
					x.Encryption.Enabled.Value(),
				)
			}
		}
		return
	},
)
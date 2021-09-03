package elasticsearch

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckEncryptReplicationGroup = rules.Register(
	rules.Rule{
		Provider:    provider.AWSProvider,
		Service:     "elastic-search",
		ShortCode:   "encrypt-replication-group",
		Summary:     "Unencrypted Elasticache Replication Group.",
		Impact:      "Data in the replication group could be readable if compromised",
		Resolution:  "Enable encryption for replication group",
		Explanation: `You should ensure your Elasticache data is encrypted at rest to help prevent sensitive information from being read by unauthorised users.`,
		Links: []string{ 
			"https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/at-rest-encryption.html",
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
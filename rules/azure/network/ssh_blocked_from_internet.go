package network

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckSshBlockedFromInternet = rules.Register(
	rules.Rule{
		Provider:    provider.AzureProvider,
		Service:     "network",
		ShortCode:   "ssh-blocked-from-internet",
		Summary:     "SSH access should not be accessible from the Internet, should be blocked on port 22",
		Impact:      "Its dangerous to allow SSH access from the internet",
		Resolution:  "Block port 22 access from the internet",
		Explanation: `SSH access can be configured on either the network security group or in the network security group rule. 

SSH access should not be permitted from the internet (*, 0.0.0.0, /0, internet, any)`,
		Links: []string{ 
		},
		Severity: severity.Critical,
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
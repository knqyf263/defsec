package compute

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckNoPlaintextPassword = rules.Register(
	rules.Rule{
		Provider:    provider.OpenStackProvider,
		Service:     "compute",
		ShortCode:   "no-plaintext-password",
		Summary:     "No plaintext password for compute instance",
		Impact:      "Including a plaintext password could lead to compromised instance",
		Resolution:  "Do not use plaintext passwords in terraform files",
		Explanation: `Assigning a password to the compute instance using plaintext could lead to compromise; it would be preferable to use key-pairs as a login mechanism`,
		Links: []string{ 
		},
		Severity: severity.Medium,
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
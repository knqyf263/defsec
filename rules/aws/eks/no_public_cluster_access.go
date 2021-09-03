package eks

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckNoPublicClusterAccess = rules.Register(
	rules.Rule{
		Provider:    provider.AWSProvider,
		Service:     "eks",
		ShortCode:   "no-public-cluster-access",
		Summary:     "EKS Clusters should have the public access disabled",
		Impact:      "EKS can be access from the internet",
		Resolution:  "Don't enable public access to EKS Clusters",
		Explanation: `EKS clusters are available publicly by default, this should be explicitly disabled in the vpc_config of the EKS cluster resource.`,
		Links: []string{ 
			"https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html",
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
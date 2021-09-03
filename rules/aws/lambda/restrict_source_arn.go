package lambda

import (
	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/defsec/state"
)

var CheckRestrictSourceArn = rules.Register(
	rules.Rule{
		Provider:    provider.AWSProvider,
		Service:     "lambda",
		ShortCode:   "restrict-source-arn",
		Summary:     "Ensure that lambda function permission has a source arn specified",
		Impact:      "Not providing the source ARN allows any resource from principal, even from other accounts",
		Resolution:  "Always provide a source arn for Lambda permissions",
		Explanation: `When the principal is an AWS service, the ARN of the specific resource within that service to grant permission to. 

Without this, any resource from principal will be granted permission – even if that resource is from another account. 

For S3, this should be the ARN of the S3 Bucket. For CloudWatch Events, this should be the ARN of the CloudWatch Events Rule. For API Gateway, this should be the ARN of the API`,
		Links: []string{ 
			"https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html",
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
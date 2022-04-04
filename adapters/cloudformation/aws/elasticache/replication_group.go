package elasticache

import (
	"github.com/aquasecurity/defsec/parsers/cloudformation/parser"
	"github.com/aquasecurity/defsec/providers/aws/elasticache"
)

func getReplicationGroups(ctx parser.FileContext) (replicationGroups []elasticache.ReplicationGroup) {

	replicationGroupResources := ctx.GetResourcesByType("AWS::ElastiCache::ReplicationGroup")

	for _, r := range replicationGroupResources {
		replicationGroup := elasticache.ReplicationGroup{
			Metadata:                 r.Metadata(),
			TransitEncryptionEnabled: r.GetBoolProperty("TransitEncryptionEnabled"),
			AtRestEncryptionEnabled:  r.GetBoolProperty("AtRestEncryptionEnabled"),
		}

		replicationGroups = append(replicationGroups, replicationGroup)
	}

	return replicationGroups
}
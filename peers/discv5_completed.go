package peers

import (
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/signal"
)

// MailServerDiscoveryTopic topic name for mailserver discovery.
const MailServerDiscoveryTopic = "mailserver.discovery"

// MailServerDiscoveryLimits default limits to discover mail servers.
var MailServerDiscoveryLimits = params.Limits{Max: 3, Min: 2}

var (
	sendMailServerDiscoveryCompleted = signal.SendMailServerDiscoveryCompleted
	sendMailServerDiscoveryErrored   = signal.SendMailServerDiscoveryErrored
)

// HandleCompletedDiscovery handles completed discovery TopicPools.
func HandleCompletedDiscovery(t *TopicPool) {
	if t.topic != MailServerDiscoveryTopic {
		return
	}

	if t.maxCachedPeersReached() {
		list := []string{}
		for _, p := range t.connectedPeers {
			list = append(list, p.node.String())
		}
		if len(list) > 0 {
			sendMailServerDiscoveryCompleted(list)
			return
		}
	}
	sendMailServerDiscoveryErrored()
}

package peers

import (
	"net"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cache(t *testing.T) *Cache {
	peersDB, err := newInMemoryCache()
	require.NoError(t, err)
	return peersDB
}

func TestHandleCompletedDiscovery(t *testing.T) {
	ni := discv5.NodeID{1}
	pi := &peerInfo{
		node: discv5.NewNode(discv5.NodeID{3}, net.IPv4(100, 100, 0, 3), 32311, 32311),
	}
	testCases := []struct {
		name        string
		topic       *TopicPool
		completed   bool
		errored     bool
		mailservers []string
	}{
		{
			name: "happy path",
			topic: &TopicPool{
				topic: MailServerDiscoveryTopic,
				connectedPeers: map[discv5.NodeID]*peerInfo{
					ni: pi,
				},
			},
			mailservers: []string{"enode://03000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000@100.100.0.3:32311"},
			completed:   true,
		},
		{
			name: "non mailserver discovery topic",
			topic: &TopicPool{
				topic: "non.mailserver.discovery",
			},
		},
		{
			name: "no mailservers found",
			topic: &TopicPool{
				topic: MailServerDiscoveryTopic,
			},
			errored: true,
		},
		{
			name: "discovery errored",
			topic: &TopicPool{
				topic:          MailServerDiscoveryTopic,
				maxCachedPeers: 100,
				cache:          cache(t),
			},
			errored: true,
		},
	}

	var completed, errored bool
	var mailservers []string

	sendMailServerDiscoveryCompleted = func(list []string) {
		completed = true
		mailservers = list
	}
	sendMailServerDiscoveryErrored = func() {
		errored = true
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mailservers = nil
			completed = false
			errored = false
			HandleCompletedDiscovery(tc.topic)
			assert.Equal(t, tc.completed, completed)
			assert.Equal(t, tc.errored, errored)
			assert.Equal(t, tc.mailservers, mailservers)
		})
	}
}

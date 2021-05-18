package minersc

import (
	"testing"

	"0chain.net/chaincore/node"
	"github.com/stretchr/testify/assert"
)

func createTestSimpleNodesAndNodePool() (SimpleNodes, *node.Pool) {

	sn := NewSimpleNodes()
	sn["0"] = &SimpleNode{ID: "0", TotalStaked: 12}
	sn["1"] = &SimpleNode{ID: "1", TotalStaked: 10}
	sn["2"] = &SimpleNode{ID: "2", TotalStaked: 8}
	sn["3"] = &SimpleNode{ID: "3", TotalStaked: 5}
	sn["4"] = &SimpleNode{ID: "4", TotalStaked: 3}
	sn["5"] = &SimpleNode{ID: "5", TotalStaked: 3}
	sn["6"] = &SimpleNode{ID: "6", TotalStaked: 2}
	sn["7"] = &SimpleNode{ID: "7", TotalStaked: 2}
	sn["8"] = &SimpleNode{ID: "8", TotalStaked: 2}
	sn["9"] = &SimpleNode{ID: "9", TotalStaked: 1}

	np := node.NewPool(node.NodeTypeMiner)

	var n *node.Node

	n = &node.Node{}
	n.ID = sn["6"].ID
	np.AddNode(n)

	n = &node.Node{}
	n.ID = sn["9"].ID
	np.AddNode(n)

	n = &node.Node{}
	n.ID = sn["4"].ID
	np.AddNode(n)

	n = &node.Node{}
	n.ID = sn["2"].ID
	np.AddNode(n)

	return sn, np
}

func TestSimpleNodesReduce(t *testing.T) {
	var pmbrss int64 = 123456789

	// select up to 5 of the existing nodes
	sn, np := createTestSimpleNodesAndNodePool()
	sn.reduce(7, 0.7, pmbrss, np)
	for _, n := range sn {
		assert.Contains(t, []string{"2", "4", "6", "9", "0", "1", "3"}, n.ID)
	}

	// select up to 3 nodes from previous set and rest by desc stake
	sn, np = createTestSimpleNodesAndNodePool()
	sn.reduce(5, 0.6, pmbrss, np)
	for _, n := range sn {
		assert.Contains(t, []string{"2", "4", "6", "0", "1"}, n.ID)
	}

	// select up to 5 nodes from previous set and rest by desc stake
	sn, np = createTestSimpleNodesAndNodePool()
	sn.reduce(8, 0.6, pmbrss, np)
	for _, n := range sn {
		assert.Contains(t, []string{"2", "4", "6", "9", "0", "1", "3", "5"}, n.ID)
	}

	// select up to 6 nodes form previous set (4), and rest by desc stake
	// resolve equal stake (7:2, 8:2) using pmbrss
	sn, np = createTestSimpleNodesAndNodePool()
	sn.reduce(9, 0.6, pmbrss, np)
	for _, n := range sn {
		assert.Contains(t, []string{"2", "4", "6", "9", "0", "1", "3", "5", "8"}, n.ID)
	}

	// select up to 6 nodes form previous set (4), and rest by desc stake
	// resolve equal stake (7:2, 8:2) using pmbrss+2
	sn, np = createTestSimpleNodesAndNodePool()
	sn.reduce(9, 0.6, pmbrss+2, np)
	for _, n := range sn {
		assert.Contains(t, []string{"2", "4", "6", "9", "0", "1", "3", "5", "7"}, n.ID)
	}

}

func TestQuickFixDuplicateHosts(t *testing.T) {
	node := func(n2nhost, host string, port int) *MinerNode {
		return &MinerNode{SimpleNode: &SimpleNode{N2NHost: n2nhost, Host: host, Port: port}}
	}
	nodes := func() []*MinerNode {
		return []*MinerNode{
			{SimpleNode: &SimpleNode{N2NHost: "abc.com", Host: "lmn.com", Port: 0}},
		}
	}
	assert.EqualError(t, quickFixDuplicateHosts(node("", "", 0), nodes()), "invalid n2nhost: ")
	assert.EqualError(t, quickFixDuplicateHosts(node("localhost", "", 0), nodes()), "invalid n2nhost: localhost")
	assert.EqualError(t, quickFixDuplicateHosts(node("127.0.0.1", "", 0), nodes()), "invalid n2nhost: 127.0.0.1")
	assert.NoError(t, quickFixDuplicateHosts(node("xyz.com", "", 0), nodes()))
	assert.NoError(t, quickFixDuplicateHosts(node("xyz.com", "localhost", 0), nodes()))
	assert.NoError(t, quickFixDuplicateHosts(node("xyz.com", "127.0.0.1", 0), nodes()))
	assert.NoError(t, quickFixDuplicateHosts(node("xyz.com", "prq.com", 0), nodes()))
	assert.EqualError(t, quickFixDuplicateHosts(node("abc.com", "", 0), nodes()), "n2nhost:port already exists: abc.com:0")
	assert.NoError(t, quickFixDuplicateHosts(node("abc.com", "", 1), nodes()))
	assert.EqualError(t, quickFixDuplicateHosts(node("lmn.com", "", 0), nodes()), "host:port already exists: lmn.com:0")
	assert.NoError(t, quickFixDuplicateHosts(node("lmn.com", "", 1), nodes()))

}

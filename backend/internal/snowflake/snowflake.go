package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch     int64 = 1577836800000 // Custom epoch (e.g., 2020-01-01)
	nodeBits  uint8 = 10
	stepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

type Node struct {
	mu        sync.Mutex
	timestamp int64
	nodeID    int64
	step      int64
}

func NewNode(nodeID int64) (*Node, error) {
	if nodeID < 0 || nodeID > nodeMax {
		return nil, errors.New("node ID out of range")
	}
	return &Node{
		timestamp: 0,
		nodeID:    nodeID,
		step:      0,
	}, nil
}

func (n *Node) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == n.timestamp {
		n.step = (n.step + 1) & stepMax
		if n.step == 0 {
			for now <= n.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		n.step = 0
	}

	n.timestamp = now

	return (now-epoch)<<timeShift | (n.nodeID << nodeShift) | n.step
}

// Copyright 2019 The readailib Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radix

import "strings"

// refer: https://en.wikipedia.org/wiki/Radix_tree
type Edge struct {
	containKey string
	sourceNode *Node
	destNode   *Node
}

// 真实存储的节点
type leafNode struct {
	key   string
	value interface{}
}

type Node struct {
	leaf  *leafNode // 叶子节点
	edges []Edge    // 每个节点对应有很多条边
}

func (n *Node) isLeafNode() bool {
	return n.leaf != nil && len(n.edges) == 0
}

func (n *Node) insertLeafNode(containKey, totalKey string, value interface{}) {
	newNode := &Node{leaf: &leafNode{key: totalKey, value: value}}
	newEdge := Edge{containKey: containKey, sourceNode: n, destNode: newNode}
	n.edges = append(n.edges, newEdge)
}

func (n *Node) insertSplitNode(splitKey string, edgeKey string) *Node {

	if n.isLeafNode() {
		// node is leaf node could not split, return nil
		return nil
	}

	for edgeIndex := range n.edges {
		if n.edges[edgeIndex].containKey == edgeKey {
			// backup for split
			originTargetNode := n.edges[edgeIndex].destNode
			// insert split node
			splitNode := &Node{}
			n.edges[edgeIndex] = Edge{containKey: splitKey, sourceNode: n, destNode: splitNode}

			// connect to origin node
			remainKey := strings.TrimPrefix(edgeKey, splitKey)
			edgeFromSplitToOri := Edge{containKey: remainKey, sourceNode: splitNode, destNode: originTargetNode}
			splitNode.edges = append(splitNode.edges, edgeFromSplitToOri)

			return splitNode
		}
	}
	return nil
}

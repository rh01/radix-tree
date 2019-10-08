package radix

import (
	"fmt"
	"strings"
)

// tier tree
type RadixTree struct {
	root Node
}

func contrainPredix(str1, str2 string) bool {
	if sub, find := getSubsetPrefix(str1, str2); find && sub == str2 {
		return true
	}
	return false
}

func getSubsetPrefix(str1, str2 string) (string, bool) {
	findSubset := false
	for i := 0; i < len(str1) && i < len(str2); i++ {
		if str1[i] != str2[i] {
			retStr := str1[:i]
			return retStr, findSubset
		}
		findSubset = true
	}

	if len(str1) > len(str2) {
		return str2, findSubset
	} else if len(str1) == len(str2) {
		// fix "" not a subset of ""
		return str1, str1 == str2
	}

	return str1, findSubset
}

// factory method
func NewRadixTree() *RadixTree {
	return &RadixTree{}
}

// recursive print tree
func (t *RadixTree) recursivePrintTree(currentNode *Node, treeLevel int) {
	indentStr := ""
	for i := 0; i < treeLevel; i++ {
		indentStr += indentStr + "\t"
	}

	if currentNode.isLeafNode() {
		// Reach the end point
		fmt.Printf("%s[%d]Leaf Key: %s value: %v\n", indentStr, treeLevel, currentNode.leaf.key, currentNode.leaf.value)
		return
	}

	fmt.Printf("%s[%d]Node has %d edges \n", indentStr, treeLevel, len(currentNode.edges))
	for _, edgeObj := range currentNode.edges {
		fmt.Printf("%s[%d]Edge[%s]\n", indentStr, treeLevel, string(edgeObj.containKey))
		t.recursivePrintTree(edgeObj.destNode, treeLevel+1)
	}
}

// PrintTree: Print out current Tree struct, it will use \t for tree level
func (t *RadixTree) PrintTree() {
	t.recursivePrintTree(&t.root, 1)
}

func (t *RadixTree) recursiveInsertTree(currentNode *Node, containKey string, targetKey string, targetValue interface{}) {
	// Reach leaf the end print, refer this case
	if currentNode.isLeafNode() {
		if targetKey == currentNode.leaf.key {
			// the same key, update the value
			currentNode.leaf.value = targetValue
			return
		} else {
			//Insert key value as new Child node of currentNode
			//original leaf node, become another leaf of currentNode
			//currentNode become not leaf node
			currentNode.insertLeafNode(containKey, targetKey, targetValue)
			currentNode.insertLeafNode("", currentNode.leaf.key, currentNode.leaf.value)
			currentNode.leaf = nil
			return
		}
	}

	for edgeIndex := range currentNode.edges {
		subStr, find := getSubsetPrefix(containKey, currentNode.edges[edgeIndex].containKey)
		if find {
			if subStr == currentNode.edges[edgeIndex].containKey {
				nextTargetKey := strings.TrimPrefix(containKey, currentNode.edges[edgeIndex].containKey)
				t.recursiveInsertTree(currentNode.edges[edgeIndex].destNode, nextTargetKey, targetKey, targetValue)
				return
			} else {
				//contain case
				//using subStr to create new node and separate two edges
				//Refer: https://goo.gl/j2MDBI
				splitNode := currentNode.insertSplitNode(subStr, currentNode.edges[edgeIndex].containKey)
				if splitNode == nil {
					panic("Unexpect error on split node")
				}

				splitContainKey := strings.TrimPrefix(containKey, subStr)
				splitNode.insertLeafNode(splitContainKey, targetKey, targetValue)
				return
			}
		}
	}

	//New edge with new key on leaf node
	//Ref: https://goo.gl/nSLTJr
	currentNode.insertLeafNode(containKey, targetKey, targetValue)
	return
}

//Insert: key and value into radix tree
//Major implement refer from Wiki: https://en.wikipedia.org/wiki/Radix_tree
func (t *RadixTree) Insert(searchKey string, value interface{}) {
	t.recursiveInsertTree(&t.root, searchKey, searchKey, value)
}

func (t *RadixTree) recursiveLoopup(searchNode *Node, searchKey string) (interface{}, bool) {
	if searchNode.isLeafNode() {
		return searchNode.leaf.value, true
	}

	for _, edgeObj := range searchNode.edges {
		if contrainPredix(searchKey, edgeObj.containKey) {
			nextSearchKey := strings.TrimPrefix(searchKey, edgeObj.containKey)
			return t.recursiveLoopup(edgeObj.destNode, nextSearchKey)
		}
	}

	return nil, false
}

//Lookup: Find if seachKey exist in current radix tree and return its value
func (t *RadixTree) Lookup(searchKey string) (interface{}, bool) {
	return t.recursiveLoopup(&t.root, searchKey)
}

func (t *RadixTree) recursiveLocateLeafNode(currentNode, parentNode *Node, containKey, locateKey string) (*Node, *Node, bool) {

	if currentNode.isLeafNode() {
		return currentNode, parentNode, currentNode.leaf.key == locateKey
	}

	for _, edgeObj := range currentNode.edges {
		if contrainPredix(containKey, edgeObj.containKey) {
			nextContainKey := strings.TrimPrefix(containKey, edgeObj.containKey)
			return t.recursiveLocateLeafNode(edgeObj.destNode, currentNode, nextContainKey, locateKey)
		}
	}

	return nil, nil, false
}

func (t *RadixTree) locateLeafNode(locateKey string) (locateNode, parentNode *Node, find bool) {
	locateNode, parentNode, find = t.recursiveLocateLeafNode(&t.root, &t.root, locateKey, locateKey)
	return locateNode, parentNode, find
}

func (t *RadixTree) recursiveFindParent(currentNode, parentNode, locateNode *Node) (*Node, bool) {
	if currentNode.isLeafNode() {
		return nil, false
	}

	if currentNode == locateNode {
		return parentNode, true
	}

	for _, edgeObj := range currentNode.edges {
		if edgeObj.destNode == locateNode {
			return currentNode, true
		}

		if pNode, find := t.recursiveFindParent(edgeObj.destNode, currentNode, locateNode); find {
			return pNode, true
		}
	}

	return nil, false
}

func (t *RadixTree) findParent(locateNode *Node) (*Node, bool) {
	return t.recursiveFindParent(&t.root, &t.root, locateNode)
}

//Delete: Delete leaf node by seachKey will return if exist
func (t *RadixTree) Delete(searchKey string) bool {

	lNode, pNode, find := t.locateLeafNode(searchKey)
	if !find {
		//leaf not exist, delete failed
		return false
	}

	for {
		//delete note from parent node
		for index, _ := range pNode.edges {
			if pNode.edges[index].destNode == lNode {
				pNode.edges = append(pNode.edges[:index], pNode.edges[index+1:]...)
				break
			}
		}

		if len(pNode.edges) != 0 || pNode == &t.root {
			//Stop loop up level condition
			//1: parent node have more than 1 edge after delete
			//2: parent node is root node
			return true
		}

		//delete lNode
		lNode = nil

		//Keep loop up level
		lNode = pNode
		pNode, _ = t.findParent(pNode)
	}

	//return false
}

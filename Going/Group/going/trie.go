package going

import (
	"fmt"
	"strings"
)

// Trie树节点 用于保存路由信息
type node struct {
	pattern  string  //待匹配路由 例如 /p/:lang
	part     string  //路由中的一部分 例如 ：lang
	children []*node //子节点，例如 [doc，tutorial，intro]
	isWild   bool    //是否精确匹配，part含有 : 或 * 时为true
}

// String()方法实现fmt.Stringer接口
func (n *node) String() string { //打印节点信息
	return fmt.Sprintf("*node{pattern=%s,part=%s,isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		//如果没有匹配到当前节点，则新建一个节点
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child) //将新节点加入子节点
	}
	child.insert(pattern, parts, height+1) //递归插入子节点
}

// 递归查找节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	//获取当前层级的part
	part := parts[height]
	children := n.matchChildren(part)

	//遍历子节点，递归查找
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 查找所有匹配节点，用于插入路由时查找冲突
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

// 匹配子节点，用于查找路由
func (n *node) matchChild(part string) *node {
	//遍历子节点，查找匹配的节点
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配所有子节点，用于插入路由时查找冲突
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

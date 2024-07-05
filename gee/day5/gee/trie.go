package gee

import "strings"

// 需要注意pattern是整个路由 如/p/:lang/doc
// 而part是路由的一部分 如:lang
// 一个node其实是路由part的封装
// 只有叶子节点的pattern才有意义 也即只有doc节点 其pattern才是/p/:lang/doc
// 非叶子节点的pattern和handler属性为空
type node struct {
	pattern  string      //待匹配的路由
	handler  HandlerFunc //当前路由的处理函数，只有末节点才有处理函数
	part     string      //路由的一部分
	children []*node     //子节点
	isWild   bool        //是否模糊匹配
}

// 初始化 假的根节点
func NewTrie() *node {
	return &node{}
}

func (node *node) matchChild(part string) *node {
	for _, child := range node.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 将pattern拆为多个part，如p/:lang/doc拆为 [p,:lang,doc]
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 注册一个路径和这个路径的处理函数
func (n *node) insert(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			child = &node{
				part:   part,
				isWild: part[0] == ':' || part[0] == '*',
			}
			n.children = append(n.children, child)
		}
		n = child
	}
	n.pattern = pattern
	n.handler = handler
}

// 搜索 根据URL搜索这个URL对应的叶节点node
// 解释下这里为什么使用matchChildren而非matchChild 假设如下：
// 我们insert一个路由 /p/c/doc，随后又insert一个路由/p/:lang/src
// 这里p节点下肯定是有两个子节点的，分别是c和:lang
// 假设此时一条请求url是/p/c/src
// 如果search的时候选择matchChild而非matchChildren，很容易匹配到/p/c这条路，但这条路叶节点是doc 因此会找不到路由
// 我们知道其实/p/:lang这条路也是符合的，因此应该走这条路，所以需要通过matchChildren查出所有符合的children
func (n *node) search(pattern string) *node {
	parts := parsePattern(pattern)
	children := []*node{n}
	for _, part := range parts {
		children = doSearch(children, part)
	}
	if len(children) == 0 {
		return nil
	}
	return children[0]
}

func doSearch(nodes []*node, part string) []*node {
	matched := make([]*node, 0)
	for _, n := range nodes {
		matched = append(matched, n.matchChildren(part)...)
	}
	return matched
}

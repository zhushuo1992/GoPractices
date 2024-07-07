package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  //全量的url /home/admin/app
	part     string  //当前节点负责的这一层是什么 是home admin 还是app
	children []*node //子层节点
	isWild   bool    //是否是通配符
}

// 打印用
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern //最后1层了
		return
	}

	// home -> admin ->app

	part := parts[height]
	child := n.macthChild(part) //假如有几个路径 /home/admin/app /home/admin/src，先插入前面的，就有可能找到了
	//home有1个child，admin有2个child
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" { //默认值，不是最后1层
			return nil
		}

		return n
	}

	part := parts[height]

	children := n.matchChildren(part) //每个子节点都要层层查找
	//为什么插入不需要匹配这么多？
	//插入是精确的，查找有多种情况，admin和*都能匹配

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil

}

// 获取所有叶子
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) macthChild(part string) *node {
	for _, child := range n.children {
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

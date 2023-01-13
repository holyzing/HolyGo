package leetcode

import (
	"testing"
)

/**
206 反转列表
*/

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	return head
}

func TestMain(t *testing.T) {
	reverseList(&ListNode{})
}

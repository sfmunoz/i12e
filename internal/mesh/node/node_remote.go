package node

import (
	"time"
)

func (n *Node) GetTsFirst() *time.Time {
	return n.tsFirst
}

func (n *Node) SetTsFirst(ts *time.Time) {
	n.tsFirst = ts
}

func (n *Node) GetTsCurr() *time.Time {
	return n.tsCurr
}

func (n *Node) SetTsCurr(ts *time.Time) {
	n.tsCurr = ts
}

func (n *Node) GetAge(tsNow *time.Time) *time.Duration {
	if tsNow == nil {
		ts := time.Now().UTC()
		tsNow = &ts
	}
	tsFirst := n.GetTsFirst()
	if tsFirst == nil {
		return nil
	}
	d := tsNow.Sub(*tsFirst)
	return &d
}

func (n *Node) GetDelta() *time.Duration {
	tsFirst := n.GetTsFirst()
	if tsFirst == nil {
		return nil
	}
	tsCurr := n.GetTsCurr()
	if tsCurr == nil {
		return nil
	}
	d := tsCurr.Sub(*tsFirst)
	return &d
}

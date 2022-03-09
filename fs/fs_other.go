//go:build !linux && !darwin
// +build !linux,!darwin

package fs

import "errors"

type Node interface {
	MountAndServe(mtpt string, debug bool) error
}

type node struct{}

func NewNode() Node {
	return  &node{}
}

func (*node) MountAndServe(mtpt string, debug bool) error {
	return errors.New("not implement")
}

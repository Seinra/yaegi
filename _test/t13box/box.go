package t13box

import "sync"

type Item[T any] struct {
	Value T
	mu    sync.Mutex
}

type Box[T any] struct {
	other *Later
	mu    sync.Mutex
	items []*Item[T]
	head  T
}

type Later struct {
	X int
}

package main

import (
	"errors"
	"hash/crc32"
	"sort"
)

type HashRing struct {
	ring  []uint32
	nodes map[uint32]string
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func NewHashRing() *HashRing {
	return &HashRing{
		ring:  []uint32{},
		nodes: make(map[uint32]string),
	}
}

func (hr *HashRing) AddNode(node string) error {
	if node == "" {
		return errors.New("node cannot be empty")
	}

	for _, existing := range hr.nodes {
		if existing == node {
			return nil
		}
	}

	if existing, ok := hr.nodes[hashKey(node)]; ok {
		return errors.New("hash collision with node: " + existing)
	}

	index := sort.Search(len(hr.ring), func(i int) bool {
		return hr.ring[i] >= hashKey(node)
	})

	hr.ring = append(hr.ring, 0)
	copy(hr.ring[index+1:], hr.ring[index:])
	hr.ring[index] = hashKey(node)
	hr.nodes[hashKey(node)] = node

	return nil
}

func (hr *HashRing) GetNode(key string) (string, error) {
	if len(hr.ring) == 0 {
		return "", errors.New("cannot lookup key: hash ring is empty")
	}

	position := hashKey(key)

	index := sort.Search(len(hr.ring), func(i int) bool {
		return hr.ring[i] >= position
	})

	if index == len(hr.ring) {
		index = 0
	}

	return hr.nodes[hr.ring[index]], nil
}

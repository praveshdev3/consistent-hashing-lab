package main

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
)

type HashRing struct {
	ring         []uint32
	nodes        map[uint32]string
	virtualNodes int
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func NewHashRing(virtualNodes int) *HashRing {
	if virtualNodes <= 0 {
		virtualNodes = 1
	}

	return &HashRing{
		ring:         []uint32{},
		nodes:        make(map[uint32]string),
		virtualNodes: virtualNodes,
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

	for i := 0; i < hr.virtualNodes; i++ {
		virtualName := fmt.Sprintf("%s#%d", node, i)
		position := hashKey(virtualName)

		if existing, ok := hr.nodes[position]; ok {
			return fmt.Errorf(
				"hash collision at position %d with server %s",
				position,
				existing,
			)
		}

		index := sort.Search(len(hr.ring), func(i int) bool {
			return hr.ring[i] >= position
		})

		hr.ring = append(hr.ring, 0)
		copy(hr.ring[index+1:], hr.ring[index:])
		hr.ring[index] = position
		hr.nodes[position] = node
	}
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

func (hr *HashRing) RemoveNode(name string) error {
	if name == "" {
		return errors.New("server name cannot be empty")
	}

	positions := make([]uint32, 0)

	for position, server := range hr.nodes {
		if server == name {
			positions = append(positions, position)
		}
	}

	if len(positions) == 0 {
		return fmt.Errorf("server %s not found in the ring", name)
	}

	for _, position := range positions {
		index := sort.Search(len(hr.ring), func(i int) bool {
			return hr.ring[i] >= position
		})

		if index < len(hr.ring) && hr.ring[index] == position {
			hr.ring = append(hr.ring[:index], hr.ring[index+1:]...)
			delete(hr.nodes, position)
		}
	}
	return nil
}

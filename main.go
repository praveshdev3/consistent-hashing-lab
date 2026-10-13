package main

import "fmt"

func main() {
	ring := NewHashRing()
	servers := []string{"server1", "server2", "server3"}
	for _, server := range servers {
		if err := ring.AddNode(server); err != nil {
			fmt.Printf("Error adding node %s: %v\n", server, err)
		}
	}

	keys := []string{"key1", "key2", "key3", "key4", "key5"}
	for _, key := range keys {
		node, err := ring.GetNode(key)
		if err != nil {
			fmt.Printf("Error getting node for key %s: %v\n", key, err)
		} else {
			fmt.Printf("Key %s is mapped to node %s\n", key, node)
		}
	}

}

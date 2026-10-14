package main

import (
	"fmt"
)

func main() {
	const (
		virtualNodes = 100
		totalKeys    = 100000
	)

	serversBefore := []string{
		"server-A",
		"server-B",
		"server-C",
	}

	serversAfter := []string{
		"server-A",
		"server-B",
		"server-C",
		"server-D",
	}

	keys := make([]string, totalKeys)
	for i := range keys {
		keys[i] = fmt.Sprintf("user:%d", i)
	}

	// Experiment 1: Modulo hashing.
	moduloBefore := make(map[string]string, totalKeys)

	for _, key := range keys {
		index := hashKey(key) % uint32(len(serversBefore))
		moduloBefore[key] = serversBefore[index]
	}

	moduloMoved := 0
	for _, key := range keys {
		index := hashKey(key) % uint32(len(serversAfter))
		newServer := serversAfter[index]

		if moduloBefore[key] != newServer {
			moduloMoved++
		}
	}

	// Experiment 2: Consistent hashing.
	ring := NewHashRing(virtualNodes)

	for _, server := range serversBefore {
		if err := ring.AddNode(server); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	consistentBefore := make(map[string]string, totalKeys)

	for _, key := range keys {
		server, err := ring.GetNode(key)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		consistentBefore[key] = server
	}

	if err := ring.AddNode("server-D"); err != nil {
		fmt.Println("Error:", err)
		return
	}

	consistentMoved := 0
	for _, key := range keys {
		newServer, err := ring.GetNode(key)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if consistentBefore[key] != newServer {
			consistentMoved++
		}
	}

	fmt.Println("Server Addition: Hashing Comparison")
	fmt.Println("------------------------------------")
	fmt.Printf("Keys tested: %d\n\n", totalKeys)

	fmt.Printf("Modulo hashing:      %d keys moved (%.2f%%)\n",
		moduloMoved,
		float64(moduloMoved)/float64(totalKeys)*100)

	fmt.Printf("Consistent hashing:  %d keys moved (%.2f%%)\n",
		consistentMoved,
		float64(consistentMoved)/float64(totalKeys)*100)

	fmt.Println("- Modulo hashing typically remaps many keys.")
	fmt.Println("- Consistent hashing typically moves keys only to the new server.")
}

package main

import (
	"fmt"
	"testing"
)

func TestEmptyRing(t *testing.T) {
	ring := NewHashRing(10)

	_, err := ring.GetNode("user:123")
	if err == nil {
		t.Fatal("expected an error for an empty ring")
	}
}

func TestSingleServer(t *testing.T) {
	ring := NewHashRing(10)

	if err := ring.AddNode("server-A"); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{
		"user:1",
		"user:2",
		"user:3",
	} {
		server, err := ring.GetNode(key)
		if err != nil {
			t.Fatal(err)
		}

		if server != "server-A" {
			t.Fatalf("expected server-A, got %s", server)
		}
	}
}

func TestDeterministicLookup(t *testing.T) {
	ring := NewHashRing(10)

	for _, server := range []string{
		"server-A",
		"server-B",
		"server-C",
	} {
		if err := ring.AddNode(server); err != nil {
			t.Fatal(err)
		}
	}

	first, err := ring.GetNode("user:123")
	if err != nil {
		t.Fatal(err)
	}

	second, err := ring.GetNode("user:123")
	if err != nil {
		t.Fatal(err)
	}

	if first != second {
		t.Fatalf("expected stable lookup, got %s and %s",
			first, second)
	}
}

func TestDuplicateServer(t *testing.T) {
	ring := NewHashRing(10)

	if err := ring.AddNode("server-A"); err != nil {
		t.Fatal(err)
	}

	if err := ring.AddNode("server-A"); err != nil {
		t.Fatal(err)
	}

	for _, key := range []string{"user:1", "user:2"} {
		server, err := ring.GetNode(key)
		if err != nil {
			t.Fatal(err)
		}

		if server != "server-A" {
			t.Fatalf("expected server-A, got %s", server)
		}
	}
}

func TestRemoveNode(t *testing.T) {
	ring := NewHashRing(10)

	for _, server := range []string{
		"server-A",
		"server-B",
		"server-C",
	} {
		if err := ring.AddNode(server); err != nil {
			t.Fatal(err)
		}
	}

	if err := ring.RemoveNode("server-B"); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user:%d", i)

		server, err := ring.GetNode(key)
		if err != nil {
			t.Fatal(err)
		}

		if server == "server-B" {
			t.Fatalf("key %s assigned to removed server", key)
		}
	}
}

func TestRemoveLastNode(t *testing.T) {
	ring := NewHashRing(10)

	if err := ring.AddNode("server-A"); err != nil {
		t.Fatal(err)
	}

	if err := ring.RemoveNode("server-A"); err != nil {
		t.Fatal(err)
	}

	if _, err := ring.GetNode("user:1"); err == nil {
		t.Fatal("expected an error when looking up a key in an empty ring")
	}
}

func TestRemoveUnknownServer(t *testing.T) {
	ring := NewHashRing(10)

	if err := ring.AddNode("server-A"); err != nil {
		t.Fatal(err)
	}

	if err := ring.RemoveNode("server-Z"); err == nil {
		t.Fatal("expected an error when removing an unknown server")
	}
}

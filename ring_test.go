package main

import "testing"

func TestEmptyRing(t *testing.T) {
	ring := NewHashRing()

	_, err := ring.GetNode("user:123")
	if err == nil {
		t.Fatal("expected an error for an empty ring")
	}
}

func TestSingleServer(t *testing.T) {
	ring := NewHashRing()

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
	ring := NewHashRing()

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
	ring := NewHashRing()

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

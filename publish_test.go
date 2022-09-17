package gost

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestPublish(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	imageBytes, err := os.ReadFile("test.png")
	if err != nil {
		t.Errorf("Failed to read test.png: %v", err)
	}
	loc, err := store.Publish(context.Background(), "test.png", "image/png", imageBytes)
	if err != nil {
		t.Errorf("Failed to publish: %v", err)
	}
	t.Logf("Location: %v", loc)
}

func TestUnpublish(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	err = store.Unpublish(context.Background(), "test.png")
	if err != nil {
		t.Errorf("Failed to unpublish: %v", err)
	}

}

func TestDenyPublic(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	err = store.DenyPublic(context.Background())
	if err != nil {
		t.Errorf("Failed to deny public: %v", err)
	}
}

func TestAllowPublic(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	err = store.AllowPublic(context.Background())
	if err != nil {
		t.Errorf("Failed to allow public: %v", err)
	}
}
func TestCheckPublic(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}

	isPublic, err := store.IsPublic(context.Background())
	if err != nil {
		t.Errorf("Failed to check public: %v", err)
	}

	fmt.Println("public: ", isPublic)
}

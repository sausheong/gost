package gost

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
)

var objects []Thingy

func setupObjects() {
	setup()
	objects = []Thingy{
		Thingy{
			Name: "Bob",
			Age:  20,
		},
		Thingy{
			Name: "Alice",
			Age:  21,
		},
	}
}

func TestPutObject(t *testing.T) {
	setupObjects()
	Register([]Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	err = store.PutObject(context.Background(), "some-things", objects)
	if err != nil {
		t.Errorf("Failed to store: %v", err)
	}
}

func TestGetObject2(t *testing.T) {
	setupObjects()
	Register([]Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}

	obj, err := store.GetObject(context.Background(), "some-things")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
		t.Log("Obj is:", obj)
	}
	things := obj.([]Thingy)

	if things[0].Name != "Bob" {
		t.Errorf("Failed to get the right thingy")
	}
	if things[1].Name != "Alice" {
		t.Errorf("Failed to get the right thingy")
	}

}

func TestGetObject(t *testing.T) {
	setupObjects()
	Register([]Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}

	obj, err := store.GetObject(context.Background(), "some-things")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
		errResp := err.(minio.ErrorResponse)
		t.Log(errResp.Code, errResp.Message)
	}
	things := obj.([]Thingy)

	if things[0].Name != "Bob" {
		t.Errorf("Failed to get the right thingy")
	}
	if things[1].Name != "Alice" {
		t.Errorf("Failed to get the right thingy")
	}

}

func TestDeleteObject(t *testing.T) {
	setupObjects()
	Register([]Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}

	err = store.DeleteObject(context.Background(), "some-things")
	if err != nil {
		t.Errorf("Failed to delete: %v", err)
	}

}

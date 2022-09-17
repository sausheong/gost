package gost

import (
	"context"
	"testing"
)

func TestBackup(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	Register(Thingy{})
	err = store.Backup(context.Background(), "sausheong")
	if err != nil {
		t.Errorf("Failed to backup: %v", err)
	}

}

func TestLoad(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	Register(Thingy{})
	data, err := store.Load(context.Background(), "sausheong")
	if err != nil {
		t.Errorf("Failed to list: %v", err)
	}
	if len(data) < 0 {
		t.Errorf("Failed to get the right number of data")
	}
	t.Logf("%#v", data)

}

func TestBackupAndRestore(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, bucket)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}
	Register(Thingy{})

	// backup the data first
	err = store.Backup(context.Background(), "sausheong")
	if err != nil {
		t.Errorf("Failed to backup: %v", err)
	}

	// put in new data
	err = store.Put(context.Background(), "sausheong", "1000", "This is new data")
	if err != nil {
		t.Errorf("Failed to store: %v", err)
	}

	// check if the data is there
	newData, err := store.Get(context.Background(), "sausheong", "1000")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}
	if newData.(string) != "This is new data" {
		t.Errorf("Failed to get the right thing")
	}

	// restore the data
	err = store.Restore(context.Background(), "sausheong")
	if err != nil {
		t.Errorf("Failed to restore: %v", err)
	}

	// data should be back to the original, newData should be gone
	newData, err = store.Get(context.Background(), "sausheong", "1000")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}
	if newData != nil {
		t.Errorf("newData should be nil:, %v", newData)
	}

}

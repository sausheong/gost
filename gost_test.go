package gost

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var policy string = `{
	"Version": "2012-10-17",
	"Statement": [
		{
		"Effect": "Allow",
		"Principal": "*",
		"Action": "s3:GetObject",
		"Resource": "arn:aws:s3:::gost/*",
		"Sid": ""
		}
	]
}`

type Thingy struct {
	Name        string
	Age         int
	DateCreated time.Time
	Length      float64
	Bunch       []OtherThingy
}

type OtherThingy struct {
	Name   string
	Number int
}

var key, secret, endpoint, region, bucket string
var useSSL bool

func setup() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load the env vars: %v", err)
	}
	key = os.Getenv("KEY")
	secret = os.Getenv("SECRET")
	endpoint = os.Getenv("ENDPOINT")
	useSSL, err = strconv.ParseBool((os.Getenv("USE_SSL")))
	if err != nil {
		log.Fatalf("Failed to parse USE_SSL: %v", err)
	}
	region = os.Getenv("REGION")
	bucket = os.Getenv("BUCKET")
}

func TestPutBasic(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	err = store.Put(context.Background(), "sausheong", "123", "hello world!")
	if err != nil {
		t.Errorf("Failed to store: %v", err)
	}
}

func TestPut(t *testing.T) {
	setup()
	thingy := Thingy{
		Name:        "Bob",
		Age:         42,
		DateCreated: time.Now(),
		Length:      1.234,
		Bunch: []OtherThingy{
			{
				Name:   "Alice",
				Number: 1,
			},
			{
				Name:   "Bob",
				Number: 2,
			},
		},
	}
	Register(thingy)
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	err = store.Put(context.Background(), "sausheong", "Bob", thingy)
	if err != nil {
		t.Errorf("Failed to store: %v", err)
	}
}

func TestGetBasic(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	thing, err := store.Get(context.Background(), "sausheong", "123")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}
	if thing.(string) != "hello world!" {
		t.Errorf("Failed to get the right thing")
	}
}

func TestGet(t *testing.T) {
	setup()
	Register(Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	thing, err := store.Get(context.Background(), "sausheong", "Bob")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}
	if thing.(Thingy).Name != "Bob" {
		t.Errorf("Failed to get the right thingy")
	}
}

func TestGetAll(t *testing.T) {
	setup()
	Register(Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	all, err := store.GetAll(context.Background(), "sausheong")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}
	if all["Bob"].(Thingy).Name != "Bob" {
		t.Errorf("Failed to get the right thingy")
	}
	if all["123"].(string) != "hello world!" {
		t.Errorf("Failed to get the right thing")
	}
}

func TestDelete(t *testing.T) {
	setup()
	Register(Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	err = store.Delete(context.Background(), "sausheong", "123")
	if err != nil {
		t.Errorf("Failed to delete: %v", err)
	}
}

func TestDeleteAll(t *testing.T) {
	setup()
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	err = store.DeleteAll(context.Background(), "sausheong")
	if err != nil {
		t.Errorf("Failed to delete all: %v", err)
	}
}

func TestPutImageFile(t *testing.T) {
	setup()
	imageBytes, err := os.ReadFile("test.png")
	if err != nil {
		t.Errorf("Failed to read test.png: %v", err)
	}
	Register(Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	err = store.Put(context.Background(), "sausheong", "test.png", imageBytes)
	if err != nil {
		t.Errorf("Failed to store: %v", err)
	}
}

func TestGetImageFile(t *testing.T) {
	setup()
	Register(Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	image, err := store.Get(context.Background(), "sausheong", "test.png")
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}

	// write image file
	err = os.WriteFile("test2.png", image.([]byte), 0644)
	if err != nil {
		t.Errorf("Failed to write test2.png: %v", err)
	}

}

func TestLargeFile(t *testing.T) {
	setup()
	Register(Thingy{})
	store, err := NewStore(key, secret, endpoint, useSSL, region, bucket)
	if err != nil {
		t.Errorf("Failed to initialise a store: %v", err)
	}
	// dataBytes, err := os.ReadFile("100MB.bin")
	// if err != nil {
	// 	t.Errorf("Failed to read big file: %v", err)
	// }
	err = store.Put(context.Background(), "large-2", "big-file-2", "123")
	if err != nil {
		t.Errorf("Failed to store: %v", err)
	}
}

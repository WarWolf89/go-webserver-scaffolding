package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fruits = []Fruit{
	{ID: "1", Name: "apple", Color: "red"},
	{ID: "2", Name: "orange", Color: "orange"},
	{ID: "3", Name: "pineapple", Color: "yellow"},
}

var expected_fruits = []Fruit{
	{ID: "1", Name: "apple", Color: "red"},
	{ID: "2", Name: "orange", Color: "orange"},
	{ID: "3", Name: "pineapple", Color: "yellow"},
}

var sample_fruit_new = Fruit{
	ID:    "4",
	Name:  "test",
	Color: "test",
}

var sample_fruit_existing = Fruit{
	ID:    "1",
	Name:  "apple",
	Color: "red",
}

func TestGetAllFruits(t *testing.T) {
	r := setupClients()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/fruits", nil)
	r.ServeHTTP(w, req)

	actual_fruits := []fruit{}
	if err := json.Unmarshal(w.Body.Bytes(), &actual_fruits); err != nil {
		fmt.Println(err)
		panic(err)
	}

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected_fruits, actual_fruits)
}

func TestAddFruit(t *testing.T) {
	r := setupClients()
	b, _ := json.Marshal(sample_fruit_new)
	req, _ := http.NewRequest("POST", "/fruits", bytes.NewBuffer(b))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := fruit{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		fmt.Println(err)
		panic(err)
	}
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, sample_fruit_new, resp)
}

func TestGetFruit(t *testing.T) {
	w := httptest.NewRecorder()
	r := setupClients()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/fruits/%s", sample_fruit_existing.ID), nil)

	r.ServeHTTP(w, req)
	resp := fruit{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		fmt.Println(err)
		panic(err)
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, sample_fruit_existing, resp)
}

func TestGetFruitNone(t *testing.T) {
	w := httptest.NewRecorder()
	r := setupClients()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/fruits/%s", "8"), nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

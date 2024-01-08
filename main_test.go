package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"

	"csaba.almasi.per/webserver/src/pkg/fruitservice"
	"csaba.almasi.per/webserver/src/pkg/fruitservice/api"
	"csaba.almasi.per/webserver/src/pkg/fruitservice/fruitstore"
)

type FruitTestSuite struct {
	suite.Suite
	api *api.Api
	w   *httptest.ResponseRecorder
	ctx context.Context
}

// Set Up testify suite
func (s *FruitTestSuite) SetupSuite() {
	validate := validator.New(validator.WithRequiredStructEnabled())
	s.w = httptest.NewRecorder()
	s.api = api.ProvideApi(gin.Default(), fruitstore.ProvideSVC(), validate)
	s.api.RegisterAPIEndpoints()
	gin.SetMode(gin.TestMode)

}

// set context before every test
func (s *FruitTestSuite) SetupTest() {
	// type assert for Redis specific client
	rc := s.api.Fsvc.(*fruitstore.RedisStore)

	// Create a gin context per test
	s.ctx = gin.CreateTestContextOnly(s.w, s.api.Gengine)

	// Flush the Redis instance before every test
	rc.Client.FlushAll(s.ctx)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFruitTestSuite(t *testing.T) {
	suite.Run(t, new(FruitTestSuite))
}

var expected_fruits = []fruitservice.Fruit{}

var sample_fruit_new = fruitservice.Fruit{
	ID:    "4",
	Name:  "test",
	Color: "test",
}

var sample_fruit_existing = fruitservice.Fruit{
	ID:    "1",
	Name:  "apple",
	Color: "red",
}

func (s *FruitTestSuite) TestGetAllFruits() {

	_, err := s.api.Fsvc.AddFruit(s.ctx, &sample_fruit_existing)
	if err != nil {
		s.Error(err)
	}

	expected_fruits = append(expected_fruits, sample_fruit_existing)

	req, err := http.NewRequest("GET", "/api/v1/fruits", nil)
	if err != nil {
		s.Error(err)
	}
	s.api.Gengine.ServeHTTP(s.w, req)

	actual_fruits := []fruitservice.Fruit{}
	if err := json.Unmarshal(s.w.Body.Bytes(), &actual_fruits); err != nil {
		fmt.Println(err)
		panic(err)
	}

	s.Equal(200, s.w.Code)
	s.Equal(expected_fruits, actual_fruits)
}

//

// func TestGetFruit(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	r := setupClients()

// 	req, _ := http.NewRequest("GET", fmt.Sprintf("/fruits/%s", sample_fruit_existing.ID), nil)

// 	r.ServeHTTP(w, req)
// 	resp := fruit{}
// 	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	assert.Equal(t, 200, w.Code)
// 	assert.Equal(t, sample_fruit_existing, resp)
// }

// func TestGetFruitNone(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	r := setupClients()

// 	req, _ := http.NewRequest("GET", fmt.Sprintf("/fruits/%s", "8"), nil)

// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, 404, w.Code)
// }

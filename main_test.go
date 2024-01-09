package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"

	"csaba.almasi.per/webserver/src/pkg/fruitservice"
	"csaba.almasi.per/webserver/src/pkg/fruitservice/api"
	"csaba.almasi.per/webserver/src/pkg/fruitservice/fruitstore"
	"csaba.almasi.per/webserver/src/pkg/util"
)

var sample_fruit_new = fruitservice.Fruit{
	Name:  "orange",
	Color: "orange",
}

var sample_fruit_existing = fruitservice.Fruit{
	Name:  "apple",
	Color: "red",
}

var faulty_fruit = fruitservice.Fruit{
	Name:  "89897899",
	Color: "67854687",
}

type FruitTestSuite struct {
	suite.Suite
	api *api.Api
	w   *httptest.ResponseRecorder
	ctx context.Context
}

// Set Up testify suite
func (s *FruitTestSuite) SetupSuite() {
	config, err := util.LoadConfig("test_config")
	if err != nil {
		s.Error(err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	s.api = api.ProvideApi(gin.Default(), fruitstore.ProvideRedisStore(config), validate)
	s.api.RegisterAPIEndpoints()
	gin.SetMode(gin.TestMode)

}

// set context before every test
func (s *FruitTestSuite) SetupTest() {
	// type assert for Redis specific client
	rc := s.api.Fsvc.(*fruitstore.RedisStore)

	// Create a gin context per test
	s.ctx = gin.CreateTestContextOnly(s.w, s.api.Gengine)

	// New to create a new recorder before every run since sporadic memory issues pop up
	s.w = httptest.NewRecorder()

	// Flush the Redis instance before every test
	rc.Client.FlushAll(s.ctx)

}

func (s *FruitTestSuite) TearDownSuite() {
	rc := s.api.Fsvc.(*fruitstore.RedisStore)

	// Flush the Redis instance after tests
	rc.Client.FlushAll(s.ctx)

	// Close connection
	if err := s.api.Fsvc.(*fruitstore.RedisStore).Client.Conn().Close(); err != nil {
		s.Error(err)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFruitTestSuite(t *testing.T) {
	suite.Run(t, new(FruitTestSuite))
}

func (s *FruitTestSuite) TestGetAllFruits() {

	expected_fruits := []fruitservice.Fruit{}
	actual_fruits := []fruitservice.Fruit{}

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

	if err := json.Unmarshal(s.w.Body.Bytes(), &actual_fruits); err != nil {
		s.Error(err)
	}

	s.Equal(200, s.w.Code)
	s.Equal(expected_fruits, actual_fruits)
}

func (s *FruitTestSuite) TestGetFruitByID() {
	id, err := s.api.Fsvc.AddFruit(s.ctx, &sample_fruit_existing)
	if err != nil {
		s.Error(err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/fruits:%s", id), nil)
	if err != nil {
		s.Error(err)
	}
	sample_fruit_existing.ID = id

	s.api.Gengine.ServeHTTP(s.w, req)

	var actual_fruit fruitservice.Fruit

	if err := json.Unmarshal(s.w.Body.Bytes(), &actual_fruit); err != nil {
		fmt.Println(err)
		panic(err)
	}

	s.Equal(200, s.w.Code)
	s.Equal(sample_fruit_existing, actual_fruit)

}

func (s *FruitTestSuite) TestAddFruit() {

	rb, err := json.Marshal(sample_fruit_new)
	if err != nil {
		s.Error(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/fruits", bytes.NewBuffer(rb))
	if err != nil {
		s.Error(err)
	}
	req.Header.Set("Content-type", "application/json")

	s.api.Gengine.ServeHTTP(s.w, req)

	// Get the ID from the header of the response
	id := strings.ReplaceAll(s.w.Result().Header.Get("Location"), "v1/api/fruits", "")

	// Make a direct call to Redis and see if fruit was actually created with the values we passed
	actual, err := s.api.Fsvc.GetFruitByID(s.ctx, id)
	if err != nil {
		s.Error(err)
	}

	s.Equal(201, s.w.Result().StatusCode)
	s.Equal(fmt.Sprintf("v1/api/fruits%s", id), s.w.Result().Header.Get("Location"))
	s.Equal(sample_fruit_new.Name, actual.Name)
	s.Equal(sample_fruit_new.Color, actual.Color)

}

func (s *FruitTestSuite) TestAddFaultyFruit() {
	rb, err := json.Marshal(faulty_fruit)
	if err != nil {
		s.Error(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/fruits", bytes.NewBuffer(rb))
	if err != nil {
		s.Error(err)
	}
	req.Header.Set("Content-type", "application/json")

	s.api.Gengine.ServeHTTP(s.w, req)

	s.Equal(400, s.w.Result().StatusCode)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/julienschmidt/httprouter"
)

func main() {
	customResolver := func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               "http://localstack:4566",
			HostnameImmutable: true,
			SigningRegion:     region,
		}, nil
	}
	s := server{
		router: httprouter.New(),
		s3: s3.NewFromConfig(aws.Config{
			Region:           "us-east-1",
			EndpointResolver: aws.EndpointResolverFunc(customResolver),
		}),
	}

	s.routes()

	log.Fatal(http.ListenAndServe(":8080", s.router))
}

type server struct {
	router *httprouter.Router // route the requests
	s3     *s3.Client         // talk to s3
}

func (s *server) routes() {
	s.router.GET("/greeting/:name", s.handleGetGreeting())
	s.router.POST("/greeting", s.handleSetGreeting())
}

func (s *server) handleGetGreeting() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Request object with key :name under the bucket greetings in s3
		input := &s3.GetObjectInput{
			Bucket: aws.String("greetings"),
			Key:    aws.String(ps.ByName("name")),
		}
		obj, err := s.s3.GetObject(r.Context(), input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Something went wrong: %s", err.Error())))
			return
		}

		// Read the file
		b, err := io.ReadAll(obj.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Something went wrong: %s", err.Error())))
			return
		}

		// Write the file contents
		w.Write(b)
	}
}

func (s *server) handleSetGreeting() httprouter.Handle {
	// Request is what we expect from the /greeting endpoint
	type Request struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Read the body to bytes
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Something went wrong: %s", err.Error())))
			return
		}

		// Unmarshal the body into the Request struct
		var req Request
		if err := json.Unmarshal(b, &req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Something went wrong: %s", err.Error())))
			return
		}

		// Store the greeting under a key of name in the greetings bucket
		input := &s3.PutObjectInput{
			Bucket: aws.String("greetings"),
			Key:    aws.String(req.Name),
			Body:   strings.NewReader(fmt.Sprintf("Hello, %s!", req.Name)),
		}

		if _, err := s.s3.PutObject(r.Context(), input); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Something went wrong: %s", err.Error())))
			return
		}
	}
}

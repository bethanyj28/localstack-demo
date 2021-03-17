package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/julienschmidt/httprouter"
)

func main() {
	s := server{
		router: httprouter.New(),
		s3:     s3.New(s3.Options{Region: "us-east-1"}),
	}

	s.routes()

	log.Fatal(http.ListenAndServe(":8080", s.router))
}

type server struct {
	router *httprouter.Router
	s3     *s3.Client
}

func (s *server) routes() {
	s.router.GET("/greeting/:name", s.handleGetGreeting())
	s.router.POST("/greeting", s.handleSetGreeting())
}

func (s *server) handleGetGreeting() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		input := &s3.GetObjectInput{
			Bucket: aws.String("greetings"),
			Key:    aws.String(ps.ByName("name")),
		}
		obj, err := s.s3.GetObject(r.Context(), input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}

		b, err := io.ReadAll(obj.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) handleSetGreeting() httprouter.Handle {
	type Request struct {
		Name string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		body, err := r.GetBody()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}
		defer body.Close()

		b, err := io.ReadAll(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}

		var req Request
		if err := json.Unmarshal(b, &req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}
	}
}

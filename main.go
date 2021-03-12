package main

import (
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
	s.router.GET("/", s.handleRenderHTML())
}

func (s *server) handleRenderHTML() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		input := &s3.GetObjectInput{
			Bucket: aws.String(""),
			Key:    aws.String(""),
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

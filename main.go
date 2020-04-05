package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var muxLambda *gorillamux.GorillaMuxAdapter

// init is automatically called by Lambda to setup global state, rather than per-request
func init() {
	r := setupRouter()
	muxLambda = gorillamux.New(r)
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello, world!"))
	})

	return r
}

// Handler implements the AWS Lambda pattern
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return muxLambda.ProxyWithContext(ctx, req)
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "test" {
		log.Println("Starting headless...")

		if err := godotenv.Load(); err != nil {
			log.Println("Failed to load .env file: ", err.Error())
		}

		r := setupRouter()

		srv := &http.Server{
			Handler:      r,
			Addr:         "127.0.0.1:8000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Println("Listening on ", srv.Addr, "...")

		log.Fatal(srv.ListenAndServe())
		return
	}

	lambda.Start(Handler)
}

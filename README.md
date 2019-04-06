# quiz-app

Server for QUIZ

## Build

    # Clone
    git clone https://github.com/maviucak/quizApp.git
    cd quizApp

    export GO111MODULE=on

    # Generate assets
    go generate ./cmd/quiz/assets/gen.go 

    # Compile
    go build ./cmd/quiz

## Run
   
    Set up env vars, then:
   
    ./quiz

## Options

    Usage of ./quiz:
      -listen string
            Listen addr (default ":8181")
      -log string
            Log level (debug, info, warn, error) (default "debug")
            
## Generate

Generated API spec (openapi.yml) is served at `/api/v1/assets/openapi.yml`

#### Go Vet

    go vet ./...
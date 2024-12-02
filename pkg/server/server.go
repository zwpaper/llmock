package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	Address string
	Port    int

	*http.Server
}

func New(address string, port int) *Server {
	server := &Server{
		Address: address,
		Port:    port,

		Server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", address, port),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/embeddings", server.handleEmbeddings)
	server.Handler = mux

	return server
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	fmt.Printf("Starting server on %s:%d\n", s.Address, s.Port)
	return s.ListenAndServe()
}

type EmbeddingsOpenAIResponse struct {
	Object string                         `json:"object"`
	Model  string                         `json:"model"`
	Data   []EmbeddingsOpenAIResponseData `json:"data"`
	Usage  EmbeddingsOpenAIResponseUsage  `json:"usage"`
}

type EmbeddingsOpenAIResponseData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingsOpenAIResponseUsage struct {
	TotalTokens  int `json:"total_tokens"`
	PromptTokens int `json:"prompt_tokens"`
}

func (s *Server) handleEmbeddings(w http.ResponseWriter, r *http.Request) {
	time.Sleep(500 * time.Millisecond)

	response := EmbeddingsOpenAIResponse{
		Object: "list",
		Model:  "gpt-3.5-turbo-instruct",
		Data: []EmbeddingsOpenAIResponseData{
			{
				Object:    "embedding",
				Embedding: generateRandomSlice(64),
				Index:     0,
			},
		},
		Usage: EmbeddingsOpenAIResponseUsage{
			TotalTokens:  64,
			PromptTokens: 32,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateRandomSlice(length int) []float64 {
	slice := make([]float64, length)
	for i := range slice {
		slice[i] = rand.Float64()
	}
	return slice
}

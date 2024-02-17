package chromem

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const BaseURLOpenAI = "https://api.openai.com/v1"

type EmbeddingModelOpenAI string

const (
	EmbeddingModelOpenAI2Ada   EmbeddingModelOpenAI = "text-embedding-ada-002"
	EmbeddingModelOpenAI3Small EmbeddingModelOpenAI = "text-embedding-3-small"
	EmbeddingModelOpenAI3Large EmbeddingModelOpenAI = "text-embedding-3-large"
)

type openAIResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// NewEmbeddingFuncDefault returns a function that creates embeddings for a document
// using OpenAI`s "text-embedding-3-small" model via their API.
// The model supports a maximum document length of 8191 tokens.
// The API key is read from the environment variable "OPENAI_API_KEY".
func NewEmbeddingFuncDefault() EmbeddingFunc {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return NewEmbeddingFuncOpenAI(apiKey, EmbeddingModelOpenAI3Small)
}

// NewEmbeddingFuncOpenAI returns a function that creates embeddings for a document
// using the OpenAI API.
func NewEmbeddingFuncOpenAI(apiKey string, model EmbeddingModelOpenAI) EmbeddingFunc {
	return NewEmbeddingFuncOpenAICompat(BaseURLOpenAI, apiKey, string(model))
}

// NewEmbeddingFuncOpenAICompat returns a function that creates embeddings for a document
// using an OpenAI compatible API. For example:
//   - Azure OpenAI: https://azure.microsoft.com/en-us/products/ai-services/openai-service
//   - LitLLM: https://github.com/BerriAI/litellm
//   - Ollama: https://github.com/ollama/ollama/blob/main/docs/openai.md
//   - etc.
func NewEmbeddingFuncOpenAICompat(baseURL, apiKey, model string) EmbeddingFunc {
	// We don't set a default timeout here, although it's usually a good idea.
	// In our case though, the library user can set the timeout on the context,
	// and it might have to be a long timeout, depending on the document size.
	client := &http.Client{}

	return func(ctx context.Context, document string) ([]float32, error) {
		// Prepare the request body.
		reqBody, err := json.Marshal(map[string]string{
			"input": document,
			"model": model,
		})
		if err != nil {
			return nil, fmt.Errorf("couldn't marshal request body: %w", err)
		}

		// Create the request. Creating it with context is important for a timeout
		// to be possible, because the client is configured without a timeout.
		req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/embeddings", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("couldn't create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// Send the request.
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("couldn't send request: %w", err)
		}
		defer resp.Body.Close()

		// Check the response status.
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("error response from the embedding API: " + resp.Status)
		}

		// Read and decode the response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("couldn't read response body: %w", err)
		}
		var embeddingResponse openAIResponse
		err = json.Unmarshal(body, &embeddingResponse)
		if err != nil {
			return nil, fmt.Errorf("couldn't unmarshal response body: %w", err)
		}

		// Check if the response contains embeddings.
		if len(embeddingResponse.Data) == 0 || len(embeddingResponse.Data[0].Embedding) == 0 {
			return nil, errors.New("no embeddings found in the response")
		}

		return embeddingResponse.Data[0].Embedding, nil
	}
}
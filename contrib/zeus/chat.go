package zeus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-kratos/blades"
)

// ChatProvider implements blades.ModelProvider for Zeus API.
type ChatProvider struct {
	client      *http.Client
	apiKey      string
	baseURL     string
	pipelineID  string
}

// NewChatProvider constructs a Zeus provider. The API key is read from
// the ZEUS_API_KEY environment variable. The pipeline ID is read from
// ZEUS_PIPELINE_ID environment variable.
func NewChatProvider() blades.ModelProvider {
	apiKey := os.Getenv("ZEUS_API_KEY")
	if apiKey == "" {
		panic("ZEUS_API_KEY environment variable is required for Zeus provider")
	}
	
	baseURL := os.Getenv("ZEUS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.zeusllm.com/v1"
	}
	
	pipelineID := os.Getenv("ZEUS_PIPELINE_ID")
	if pipelineID == "" {
		panic("ZEUS_PIPELINE_ID environment variable is required for Zeus provider")
	}

	return &ChatProvider{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:     apiKey,
		baseURL:    baseURL,
		pipelineID: pipelineID,
	}
}

// Generate executes a non-streaming chat completion request.
func (p *ChatProvider) Generate(ctx context.Context, req *blades.ModelRequest, opts ...blades.ModelOption) (*blades.ModelResponse, error) {
	// Convert Blades request to Zeus API format
	zeusReq := p.convertToZeusRequest(req)
	
	// Make HTTP request
	resp, err := p.makeRequest(ctx, zeusReq)
	if err != nil {
		return nil, err
	}
	
	// Convert Zeus response to Blades format
	return p.convertFromZeusResponse(resp)
}

// NewStream executes a streaming chat completion request.
func (p *ChatProvider) NewStream(ctx context.Context, req *blades.ModelRequest, opts ...blades.ModelOption) (blades.Streamer[*blades.ModelResponse], error) {
	// For now, implement as non-streaming since Zeus API doesn't show streaming support
	// in the provided example. This can be enhanced later if streaming is supported.
	
	pipe := blades.NewStreamPipe[*blades.ModelResponse]()
	pipe.Go(func() error {
		response, err := p.Generate(ctx, req, opts...)
		if err != nil {
			return err
		}
		pipe.Send(response)
		return nil
	})
	
	return pipe, nil
}

// convertToZeusRequest converts Blades ModelRequest to Zeus API format
func (p *ChatProvider) convertToZeusRequest(req *blades.ModelRequest) map[string]interface{} {
	messages := make([]map[string]string, 0, len(req.Messages))
	
	for _, msg := range req.Messages {
		role := string(msg.Role)
		// Ensure proper role mapping
		switch role {
		case "assistant":
			role = "assistant"
		case "user":
			role = "user"
		case "system":
			role = "system"
		default:
			role = "user" // Default to user if unknown
		}
		
		// Extract text content from parts
		content := ""
		for _, part := range msg.Parts {
			if textPart, ok := part.(blades.TextPart); ok {
				content += textPart.Text
			}
		}
		
		// Only add message if it has content
		if content != "" {
			messages = append(messages, map[string]string{
				"role":    role,
				"content": content,
			})
		}
	}
	
	return map[string]interface{}{
		"messages":     messages,
		"pipeline_id":  p.pipelineID,
	}
}

// convertFromZeusResponse converts Zeus API response to Blades ModelResponse
func (p *ChatProvider) convertFromZeusResponse(resp *ZeusResponse) (*blades.ModelResponse, error) {
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in Zeus response")
	}
	
	choice := resp.Choices[0]
	content := choice.Message.Content
	
	return &blades.ModelResponse{
		Messages: []*blades.Message{
			{
				Role:   blades.RoleAssistant,
				Status: blades.StatusCompleted,
				Parts: []blades.Part{
					blades.TextPart{Text: content},
				},
			},
		},
	}, nil
}

// makeRequest makes HTTP request to Zeus API
func (p *ChatProvider) makeRequest(ctx context.Context, req map[string]interface{}) (*ZeusResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Debug: Print request (remove in production)
	// fmt.Printf("Zeus Request: %s\n", string(jsonData))
	
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/ai", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Zeus API error: %d - %s", resp.StatusCode, string(body))
	}
	
	var zeusResp ZeusResponse
	if err := json.NewDecoder(resp.Body).Decode(&zeusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Debug: Print response (remove in production)
	// fmt.Printf("Zeus Response: %+v\n", zeusResp)
	
	return &zeusResp, nil
}

// ZeusResponse represents the response from Zeus API
type ZeusResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Message      struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Created int64 `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

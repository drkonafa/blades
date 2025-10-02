package gemini

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/go-kratos/blades"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	// ErrEmptyResponse indicates the provider returned no candidates.
	ErrEmptyResponse = errors.New("empty completion response")
	// ErrToolNotFound indicates a tool call was made to an unknown tool.
	ErrToolNotFound = errors.New("tool not found")
)

// ChatProvider implements blades.ModelProvider for Gemini models.
type ChatProvider struct {
	client *genai.Client
}

// NewChatProvider constructs a Gemini provider. The API key is read from
// the API_KEY environment variable.
func NewChatProvider() blades.ModelProvider {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required for Gemini provider")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	return &ChatProvider{client: client}
}

// Generate executes a non-streaming chat completion request.
func (p *ChatProvider) Generate(ctx context.Context, req *blades.ModelRequest, opts ...blades.ModelOption) (*blades.ModelResponse, error) {
	model := p.client.GenerativeModel(req.Model)

	// Convert messages to Gemini format
	var parts []genai.Part
	for _, msg := range req.Messages {
		for _, part := range msg.Parts {
			switch v := part.(type) {
			case blades.TextPart:
				parts = append(parts, genai.Text(v.Text))
			case blades.FilePart:
				// Handle file parts - for now, just add as text
				parts = append(parts, genai.Text("File: "+v.URI))
			case blades.DataPart:
				// Handle data parts - for now, just add as text
				parts = append(parts, genai.Text("Data: "+string(v.Bytes)))
			}
		}
	}

	// Generate content
	resp, err := model.GenerateContent(ctx, parts...)
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) == 0 {
		return nil, ErrEmptyResponse
	}

	// Convert response to Blades format
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			responseText += string(text)
		}
	}

	response := &blades.ModelResponse{
		Messages: []*blades.Message{
			{
				Role:   blades.RoleAssistant,
				Status: blades.StatusCompleted,
				Parts: []blades.Part{
					blades.TextPart{Text: responseText},
				},
			},
		},
	}

	return response, nil
}

// NewStream executes a streaming chat completion request.
func (p *ChatProvider) NewStream(ctx context.Context, req *blades.ModelRequest, opts ...blades.ModelOption) (blades.Streamer[*blades.ModelResponse], error) {
	model := p.client.GenerativeModel(req.Model)

	// Convert messages to Gemini format
	var parts []genai.Part
	for _, msg := range req.Messages {
		for _, part := range msg.Parts {
			switch v := part.(type) {
			case blades.TextPart:
				parts = append(parts, genai.Text(v.Text))
			case blades.FilePart:
				parts = append(parts, genai.Text("File: "+v.URI))
			case blades.DataPart:
				parts = append(parts, genai.Text("Data: "+string(v.Bytes)))
			}
		}
	}

	// Generate content with streaming
	iter := model.GenerateContentStream(ctx, parts...)

	pipe := blades.NewStreamPipe[*blades.ModelResponse]()
	pipe.Go(func() error {
		var fullText string
		for {
			resp, err := iter.Next()
			if err != nil {
				return err
			}

			if len(resp.Candidates) > 0 {
				for _, part := range resp.Candidates[0].Content.Parts {
					if text, ok := part.(genai.Text); ok {
						fullText += string(text)

						// Send incremental response
						response := &blades.ModelResponse{
							Messages: []*blades.Message{
								{
									Role:   blades.RoleAssistant,
									Status: blades.StatusIncomplete,
									Parts: []blades.Part{
										blades.TextPart{Text: string(text)},
									},
								},
							},
						}
						pipe.Send(response)
					}
				}
			}
		}
	})

	return pipe, nil
}

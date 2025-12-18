package service

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

type GeminiService struct {
	client    *genai.Client
	modelName string
	projectID string
	location  string
}

// NewGeminiService: クライアントを初期化します
func NewGeminiService(ctx context.Context, projectID, location, modelName string) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiService{
		client:    client,
		modelName: modelName, // "gemini-2.5-flash"
		projectID: projectID,
		location:  location,
	}, nil
}

// GenerateDescription: プロンプトを送信して生成結果を返します
func (s *GeminiService) GenerateDescription(ctx context.Context, promptText string) (string, error) {
	model := s.client.GenerativeModel(s.modelName)

	// パラメータ調整 (必要に応じて変更)
	model.SetTemperature(0.7) // 創造性 (0.0 - 1.0)

	resp, err := model.GenerateContent(ctx, genai.Text(promptText))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}

	// レスポンスからテキスト部分を抽出
	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			sb.WriteString(string(txt))
		}
	}

	return sb.String(), nil
}

// Close: アプリ終了時にクライアントを閉じます
func (s *GeminiService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

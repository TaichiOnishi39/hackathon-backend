package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"hackathon-backend/model"
	"hackathon-backend/service"
	"strings"
)

type ProductDescriptionUsecase struct {
	GeminiService *service.GeminiService
}

func NewProductDescriptionUsecase(gService *service.GeminiService) *ProductDescriptionUsecase {
	return &ProductDescriptionUsecase{
		GeminiService: gService,
	}
}

// Generate: 商品名とキーワードから説明文を生成 (テキストのみ)
func (u *ProductDescriptionUsecase) Generate(ctx context.Context, name, keywords string) (string, error) {
	prompt := fmt.Sprintf(`
あなたはプロのコピーライターです。フリマアプリに出品するための魅力的な商品説明文を書いてください。

【商品名】
%s

【特徴・キーワード】
%s

【条件】
- ターゲットが欲しくなるような文章にする
- 商品の状態や魅力が伝わるようにする
- 丁寧語（です・ます調）を使う
- 200文字以内で簡潔にまとめる

【重要：出力形式】
- 生成された説明文のみを出力してください。
- 「はい、承知いたしました」や「以下の通り作成しました」などの挨拶や前置きは一切不要です。
`, name, keywords)

	return u.GeminiService.GenerateDescription(ctx, prompt)
}

// GenerateInfoFromImage: 画像から商品情報を抽出 (マルチモーダル)
func (u *ProductDescriptionUsecase) GenerateInfoFromImage(ctx context.Context, imgData []byte, mimeType string) (*model.GenerateImageRes, error) {
	// プロンプト：JSON形式での出力を強制します
	prompt := `
この商品画像を解析し、フリマアプリ出品用の情報をJSON形式で出力してください。
以下のキーを含めてください：
- "name": 商品名 (30文字以内)
- "price": 推定価格 (数値のみ、円は不要。不明なら3000など適当な値)
- "keywords": 特徴を表すキーワード (カンマ区切り)
- "description": 魅力的な商品説明文 (150文字程度。丁寧語で)

出力はJSONのみにしてください。Markdownのコードブロックは不要です。
`

	// 1. Gemini呼び出し
	respText, err := u.GeminiService.GenerateFromImage(ctx, prompt, imgData, mimeType)
	if err != nil {
		return nil, err
	}

	// 2. JSONの整形（```json 等が含まれていた場合の除去）
	respText = strings.TrimSpace(respText)
	respText = strings.TrimPrefix(respText, "```json")
	respText = strings.TrimPrefix(respText, "```")
	respText = strings.TrimSuffix(respText, "```")

	// 3. JSONパース用の一時構造体
	// (Geminiが price を文字列で返してくる場合と数値で返してくる場合の両方に対応するため interface{} で受けます)
	type AiResponse struct {
		Name        string      `json:"name"`
		Price       interface{} `json:"price"`
		Keywords    string      `json:"keywords"`
		Description string      `json:"description"`
	}

	var aiRes AiResponse
	if err := json.Unmarshal([]byte(respText), &aiRes); err != nil {
		// パース失敗時はエラーを返す（あるいはログに出して空の結果を返すなど）
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// 4. Priceをintに安全に変換
	var priceInt int
	switch v := aiRes.Price.(type) {
	case float64:
		priceInt = int(v)
	case string:
		// "3,000" や "3000円" などの表記揺れ対策として、数字以外を除去してから変換しても良いですが
		// ここではシンプルに fmt.Sscanf で試みます
		fmt.Sscanf(v, "%d", &priceInt)
	case int:
		priceInt = v
	}

	// 5. 結果を返す
	return &model.GenerateImageRes{
		Name:        aiRes.Name,
		Price:       priceInt,
		Keywords:    aiRes.Keywords,
		Description: aiRes.Description,
	}, nil
}

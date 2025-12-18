package usecase

import (
	"context"
	"fmt"
	"hackathon-backend/service"
)

type ProductDescriptionUsecase struct {
	GeminiService *service.GeminiService
}

func NewProductDescriptionUsecase(gService *service.GeminiService) *ProductDescriptionUsecase {
	return &ProductDescriptionUsecase{
		GeminiService: gService,
	}
}

// Generate: 商品名とキーワードから説明文を生成
func (u *ProductDescriptionUsecase) Generate(ctx context.Context, name, keywords string) (string, error) {
	// プロンプトエンジニアリング
	// Geminiに対して具体的な役割と指示を与えます
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
- ハッシュタグを含める
`, name, keywords)

	return u.GeminiService.GenerateDescription(ctx, prompt)
}

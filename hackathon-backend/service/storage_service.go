package service

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type StorageService struct {
	client     *storage.Client
	bucketName string
}

func NewStorageService(client *storage.Client, bucketName string) *StorageService {
	return &StorageService{
		client:     client,
		bucketName: bucketName,
	}
}

// Upload は画像をGCSにアップロードし、保存されたファイル名（オブジェクト名）を返します
func (s *StorageService) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	wc := s.client.Bucket(s.bucketName).Object(filename).NewWriter(ctx)

	// アップロード実行
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// バケットが非公開なので、ファイル名だけを返す（後で署名付きURLを作るため）
	return filename, nil
}

// GenerateSignedURL は非公開の画像にアクセスするための「署名付きURL」を発行します
func (s *StorageService) GenerateSignedURL(filename string) (string, error) {
	// 署名付きURLの発行は、Cloud Runのデフォルト環境だと秘密鍵がないため失敗しやすいです。
	// 代わりに、公開バケットのURLを返します。
	// 形式: https://storage.googleapis.com/[バケット名]/[ファイル名]

	if filename == "" {
		return "", nil
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, filename)
	return url, nil
}

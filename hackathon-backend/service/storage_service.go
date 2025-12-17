package service

import (
	"context"
	"io"
	"time"

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
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute), // 15分間だけ有効
	}

	// バケットハンドルから署名付きURLを生成
	url, err := s.client.Bucket(s.bucketName).SignedURL(filename, opts)
	if err != nil {
		return "", err
	}
	return url, nil
}

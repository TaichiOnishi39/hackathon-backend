package service

import (
	"context"
	"fmt"
	"io"
	"strings"

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
func (s *StorageService) UploadImage(ctx context.Context, file io.Reader, filename string) (string, error) {
	// 1. GCSにアップロード
	wc := s.client.Bucket(s.bucketName).Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// 2. ★修正ポイント: 完全なURLを作成して返す
	// 公開バケット前提のURL形式: https://storage.googleapis.com/[バケット名]/[ファイル名]
	fullURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, filename)

	return fullURL, nil
}

// GenerateSignedURL は非公開の画像にアクセスするための「署名付きURL」を発行します
func (s *StorageService) GenerateSignedURL(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}

	// ★修正ポイント: すでに "https://" で始まっているなら、そのまま返す
	// これにより、昔のデータ(ファイル名のみ)と新しいデータ(URL)の両方に対応できます
	if strings.HasPrefix(filename, "https://") || strings.HasPrefix(filename, "http://") {
		return filename, nil
	}

	// ファイル名だけの場合は、URLを組み立てて返す
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, filename)
	return url, nil
}

package service

import (
	"context"
	"fmt"
	"io"
	"strings"
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
func (s *StorageService) GenerateSignedURL(objectName string) (string, error) {
	if objectName == "" {
		return "", nil
	}

	// ★修正: もし "https://storage.googleapis.com/バケット名/" で始まっていたら、そこを取り除く
	// これにより、ファイル名だけのデータも、フルURLのデータも両方扱えるようになります
	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", s.bucketName)
	if strings.HasPrefix(objectName, prefix) {
		objectName = strings.TrimPrefix(objectName, prefix)
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := s.client.Bucket(s.bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", err
	}

	return u, nil
}

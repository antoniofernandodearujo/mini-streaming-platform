package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	BucketName string
	S3Service  *s3.S3
}

// Função para criar um novo cliente S3
func NewS3Client(bucketName, region string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return &S3Client{
		BucketName: bucketName,
		S3Service:  s3.New(sess),
	}, nil
}

// Função para fazer upload de um arquivo para o S3
func (s *S3Client) UploadFile(ctx context.Context, file multipart.File, fileName string) error {
	// Cria um buffer para armazenar o arquivo
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return err
	}

	// Faz o upload do arquivo para o S3
	_, err := s.S3Service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	return err
}

// Função para baixar um arquivo do S3
func (s *S3Client) DownloadFile(ctx context.Context, s3Key string) ([]byte, error) {
	result, err := s.S3Service.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// Função para fazer upload de um arquivo de um caminho para o S3
func (s *S3Client) UploadFileFromPath(ctx context.Context, s3Key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Faz o upload do arquivo para o S3
	_, err = s.S3Service.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s3Key),
		Body:   file,
	})
	return err
}

// ListFiles lista os arquivos dentro de um prefixo específico no bucket S3.
func (s *S3Client) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	var fileKeys []string
	err := s.S3Service.ListObjectsV2PagesWithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, item := range page.Contents {
			fileKeys = append(fileKeys, *item.Key)
		}
		return !lastPage
	})
	if err != nil {
		return nil, err
	}
	return fileKeys, nil
}

func (s *S3Client) ListDirectories(ctx context.Context, prefix string) ([]string, error) {
	objects, err := s.S3Service.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.BucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})

	if err != nil {
		return nil, err
	}

	var directories []string
	for _, commonPrefix := range objects.CommonPrefixes {
		// Remove o prefixo base e "/" final para obter apenas o nome do diretório
		trimmed := strings.TrimSuffix(strings.TrimPrefix(*commonPrefix.Prefix, prefix), "/")
		directories = append(directories, trimmed)
	}

	return directories, nil
}

func (s *S3Client) ListResolutions(ctx context.Context, videoID string) ([]string, error) {
	// Prefixo para as resoluções dentro do vídeo
	prefix := "videos-transcoded/" + videoID + "/"

	objects, err := s.S3Service.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.BucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	var resolutions []string
	for _, commonPrefix := range objects.CommonPrefixes {
		// Remove o prefixo base e "/" final para obter apenas o nome da resolução
		trimmed := strings.TrimSuffix(strings.TrimPrefix(*commonPrefix.Prefix, prefix), "/")
		resolutions = append(resolutions, trimmed)
	}

	return resolutions, nil
}


func (s *S3Client) GetVideoFile(ctx context.Context, videoID string, resolution string) (string, error) {
	// Caminho específico para a resolução escolhida
	prefix := "videos-transcoded/" + videoID + "/" + resolution + "/"

	objects, err := s.S3Service.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return "", err
	}

	// Procurar pelo arquivo "video.m3u8"
	for _, object := range objects.Contents {
		if strings.HasSuffix(*object.Key, "video.m3u8") {
			return *object.Key, nil // Caminho completo do arquivo
		}
	}

	return "", fmt.Errorf("arquivo 'video.m3u8' não encontrado em %s", prefix)
}

func (s *S3Client) GetFileURL(fileKey string) (string, error) {
	// Gerar a URL pública do arquivo
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.BucketName, *s.S3Service.Config.Region, fileKey)
	return url, nil
}



package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIO client
var minioClient *minio.Client

// InitMinio initializes MinIO client
func InitMinio() {
	var err error
	minioClient, err = minio.New("10.162.14.111:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false, // Set to true if using HTTPS
	})
	if err != nil {
		log.Fatalf("❌ Failed to initialize MinIO client: %v", err)
	}
	log.Println("✅ MinIO Client Initialized Successfully!")
}

// UploadFileToMinio uploads a file to the 'user-files' bucket and returns its URL
// UploadFileToMinio uploads a file to the 'user-files' bucket and returns its URL
func UploadFileToMinio(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	bucketName := "user-files" // MinIO bucket name
	objectName := fileHeader.Filename
	contentType := fileHeader.Header.Get("Content-Type")

	log.Printf("🚀 Checking if bucket '%s' exists...", bucketName)
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("⚠️ Error checking bucket existence: %v", err)
		return "", err
	}

	if !exists {
		log.Printf("🚀 Bucket '%s' does not exist, creating it...", bucketName)
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("❌ Failed to create bucket: %v", err)
			return "", err
		}
		log.Printf("✅ Bucket '%s' created successfully!", bucketName)
	} else {
		log.Printf("✅ Bucket '%s' already exists!", bucketName)
	}

	// Read file into buffer
	log.Println("📂 Reading file into memory...")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("❌ Failed to read file: %v", err)
		return "", err
	}
	fileSize := int64(len(fileBytes))
	log.Printf("📏 File size: %d bytes", fileSize)

	// Upload file
	log.Printf("🚀 Uploading file: %s (%d bytes)", objectName, fileSize)
	n, err := minioClient.PutObject(ctx, bucketName, objectName, bytes.NewReader(fileBytes), fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("❌ Failed to upload file: %v", err)
		return "", err
	}

	log.Printf("✅ File uploaded successfully: %s (%d bytes written)", objectName, n.Size)
	log.Printf("🚀 Uploading to MinIO: Bucket='%s', Filename='%s', Content-Type='%s'", bucketName, objectName, contentType)

	// Construct file path
	filePath := fmt.Sprintf("http://10.162.14.111:9000/%s/%s", bucketName, objectName)
	log.Printf("🔗 File can be accessed at: %s", filePath)
	return filePath, nil

}

// GetFileFromMinio fetches a file from MinIO and returns its reader
func GetFileFromMinio(objectName string) (io.ReadCloser, error) {
	ctx := context.Background()
	bucketName := "user-files"

	object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("❌ Failed to fetch file: %v", err)
		return nil, err
	}
	log.Printf("✅ File fetched successfully: %s", objectName)
	return object, nil
}

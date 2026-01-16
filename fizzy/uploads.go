package fizzy

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
)

// CreateDirectUpload requests a direct upload URL for a file.
//
// This is the first step in uploading a file to Fizzy.
// After getting the upload URL, use the returned URL and headers to upload the file.
func (s *UploadsService) CreateDirectUpload(ctx context.Context, req *DirectUploadRequest) (*DirectUploadResponse, error) {
	path := fmt.Sprintf("/%s/rails/active_storage/direct_uploads", s.client.accountSlug)

	var response DirectUploadResponse
	if err := s.client.doRequest(ctx, "POST", path, req, &response); err != nil {
		return nil, fmt.Errorf("failed to create direct upload: %w", err)
	}

	return &response, nil
}

// UploadFile is a convenience method that handles the complete file upload flow.
//
// It:
// 1. Reads the file and calculates MD5 checksum
// 2. Requests a direct upload URL
// 3. Uploads the file to the direct upload URL
// 4. Returns the blob ID for use in rich text fields
//
// Example:
//
//	blobID, err := client.Uploads.UploadFile(ctx, "/path/to/image.png", "image/png")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Use blobID in card body or comment
func (s *UploadsService) UploadFile(ctx context.Context, filePath, contentType string) (string, error) {
	// Open and read file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}

	// Calculate MD5 checksum
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}
	checksum := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file: %w", err)
	}

	// Request direct upload
	uploadResp, err := s.CreateDirectUpload(ctx, &DirectUploadRequest{
		Blob: DirectUploadBlob{
			Filename:    fileInfo.Name(),
			ByteSize:    fileInfo.Size(),
			Checksum:    checksum,
			ContentType: contentType,
		},
	})
	if err != nil {
		return "", err
	}

	// Upload file to direct upload URL
	req, err := http.NewRequestWithContext(ctx, "PUT", uploadResp.DirectUploadURL, file)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	// Set headers from response
	for key, value := range uploadResp.Headers {
		req.Header.Set(key, value)
	}

	// Execute upload
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return uploadResp.BlobID, nil
}

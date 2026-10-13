package domain

import "context"

// Interface for downloading file via link in Localsend
type FileDownloader interface {
	GetFileList(ctx context.Context) (*FileMetadata, error)
	DownloadFile(ctx context.Context, file *File) error
}
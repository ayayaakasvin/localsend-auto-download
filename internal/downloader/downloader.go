package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ayayaakasvin/localsend-auto-download/internal/domain"
)

var _ domain.FileDownloader = &Downloader{}

const (
	localsendPort         = "53317"
	v2DownloadPath        = "/api/localsend/v2/download"
	v2PrepareDownloadPath = "/api/localsend/v2/prepare-download"
)

type Downloader struct {
	client *http.Client

	basePath string // Base path for files

	host      string // ip of Localsend host
	sessionID string // ip address of device inside netword, can be speficied
	pinCode   string // ilf localsend host requires pin

	timeout     time.Duration // 10s default
	idleTimeout time.Duration // 30s default
}

// NewDownloader creates a downloader with 10-second request and 30-second
// idle-connection timeouts. Optional arguments override those defaults, in
// order: timeout, then idleTimeout.
// basePath is the directory where the files will be downloaded
// host for target, ip of device that hosting Localsend share link
// sessionID for LocalSend is device ip inside network, if now specified will be fetched. e.g. wlanX inet <IP>
func NewDownloader(basePath string, slc *domain.ShareLinkConfig, timeoutCfg *domain.TimeoutConfig) (*Downloader, error) {
	if basePath == "" {
		exePath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %w", err)
		}

		basePath = filepath.Dir(exePath)
	}

	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to init directory: %w", err)
	}

	timeout := 10 * time.Second
	idleTimeout := 30 * time.Second
	if timeoutCfg != nil {
		if timeoutCfg.Timeout > 0 {
			timeout = timeoutCfg.Timeout
		}
		if timeoutCfg.IdleTimeout > 0 {
			idleTimeout = timeoutCfg.IdleTimeout
		}
	}

	d := &Downloader{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				IdleConnTimeout: idleTimeout,
			},
		},
		timeout:     timeout,
		idleTimeout: idleTimeout,
		basePath:    basePath,
	}
	if slc != nil {
		d.host = slc.Host
		d.sessionID = slc.SessionID
		d.pinCode = slc.PinCode
	}
	if d.sessionID != "" {
		return d, nil
	}

	if err := d.establishIPInsideNetwork(d.host); err != nil {
		return nil, err
	}
	return d, nil
}

// DownloadFiles implements [domain.FileDownloader].
func (d *Downloader) DownloadFile(ctx context.Context, file *domain.File) error {
	if d.host == "" {
		return fmt.Errorf("host is required")
	}

	if file.ID == "" {
		return fmt.Errorf("fileID is required")
	}

	if file.FileName == "" {
		return fmt.Errorf("filename is required")
	}

	query := url.Values{
		"sessionId": {d.sessionID},
		"fileId":    {file.ID},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, d.requestURL(v2DownloadPath, query), nil)
	if err != nil {
		return fmt.Errorf("create get request: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("perform get request: %w", err)
	}
	defer resp.Body.Close()

	counter := &countingReader{r: resp.Body}

	targetPath := filepath.Join(d.basePath, file.FileName)
	_, err = d.saveStream(counter, targetPath)
	if err != nil {
		err = d.delete(targetPath)
		if err != nil {
			return fmt.Errorf("failed to delete failed file: %s", err.Error())
		}

		return err
	}

	if file.Size > 0 && counter.n != file.Size {
		err = d.delete(targetPath)
		if err != nil {
			return fmt.Errorf("failed to delete file with mismatched size: %s", err.Error())
		}
		return fmt.Errorf("downloaded file size mismatch: expected %d, got %d", file.Size, counter.n)
	}

	return nil
}

// GetFileList implements [domain.FileDownloader].
func (d *Downloader) GetFileList(ctx context.Context) (*domain.FileMetadata, error) {
	if d.host == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	query := make(url.Values)
	if d.pinCode != "" {
		query.Set("pin", d.pinCode)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.requestURL(v2PrepareDownloadPath, query), nil)
	if err != nil {
		return nil, fmt.Errorf("create get request: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("get file list failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var files *domain.FileMetadata
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("decode file list response: %w", err)
	}

	return files, nil
}

func (d *Downloader) requestURL(path string, query url.Values) string {
	return (&url.URL{
		Scheme:   "http",
		Host:     net.JoinHostPort(d.host, localsendPort),
		Path:     path,
		RawQuery: query.Encode(),
	}).String()
}

func (d *Downloader) resolvePath(name string) string {
	if filepath.IsAbs(name) {
		return filepath.Clean(name)
	}
	return filepath.Join(d.basePath, name)
}

func (d *Downloader) String() string {
	return fmt.Sprintf(
		"base-path=%q session-id=%q timeout=%s idle-timeout=%s",
		d.basePath,
		d.sessionID,
		d.timeout,
		d.idleTimeout,
	)
}

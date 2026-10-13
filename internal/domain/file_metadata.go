package domain

import "time"

/*
FileMetadata based on JSON data returned by: `curl -X POST 'http://<IP>:<PORT>/api/localsend/v2/prepare-download'`

	{
		"info": {
			"alias": "Good Melon",
			"version": "2.2",
			"deviceModel": "iPhone",
			"deviceType": "mobile",
			"fingerprint": "4B1AED39C0BDE1AC8CCD6B8636FDB70220038A8EB6BE19025C3F76758194E0A6",
			"download": true
		},
		"sessionId": "192.168.8.9",
		"files": {
			"233b0a48-c824-4701-9758-9f3bb968f178": {
				"id": "233b0a48-c824-4701-9758-9f3bb968f178",
				"fileName": "IMG_8616.GIF",
				"size": 65138,
				"fileType": "image/gif",
				"metadata": {
					"modified": "2026-08-02T14:38:44.98111Z",
					"accessed": "2026-08-02T14:38:44.98111Z"
				}
			}
		}
	}
*/
type FileMetadata struct {
	Info      FileInfo        `json:"info"`
	SessionID string          `json:"sessionId"`
	Files     map[string]File `json:"files"`
}

type FileInfo struct {
	Alias       string `json:"alias"`
	Version     string `json:"version"`
	DeviceModel string `json:"deviceModel"`
	DeviceType  string `json:"deviceType"`
	Fingerprint string `json:"fingerprint"`
	Download    bool   `json:"download"`
}

type File struct {
	ID       string              `json:"id"`
	FileName string              `json:"fileName"`
	Size     int64               `json:"size"`
	FileType string              `json:"fileType"`
	Metadata FileMetadataDetails `json:"metadata"`
}

type FileMetadataDetails struct {
	Modified time.Time `json:"modified"`
	Accessed time.Time `json:"accessed"`
}

package options

import (
	"errors"
	"flag"
	"time"

	"github.com/ayayaakasvin/localsend-auto-download/internal/domain"
	"github.com/ayayaakasvin/localsend-auto-download/internal/downloader"
)

type downloaderFlags struct {
	host string

	pinCode     string
	basePath    string
	sessionID   string
	timeout     time.Duration
	idleTimeout time.Duration
	help        bool
}

func (d *downloaderFlags) Register(fs *flag.FlagSet) {
	fs.StringVar(&d.basePath, "base-path", "./", "Directory where received files are saved")
	fs.StringVar(&d.host, "host", "", "IP address of the LocalSend device to download from")

	fs.StringVar(&d.pinCode, "pin", "", "Set if Localsend Host requires pin")

	fs.StringVar(&d.sessionID, "session-id", "", "LocalSend session ID (discovered from the host when omitted)")
	fs.DurationVar(&d.timeout, "timeout", 10*time.Second, "HTTP request timeout")
	fs.DurationVar(&d.idleTimeout, "idle-timeout", 30*time.Second, "HTTP idle connection timeout")
	fs.BoolVar(&d.help, "help", false, "Show this help message")
}

func (d *downloaderFlags) Validate() error {
	if d.host == "" {
		return errors.New("host is required when session ID is not specified")
	}
	return nil
}

func (d *downloaderFlags) Build() (*downloader.Downloader, error) {
	slc := domain.NewShareLinkConfig(d.host, d.sessionID, d.pinCode)
	tts := domain.NewTimeoutConfig(d.timeout, d.idleTimeout)
	return downloader.NewDownloader(d.basePath, slc, tts)
}

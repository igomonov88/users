package content_uploader

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

type Uploader struct {
	cfg      Config
	uploader *s3manager.Uploader
}

type Config struct {
	Region string
	Bucket string
}

// New returns new Content Uploader
func New(cfg Config) (*Uploader, error) {

	// create new session with given configuration
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating content_uploader client")
	}

	//
	u := Uploader{
		cfg: cfg,
		uploader: s3manager.NewUploader(s),
	}

	return &u, nil
}

// Upload uploads file with given name to s3 bucket
func (u *Uploader) Upload(fileName string, file io.Reader) error {
	ui := s3manager.UploadInput{
		Body:                      file,
		Bucket:                    aws.String(u.cfg.Bucket),
		Key:                       aws.String(fileName),
	}
	
	if _, err := u.uploader.Upload(&ui); err != nil {
		return err
	}

	return nil
}

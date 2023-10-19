package mongo

import (
	"context"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/ioc/apps/oss"
)

func (s *service) UploadFile(ctx context.Context, req *oss.UploadFileRequest) error {
	s.log.Debug().Msgf("bucket name: %s, db file name: %s", req.BucketName, req.FileName)

	if err := req.Validate(); err != nil {
		return exception.NewBadRequest("valiate upload file request error, %s", err)
	}

	bucket, err := s.getBucket(req.BucketName)
	if err != nil {
		return err
	}

	opts := options.GridFSUpload()
	opts.Metadata = req.Meta()

	// 清除已有文件
	bucket.Delete(req.FileName)

	// 上传新文件
	uploadStream, err := bucket.OpenUploadStreamWithID(req.FileName, req.FileName, opts)
	if err != nil {
		return err
	}
	defer uploadStream.Close()

	fileSize, err := io.Copy(uploadStream, req.ReadCloser())
	if err != nil {
		return err
	}

	s.log.Debug().Msgf("Write file %s to DB was successful. File size: %d M", req.FileName, fileSize/1024/1024)
	return nil
}

func (s *service) Download(ctx context.Context, req *oss.DownloadFileRequest) error {
	if err := req.Validate(); err != nil {
		return exception.NewBadRequest("valiate upload file request error, %s", err)
	}

	bucket, err := s.getBucket(req.BucketName)
	if err != nil {
		return err
	}

	s.log.Debug().Msgf("start download file: %s ...", req.FileID)
	// 下载文件
	size, err := bucket.DownloadToStream(req.FileID, req.Writer())
	if err != nil {
		return err
	}

	s.log.Debug().Msgf("download file: %s complete, size: %d", req.FileID, size)
	return nil
}

func (s *service) getBucket(name string) (*gridfs.Bucket, error) {
	opts := options.GridFSBucket()
	opts.SetName(name)

	bucket, err := gridfs.NewBucket(s.db, opts)
	if err != nil {
		return nil, fmt.Errorf("new bucket error, %s", err)
	}

	return bucket, nil
}

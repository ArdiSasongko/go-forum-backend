package cld

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/sirupsen/logrus"
)

type CldService struct {
	Client *cloudinary.Cloudinary
}

func Init(url string) (*CldService, error) {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		logrus.WithField("connect cloudinary", err.Error()).Error(err.Error())
		return nil, err
	}

	return &CldService{Client: cld}, nil
}

func UploadImageByte(ctx context.Context, file []byte, url, folder string) (string, string, error) {
	cld, err := Init(url)
	if err != nil {
		logrus.WithField("init cloudinary", err.Error()).Error(err.Error())
		return "", "", err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := cld.Client.Upload.Upload(ctx, bytes.NewReader(file), uploader.UploadParams{Folder: folder})
	if err != nil {
		logrus.WithField("upload image", err.Error()).Error(err.Error())
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}

func UploadImage(ctx context.Context, file *multipart.FileHeader, url, folder string) (string, string, error) {
	fileName := file.Filename

	src, err := file.Open()
	if err != nil {
		logrus.WithField("open file", err.Error()).Error(err.Error())
		return "", "", err
	}
	defer src.Close()

	fileUpload := "./temp/upload" + fileName
	dst, err := os.Create(fileUpload)
	if err != nil {
		logrus.WithField("create file", err.Error()).Error(err.Error())
		return "", "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		logrus.WithField("copy file", err.Error()).Error(err.Error())
		return "", "", err
	}

	defer func() {
		os.Remove(fileUpload)
	}()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cld, err := Init(url)
	if err != nil {
		logrus.WithField("init cloudinary", err.Error()).Error(err.Error())
		return "", "", err
	}

	result, err := cld.Client.Upload.Upload(ctx, fileUpload, uploader.UploadParams{Folder: folder})
	if err != nil {
		logrus.WithField("upload image", err.Error()).Error(err.Error())
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}

func GetPublicID(imageUrl, folder string) (string, error) {
	validFolder := fmt.Sprintf("/%s/", folder)
	filePath := strings.Split(imageUrl, validFolder)[1]
	publicID := strings.TrimSuffix(filePath, path.Ext(filePath))
	validPublicID := fmt.Sprintf("%s/%s", folder, publicID)
	return validPublicID, nil
}

func DestroyImage(ctx context.Context, url, publicID string) error {
	cld, err := Init(url)
	if err != nil {
		logrus.WithField("init cloudinary", err.Error()).Error(err.Error())
		return err
	}

	logrus.Info(publicID)
	_, err = cld.Client.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		logrus.WithField("delete image", err.Error()).Error(err.Error())
		return err
	}

	return nil
}

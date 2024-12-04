package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func supportedExt(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, supp := range extensions {
		if ext == supp {
			return true
		}
	}
	return false
}

func GetProfile() ([]byte, string, error) {
	dir := "./temp/default_profile"
	files, err := os.ReadDir(dir)
	if err != nil {
		logrus.Error(err.Error())
		return nil, "", fmt.Errorf("failed read folder : %v", err)
	}

	var images []string
	for _, v := range files {
		if !v.IsDir() && supportedExt(v.Name()) {
			images = append(images, filepath.Join(dir, v.Name()))
		}
	}

	if len(images) == 0 {
		logrus.Error("empty images")
		return nil, "", fmt.Errorf("empty images in folder")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rng.Intn(len(images))
	randomImagePath := images[randomIndex]

	imageData, err := os.ReadFile(randomImagePath)
	if err != nil {
		logrus.Error(err.Error())
		return nil, "", fmt.Errorf("failed read file : %v", err)
	}

	return imageData, filepath.Base(randomImagePath), nil
}

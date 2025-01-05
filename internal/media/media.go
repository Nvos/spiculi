package media

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"io"
)

type Store struct {
	fs afero.Fs
}

func NewStore(fs afero.Fs) *Store {
	return &Store{fs: fs}
}

func (s *Store) UploadImage(file io.Reader) (string, error) {
	fileName := fmt.Sprintf("%s.webp", uuid.New().String())
	create, err := s.fs.Create(fileName)
	if err != nil {
		return "", err
	}
	defer create.Close()

	err = optimizeAndWriteImage(file, create, optimizeImageOpts{
		quality: 80,
	})
	if err != nil {
		return "", fmt.Errorf("upload iamge (optimizing): %w", err)
	}

	if err := create.Sync(); err != nil {
		return "", fmt.Errorf("upload image (file sync): %w", err)
	}

	return fileName, nil
}

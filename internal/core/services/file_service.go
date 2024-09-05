package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"wanny-web-services/internal/core/domain"
	"wanny-web-services/internal/core/ports"
)

type FileService struct {
	repo              ports.FileRepository
	maxStoragePerUser int64
	encryptionKey     []byte
}

func NewFileService(repo ports.FileRepository, maxStoragePerUser int64, encryptionKey string) *FileService {
	return &FileService{
		repo:              repo,
		maxStoragePerUser: maxStoragePerUser,
		encryptionKey:     []byte(encryptionKey),
	}
}

func (s *FileService) Upload(userID int64, filename string, data []byte) error {
	if int64(len(data)) > s.maxStoragePerUser {
		return errors.New("file size exceeds maximum allowed storage")
	}

	encryptedData, err := s.encrypt(data)
	if err != nil {
		return err
	}

	file := &domain.File{
		UserID:   userID,
		Filename: filename,
		Size:     int64(len(encryptedData)),
	}

	if err := s.repo.Create(file); err != nil {
		return err
	}

	return nil
}

func (s *FileService) Download(userID int64, filename string) ([]byte, error) {
	//TODO return original file
	_, err := s.repo.GetByUserAndFilename(userID, filename)
	if err != nil {
		return nil, err
	}

	// Here you would typically read the encrypted data from your storage system
	// For this example, we'll just return some dummy data
	dummyEncryptedData := []byte("encrypted data")

	decryptedData, err := s.decrypt(dummyEncryptedData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func (s *FileService) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *FileService) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"uuid"

	"github.com/winstonajr/dashdrop/internal/repository"
)

type FileService struct {
	repo      *repository.Queries
	uploadDir string
}

func NewFileService(repo *repository.Queries, uploadDir string) *FileService {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("falha ao criar pasta de uploads: %v", err))
	}

	return &FileService{
		repo:      repo,
		uploadDir: uploadDir,
	}
}

func (s *FileService) ListFiles(ctx context.Context) ([]repository.File, error) {
	return s.repo.ListFiles(ctx)
}

func (s *FileService) UploadFiles(ctx context.Context, originalName string, size int64, content io.Reader) (string, error) {
	id := uuid.New().String()
	savedName := fmt.Sprintf("%s_%s", id, originalName)
	fullPath := filepath.Join(s.uploadDir, savedName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("erro ao criar arquivo em disco: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, content); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("erro ao gravar arquivo: %w", err)
	}

	params := repository.CreateFileParams{
		ID:           id,
		OriginalName: originalName,
		SavedName:    savedName,
		FileSize:     size,
	}

	if err := s.repo.CreateFile(ctx, params); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("erro ao persistir metadados no banco: %w", err)
	}
	return id, nil
}

func (s *FileService) GetFileDownload(ctx context.Context, id string) (filePath, originalName string, err error) {
	file, err := s.repo.GetFileByID(ctx, id)
	if err != nil {
		return "", "", fmt.Errorf("arquivo não encontrado no banco: %w", err)
	}

	fullPath := filepath.Join(s.uploadDir, file.SavedName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("arquivo físico não encontrado no servidor")
	}

	return fullPath, file.OriginalName, nil
}

func (s *FileService) DeleteFile(ctx context.Context, id string) error {
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		return fmt.Errorf("arquivo não encontrado: %w", err)
	}

	if err := s.repo.DeleteFile(ctx, id); err != nil {
		return fmt.Errorf("falha ao remover registro do banco: %w", err)
	}

	fullPath := filepath.Join(s.uploadDir, file.SavedName)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("registro removido, mas falha ao apagar arquivo físico: %w", err)
	}

	return nil
}

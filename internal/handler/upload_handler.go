package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maynguyen24/sever/pkg/response"
	"github.com/maynguyen24/sever/pkg/upload"
)

type UploadHandler struct {
	uploader upload.Uploader
}

func NewUploadHandler(uploader upload.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

// Upload handles single file upload via multipart/form-data
func (h *UploadHandler) Upload(c *fiber.Ctx) error {
	// Parse the multipart form:
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File is required")
	}

	// Get the folder/prefix from query or default
	folder := c.Query("folder", "general")

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
	}
	defer fileContent.Close()

	// Upload using the service
	info, err := h.uploader.UploadFile(c.Context(), fileContent, file.Size, file.Header.Get("Content-Type"), folder)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}

	return response.Success(c, 2000, "File uploaded successfully", info)
}

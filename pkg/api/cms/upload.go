package cmsApi

import (
	"fmt"
	"path/filepath"

	"github.com/dtekltd/common/api"
	"github.com/gofiber/fiber/v2"
)

var uploadDir = "./uploads"

// File upload route for CkEditor
func ckUpload(ctx *fiber.Ctx) error {
	// Get uploaded file
	file, err := ctx.FormFile("upload")
	if err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	// Define file path
	filePath := filepath.Join(uploadDir, file.Filename)

	// Save file to uploads folder
	if err := ctx.SaveFile(file, filePath); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	// Return CKEditor-friendly response
	return api.SuccessResp(ctx, fmt.Sprintf("/uploads/%s", file.Filename))
}

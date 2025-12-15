package mediaApi

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dtekltd/common/api"
	"github.com/gofiber/fiber/v2"
)

var (
	rootDir = "./uploads"
	rootUrl = "/uploads"
)

type FileInfo struct {
	IsDir     bool   `json:"isDir"`
	Name      string `json:"name"`
	Ext       string `json:"ext"`
	Url       string `json:"url"`
	UpdatedAt int64  `json:"updatedAt"`
}

func createFolder(ctx *fiber.Ctx) error {
	type CreateFolderReq struct {
		Name string `json:"name"`
	}
	req := &CreateFolderReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	path := ctx.Query("p", "/")
	dir := filepath.Join(rootDir, path, req.Name)
	if err := os.Mkdir(dir, 0o755); err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}

	return api.SuccessResp(ctx, true)
}

func deleteMedia(ctx *fiber.Ctx) error {
	type DeleteMediaReq struct {
		Path  string   `json:"path"`
		Files []string `json:"files"`
	}
	req := &DeleteMediaReq{}
	if err := ctx.BodyParser(req); err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	if len(req.Files) == 0 {
		return api.ErrorBadRequestResp(ctx, "Missing files")
	}

	count := 0
	for _, fileName := range req.Files {
		filePath := filepath.Join(rootDir, req.Path, fileName)
		_, err := os.Stat(filePath)

		if os.IsNotExist(err) {
			return api.ErrorBadRequestResp(ctx, "File does not exist")
		} else if err != nil {
			return api.ErrorBadRequestResp(ctx, err.Error())
		} else {
			if err := os.Remove(filePath); err != nil {
				return api.ErrorInternalServerErrorResp(ctx, err.Error())
			}
		}
		count += 1
	}

	return api.SuccessResp(ctx, &count)
}

func uploadMedia(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return api.ErrorBadRequestResp(ctx, err.Error())
	}

	fileInfos := []FileInfo{}
	for _, files := range form.File {
		file := files[0]
		path := ctx.Query("p", "/")
		uploadDir := filepath.Join(rootDir, path)
		filePath := filepath.Join(uploadDir, file.Filename)

		if err := ctx.SaveFile(file, filePath); err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}

		info, err := os.Stat(filePath)
		if err != nil {
			return api.ErrorInternalServerErrorResp(ctx, err.Error())
		}

		url := rootUrl + "/" + file.Filename
		if runtime.GOOS == "windows" {
			url = strings.ReplaceAll(url, "\\", "/")
		}

		fileInfos = append(fileInfos, FileInfo{
			Url:       url,
			IsDir:     false,
			Name:      file.Filename,
			Ext:       filepath.Ext(file.Filename),
			UpdatedAt: info.ModTime().Unix(),
		})
	}

	return api.SuccessResp(ctx, fileInfos)
}

func getMedia(ctx *fiber.Ctx) error {
	path := ctx.Query("p", "/")
	files, err := getFiles(path)
	if err != nil {
		return api.ErrorInternalServerErrorResp(ctx, err.Error())
	}
	return api.SuccessResp(ctx, files)
}

func getFiles(path string) ([]FileInfo, error) {
	var files []FileInfo

	entries, err := os.ReadDir(filepath.Join(rootDir, path))
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		files = append(files, newFileInfo(path, entry))
	}

	return files, nil
}

func newFileInfo(path string, entry os.DirEntry) FileInfo {
	url := rootUrl
	if path != "/" {
		url += path
	}
	url += "/" + entry.Name()
	if runtime.GOOS == "windows" {
		url = strings.ReplaceAll(url, "\\", "/")
	}
	info, _ := entry.Info()
	return FileInfo{
		Url:       url,
		IsDir:     entry.IsDir(),
		Name:      filepath.Base(entry.Name()),
		Ext:       filepath.Ext(entry.Name()),
		UpdatedAt: info.ModTime().Unix(),
	}
}

package utils

import (
	"encoding/json"
	"net/http"
	"path/filepath"
)

func GetContentType(path string) string {
	switch filepath.Ext(path) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		return ""
	}
}

func CheckAllowedImageExtension(extension string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range allowedExtensions {
		if extension == ext {
			return true
		}
		return true
	}
	return false
}

func CheckModel(inputModel string) bool {
	models := []string{
		"4x_NMKD-Superscale-SP_178000_G",
		"realesrgan-x4plus-anime",
		"realesrgan-x4plus",
		"RealESRGAN_General_x4_v3",
		"remacri",
		"ultramix_balanced",
		"ultrasharp",
	}
	for _, model := range models {
		if inputModel == model {
			return true
		}
	}
	return false
}

func WriteJSON(w http.ResponseWriter, code int, message map[string]any) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(message)

	w.WriteHeader(code)
	_, _ = w.Write(j)

}

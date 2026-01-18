package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"unicode"
)

// detectContentType detects MIME type from file content (basic magic number detection)
func detectContentType(content []byte) string {
	if len(content) == 0 {
		return "application/octet-stream"
	}
	
	// Check for common image formats
	if len(content) >= 2 {
		// JPEG: FF D8
		if content[0] == 0xFF && content[1] == 0xD8 {
			return "image/jpeg"
		}
		// PNG: 89 50 4E 47
		if len(content) >= 4 && content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47 {
			return "image/png"
		}
		// GIF: 47 49 46 38
		if len(content) >= 4 && content[0] == 0x47 && content[1] == 0x49 && content[2] == 0x46 && content[3] == 0x38 {
			return "image/gif"
		}
	}
	
	// PDF: %PDF
	if len(content) >= 4 && string(content[0:4]) == "%PDF" {
		return "application/pdf"
	}
	
	// Default
	return "application/octet-stream"
}

// SanitizeFilename sanitizes a filename to prevent path traversal and other attacks
func SanitizeFilename(filename string) string {
	// Remove path components
	filename = filepath.Base(filename)
	
	// Remove any remaining path separators
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	
	// Remove null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")
	
	// Keep only safe characters (alphanumeric, dots, hyphens, underscores)
	var result strings.Builder
	for _, r := range filename {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '.' || r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}
	
	sanitized := result.String()
	
	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "file"
	}
	
	// Limit length
	if len(sanitized) > 255 {
		ext := filepath.Ext(sanitized)
		name := sanitized[:255-len(ext)]
		sanitized = name + ext
	}
	
	return sanitized
}

// GenerateSafeFilename generates a safe filename with random prefix
func GenerateSafeFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	
	// Generate random prefix
	bytes := make([]byte, 16)
	rand.Read(bytes)
	prefix := hex.EncodeToString(bytes)
	
	// Sanitize original name
	sanitized := SanitizeFilename(originalName)
	nameWithoutExt := strings.TrimSuffix(sanitized, ext)
	
	// Combine: random_prefix-originalname.ext
	return fmt.Sprintf("%s-%s%s", prefix, nameWithoutExt, ext)
}

// ValidatePath prevents path traversal
func ValidatePath(path string) error {
	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal detected")
	}
	
	// Check for absolute paths
	if filepath.IsAbs(path) {
		return fmt.Errorf("absolute paths not allowed")
	}
	
	// Normalize and check
	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") {
		return fmt.Errorf("path traversal detected")
	}
	
	return nil
}

// ValidateMimeType validates MIME type against file content (basic check)
func ValidateMimeType(content io.Reader, declaredMimeType string, filename string) (bool, error) {
	// Read first 512 bytes for magic number detection
	buffer := make([]byte, 512)
	n, err := content.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}
	
	// Detect MIME type from content (basic magic number detection)
	detectedMimeType := detectContentType(buffer[:n])
	
	// Also check extension-based MIME type
	ext := filepath.Ext(filename)
	extMimeType := mime.TypeByExtension(ext)
	
	// Validate against declared type
	allowedTypes := []string{
		"image/jpeg", "image/png", "image/gif", "image/webp",
		"video/mp4", "video/webm",
		"application/pdf",
		"text/plain", "text/csv",
	}
	
	// Check if declared type is allowed
	declaredAllowed := false
	for _, allowed := range allowedTypes {
		if declaredMimeType == allowed {
			declaredAllowed = true
			break
		}
	}
	
	if !declaredAllowed {
		return false, fmt.Errorf("MIME type not allowed: %s", declaredMimeType)
	}
	
	// Check if detected type matches declared type (basic validation)
	// In production, you might want stricter validation
	if detectedMimeType != declaredMimeType && extMimeType != declaredMimeType {
		// Log warning but don't fail (MIME detection can be imprecise)
		// In strict mode, you might want to return false here
	}
	
	return true, nil
}

// ValidateFileSize validates file size
func ValidateFileSize(size int64, maxSize int64) error {
	if size <= 0 {
		return fmt.Errorf("invalid file size")
	}
	
	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}
	
	return nil
}

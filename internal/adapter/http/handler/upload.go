package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

// imageTypes maps each accepted extension to the content type its bytes
// must sniff as, so a renamed HTML or script file is rejected.
var imageTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

// uploadError carries a client-facing message and status code.
type uploadError struct {
	msg    string
	status int
}

func (e *uploadError) Error() string { return e.msg }

// saveImage stores the multipart file in field under images/<subdir>/ and
// returns its path relative to the working directory and its file name.
func saveImage(r *http.Request, field, subdir string) (path, fileName string, err error) {
	file, header, err := r.FormFile(field)
	if err != nil {
		return "", "", &uploadError{field + " is required", http.StatusBadRequest}
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	wantType, ok := imageTypes[ext]
	if !ok {
		return "", "", &uploadError{"only jpg, jpeg, png, and webp files are allowed", http.StatusBadRequest}
	}

	head := make([]byte, 512)
	n, err := io.ReadFull(file, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", "", &uploadError{"could not read image", http.StatusBadRequest}
	}
	if http.DetectContentType(head[:n]) != wantType {
		return "", "", &uploadError{"file content does not match its extension", http.StatusBadRequest}
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", err
	}

	uploadDir := filepath.Join("images", subdir)
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", "", fmt.Errorf("create upload directory: %w", err)
	}

	fileName = fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.NewString(), ext)
	path = filepath.Join(uploadDir, fileName)

	dst, err := os.Create(path)
	if err != nil {
		return "", "", fmt.Errorf("create image file: %w", err)
	}
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		os.Remove(path) // don't leave a truncated file behind
		return "", "", fmt.Errorf("write image: %w", err)
	}
	if err := dst.Close(); err != nil {
		os.Remove(path)
		return "", "", fmt.Errorf("close image: %w", err)
	}
	return path, fileName, nil
}

// sendUploadError writes client errors as-is and hides internal ones.
func sendUploadError(w http.ResponseWriter, err error) {
	var ue *uploadError
	if errors.As(err, &ue) {
		util.SendData(w, ue.msg, ue.status)
		return
	}
	log.Println("image upload failed:", err)
	util.SendData(w, "failed to save image", http.StatusInternalServerError)
}

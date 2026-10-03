package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

func (h *AgencyHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	agencyID, err := uuid.Parse(r.PathValue("agency_id"))
	if err != nil {
		util.SendData(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	replaceImage(w, r, "agencies", func(path string) (string, error) {
		return h.imageUC.Execute(r.Context(), actor, agencyID, path)
	})
}

func (h *TourHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	agencyID, err := uuid.Parse(r.PathValue("agency_id"))
	if err != nil {
		util.SendData(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	tourID, err := uuid.Parse(r.PathValue("tour_id"))
	if err != nil {
		util.SendData(w, "invalid tour id", http.StatusBadRequest)
		return
	}
	replaceImage(w, r, "tours", func(path string) (string, error) {
		return h.imageUC.Execute(r.Context(), actor, agencyID, tourID, path)
	})
}

const maxImageBytes = 10 << 20

func replaceImage(w http.ResponseWriter, r *http.Request, subdir string, update func(string) (string, error)) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImageBytes+(1<<20))
	err := r.ParseMultipartForm(maxImageBytes)
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			util.SendData(w, "image must be at most 10 MB", http.StatusRequestEntityTooLarge)
		} else {
			util.SendData(w, "invalid multipart form data", http.StatusBadRequest)
		}
		return
	}
	files := r.MultipartForm.File["image"]
	if len(files) != 1 {
		util.SendData(w, "select one image", http.StatusBadRequest)
		return
	}
	if files[0].Size > maxImageBytes {
		util.SendData(w, "image must be at most 10 MB", http.StatusRequestEntityTooLarge)
		return
	}
	path, _, err := saveImage(r, "image", subdir)
	if err != nil {
		sendUploadError(w, err)
		return
	}
	imagePath := filepath.ToSlash(path)
	oldPath, err := update(imagePath)
	if err != nil {
		_ = os.Remove(path)
		switch {
		case errors.Is(err, domain.ErrImageTargetNotFound):
			util.SendData(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domain.ErrImageAccessDenied):
			util.SendData(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, domain.ErrTourCancelled):
			util.SendData(w, err.Error(), http.StatusConflict)
		default:
			log.Println("replace image:", err)
			util.SendData(w, "failed to update image", http.StatusInternalServerError)
		}
		return
	}
	removeReplacedImage(oldPath, subdir)
	util.SendData(w, map[string]string{"message": "Image updated successfully", "image_path": imagePath, "image_url": "/" + imagePath}, http.StatusOK)
}

// Only remove direct files in the expected upload directory. Older tour paths
// can omit the images/ prefix; arbitrary database paths must never be deleted.
func removeReplacedImage(path, subdir string) {
	if path == "" {
		return
	}
	if strings.HasPrefix(path, subdir+"/") {
		path = filepath.Join("images", path)
	}
	path = filepath.Clean(filepath.FromSlash(path))
	if filepath.Dir(path) != filepath.Join("images", subdir) {
		return
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Println("remove replaced image:", err)
	}
}

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/azzamt11/todoapp_backend/app/model"
)

func GetAllProjects(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
    if err != nil || pageSize < 1 {
        pageSize = 10
    }

    sortBy := r.URL.Query().Get("sort_by")
    sortOrder := r.URL.Query().Get("sort_order")

    // Parse search query parameter
    searchQuery := r.URL.Query().Get("search")

    query := db.Offset((page - 1) * pageSize).Limit(pageSize)
    if sortBy != "" {
        switch sortBy {
		case "title":
            query = query.Order("title " + sortOrder)
        case "updated":
            query = query.Order("updated_at " + sortOrder)
        case "created":
            query = query.Order("created_at " + sortOrder)
        default:
            http.Error(w, "Unsupported sorting criteria", http.StatusBadRequest)
            return
        }
    }

    // Apply search filter if search query is provided
    if searchQuery != "" {
        // Use ILIKE for case-insensitive search in PostgreSQL, use LIKE for case-sensitive search in MySQL
        //query = query.Where("name ILIKE ?", "%"+searchQuery+"%") // Use ILIKE for PostgreSQL
        query = query.Where("name LIKE ?", "%"+searchQuery+"%") // Use LIKE for MySQL
    }

    projects := []model.Project{}
    query.Find(&projects)

    respondJSON(w, http.StatusOK, projects)
}



func CreateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	project := model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, project)
}

func GetProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _:= strconv.Atoi(vars["id"])
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func UpdateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _:= strconv.Atoi(vars["id"])
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func DeleteProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	if err := db.Delete(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func ArchiveProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _:= strconv.Atoi(vars["id"])
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	project.Archive()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func RestoreProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _:= strconv.Atoi(vars["id"])
	project := getProjectOr404(db, id, w, r)
	if project == nil {
		return
	}
	project.Restore()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

// getProjectOr404 gets a project instance if exists, or respond the 404 error otherwise
func getProjectOr404(db *gorm.DB, projectID int, w http.ResponseWriter, r *http.Request) *model.Project {
	project := model.Project{}
	if err := db.First(&project, projectID).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &project
}

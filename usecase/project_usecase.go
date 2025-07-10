package usecase

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"taskify/config"
	"taskify/models"
	"taskify/utils"
)

// InputProject menerima CreatedByID secara eksplisit untuk pembuatan proyek.
type InputProject struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	CreatedByID uuid.UUID `json:"created_by" binding:"required"` // Field ini wajib diisi di body saat Create Project
}

// CreateProject: Membuat proyek baru. Memerlukan otentikasi (JWT) dan CreatedByID manual.
func CreateProject(c *gin.Context) {
	var input InputProject
	// Melakukan binding JSON dari request body ke struct InputProject
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Jika binding gagal (misal field required kosong), kirim 400
		return
	}

	// Verifikasi apakah CreatedByID yang diberikan di body request adalah user yang valid di database.
	var user models.User
	if err := config.DB.Where("id = ?", input.CreatedByID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika user dengan ID tersebut tidak ditemukan
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provided created_by ID does not correspond to a valid user"})
		} else {
			// Jika ada error database lain saat mencari user
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error validating created_by ID"})
		}
		return
	}

	// Buat objek Project baru
	project := models.Project{
		ID:          uuid.New(), // Generate UUID baru untuk ID proyek
		Name:        input.Name,
		Description: input.Description,
		CreatedByID: input.CreatedByID, // Set CreatedByID dari input request
	}

	// Simpan proyek ke database
	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// PRELOAD USER UNTUK RESPONSE: Agar detail User muncul di JSON respons
	if err := config.DB.Preload("User").First(&project, "id = ?", project.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload user for project response"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully", "project": project})
}

// GetProjects: Mengambil semua proyek yang dimiliki oleh user yang sedang login.
func GetProjects(c *gin.Context) {
	// AMBIL USER ID DARI JWT TOKEN: Ini adalah kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		// Jika userID tidak ditemukan di konteks (misal token tidak ada/invalid), error sudah ditangani di middleware/utils
		return
	}

	var projects []models.Project
	// Ambil proyek di mana created_by_id sama dengan userID dari token
	// PRELOAD USER UNTUK RESPONSE: Agar detail User muncul di JSON respons
	if err := config.DB.Preload("User").Where("created_by_id = ?", userID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projects retrieved successfully", "projects": projects})
}

// GetProjectByID: Mengambil satu proyek berdasarkan ID, memastikan user yang login adalah pemiliknya.
func GetProjectByID(c *gin.Context) {
	projectIDStr := c.Param("id")              // Ambil projectID dari parameter URL
	projectID, err := uuid.Parse(projectIDStr) // Parse string ke UUID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// AMBIL USER ID DARI JWT TOKEN: Kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var project models.Project
	// Cari proyek berdasarkan ID DAN pastikan created_by_id cocok dengan userID dari token.
	// PRELOAD USER UNTUK RESPONSE: Agar detail User muncul di JSON respons
	if err := config.DB.Preload("User").Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Jika proyek tidak ditemukan ATAU user yang login bukan pemiliknya
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project retrieved successfully", "project": project})
}

// UpdateProject: Mengupdate proyek yang sudah ada, memastikan user yang login adalah pemiliknya.
func UpdateProject(c *gin.Context) {
	projectIDStr := c.Param("id") // Ambil projectID dari parameter URL
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// AMBIL USER ID DARI JWT TOKEN: Kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var input struct { // Struct anonim untuk input update (hanya name & description)
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var project models.Project
	// Cari proyek untuk diupdate berdasarkan ID DAN pastikan user yang login adalah pemiliknya.
	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project for update"})
		}
		return
	}

	// Lakukan update pada field yang diberikan
	updates := models.Project{Name: input.Name, Description: input.Description}
	if err := config.DB.Model(&project).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// PRELOAD USER UNTUK RESPONSE: Agar detail User muncul di JSON respons setelah update
	if err := config.DB.Preload("User").First(&project, "id = ?", project.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload user after project update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully", "project": project})
}

// DeleteProject: Menghapus proyek, memastikan user yang login adalah pemiliknya.
func DeleteProject(c *gin.Context) {
	projectIDStr := c.Param("id") // Ambil projectID dari parameter URL
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// AMBIL USER ID DARI JWT TOKEN: Kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var project models.Project
	// Cari proyek untuk dihapus berdasarkan ID DAN pastikan user yang login adalah pemiliknya.
	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project for deletion"})
		}
		return
	}

	// Hapus tugas-tugas terkait. (Secara teori, ON DELETE CASCADE di DB sudah cukup, tapi ini sebagai fallback)
	if err := config.DB.Where("project_id = ?", projectID).Delete(&models.Task{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated tasks"})
		return
	}

	// Hapus proyek
	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

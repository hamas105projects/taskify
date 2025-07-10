package usecase

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"taskify/config"
	"taskify/models"
	"taskify/utils"
)

// InputTask: Struktur input untuk membuat atau mengupdate tugas.
type InputTask struct {
	Title       string            `json:"title" binding:"required"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status" binding:"required,oneof=todo in_progress done"` // Status wajib dan harus salah satu dari enum
	Deadline    *CustomDate       `json:"deadline"`                                              // Format tanggal YYYY-MM-DD
}
type CustomDate struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // remove quotes
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.Time.Format(dateLayout))
}

func CreateTask(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"}) // Jika format ID salah
		return
	}

	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var project models.Project

	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access to this project"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project for task creation"})
		}
		return
	}

	var input InputTask
	// Melakukan binding JSON dari request body ke struct InputTask
	if err := c.ShouldBindJSON(&input); err != nil {
		// Jika binding gagal (misal field required kosong, status salah, format tanggal salah)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat objek Task baru
	task := models.Task{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
	}

	if input.Deadline != nil {
		t := input.Deadline.Time
		task.Deadline = &t
	}

	// Simpan tugas ke database
	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// PRELOAD PROJECT DAN USER UNTUK RESPONSE: Agar detail Project dan User muncul di JSON respons
	if err := config.DB.Preload("Project.User").First(&task, "id = ?", task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload project/user for task response"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": task})
}

// GetTasksByProject: Mengambil semua tugas untuk proyek spesifik, memastikan user yang login punya akses.
func GetTasksByProject(c *gin.Context) {
	projectIDStr := c.Param("project_id") // Ambil projectID dari parameter URL
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
	// VERIFIKASI KEPEMILIKAN PROYEK: Pastikan projectID di URL dimiliki oleh userID yang login.
	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access to this project"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		}
		return
	}

	var tasks []models.Task
	// Ambil semua tugas yang ProjectID-nya cocok dengan projectID dari URL.
	// PRELOAD PROJECT DAN USER UNTUK RESPONSE: Agar detail Project dan User muncul di JSON respons
	if err := config.DB.Preload("Project.User").Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tasks retrieved successfully", "tasks": tasks})
}

// GetTaskByID: Mengambil satu tugas berdasarkan ID dalam proyek spesifik, memastikan user punya akses.
func GetTaskByID(c *gin.Context) {
	projectIDStr := c.Param("project_id") // Ambil projectID dari parameter URL
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	taskIDStr := c.Param("task_id") // Ambil taskID dari parameter URL
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// AMBIL USER ID DARI JWT TOKEN: Kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var project models.Project
	// VERIFIKASI KEPEMILIKAN PROYEK: Penting sebelum mencari task
	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access to this project"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		}
		return
	}

	var task models.Task
	// Cari tugas berdasarkan taskID DAN pastikan projectID cocok DAN project tersebut dimiliki oleh userID yang login
	// PRELOAD PROJECT DAN USER UNTUK RESPONSE: Agar detail Project dan User muncul di JSON respons
	if err := config.DB.Preload("Project.User").Where("id = ? AND project_id = ?", taskID, projectID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this project or you don't have access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task retrieved successfully", "task": task})
}

// UpdateTask: Mengupdate tugas dalam proyek spesifik, memastikan user punya akses.
func UpdateTask(c *gin.Context) {
	projectIDStr := c.Param("project_id") // Ambil projectID dari parameter URL
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	taskIDStr := c.Param("task_id") // Ambil taskID dari parameter URL
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// AMBIL USER ID DARI JWT TOKEN: Kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var project models.Project
	// VERIFIKASI KEPEMILIKAN PROYEK: Penting sebelum mengupdate task
	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access to this project"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project for task update"})
		}
		return
	}

	var input InputTask
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	// Cari tugas untuk diupdate berdasarkan taskID DAN pastikan projectID cocok
	if err := config.DB.Where("id = ? AND project_id = ?", taskID, projectID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this project or you don't have access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task for update"})
		}
		return
	}

	// Lakukan update pada field yang diberikan
	updates := models.Task{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
	}

	if input.Deadline != nil {
		t := input.Deadline.Time
		task.Deadline = &t
	}

	if err := config.DB.Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// PRELOAD PROJECT DAN USER UNTUK RESPONSE: Agar detail Project dan User muncul di JSON respons setelah update
	if err := config.DB.Preload("Project.User").First(&task, "id = ?", task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to preload project/user after task update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}

// DeleteTask: Menghapus tugas dalam proyek spesifik, memastikan user punya akses.
func DeleteTask(c *gin.Context) {
	projectIDStr := c.Param("project_id") // Ambil projectID dari parameter URL
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	taskIDStr := c.Param("task_id") // Ambil taskID dari parameter URL
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// AMBIL USER ID DARI JWT TOKEN: Kunci otorisasi!
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		return
	}

	var project models.Project
	// VERIFIKASI KEPEMILIKAN PROYEK: Penting sebelum menghapus task
	if err := config.DB.Where("id = ? AND created_by_id = ?", projectID, userID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or you don't have access to this project"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project for task deletion"})
		}
		return
	}

	var task models.Task
	// Cari tugas untuk dihapus berdasarkan taskID DAN pastikan projectID cocok
	if err := config.DB.Where("id = ? AND project_id = ?", taskID, projectID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found in this project or you don't have access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task for deletion"})
		}
		return
	}

	// Hapus tugas
	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

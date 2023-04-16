package main

import (
	"fmt"
	"gin_webserver/web_structs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var notes = make(map[int]*web_structs.Note)
var notesMutex sync.Mutex
var noteIDCounter int

const uploadsDir = "./uploads"

func createNote(c *gin.Context) {
	var newNote web_structs.Note

	if err := c.BindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notesMutex.Lock()
	newNote.ID = noteIDCounter
	notes[noteIDCounter] = &newNote
	noteIDCounter++
	notesMutex.Unlock()

	c.JSON(http.StatusCreated, newNote)
}

func getAllNotes(c *gin.Context) {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	var notesList []*web_structs.Note
	for _, note := range notes {
		notesList = append(notesList, note)
	}

	c.JSON(http.StatusOK, notesList)
}

func getNoteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	notesMutex.Lock()
	defer notesMutex.Unlock()

	note, ok := notes[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func updateNoteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var updatedNote web_structs.Note
	if err := c.BindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notesMutex.Lock()
	defer notesMutex.Unlock()

	note, ok := notes[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	updatedNote.ID = id
	*note = updatedNote
	c.JSON(http.StatusOK, updatedNote)
}

func deleteNoteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	notesMutex.Lock()
	defer notesMutex.Unlock()

	if _, ok := notes[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	delete(notes, id)
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dst := filepath.Join(uploadsDir, file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}

func listFiles(c *gin.Context) {
	files, err := ioutil.ReadDir(uploadsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name()
	}

	c.JSON(http.StatusOK, filenames)
}

func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(uploadsDir, filename)

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.File(filePath)
}

func deleteFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(uploadsDir, filename)

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func generateString(length int, base string) string {
	return fmt.Sprintf("%-*.*s", length, length, base)
}

func initNotes(c *gin.Context) {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	notes = make(map[int]*web_structs.Note)
	for i := 1; i <= 10000; i++ {
		note := &web_structs.Note{
			ID:      i,
			Title:   generateString(10, "Title "+strconv.Itoa(i)),
			Content: generateString(100, "Content "+strconv.Itoa(i)),
		}
		notes[i] = note
	}
	c.JSON(http.StatusOK, gin.H{"message": "10,000 notes initialized"})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// Register the CRUD operation endpoints
	r.POST("/init_notes", initNotes)
	r.POST("/notes", createNote)
	r.GET("/notes", getAllNotes)
	r.GET("/notes/:id", getNoteByID)
	r.PUT("/notes/:id", updateNoteByID)
	r.DELETE("/notes/:id", deleteNoteByID)
	r.POST("/files", uploadFile)
	r.GET("/files", listFiles)
	r.GET("/files/:filename", downloadFile)
	r.DELETE("/files/:filename", deleteFile)

	// Start the server
	port := "8080"
	fmt.Printf("Server running on port %s\n", port)
	err := r.Run(":" + port)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"registration_backend/database"
	"registration_backend/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct for registration
type User struct {
	RegNo          string `json:"regno"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Password       string `json:"password"`
	UserImagePath  string `json:"userImagePath"`
	AadharCardPath string `json:"aadharCardPath"` // JSON array of multiple PDFs
	CreatedAt      string `json:"created_at"`
}

// Generate unique registration number
func generateRegNo() string {
	var maxRegNo *int
	currentYear := time.Now().Year()
	startRegNo := currentYear*10000 + 1000

	query := `SELECT MAX(regno) FROM user_registration`
	err := database.DB.QueryRow(context.Background(), query).Scan(&maxRegNo)
	if err != nil {
		log.Printf("⚠️ Error fetching max regno: %v", err)
		return fmt.Sprintf("%d", startRegNo)
	}

	if maxRegNo == nil {
		return fmt.Sprintf("%d", startRegNo)
	}

	return fmt.Sprintf("%d", *maxRegNo+1)
}

// Hash password
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// RegisterUserHandler
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // 50 MB limit
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	// Check if email or phone exists
	var existingEmail, existingPhone string
	checkUserQuery := `SELECT email, phone FROM user_registration WHERE email = $1 OR phone = $2`
	err = database.DB.QueryRow(context.Background(), checkUserQuery, email, phone).Scan(&existingEmail, &existingPhone)
	if existingEmail == email {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}
	if existingPhone == phone {
		http.Error(w, "Phone number already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	regNo := generateRegNo()

	// Handle file uploads
	var userImagePath string
	var aadharFiles []string // Slice for Aadhar and other PDFs

	// Log received files for debugging
	log.Println("📥 Checking uploaded files...")

	// 📌 **Handle Multiple PDFs**
	if pdfFiles, ok := r.MultipartForm.File["pdfFiles"]; ok {
		log.Printf("📥 Received %d PDF files", len(pdfFiles))

		for _, pdfFileHeader := range pdfFiles {
			pdfFile, err := pdfFileHeader.Open()
			if err != nil {
				log.Println("❌ Failed to open PDF file:", err)
				http.Error(w, "Failed to process PDF files", http.StatusInternalServerError)
				return
			}
			defer pdfFile.Close()

			// Upload PDF
			pdfPath, err := storage.UploadFileToMinio(pdfFile, pdfFileHeader)
			if err != nil {
				log.Println("❌ Failed to upload PDF:", err)
				http.Error(w, "Failed to upload PDF file", http.StatusInternalServerError)
				return
			}

			log.Println("✅ Successfully uploaded PDF:", pdfPath)
			aadharFiles = append(aadharFiles, pdfPath)
		}
	} else {
		log.Println("⚠️ No PDF files received!")
	}

	// 📌 **Handle User Image Upload**
	if userImage, userImageHeader, err := r.FormFile("userImage"); err == nil {
		defer userImage.Close()
		userImagePath, err = storage.UploadFileToMinio(userImage, userImageHeader)
		if err != nil {
			log.Println("❌ Failed to upload user image:", err)
			http.Error(w, "Failed to upload user image", http.StatusInternalServerError)
			return
		}
		log.Println("✅ Successfully uploaded user image:", userImagePath)
	}

	// 📌 **Handle Aadhar Card Upload**
	// 📌 **Handle Aadhar Card Upload** - Ensure aadharCard is included
	if aadharCard, aadharCardHeader, err := r.FormFile("aadharCard"); err == nil {
		defer aadharCard.Close()
		aadharCardPath, err := storage.UploadFileToMinio(aadharCard, aadharCardHeader)
		if err != nil {
			log.Println("❌ Failed to upload Aadhar card:", err)
			http.Error(w, "Failed to upload Aadhar card", http.StatusInternalServerError)
			return
		}
		log.Println("✅ Successfully uploaded Aadhar card:", aadharCardPath)
		aadharFiles = append(aadharFiles, aadharCardPath) // Append to PDF array
	}

	// Convert all file paths (Aadhar + PDFs) to JSON
	aadharPathsJSON, err := json.Marshal(aadharFiles)
	if err != nil {
		log.Println("❌ Failed to convert file paths to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("📄 Stored in Aadhar Field: %s", string(aadharPathsJSON))

	// 📌 **Insert User Into Database**
	query := `INSERT INTO user_registration (regno, user_name, email, phone, password, photo, aadhar_card, pdf_files, created_at) 
          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`

	_, err = database.DB.Exec(context.Background(), query, regNo, name, email, phone, hashedPassword, userImagePath, string(aadharPathsJSON), string(aadharPathsJSON))
	if err != nil {
		log.Println("❌ Error inserting user:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Publish User Data to RabbitMQ
	user := User{
		RegNo:          regNo,
		Name:           name,
		Email:          email,
		Phone:          phone,
		Password:       hashedPassword,
		UserImagePath:  userImagePath,
		AadharCardPath: string(aadharPathsJSON),
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	_ = PublishToQueue(user)

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully. Processing in background.",
		"regNo":   regNo,
	})
}

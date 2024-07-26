package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"jobApps/internal/database"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"jobApps/authentication"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Controllers interface {
	SignUp(g *gin.Context)
	Login(g *gin.Context)
	Careerpost(g *gin.Context)
	Profile(g *gin.Context)
}

type DbConnection struct {
	Query *database.Queries
}

func ControllerInstance(q *database.Queries) *DbConnection {
	return &DbConnection{
		Query: q,
	}
}

func (db DbConnection) SignUp(g *gin.Context) {
	var users database.CreateUserParams
	if err := g.BindJSON(&users); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  err.Error(),
		})
		return
	}

	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(users.Email) {
		g.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  "invalid email format",
		})
		return
	}
	if users.Username == "" {
		g.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"Error":  "Username field should not be empty",
		})
		return
	}

	if users.Role != "user" && users.Role != "admin" {
		g.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"Error":  "Invalid value for role field.Only 'user' and 'admin' are allowed.",
		})
		return
	}
	//password should have minimum 8 character
	if len(users.Password) < 8 {
		g.JSON(http.StatusBadRequest, gin.H{
			"Error":  "Password should be more than 8 characters",
			"status": 400,
		})
		return
	}

	//passwords are stored in hashing method in the database
	password, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"Error":  "failed to hashing the password",
			"status": 500,
		})
		return
	}
	users.Password = string(password)

	// Validate phone number
	phoneNumber := strings.TrimSpace(users.Phonenumber)
	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		g.JSON(http.StatusBadRequest, gin.H{
			"Error":  "Invalid phone number format",
			"status": 400,
		})
		return
	}

	_, err = db.Query.GetUserByEmail(context.Background(), users.Email)
	if err == nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"Error":  "email ID already exist",
			"status": 500,
		})
		return
	}

	_, err = db.Query.GetUserByPhoneNumber(context.Background(), users.Phonenumber)
	if err == nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"Error":  "user's Phonenumber already exist",
			"status": 500,
		})
		return
	}

	usersData, err := db.Query.CreateUser(context.Background(), users)
	if err != nil {
		fmt.Println("Error", err)
		g.JSON(http.StatusInternalServerError, gin.H{
			"Error":  "failed to create user",
			"status": 500,
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{"Inserted details": usersData})
}

func (db DbConnection) Login(g *gin.Context) {
	var users database.CreateUserParams
	if err := g.BindJSON(&users); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"error":  err.Error(),
		})
		return
	}

	// Retrieve user data based on existing email
	userData, err := db.Query.GetUserByEmail(context.Background(), users.Email)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"error":  "email does not exists",
		})
		return
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(users.Password))
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"error":  "password not matching",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": users.Email,
		"role":  userData.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"status": 500,
			"error":  "Failed to generate token",
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status":    200,
		"message":   "Login successful",
		"user_data": userData,
		"token":     tokenString,
	})
}

func (db DbConnection) GetAllUsersEmail(g *gin.Context) {
	err := authentication.AdminAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can access this endpoint"})
		return
	}

	usersEmail, err := db.Query.GetallusersEmail(context.Background())
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to get users Email ",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "All users email details were retrieved successfully",
		"data":    usersEmail,
	})
}

func (db DbConnection) CreateCareer(g *gin.Context) {
	err := authentication.AdminAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can access this endpoint"})
		return
	}

	var career database.CreateCareerParams

	if err := g.BindJSON(&career); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   err.Error(),
			"message": "Failed to bind JSON data",
		})
		return
	}

	// Validate required fields
	missingFields := []string{}
	if career.Company == "" {
		missingFields = append(missingFields, "Company")
	}
	if career.Position == "" {
		missingFields = append(missingFields, "Position")
	}
	if career.Jobtype == "" {
		missingFields = append(missingFields, "Jobtype")
	}
	if career.Description == "" {
		missingFields = append(missingFields, "Description")
	}
	if career.Startdate.IsZero() {
		missingFields = append(missingFields, "Startdate")
	}
	if career.Enddate.IsZero() {
		missingFields = append(missingFields, "Enddate")
	}

	if len(missingFields) > 0 {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":        400,
			"error":         "Missing fields",
			"missingFields": missingFields,
		})
		return
	}

	if career.Startdate.After(career.Enddate) {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Invalid date range",
			"message": "Start date cannot be greater than end date",
		})
		return
	}
	_, err = db.Query.CreateCareer(context.Background(), career)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"error":   "Failed to create career post",
			"message": err.Error(),
		})
		return
	}

	byteData, err := json.Marshal(&career)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send POST request to webhook server
	_, err = http.Post("http://localhost:9000/webhook", "application/json", bytes.NewBuffer(byteData))
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": " Error sending webhook notification"})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Career post created successfully",
		"data":    career,
	})
}

func (db DbConnection) GetCareerByJobId(g *gin.Context) {
	err := authentication.CommonAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "either admin or user can access this endpoint"})
		return
	}
	jobId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "error on converting string to integer"})
		return
	}
	career, err := db.Query.GetCareerByJobId(context.Background(), int64(jobId))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to get a career details",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "career detail retrieved successfully",
		"data":    career,
	})
}

func (db DbConnection) GetAllCareers(g *gin.Context) {
	err := authentication.CommonAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "either admin or user can access this endpoint"})
		return
	}
	careers, err := db.Query.GetAllCareerDetails(context.Background())
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to get all career details",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "All career details were retrieved successfully",
		"data":    careers,
	})
}

func (db DbConnection) UpdateCareerById(g *gin.Context) {
	err := authentication.AdminAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "only admin can access this endpoint"})
		return
	}

	var career database.UpdateCareerByJobIdParams
	if err := g.BindJSON(&career); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   err.Error(),
			"message": "Failed to bind JSON data",
		})
		return
	}

	jobId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "error on converting string to integer"})
		return
	}

	career.Jobid = int64(jobId)

	existingCareerDetail, _ := db.Query.GetCareerByJobId(context.Background(), career.Jobid)

	if career.Company == "" {
		career.Company = existingCareerDetail.Company
	}
	if career.Position == "" {
		career.Position = existingCareerDetail.Position
	}
	if career.Jobtype == "" {
		career.Jobtype = existingCareerDetail.Jobtype
	}
	if career.Description == "" {
		career.Description = existingCareerDetail.Description
	}
	careerDetail, err := db.Query.UpdateCareerByJobId(context.Background(), career)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to Update career details",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": " career details were retrieved successfully",
		"data":    careerDetail,
	})
}

func (db DbConnection) DeleteCareerById(g *gin.Context) {
	err := authentication.AdminAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "only admin can access this endpoint"})
		return
	}

	jobId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "error on converting string to integer"})
		return
	}

	career, err := db.Query.DeleteCareerByJobId(context.Background(), int64(jobId))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to delete a career detail",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":       200,
		"message":      "career detail deleted successfully",
		"deleted data": career,
	})
}

func (db DbConnection) CreateProfile(g *gin.Context) {
	err := authentication.UserAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "Only users can access this endpoint"})
		return
	}

	var profile database.CreateProfileParams
	if err := g.BindJSON(&profile); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   err.Error(),
			"message": "Failed to bind JSON data",
		})
		return
	}
	// Validate required fields
	missingFields := []string{}
	if profile.Fullname == "" {
		missingFields = append(missingFields, "Fullname")
	}
	if profile.Address == "" {
		missingFields = append(missingFields, "Address")
	}
	if profile.Gender == "" {
		missingFields = append(missingFields, "Gender")
	}
	if profile.Age == 0 {
		missingFields = append(missingFields, "Age")
	}

	if len(missingFields) > 0 {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":        400,
			"error":         "Missing fields",
			"missingFields": missingFields,
		})
		return
	}

	_, err = db.Query.CreateProfile(context.Background(), profile)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"error":   "Failed to create profile",
			"message": err.Error(),
		})
		return
	}

	// Respond with success message
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "profile created successfully",
		"data":    profile,
	})
}

func (db DbConnection) GetProfileById(g *gin.Context) {
	err := authentication.CommonAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "either admin or user can access this endpoint"})
		return
	}
	profileid, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "error on converting string to integer"})
		return
	}
	profile, err := db.Query.GetProfileByuserId(context.Background(), int64(profileid))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to get a Profile details",
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "profile detail retrieved successfully",
		"data":    profile,
	})
}

func (db DbConnection) GetAllProfiles(g *gin.Context) {
	err := authentication.CommonAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "either admin or user can access this endpoint"})
		return
	}
	profile, err := db.Query.GetAllProfileDetails(context.Background())
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to get all profile details",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "All Profile details were retrieved successfully",
		"data":    profile,
	})
}

func (db DbConnection) DeleteProfileById(g *gin.Context) {
	err := authentication.UserAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "only user can access this endpoint"})
		return
	}

	jobId, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "error on converting string to integer"})
		return
	}
	career, err := db.Query.DeleteProfileByUserId(context.Background(), int64(jobId))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to delete a profile detail",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":       200,
		"message":      "profile detail deleted successfully",
		"deleted data": career,
	})
}

func (db DbConnection) UpdateProfileById(g *gin.Context) {
	err := authentication.UserAuth(g)
	if err != nil {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "only user can access this endpoint"})
		return
	}

	var profile database.UpdateProfileByuserIdParams
	if err := g.BindJSON(&profile); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   err.Error(),
			"message": "Failed to bind JSON data",
		})
		return
	}
	userid, err := strconv.Atoi(g.Param("id"))
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "error on converting string to integer"})
		return
	}
	profile.Userid = int64(userid)
	existingProfileDetail, _ := db.Query.GetProfileByuserId(context.Background(), profile.Userid)
	if profile.Fullname == "" {
		profile.Fullname = existingProfileDetail.Fullname
	}
	if profile.Address == "" {
		profile.Address = existingProfileDetail.Address

	}
	if profile.Gender == "" {
		profile.Gender = existingProfileDetail.Gender

	}
	if profile.Age == 0 {
		profile.Age = existingProfileDetail.Age

	}
	profileDetail, err := db.Query.UpdateProfileByuserId(context.Background(), profile)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"error":   "Failed to Update profile details",
			"message": err.Error(),
		})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "profile details were retrieved successfully",
		"data":    profileDetail,
	})
}

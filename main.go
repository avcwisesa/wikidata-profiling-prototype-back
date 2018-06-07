package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type profile struct {
	gorm.Model
	Name       string `json:"name"`
	Class      string `json:"class"`      // JSON string containing name and code for class
	Facets     string `json:"facets"`     // JSON string containing list of facets, a facet got name and code
	Attributes string `json:"attributes"` // JSON string containing list of attributes, an attribute got name and code
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system ENV only.")
	}

	preDBEnv := strings.Split(os.Getenv("DATABASE_URL"), "://")
	dbEnv := strings.Split(preDBEnv[1], "@")
	dbEnv1 := strings.Split(dbEnv[0], ":")
	dbEnv2 := strings.Split(dbEnv[1], "/")
	dbEnv3 := strings.Split(dbEnv2[0], ":")

	client, err := gorm.Open("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
		dbEnv3[0],
		dbEnv2[1],
		dbEnv1[0],
		dbEnv1[1],
	))

	if err != nil {
		log.Println("Error connecting to DB")
		panic(err)
	}

	// Migrate necessary model
	client.AutoMigrate(&profile{})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Ping handler
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, "ProWD prototype database API")
	})

	// new profile handler
	router.POST("profile/new", func(ctx *gin.Context) {

		var newProfile profile
		err := ctx.ShouldBindJSON(&newProfile)
		if err != nil {
			ctx.JSON(400, err)
			return
		}

		profile, err := createProfile(client, newProfile)
		if err != nil {
			ctx.JSON(500, err)
			return
		}

		ctx.JSON(200, profile)
	})

	// get all profile handler
	router.GET("profile", func(ctx *gin.Context) {

		profiles, err := retrieveAllProfile(client)
		if err != nil {
			ctx.JSON(500, err)
		}

		ctx.JSON(200, profiles)
	})

	// get profile by ID handler
	router.GET("profile/:id", func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))

		profile, err := retrieveProfileByID(client, uint(id))
		if err != nil {
			ctx.JSON(500, err)
		}

		ctx.JSON(200, profile)
	})

	router.Run(":" + port)
}

// retrieveAllProfile is a function for retrieving all profile from database
func retrieveAllProfile(client *gorm.DB) ([]profile, error) {

	var retrievedProfiles []profile
	if err := client.Find(&retrievedProfiles).Error; err != nil {
		return nil, err
	}

	return retrievedProfiles, nil
}

// retrieveProfileByID is a function for retrieving profile from database within specified id
func retrieveProfileByID(client *gorm.DB, id uint) (profile, error) {

	var retrievedProfile profile
	if err := client.First(&retrievedProfile, id).Error; err != nil {
		return profile{}, err
	}

	return retrievedProfile, nil
}

// createProfile is a function for storing profile into database
func createProfile(client *gorm.DB, newProfile profile) (profile, error) {

	if err := client.Where(&profile{Name: newProfile.Name}).First(&profile{}).Error; err != nil {
		client.Create(&profile{
			Name:       newProfile.Name,
			Class:      newProfile.Class,
			Facets:     newProfile.Facets,
			Attributes: newProfile.Attributes,
		})
	}

	client.Where(&profile{Name: newProfile.Name}).First(&newProfile)
	return newProfile, nil
}

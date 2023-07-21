package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
)




var Imagemap = make(map[string]string)

// Define a struct to hold the image data
type ImageData struct {
	Hash     string `json:"hash"`
	FileName string `json:"fileName"`
}
const DOWNLOADS_PATH = "pictures/"

// Load the image data from the JSON file at the start of the application
func loadImageDataFromJSON() error {
	data, err := ioutil.ReadFile("imageData.json")
	if err != nil {
		return err
	}

	var images []ImageData
	err = json.Unmarshal(data, &images)
	if err != nil {
		return err
	}

	// Populate the imagemap with the loaded data
	for _, img := range images {
		Imagemap[img.Hash] = img.FileName
	}

	return nil
}

 //saveImageData saves the image to the JSON file
 func saveImageDataToJSON(imagemap map[string]string) error {
	var images []ImageData
	for hash, fileName := range imagemap {
		images = append(images, ImageData{
			Hash:     hash,
			FileName: fileName,
		})
	}

	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("imageData.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

//AddNewImage adds the image to the pictures folder
func AddNewImage(key string, value string) {
	Imagemap[key] = value 
	fmt.Println("New Image Added:", key, value)

	// Save the updated image data to JSON
	err := saveImageDataToJSON(Imagemap)
	if err != nil {
		fmt.Println("Failed to save image data:", err)
	}
}


//DeleteImage is a function for deleting images out of the folder and out of the JSON
func DeleteImage(Imagemap map[string]string, hash string) error {
	// Check if the hash exists in the imagemap
	fmt.Println("Called me 2")
	fileName := Imagemap[hash] 
	fmt.Println("Deleting image with hash:", hash)
	fmt.Println("Corresponding file name:", fileName)
     
	fmt.Println("Hashmap before deleting the file,",Imagemap)
	// Remove the hash from the imagemap using the correct key
	delete(Imagemap, hash)
	fmt.Println("ext is")
	fmt.Println("Hashmap after deleting the file,",Imagemap)
	// Get the file extension from the file name

	// Delete the corresponding image file from the folder
	folderPath := "pictures"
	filePathh := filepath.Join(folderPath, hash)
	if err := os.Remove(filePathh); err != nil {
		return err
	}

	// Save the updated image data to JSON
	err := saveImageDataToJSON(Imagemap)
	if err != nil {
		return err
	}

	return nil
}



//DownloadFile Downloads the image from the 'pictures' folder
func DownloadFile(ctx *gin.Context,filename string){
	ctx.Header("Content-Disposition","attachment; filename=" + filename)
	ctx.Header("Content-Type","application/text/plain")
}



//calculateSHA256Hash converts title to a hash function
func calculateSHA256Hash(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

//saveImageToFolder saves the image to the pictures folder
func saveImageToFolder(img *imageupload.Image, filename string) (*ImageData, error) {
	folderPath := "pictures"
	err := os.MkdirAll(folderPath, 0775)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(folderPath, filename)

	// Check if the file already exists in the folder
	if _, err := os.Stat(filePath); err == nil {
		// File with the same filename already exists
		fmt.Println("Image with the same hash already exists in the folder.")
		return nil, nil
	}

	// File does not exist, save the image to the folder
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Write(img.Data)
	if err != nil {
		return nil, err
	}

	// Extract the file extension
	//ext := filepath.Ext(filename)

	// Create and return the ImageData struct
	imageData := &ImageData{
		Hash:      filename,
		FileName:  filename,
	}

	return imageData, nil
}
//populateImageDataJSON checks if there has been a change in the JSON file and updates it accordingly
func populateImageDataJSON() error {
	folderPath := "pictures"
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	var images []ImageData
	for _, file := range files {
		// Read the image file
		data, err := ioutil.ReadFile(filepath.Join(folderPath, file.Name()))
		if err != nil {
			return err
		}

		// Calculate the hash for the image data
		hash := calculateSHA256Hash(data)

		// Append the image data to the images slice
		images = append(images, ImageData{
			Hash:     hash,
			FileName: file.Name(),
		})
	}

	// Write the image data to the JSON file
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("imageData.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}



//main function handles routes using gin
func main() {
	r := gin.Default()

	// Check if the "pictures" folder exists; if not, create it
	if _, err := os.Stat(DOWNLOADS_PATH); os.IsNotExist(err) {
		fmt.Println("Pictures folder not found, creating...")
		err = os.Mkdir(DOWNLOADS_PATH, 0755)
		if err != nil {
			fmt.Println("Failed to create pictures folder:", err)
			return
		}
		fmt.Println("Pictures folder created.")
	}

	
	// Populate the imageData.json file with current images
	err := populateImageDataJSON()
	if err != nil {
		fmt.Println("Failed to populate image data:", err)
		return
	}

	// Load the image data from JSON
	err = loadImageDataFromJSON()
	if err != nil {
		fmt.Println("Failed to load image data:", err)
		return
	}

	r.GET("/", func(c *gin.Context) {
		r.LoadHTMLGlob("index.html")
	
		// Get the 'sort' query parameter from the request
		sortOption := c.DefaultQuery("sort", "")
	
		// Get the 'images' from the Imagemap
		var images []string
		for _, fileName := range Imagemap {
			images = append(images, fileName)
		}
	
		// Sort the 'images' based on the selected sort option
		switch sortOption {
		case "asc":
			sort.Strings(images)
		case "desc":
			sort.Sort(sort.Reverse(sort.StringSlice(images)))
		}
	
		// Pass the 'images' slice and 'sortOption' to the context for use in the template
		c.HTML(http.StatusOK, "index.html", gin.H{
			"images":    images,
			"sortOption": sortOption,
		})
	})
	


	r.POST("/upload", func(c *gin.Context) {
		img, err := imageupload.Process(c.Request, "file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
			return
		}
	
		hash := calculateSHA256Hash(img.Data)
		fmt.Println("Hash:", hash)
	
		// Check if the image exists in the hashmap
		_, imageExists := Imagemap[hash]
	
		if !imageExists {
			// Image does not exist, upload it and add the hash to the hashmap
			imageData, err := saveImageToFolder(img, hash)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
				return
			}
			AddNewImage(hash, hash)
	
			// Save the updated image data with extension to JSON
			err = saveImageDataToJSON(Imagemap)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image data"})
				return
			}
	
			// Return success response if the image is uploaded successfully
			c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "imageData": imageData})
		} else {
			// Return an alert message to indicate that the image already exists
			c.JSON(http.StatusConflict, gin.H{"message": "Image with the same hash already exists"})
		}
	})
	
	

	r.DELETE("/delete/:hash", func(c *gin.Context) {
		// Get the hash parameter from the URL
		hash := c.Param("hash")
		fmt.Println("Hash received from URL:", hash)
	
		// Delete the image and update the JSON file
		err := DeleteImage(Imagemap, hash)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Image deletion failed", "error": err.Error()})
			return
		}
	
		// Save the updated image data to JSON after successful deletion
		err = saveImageDataToJSON(Imagemap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Image deletion successful, but failed to update JSON", "error": err.Error()})
			return
		}
	
		// Redirect to the home page after successful deletion and JSON update
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/download-user-file/:filename", func(ctx *gin.Context) {
		fileName := ctx.Param("filename")
		ext := filepath.Ext(fileName)
		targetPath := filepath.Join(DOWNLOADS_PATH, fileName + ext)
	
		// Set the appropriate Content-Type based on the file extension
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			// If the Content-Type is unknown, set a default value for binary data
			contentType = "application/octet-stream"
		}
	
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Transfer-Encoding", "binary")
		ctx.Header("Content-Disposition", "attachment; filename="+fileName)
		ctx.Header("Content-Type", contentType)
	
		ctx.File(targetPath)
	})
	
	
	
	r.Run("localhost:8080")
}

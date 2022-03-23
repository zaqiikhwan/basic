package main

import (
	"main.go/authmiddleware"
	"main.go/user"
	"fmt"
	"net/http"
	"strconv"
	"time"
	

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type postBiodataBody struct {
	Nama_Hewan    string `json:"nama_hewan"`
	Umur_Hewan    string `json:"umur_hewan"`
	Jenis_Kelamin string `json:"jenis_kelamin"`
	Jenis_Hewan   string `json:"jenis_hewan"`
	Warna_Hewan   string `json:"warna_hewan"`
}
type Doctor struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Jadwal       string `json:"jadwal"`
	Lokasi_Kerja string `gorm:"lokasi_kerja" json:"lokasi_kerja"`
	Meet         string `gorm:"meet" json:"meet"`
	Picture      string `gorm:"picture" json:"picture"`
	Pengalaman   uint   `gorm:"lama_pengalaman" json:"pengalaman"`
	Price        string `gorm:"price" json:"price"`
}

type selectDoctor struct {
	ID uint
}
type postRegisterBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type postLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type patchUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Scrape struct {
	ID               uint   `gorm:"primarykey" json:"id"`
	Location         string `gorm:"location" json:"location"`
	Name             string `gorm:"name" json:"name"`
	Address          string `gorm:"address" json:"address"`
	Phone_Number     string `gorm:"phone_number" json:"phone_number"`
	Link_Google_Maps string `gorm:"link_google_maps" json:"link_google_maps"`
}

type Article struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Image    string `json:"image"`
	Category string `json:"category"`
}

type Transaction struct {
	ID                uint   `gorm:"primarykey" json:"id"`
	Tanggal_Pemesanan string `gorm:"tanggal_pemesanan" json:"tgl_pesan"`
	Jam_Konsultasi    string `gorm:"jam_konsultasi" json:"jam_konsultasi"`
	Bukti_Pembayaran  string `gorm:"bukti_pembayaran" json:"bukti_pembayaran"`
}

type postTransactionBody struct {
	Tanggal_Pemesanan string `gorm:"tanggal_pemesanan" json:"tgl_pesan"`
	Jam_Konsultasi    string `gorm:"jam_konsultasi" json:"jam_konsultasi"`
}

type searchClinic struct {
	Location string
}

type searchArticle struct {
	ID       uint
	Category string
}

func StartServer() error {
	return r.Run()
}

var db *gorm.DB
var r *gin.Engine

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/project_1?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	err = db.AutoMigrate(&user.User{}, &user.Biodata{}, &Doctor{}, &Transaction{})
	if err != nil {
		return err
	}
	db.Exec("INSERT INTO biodata (id) VALUES (?)", 1)
	return nil
}

func InitGin() {
	r = gin.Default()
	r.Use(CORSPreflightMiddleware())
}

func CORSPreflightMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatus(204)
		} else {
			c.Next()
		}
	}
}

func InitRouter() {
	r.POST("/user/register", func(c *gin.Context) {
		var body postRegisterBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := user.User{
			Name:      body.Name,
			Email:     body.Email,
			Password:  body.Password,
			Username:  body.Username,
			BiodataID: 1,
		}
		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User Registered Successfully",
			"data": gin.H{
				"id": user.ID,
			},
		})
	})

	r.POST("/user/login", func(c *gin.Context) {
		var body postLoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := user.User{}
		if result := db.Where("email = ? ", body.Email).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if user.Password == body.Password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
			})
			tokenString, err := token.SignedString([]byte("passwordBuatSigning"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when generating the token.",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Password is correct.",
				"data": gin.H{
					"id":       user.ID,
					"name":     user.Name,
					"username": user.Username,
					"token":    tokenString,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Password is incorrect.",
			})
			return
		}
	})

	r.GET("/user", authmiddleware.AuthMiddleware(), func(c *gin.Context) {
	
		id, _ := c.Get("id")
		user := user.User{}
		if result := db.Where("id = ?", id).Preload("Biodata").Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		user.Password = ""
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful",
			"data":    user,
		})
	})

	r.GET("/doctor", func(c *gin.Context) {
		var doctors []Doctor
		if result := db.Find(&doctors); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful",
			"data":    doctors,
		})
	})

	r.GET("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		user := user.User{}
		if result := db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    user,
		})
	})

	r.PATCH("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body patchUserBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := user.User{
			ID:       uint(parsedId),
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
			Username: body.Username,
		}
		result := db.Model(&user).Updates(user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", parsedId).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    user,
		})
	})

	r.GET("/user/search", func(c *gin.Context) {
		name, isNameExists := c.GetQuery("name")
		email, isEmailExists := c.GetQuery("email")
		username, isUsernameExists := c.GetQuery("username")
		if !isNameExists && !isEmailExists && !isUsernameExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
			})
			return
		}

		var queryResults []user.User
		trx := db
		if isNameExists {
			trx = trx.Where("name LIKE ?", "%"+name+"%")
		}
		if isEmailExists {
			trx = trx.Where("email LIKE ?", "%"+email+"%")
		}
		if isUsernameExists {
			trx = trx.Where("username LIKE ?", "%"+username+"%")
		}

		if result := trx.Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data": gin.H{
				"query": gin.H{
					"name":     name,
					"email":    email,
					"username": username,
				},
				"result": queryResults,
			},
		})
	})

	r.DELETE("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		user := user.User{
			ID: uint(parsedId),
		}
		if result := db.Delete(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
		})
	})

	r.DELETE("/doctor/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		doctor := Doctor{
			ID: uint(parsedId),
		}
		if result := db.Delete(&doctor); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
		})
	})

	r.GET("/clinic", func(c *gin.Context) {
		var queryResults []Scrape
		trx := db
		if result := trx.Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data":    queryResults,
		})
	})

	r.POST("/clinic/search", func(c *gin.Context) {
		var body searchClinic
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Location is invalid.",
				"error":   err.Error(),
			})
			return
		}
		var queryResults []Scrape
		trx := db
		if result := trx.Where("Location = ?", body.Location).Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data":    queryResults,
		})

	})

	r.GET("/article", func(c *gin.Context) {
		var queryResults []Article
		trx := db
		if result := trx.Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data":    queryResults,
		})
	})

	r.POST("/article/search", func(c *gin.Context) {
		var body searchArticle
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Category is invalid.",
				"error":   err.Error(),
			})
			return
		}
		var queryResults []Article
		trx := db
		if result := trx.Where("ID = ?", body.ID).Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data":    queryResults,
		})
	})

	r.POST("/article/category", func(c *gin.Context) {
		var body searchArticle
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Category is invalid.",
				"error":   err.Error(),
			})
			return
		}
		var queryResults []Article
		trx := db
		if result := trx.Where("Category = ?", body.Category).Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data":    queryResults,
		})

	})

	r.POST("/biodata", func(c *gin.Context) {
		var body postBiodataBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		id, _ := c.Get("id")
		biodata := user.Biodata{
			Nama_Hewan:    body.Nama_Hewan,
			Umur_Hewan:    body.Umur_Hewan,
			Jenis_Kelamin: body.Jenis_Kelamin,
			Jenis_Hewan:   body.Jenis_Hewan,
			Warna_Hewan:   body.Warna_Hewan,
		}
		if result := db.Create(&biodata); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result := db.Model(&user.User{}).Where("id = ?", id).Update("biodata_id", biodata.ID); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Biodata Created Successfully",
			"data":    biodata,
		})
	})

	r.POST("/doctor/search", func(c *gin.Context) {
		var body selectDoctor
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Category is invalid.",
				"error":   err.Error(),
			})
			return
		}
		var queryResults []Doctor
		trx := db
		if result := trx.Where("ID = ?", body.ID).Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data":    queryResults,
		})
	})

	r.Static("/assets", "./assets")
	r.POST("/upload", authmiddleware.AuthMiddleware(), func(c *gin.Context) {
		//Upload file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		if err := c.SaveUploadedFile(file, "./assets/"+file.Filename); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", "./assets/"+file.Filename))
	})

	r.POST("/order/date", authmiddleware.AuthMiddleware(), func(c *gin.Context) {
		var body postTransactionBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Location is invalid.",
				"error":   err.Error(),
			})
			return
		}
		fmt.Println(body.Tanggal_Pemesanan)
		transaction := Transaction{
			Tanggal_Pemesanan: body.Tanggal_Pemesanan,
		}
		if result := db.Create(&transaction); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Date Registered Successfully",
			"data": gin.H{
				"tanggal": transaction.Tanggal_Pemesanan,
			},
		})

	})

	r.POST("/order/time", authmiddleware.AuthMiddleware(), func(c *gin.Context) {
		var body postTransactionBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Location is invalid.",
				"error":   err.Error(),
			})
			return
		}
		transaction := Transaction{
			Jam_Konsultasi: body.Jam_Konsultasi,
		}
		if result := db.Create(&transaction); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Time Registered Successfully",
			"data": gin.H{
				"jam": transaction.Jam_Konsultasi,
			},
		})

	})

}

func main() {
	if err := InitDB(); err != nil {
		fmt.Println("Database error on init!")
		fmt.Println(err.Error())
		return
	}
	InitGin()
	InitRouter()
	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}

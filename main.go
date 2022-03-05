package main

import (
	// "database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	// "strconv"
	// "time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"primarykey" json:"id"`
	// gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Doctor struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
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

type patchDoctorBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Scrape struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	Location     string `gorm:"location" json:"location"`
	Name         string `gorm:"name" json:"name"`
	Address      string `gorm:"address" json:"address"`
	Phone_Number string `gorm:"phone_number" json:"phone_number"`
	Website      string `gorm:"website" json:"website"`
}

type Hospital struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	NameCity     string `json:"name_city"`
	HospitalName string `json:"hospital_name"`
	Contact      string `json:"contact"`
	Address      string `json:"address"`
	Link         string `json:"link"`
}

type postHospitalBody struct {
	NameCity     string `json:"name_city"`
	HospitalName string `json:"hospital_name"`
	Contact      string `json:"contact"`
	Address      string `json:"address"`
	Link         string `json:"link"`
}

type patchHospitalBody struct {
	NameCity     string `json:"name_city"`
	HospitalName string `json:"hospital_name"`
	Contact      string `json:"contact"`
	Address      string `json:"address"`
	Link         string `json:"link"`
}

type Article struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Link     string `json:"link"`
}

type postArticleBody struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Link     string `json:"link"`
}

type searchClinic struct {
	Location string
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
	err = db.AutoMigrate(&Hospital{}, &User{}, &Doctor{}, &Article{})
	if err != nil {
		return err
	}
	return nil
}

func InitGin() {
	r = gin.Default()
	r.Use(cors.Default())
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]
		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte("passwordBuatSigning"), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "JWT validation error.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", claims["id"])
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "JWT invalid.",
				"error":   err.Error(),
			})
			c.Abort()
			return
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
		user := User{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
			Username: body.Username,
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
			"message": "User Registered successfully",
			"data": gin.H{
				"id": user.ID,
			},
		})
	})

	r.POST("/doctor/register", func(c *gin.Context) {
		var body postRegisterBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		doctor := Doctor{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
			Username: body.Username,
		}
		if result := db.Create(&doctor); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User Registered successfully",
			"data": gin.H{
				"id": doctor.ID,
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
		user := User{}
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

	r.POST("/doctor/login", func(c *gin.Context) {
		var body postLoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		doctor := Doctor{}
		if result := db.Where("email = ? ", body.Email).Take(&doctor); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if doctor.Password == body.Password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  doctor.ID,
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
					"id":       doctor.ID,
					"name":     doctor.Name,
					"username": doctor.Username,
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

	r.GET("/user", AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
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
			"message": "Query successful",
			"data":    user,
		})
	})

	r.GET("/doctor", AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		doctor := Doctor{}
		if result := db.Where("id = ?", id).Take(&doctor); result.Error != nil {
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
			"data":    doctor,
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
		user := User{}
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

	r.GET("/doctor/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		doctor := Doctor{}
		if result := db.Where("id = ?", id).Take(&doctor); result.Error != nil {
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
			"data":    doctor,
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
		user := User{
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

	r.PATCH("/doctor/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body patchDoctorBody
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
		doctor := Doctor{
			ID:       uint(parsedId),
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
			Username: body.Username,
		}
		result := db.Model(&doctor).Updates(doctor)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", parsedId).Take(&doctor); result.Error != nil {
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
			"data":    doctor,
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

		var queryResults []User
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

	r.GET("/doctor/search", func(c *gin.Context) {
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

		var queryResults []Doctor
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
		user := User{
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

	r.POST("/hospital/register", func(c *gin.Context) {
		var body postHospitalBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		hospital := Hospital{
			NameCity:     body.NameCity,
			HospitalName: body.HospitalName,
			Contact:      body.Contact,
			Address:      body.Address,
			Link:         body.Link,
		}
		if result := db.Create(&hospital); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Hospital registered successfully.",
			"data": gin.H{
				"id": hospital.ID,
			},
		})
	})

	r.GET("/hospital/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		hospital := Hospital{}
		if result := db.Where("id = ?", id).Take(&hospital); result.Error != nil {
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
			"data":    hospital,
		})
	})

	r.GET("/hospital", func(c *gin.Context) {
		var queryResults []Hospital
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
			"data": gin.H{
				"result": queryResults,
			},
		})
	})

	// var body searchClinic
	// if err := c.BindJSON(&body); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"success": false,
	// 		"message": "Location is invalid.",
	// 		"error":   err.Error(),
	// 	})
	// 	return
	//}

	r.POST("/clinic", func(c *gin.Context) {
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

		// if result := db.Where("Location LIKE ?", "%"+body.Location+"%"); result.Error != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"success": false,
		// 		"message": "Error when inserting into the database.",
		// 		"error":   result.Error.Error(),
		// 	})
		// 	return
		// }

		trx := db
		// if isLocationExists {
		// 	trx = trx.Where("Location LIKE ?", "%"+location+"%")
		// }

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

	r.GET("/clinic/search", func(c *gin.Context) {
		location, isLocationExists := c.GetQuery("Location")
		if !isLocationExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
			})
			return
		}

		var queryResults []Scrape
		trx := db
		if isLocationExists {
			trx = trx.Where("Location LIKE ?", "%"+location+"%")
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
			"data":    queryResults,
		})
	})

	r.GET("/hospital/search", func(c *gin.Context) {
		namakota, isNamaKotaExists := c.GetQuery("name_city")
		if !isNamaKotaExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
			})
			return
		}

		var queryResults []Hospital
		trx := db
		if isNamaKotaExists {
			trx = trx.Where("name_city LIKE ?", "%"+namakota+"%")
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
				"result": queryResults,
			},
		})
	})

	r.PATCH("/hospital/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body patchHospitalBody
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
		hospital := Hospital{
			ID:           uint(parsedId),
			NameCity:     body.NameCity,
			HospitalName: body.HospitalName,
			Contact:      body.Contact,
			Address:      body.Address,
			Link:         body.Link,
		}
		result := db.Model(&hospital).Updates(hospital)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", parsedId).Take(&hospital); result.Error != nil {
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
				"message": "Hospital not found.",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    hospital,
		})
	})

	r.POST("/article", func(c *gin.Context) {
		var body postArticleBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		article := Article{
			Title:    body.Title,
			Content:  body.Content,
			Category: body.Category,
			Link:     body.Link,
		}
		if result := db.Create(&article); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Article Posted successfully",
			"data": gin.H{
				"id": article.ID,
			},
		})
	})

	r.GET("/article/search", func(c *gin.Context) {
		title, isTitleExists := c.GetQuery("title")
		category, isCategoryExists := c.GetQuery("category")
		var queryResults []Article
		trx := db
		if isTitleExists {
			trx = trx.Where("title LIKE ?", "%"+title+"%")
		}
		if isCategoryExists {
			trx = trx.Where("category LIKE ?", "%"+category+"%")
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
			"data":    queryResults,
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
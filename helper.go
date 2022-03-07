package main

import (
	"backendtest/database"
	"backendtest/models"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

//jwt için secret key
var SecretKey = []byte("secret")

// Register fonksiyonu
func Register(c *fiber.Ctx) error {
	//body içindeki verileri alıyoruz
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	//db bağlantısını oluşturuyoruz
	db := database.Connect()

	//body içindeki password adlı veriyi alıyoruz ve onu hashleyip veritabanına kaydediyoruz
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	//body içindeki username adlı veriyi alıyoruz
	userData := models.User{
		UserName: data["username"],
		Password: string(password),
	}

	//veritabanına kayıt işlemi yapıyoruz
	db.Save(&userData)

	return c.JSON(userData)

}

// Setup fonksiyonu
func Login(c *fiber.Ctx) error {
	//body içindeki verileri alıyoruz
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// body içindeki username adlı veriyi alıyoruz
	user := models.User{
		UserName: data["username"],
	}

	// db bağlantısını oluşturuyoruz
	db := database.Connect()

	var userDB models.User

	// veritabanından username ile kayıt getiriyoruz
	db.Model(&userDB).Where("user_name = ?", user.UserName).First(&userDB)

	// body içindeki password adlı verinin hashi dbdeki veri ile karşılaştırılıyor
	err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(data["password"]))

	// karşılaştırma başarısız ise
	if err != nil {
		// giriş yapamadınız döndürüyoruz
		return c.Status(401).SendString("Invalid username or password")
	}
	//karşılaştırma başarılı ise

	//token oluşturuyoruz
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), //tokenın expire tarihi 1 saat sonra olacak
		Id:        fmt.Sprintf("%d", userDB.ID),         //tokenın id'si userDB.ID olarak ayarlıyoruz
	})

	//tokenın string halini alıyoruz
	token, err := jwtToken.SignedString(SecretKey)
	if err != nil {
		return err

	}

	//tokenı cookie olarak belirliyoruz
	cookie := fiber.Cookie{
		Name:     "jwt",                         //cookie adı
		Value:    token,                         //tokenın string halini
		Expires:  time.Now().Add(time.Hour * 1), //tokenın expire tarihi 1 saat sonra olacak
		HTTPOnly: true,                          //cookie httponly olarak ayarlıyoruz
	}

	//cookieyi ctx'e ekliyoruz
	c.Cookie(&cookie)

	return c.Status(202).JSON(userDB)

}

// // main
// count := 0

// //login
// count = 1
// //logout
// count = 0

// User fonksiyonu
func User(c *fiber.Ctx) error {
	//db bağlantısını oluşturuyoruz
	db := database.Connect()

	//tokenı cookie olarak alıyoruz
	cookie := c.Cookies("jwt")

	//tokenı parse ediyoruz
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	//token parse edilirken hata alınırsa
	if err != nil {
		return c.Status(403).SendString("Giriş yapamadınız lütfen önce giriş yapınız")
	}

	//token parse edilirken hata alınmadıysa

	// token içindeki id adlı veriyi alıyoruz
	claims := token.Claims.(*jwt.StandardClaims)

	//veritabanından id ile kayıt getiriyoruz jwtdeki token paremetresine göre
	userJwt := models.User{}

	db.Where("id = ?", claims.Id).First(&userJwt)

	return c.JSON(userJwt)

}

// Logout fonksiyonu
func Logout(c *fiber.Ctx) error {
	//tokenı cookie olarak alıyoruz
	cookieJwt := c.Cookies("jwt")

	//tokenı parse ediyoruz
	_, err := jwt.ParseWithClaims(cookieJwt, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	//token parse edilirken hata alınırsa

	// hata mesajı
	if err != nil {
		return c.Status(403).SendString("Giriş yapamadınız lütfen önce giriş yapınız")
	}

	//token parse edilirken hata alınmadıysa
	//cookie sonlanıyor
	cookie := fiber.Cookie{
		Name:    "jwt",                          //cookie adı
		Value:   "",                             //cookie değeri sıfırlıyoruz
		Expires: time.Now().Add(-1 * time.Hour), //cookie sonlandırmak için expiry tarihi 1 saat önceye ayarlıyoruz
	}

	//cookieyi ctx'e ekliyoruz
	c.Cookie(&cookie)

	return c.Status(202).SendString("Success")
}

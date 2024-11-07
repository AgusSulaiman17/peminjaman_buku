package controllers

import (
    "net/http"
    "time"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
    "github.com/labstack/echo/v4"
    "buku/config"
    "buku/models"
    "buku/utils"
)

// Fungsi Register menangani registrasi pengguna
func Register(c echo.Context) error {
    var user models.User

    // Bind data input ke variabel user
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
    }

    // Cek apakah email sudah terdaftar
    var existingUser models.User
    if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
        return c.JSON(http.StatusConflict, map[string]string{"message": "Email sudah terdaftar"})
    }

    // Hash password menggunakan bcrypt untuk keamanan
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.KataSandi), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Tidak dapat hash password"})
    }
    user.KataSandi = string(hashedPassword)

    // Tentukan timestamp dibuat dan diperbarui
    user.DibuatPada = time.Now()
    user.DiperbaruiPada = time.Now()

    // Simpan data pengguna ke dalam database
    if err := config.DB.Create(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Tidak dapat membuat pengguna"})
    }

    return c.JSON(http.StatusCreated, user)
}

// Fungsi Login menangani login pengguna
func Login(c echo.Context) error {
    var user models.User
    var dbUser models.User

    // Bind data input ke variabel user
    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Input tidak valid"})
    }

    // Cari pengguna berdasarkan email
    if err := config.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Kredensial tidak valid"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan saat mencari pengguna"})
    }

    // Bandingkan password hash
    err := bcrypt.CompareHashAndPassword([]byte(dbUser.KataSandi), []byte(user.KataSandi))
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Kredensial tidak valid"})
    }

    // Buat token JWT untuk autentikasi
    token, err := utils.GenerateJWT(dbUser)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan saat membuat token"})
    }

    return c.JSON(http.StatusOK, map[string]string{"token": token})
}

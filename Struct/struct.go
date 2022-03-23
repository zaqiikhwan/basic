package Struct

type User struct {
	ID uint `gorm:"primarykey" json:"id"`
	// gorm.Model
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Username  string  `json:"username"`
	Biodata   Biodata `json:"biodata"`
	BiodataID uint
}
type Biodata struct {
	ID            uint   `gorm:"primarykey" json:"id"`
	Nama_Hewan    string `json:"nama_hewan"`
	Umur_Hewan    string `json:"umur_hewan"`
	Jenis_Kelamin string `json:"jenis_kelamin"`
	Jenis_Hewan   string `json:"jenis_hewan"`
	Warna_Hewan   string `json:"warna_hewan"`
}
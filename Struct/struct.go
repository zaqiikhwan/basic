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
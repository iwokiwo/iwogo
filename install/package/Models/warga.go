package Models

type Warga struct {
	Nama         string `json:"nama"`
	NoKtp        string `json:"no_ktp"`
	Alamat       string `json:"alamat"`
	GolDarah     string `json:"gol_darah"`
	TglLahir     string `json:"tgl_lahir"`
	HubKeluarga  string `json:"hub_keluarga"`
	CodeHub      int    `json:"code_hub"`
	Blok         string `json:"blok"`
	Tlp          string `json:"tlp"`
	CodeUser     int    `json:"code_user"`
	AksesLogin   int    `json:"akses_login"`
	JenisKelamin string `json:"jenis_kelamin"`
	TglTagihan   string `json:"tgl_tagihan"`
	User         User   `json:"user" gorm:"foreignKey:code_user;references:id"`
	Base
}

func (b *Warga) TableName() string {
	return "wargas"
}

// id                  Int       @id @default(autoincrement())
// nama                String    @db.VarChar(255)
// no_ktp              String?   @unique @db.VarChar(100)
// alamat              String?   @db.VarChar(255)
// gol_darah           String?   @db.VarChar(2)
// tgl_lahir           DateTime? @db.Date
// hub_keluarga        String?    @db.VarChar(50)
// code_hub            Int?
// blok                String?    @db.VarChar(50)
// tlp                 String?    @db.VarChar(50)
// code_user           Int?
// akses_login         Int?
// jenis_kelamin       String?    @db.VarChar(10)
// tgl_tagihan         DateTime?  @db.Date
// createdAt           DateTime   @default(now())
// updatedAt           DateTime   @updatedAt
// rfids  rfids[]

package entities

type Faculty struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type StudyPrograms struct {
	Faculty string `json:"faculty"`
	Value   string `json:"value"`
	Name    string `json:"name"`
}

type RowData struct {
	Hari     string `json:"hari"`
	Kode     string `json:"kode"`
	Matkul   string `json:"matkul"`
	Kelas    string `json:"kelas"`
	Sks      string `json:"sks"`
	Jam      string `json:"jam"`
	Semester string `json:"semester"`
	Dosen    string `json:"dosen"`
	Ruang    string `json:"ruang"`
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Tabungan struct {
	Penghasilan float64            `json:"penghasilan"`
	Konsumtif   map[string]float64 `json:"konsumtif"`
	Alokasi     map[string]float64 `json:"alokasi"`
}

type TabunganRequest struct {
	Penghasilan float64            `json:"penghasilan"`
	Konsumtif   map[string]float64 `json:"konsumtif"`
	Alokasi     map[string]float64 `json:"alokasi"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"pesan"`
	Data    Data
}

type Data struct {
	Penghasilan             string                 `json:"penghasilan"`
	TotalPengeluaranHarian  string                 `json:"totalpengeluaranperhari"`
	TotalPengeluaranBulanan string                 `json:"totalpengeluaranperbulan"`
	TotalPengeluaranTahunan string                 `json:"totalpengeluaranpertahun"`
	SisaUang                string                 `json:"sisauang"`
	Konsumtif               map[string]interface{} `json:"konsumtif"`
}

func NewTabungan(req *TabunganRequest) (*Tabungan, error) {
	if err := InputOuputValidasi(req); err != nil {
		return nil, err
	}
	tabungan := &Tabungan{
		Penghasilan: req.Penghasilan,
		Konsumtif:   req.Konsumtif,
		Alokasi:     req.Alokasi,
	}
	return tabungan, nil
}

func InputOuputValidasi(req *TabunganRequest) error {
	if req.Penghasilan <= 0 {
		return errors.New("penghasilan is required and must be greater than 0")
	}
	for key, value := range req.Konsumtif {
		if value <= 0 {
			return fmt.Errorf("%s must be greater than or equal to 0", key)
		}
	}
	for l, v := range req.Alokasi {
		if v <= 0 || v >= 100 {
			return fmt.Errorf("%s harus antara 0 sampai 100", l)
		}
	}
	return nil
}

func GetTotalBulanan(tabungan *Tabungan) float64 {
	total := 0.0
	for _, v := range tabungan.Konsumtif {
		total += v
	}
	return total
}

func GetTotalHarian(tabungan *Tabungan) float64 {
	totalHarian := GetTotalBulanan(tabungan)
	return totalHarian / 30
}

func GetTotalTahunan(tabungan *Tabungan) float64 {
	totalTahuan := GetTotalBulanan(tabungan)
	return totalTahuan * 120
}

func FormatRupiah(amount float64) string {
	amountStr := fmt.Sprintf("%.0f", amount)
	var result strings.Builder
	length := len(amountStr)
	for i, char := range amountStr {
		if (length-i)%3 == 0 && i != 0 {
			result.WriteString(".")
		}
		result.WriteRune(char)
	}
	return "Rp" + result.String()
}

func HandleHitung(w http.ResponseWriter, r *http.Request) {
	var req TabunganRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tabungan, err := NewTabungan(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	totalBulanan := GetTotalBulanan(tabungan)

	if totalBulanan >= req.Penghasilan {
		response := Response{
			Message: "Gaji bulanan Anda tidak cukup untuk menabung untuk nikah.",
			Status:  "error",
			Data:    Data{},
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	totalHarian := GetTotalHarian(tabungan)
	totalTahunan := GetTotalTahunan(tabungan)

	sisaUang := tabungan.Penghasilan - totalBulanan

	// semua alokasi nya dibagikan dengan berapapun persentasenya
	alokasi := make(map[string]string)
	for key, persen := range tabungan.Alokasi {
		alokasi[key] = FormatRupiah(sisaUang * persen / 100)
	}

	// ambil semua kosumtifnya
	konsumtif := make(map[string]interface{})
	for key, value := range tabungan.Konsumtif {
		konsumtif[key] = FormatRupiah(value)
	}

	konsumtif["alokasi"] = alokasi

	data := Data{
		Penghasilan:             FormatRupiah(tabungan.Penghasilan),
		TotalPengeluaranBulanan: FormatRupiah(totalBulanan),
		TotalPengeluaranHarian:  FormatRupiah(totalHarian),
		TotalPengeluaranTahunan: FormatRupiah(totalTahunan),
		SisaUang:                FormatRupiah(sisaUang),
		Konsumtif:               konsumtif,
	}

	response := Response{
		Message: "data berhasil di hitung",
		Status:  "success",
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func main() {
	http.HandleFunc("/hitung", HandleHitung)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

package entities

import (
	"errors"
	"fmt"
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

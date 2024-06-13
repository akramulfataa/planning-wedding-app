package handlers

import (
	"encoding/json"
	"net/http"
	"wedding-planner/entities"
)

func HandleHitung(w http.ResponseWriter, r *http.Request) {
	var req entities.TabunganRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	tabungan, err := entities.NewTabungan(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	totalBulanan := entities.GetTotalBulanan(tabungan)
	if totalBulanan >= req.Penghasilan {
		response := entities.Response{
			Message: "Gaji bulanan Anda tidak cukup untuk menabung untuk nikah.",
			Status:  "error",
			Data:    entities.Data{},
		}
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	totalHarian := entities.GetTotalHarian(tabungan)
	totalTahunan := entities.GetTotalTahunan(tabungan)

	sisaUang := tabungan.Penghasilan - totalBulanan

	// semua alokasi nya dibagikan dengan berapapun persentasenya
	alokasi := make(map[string]string)
	for key, persen := range tabungan.Alokasi {
		alokasi[key] = entities.FormatRupiah(sisaUang * persen / 100)
	}

	// ambil semua kosumtifnya
	konsumtif := make(map[string]interface{})
	for key, value := range tabungan.Konsumtif {
		konsumtif[key] = entities.FormatRupiah(value)
	}

	konsumtif["alokasi"] = alokasi

	data := entities.Data{
		Penghasilan:             entities.FormatRupiah(tabungan.Penghasilan),
		TotalPengeluaranBulanan: entities.FormatRupiah(totalBulanan),
		TotalPengeluaranHarian:  entities.FormatRupiah(totalHarian),
		TotalPengeluaranTahunan: entities.FormatRupiah(totalTahunan),
		SisaUang:                entities.FormatRupiah(sisaUang),
		Konsumtif:               konsumtif,
	}

	response := entities.Response{
		Message: "data berhasil di hitung",
		Status:  "success",
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

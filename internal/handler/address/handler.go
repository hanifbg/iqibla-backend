package address

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func SearchAddress(c echo.Context) error {
	kecamatan := c.QueryParam("kecamatan")
	kota := c.QueryParam("kota")

	// Hardcoded data
	addresses := []map[string]string{
		{
			"kecamatan": "Cempaka Putih",
			"kota":      "Jakarta Pusat",
			"alamat":    "Jl. Cempaka Putih Timur No. 1",
		},
		{
			"kecamatan": "Wiyung",
			"kota":      "Surabaya",
			"alamat":    "Jl. Babatan Indah",
		},
		{
			"kecamatan": "Kramat Jati",
			"kota":      "Jakarta Timur",
			"alamat":    "Jl. Batu Ampar 3 No. 1",
		},
		{
			"kecamatan": "Krembangan",
			"kota":      "Surabaya",
			"alamat":    "Jl. Tambak Asri",
		},
	}

	// Validasi input
	if kecamatan == "" && kota == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "At least one of kecamatan or kota is required"})
	}

	// Filter hasil berdasarkan input
	var results []map[string]string
	for _, addr := range addresses {
		matchKecamatan := true
		matchKota := true

		// Cek kecamatan (jika ada input)
		if kecamatan != "" {
			matchKecamatan = strings.Contains(strings.ToLower(addr["kecamatan"]), strings.ToLower(kecamatan))
		}

		// Cek kota (jika ada input)
		if kota != "" {
			matchKota = strings.Contains(strings.ToLower(addr["kota"]), strings.ToLower(kota))
		}

		// Tambahkan ke hasil jika memenuhi kriteria
		if matchKecamatan && matchKota {
			results = append(results, addr)
		}
	}

	if len(results) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No addresses found"})
	}

	return c.JSON(http.StatusOK, results)
}

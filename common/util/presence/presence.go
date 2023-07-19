package presence

import "time"

func GetCurrentMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// Fungsi untuk mengecek apakah tanggal dalam range start dan end date
func IsDateInRange(inputDate, startDate, endDate string) (bool, error) {
	// Mengonversi string menjadi time.Time
	date, err := parseDate(inputDate)
	if err != nil {
		return false, err
	}

	start, err := parseDate(startDate)
	if err != nil {
		return false, err
	}

	end, err := parseDate(endDate)
	if err != nil {
		return false, err
	}

	// Melakukan pengecekan apakah tanggal berada dalam range start dan end date
	return date.After(start) && date.Before(end) || date.Equal(start) || date.Equal(end), nil
}

package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/zsefvlol/timezonemapper"
)

func GetTimeZone(latitude, longitude float64) (timeZoneCode int) {
	timezone := timezonemapper.LatLngToTimezoneString(latitude, longitude)
	loc, _ := time.LoadLocation(timezone)
	// Parse time string with location
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), loc)
	_, offset := t.Zone()
	hoursGMT := int(offset / 3600)
	return hoursGMT
}

func MillisToTimeString(timeMillis int64, timeZone int) (timeString string) {
	if timeMillis <= 0 {
		return "--:--"
	}
	if timeZone == 0 {
		dateTime := time.Unix(0, timeMillis*int64(time.Millisecond))
		return dateTime.Format("15:04")
	} else {
		millisTimeZone := int64(timeZone) * 3600000
		millisFinal := timeMillis + millisTimeZone
		dateTime := time.Unix(0, millisFinal*int64(time.Millisecond))
		return dateTime.UTC().Format("15:04")
	}
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func GetDayName(myTime time.Time) (dayName string) {
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	return days[int(myTime.Weekday())]
}

func GetDayNameFromDateString(date string) (dayName string) {
	myTime, _ := time.Parse("2006-01-02", GetOnlyDateString(date))
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	return days[int(myTime.Weekday())]
}

func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

func GetOnlyDateString(dateTimeString string) (dateString string) {
	if isContains := strings.Contains(dateTimeString, "T"); isContains {
		datas := strings.Split(dateTimeString, "T")
		return datas[0]
	} else {
		return dateTimeString
	}
}

func GetDatesArray(bulan, tahun int) []string {
	var dates []string

	// Membuat tanggal pertama pada bulan dan tahun yang diberikan
	tanggalPertama := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC)

	// Membuat tanggal terakhir pada bulan dan tahun yang diberikan
	// Dapatkan bulan berikutnya dan kurangi 1 hari untuk mendapatkan tanggal terakhir pada bulan yang diberikan
	tanggalTerakhir := tanggalPertama.AddDate(0, 1, -1)

	// Iterasi dari tanggal pertama hingga tanggal terakhir untuk membuat array tanggal
	for d := tanggalPertama; d.Before(tanggalTerakhir) || d.Equal(tanggalTerakhir); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates
}

func FormatTanggalIndonesia(tanggal string) (string, error) {
	// Parsing string tanggal menjadi tipe time.Time
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return "", err
	}

	// Array nama-nama hari dalam Bahasa Indonesia
	hariIndonesia := [...]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

	// Array nama-nama bulan dalam Bahasa Indonesia
	bulanIndonesia := [...]string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

	// Mendapatkan nama hari dan bulan dalam Bahasa Indonesia
	namaHari := hariIndonesia[t.Weekday()]
	namaBulan := bulanIndonesia[t.Month()]

	// Mendapatkan tanggal, bulan, dan tahun
	tanggalIndonesia := fmt.Sprintf("%s, %d %s %d", namaHari, t.Day(), namaBulan, t.Year())

	return tanggalIndonesia, nil
}

func GetEnglishDayName(tanggal string) (string, error) {
	// Parsing string tanggal menjadi tipe time.Time
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return "", err
	}

	// Array nama-nama hari dalam Bahasa Inggris lowercase
	hariInggris := [...]string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	// Mendapatkan indeks hari dalam tipe time.Weekday (mulai dari 0 untuk Minggu)
	indeksHari := int(t.Weekday())

	// Mendapatkan nama hari dalam Bahasa Inggris lowercase
	namaHari := hariInggris[indeksHari]

	return namaHari, nil
}

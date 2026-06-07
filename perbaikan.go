package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const MAXDATA = 1000

// Tipe bentukan dalam program
type Milestone struct {
	ID             int
	NamaTugas      string
	Deskripsi      string
	Prioritas      int
	SkorStres      int
	SkorMood       int
	CatatanRasa    string
	TanggalSelesai string
	Selesai        bool
}

// Tipe bentukan untuk array statis pada tugas
type ArrayTugas struct {
	Data    [MAXDATA]Milestone
	Jumlah  int
	Counter int
}

var tugasDB ArrayTugas //  VARIABEL GLOBAL — untuk array utama saja

const ( //  Warna Tampilan
	warnaNormal = "\033[0m"
	warnaBold   = "\033[1m"
	warnaCyan   = "\033[36m"
	warnaHijau  = "\033[32m"
	warnaKuning = "\033[33m"
	warnaMerah  = "\033[31m"
	warnaAbu    = "\033[90m"
	warnaBiru   = "\033[34m"
)

//  SUBPROGRAM: I/O Program
//  Parameter: pertanyaan dengan tipe data string, reader dengan menggunakan *bufio.Reader
//  Spesifikasi: Membaca input teks dari pengguna

func inputStr(pertanyaan string, reader *bufio.Reader) string {
	fmt.Print(pertanyaan)
	teks, _ := reader.ReadString('\n')
	return strings.TrimSpace(teks)
}

// Spesifikasi: Membaca input angka, akan meminta ulangi input jika bukan angka valid
func inputAngka(pertanyaan string, reader *bufio.Reader) int {
	valid := false
	hasil := 0
	for !valid {
		s := inputStr(pertanyaan, reader)
		n, err := strconv.Atoi(s)
		if err == nil {
			hasil = n
			valid = true
		} else {
			fmt.Println(" Coba Lagi Nyak Angka Tidak Valid")
		}
	}
	return hasil
}

// Spesifikasi: Membaca angka dalam rentang [min, max], ulangi jika di luar rentang
func inputAngkaRentang(pertanyaan string, min, max int, reader *bufio.Reader) int {
	valid := false
	hasil := 0
	for !valid {
		n := inputAngka(pertanyaan, reader)
		if n >= min && n <= max {
			hasil = n
			valid = true
		} else {
			fmt.Printf(" Angkanya Harus Antara %d Sampai %d!\n", min, max)
		}
	}
	return hasil
}

//  SUBPROGRAM: Tampilan
//  Spesifikasi: Mencetak berbagai elemen UI ke terminal

func judulBesar(teks string) {
	garis := strings.Repeat("─", 60)
	fmt.Printf("\n%s%s%s\n", warnaCyan, garis, warnaNormal)
	fmt.Printf("%s  ◆ %s%s\n", warnaBold, teks, warnaNormal)
	fmt.Printf("%s%s%s\n", warnaCyan, garis, warnaNormal)
}

func judulKecil(teks string) {
	fmt.Printf("\n%s── %s%s\n", warnaAbu, teks, warnaNormal)
}

func pesanOke(pesan string)        { fmt.Printf("%s  %s%s\n", warnaHijau, pesan, warnaNormal) }
func pesanPeringatan(pesan string) { fmt.Printf("%s  %s%s\n", warnaKuning, pesan, warnaNormal) }
func pesanInfo(pesan string)       { fmt.Printf("%s    %s%s\n", warnaBiru, pesan, warnaNormal) }
func tidakKetemu()                 { pesanPeringatan("Hmm, datanya tidak ketemu nih.") }

// Spesifikasi: Mengembalikan label prioritas berwarna
func labelPrioritas(p int) string {
	label := "?"
	if p == 1 {
		label = warnaAbu + "Santai aja" + warnaNormal
	} else if p == 2 {
		label = warnaKuning + "Lumayan penting" + warnaNormal
	} else if p == 3 {
		label = warnaMerah + "Penting banget!" + warnaNormal
	}
	return label
}

// Spesifikasi: Mengembalikan bar visual skor 1-10
func barMood(skor int) string {
	isi := skor
	kosong := 10 - skor
	return warnaHijau + strings.Repeat("█", isi) + warnaAbu + strings.Repeat("░", kosong) + warnaNormal
}

// Spesifikasi: Mengembalikan label stres berwarna berdasarkan nilai
func labelStres(s int) string {
	label := ""
	if s <= 3 {
		label = warnaHijau + strconv.Itoa(s) + "/10 (Tenang)" + warnaNormal
	} else if s <= 6 {
		label = warnaKuning + strconv.Itoa(s) + "/10 (Lumayan stres)" + warnaNormal
	} else {
		label = warnaMerah + strconv.Itoa(s) + "/10 (Stres tinggi!)" + warnaNormal
	}
	return label
}

// Spesifikasi: Mencetak satu tugas secara lengkap ke terminal
func tampilkanTugas(m Milestone) {
	status := warnaAbu + "○ Belum selesai" + warnaNormal
	if m.Selesai {
		status = warnaHijau + "✔ Sudah selesai!" + warnaNormal
	}
	fmt.Printf("\n  %s[ID: %d]%s %s%s%s\n", warnaBold, m.ID, warnaNormal, warnaBold, m.NamaTugas, warnaNormal)
	fmt.Printf("      Keterangan   : %s\n", m.Deskripsi)
	fmt.Printf("      Prioritas    : %s\n", labelPrioritas(m.Prioritas))
	fmt.Printf("      Tingkat Stres: %s\n", labelStres(m.SkorStres))
	fmt.Printf("      Mood/Perasaan: %s (%d/10)\n", barMood(m.SkorMood), m.SkorMood)
	fmt.Printf("      Curhatan     : %s\n", m.CatatanRasa)
	fmt.Printf("      Tanggal Kelar: %s\n", m.TanggalSelesai)
	fmt.Printf("      Status       : %s\n", status)
}

// Spesifikasi: Mencetak daftar tugas dari array statis sejumlah n elemen
func tampilkanArrayTugas(arr [MAXDATA]Milestone, n int) {
	i := 0
	for i < n {
		tampilkanTugas(arr[i])
		i++
	}
}

// Spesifikasi: Mencari tugas berdasarkan nama secara linear/sequential
func cariSequential(db ArrayTugas, katakunci string) ([MAXDATA]Milestone, int) {
	var hasil [MAXDATA]Milestone
	jumlahHasil := 0
	katakunci = strings.ToLower(katakunci)
	i := 0
	for i < db.Jumlah {
		if strings.Contains(strings.ToLower(db.Data[i].NamaTugas), katakunci) {
			hasil[jumlahHasil] = db.Data[i]
			jumlahHasil++
		}
		i++
	}
	return hasil, jumlahHasil
}

// Spesifikasi: Mengurutkan salinan array berdasarkan tanggal agar bisa di-binary search
func urutTanggalUntukBinarySearch(arr [MAXDATA]Milestone, n int) [MAXDATA]Milestone {
	salinan := arr
	// Gunakan insertion sort untuk mengurutkan tanggal
	i := 1
	for i < n {
		kunci := salinan[i]
		j := i - 1
		for j >= 0 && salinan[j].TanggalSelesai > kunci.TanggalSelesai {
			salinan[j+1] = salinan[j]
			j--
		}
		salinan[j+1] = kunci
		i++
	}
	return salinan
}

// Spesifikasi: Mencari tugas berdasarkan tanggal dengan metode biner (data diurutkan dulu)
func cariBinary(db ArrayTugas, tanggal string) ([MAXDATA]Milestone, int) {
	var hasil [MAXDATA]Milestone
	jumlahHasil := 0

	terurut := urutTanggalUntukBinarySearch(db.Data, db.Jumlah)

	kiri := 0
	kanan := db.Jumlah - 1
	idxTengah := -1

	for kiri <= kanan && idxTengah == -1 {
		tengah := (kiri + kanan) / 2
		if terurut[tengah].TanggalSelesai == tanggal {
			idxTengah = tengah
		} else if terurut[tengah].TanggalSelesai < tanggal {
			kiri = tengah + 1
		} else {
			kanan = tengah - 1
		}
	}

	if idxTengah != -1 {
		k := idxTengah
		for k >= 0 && terurut[k].TanggalSelesai == tanggal {
			k--
		}
		awal := k + 1
		idx := awal
		for idx < db.Jumlah && terurut[idx].TanggalSelesai == tanggal {
			hasil[jumlahHasil] = terurut[idx]
			jumlahHasil++
			idx++
		}
	}

	return hasil, jumlahHasil
}

// Spesifikasi: Mencari posisi tugas berdasarkan ID secara sequential
func cariIndexByID(db ArrayTugas, id int) int {
	indeks := -1
	ditemukan := false
	i := 0
	for i < db.Jumlah && !ditemukan {
		if db.Data[i].ID == id {
			indeks = i
			ditemukan = true
		}
		i++
	}
	return indeks
}

// Spesifikasi: Mengurutkan tugas berdasarkan prioritas dengan Selection Sort
// ascending = true → prioritas rendah ke tinggi jika ascending = false → prioritas tinggi ke rendah
func selectionSortPrioritas(arr [MAXDATA]Milestone, n int, ascending bool) [MAXDATA]Milestone {
	salinan := arr
	i := 0
	for i < n-1 {
		idxEkstrem := i
		j := i + 1
		for j < n {
			pilih := false
			if ascending {
				pilih = salinan[j].Prioritas < salinan[idxEkstrem].Prioritas
			} else {
				pilih = salinan[j].Prioritas > salinan[idxEkstrem].Prioritas
			}
			if pilih {
				idxEkstrem = j
			}
			j++
		}
		salinan[i], salinan[idxEkstrem] = salinan[idxEkstrem], salinan[i]
		i++
	}
	return salinan
}

// Spesifikasi: Mengurutkan tugas berdasarkan skor mood dengan Insertion Sort
// ascending = true → mood rendah ke tinggi jika  ascending = false → mood tinggi ke rendah
func insertionSortMood(arr [MAXDATA]Milestone, n int, ascending bool) [MAXDATA]Milestone {
	salinan := arr
	i := 1
	for i < n {
		kunci := salinan[i]
		j := i - 1
		geser := true
		for j >= 0 && geser {
			harus := false
			if ascending {
				harus = salinan[j].SkorMood > kunci.SkorMood
			} else {
				harus = salinan[j].SkorMood < kunci.SkorMood
			}
			if harus {
				salinan[j+1] = salinan[j]
				j--
			} else {
				geser = false
			}
		}
		salinan[j+1] = kunci
		i++
	}
	return salinan
}

// Spesifikasi: Menambahkan satu tugas baru ke dalam array
func tambahTugas(db *ArrayTugas, reader *bufio.Reader) {
	judulBesar("TAMBAH TUGAS BARU")

	if db.Jumlah >= MAXDATA {
		pesanPeringatan(fmt.Sprintf("Array Sudah Penuh! Maksimal %d Tugas.", MAXDATA))
		return
	}

	m := Milestone{}
	m.ID = db.Counter
	db.Counter++

	fmt.Println("Yuk Isi Info Tugasnya!")
	fmt.Println()
	m.NamaTugas = inputStr("  Nama Tugasnya Apa?        : ", reader)
	m.Deskripsi = inputStr("  Ceritain Dikit Tugasnya   : ", reader)
	m.Prioritas = inputAngkaRentang("  Seberapa Penting? (1=Santai / 2=Lumayan / 3=Penting Banget): ", 1, 3, reader)
	m.SkorStres = inputAngkaRentang("  Seberapa Stres Ngerjainnya? (1=Tenang - 10=Panik): ", 1, 10, reader)
	m.SkorMood = inputAngkaRentang("  Mood Kamu Waktu Ngerjain? (1=Sedih - 10=Happy banget): ", 1, 10, reader)
	m.CatatanRasa = inputStr("  Mau Curhat Soal Tugas Ini?   : ", reader)
	m.TanggalSelesai = inputStr("  Kapan Selesainya? (contoh: 20-08-2007) : ", reader)

	jwb := strings.ToLower(inputStr("  Udah selesai belum? (y = sudah / n = belum): ", reader))
	m.Selesai = jwb == "y"

	db.Data[db.Jumlah] = m
	db.Jumlah++

	pesanOke(fmt.Sprintf("Yeay! Tugas '%s' Berhasil Disimpan (ID: %d).", m.NamaTugas, m.ID))
}

// Spesifikasi: Menampilkan semua tugas yang tersimpan di array
func lihatSemuaTugas(db ArrayTugas) {
	judulBesar("SEMUA TUGAS KAMU")
	if db.Jumlah == 0 {
		pesanPeringatan("Belum ada tugas nih. Yuk tambah dulu!")
		return
	}
	fmt.Printf("  Total tugas: %d\n", db.Jumlah)
	tampilkanArrayTugas(db.Data, db.Jumlah)
	fmt.Println()
}

// Spesifikasi: Mengubah data tugas yang dicari berdasarkan ID (Sequential Search)
func ubahTugas(db *ArrayTugas, reader *bufio.Reader) {
	judulBesar("UBAH DATA TUGAS")
	if db.Jumlah == 0 {
		tidakKetemu()
		return
	}

	id := inputAngka("  Masukkan ID tugas yang mau diubah: ", reader)
	idx := cariIndexByID(*db, id)

	if idx == -1 {
		tidakKetemu()
		return
	}

	m := db.Data[idx]
	fmt.Printf("  Ketemu! Ini tugasnya: %s%s%s\n", warnaBold, m.NamaTugas, warnaNormal)
	fmt.Println("  Tekan Enter aja kalau nggak mau ubah bagian itu.")
	fmt.Println()

	v := inputStr("  Nama baru (Enter = skip): ", reader)
	if v != "" {
		db.Data[idx].NamaTugas = v
	}

	v = inputStr("  Keterangan baru (Enter = skip): ", reader)
	if v != "" {
		db.Data[idx].Deskripsi = v
	}

	v = inputStr("  Prioritas baru 1/2/3 (Enter = skip): ", reader)
	if v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 1 && n <= 3 {
			db.Data[idx].Prioritas = n
		}
	}

	v = inputStr("  Skor stres baru 1-10 (Enter = skip): ", reader)
	if v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 1 && n <= 10 {
			db.Data[idx].SkorStres = n
		}
	}

	v = inputStr("  Skor mood baru 1-10 (Enter = skip): ", reader)
	if v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 1 && n <= 10 {
			db.Data[idx].SkorMood = n
		}
	}

	v = inputStr("  Curhatan Baru : ", reader)
	if v != "" {
		db.Data[idx].CatatanRasa = v
	}

	v = inputStr("  Tanggal Selesai Baru DD-MM-YYYY : ", reader)
	if v != "" {
		db.Data[idx].TanggalSelesai = v
	}

	jwb := strings.ToLower(inputStr("  Ubah Status Selesai? (y=sudah / n=belum / Enter=skip): ", reader))
	if jwb == "y" {
		db.Data[idx].Selesai = true
	} else if jwb == "n" {
		db.Data[idx].Selesai = false
	}

	pesanOke("Data Tugasnya Berhasil Diperbarui!")
}

// Spesifikasi: Menghapus tugas dari array berdasarkan ID (Sequential Search) ##elemen digeser ke kiri untuk mengisi kekosongan
func hapusTugas(db *ArrayTugas, reader *bufio.Reader) {
	judulBesar("HAPUS TUGAS")
	if db.Jumlah == 0 {
		tidakKetemu()
		return
	}

	id := inputAngka("  Masukkan ID Tugas Yang Mau Dihapus: ", reader)
	idx := cariIndexByID(*db, id)

	if idx == -1 {
		tidakKetemu()
		return
	}

	nama := db.Data[idx].NamaTugas
	konfirmasi := strings.ToLower(inputStr(
		fmt.Sprintf("  Yakin Mau Hapus Tugas '%s'? (y = yakin / n = batal): ", nama),
		reader,
	))

	if konfirmasi == "y" {
		i := idx
		for i < db.Jumlah-1 {
			db.Data[i] = db.Data[i+1]
			i++
		}
		db.Data[db.Jumlah-1] = Milestone{}
		db.Jumlah--
		pesanOke("Oke, Tugasnya Udah Dihapus!")
	} else {
		pesanInfo("Oke, Nggak Jadi Hapus Ya.")
	}
}

// Spesifikasi: Menu pencarian dengan pilihan Sequential (nama) atau Binary (tanggal)
func menuCari(db ArrayTugas, reader *bufio.Reader) {
	judulBesar("CARI TUGAS")
	fmt.Println("  Mau Cari Pakai Para Apa?")
	fmt.Println()
	fmt.Println("  1. Cari Berdasarkan Nama Tugas  ")
	fmt.Println("  2. Cari Berdasarkan Tanggal    ")
	fmt.Println()
	pilihan := inputAngkaRentang("  Pilih (1 atau 2): ", 1, 2, reader)

	if pilihan == 1 {
		katakunci := inputStr("  Ketik Kata Kunci Nama Tugas: ", reader)
		hasil, jumlah := cariSequential(db, katakunci)
		judulKecil(fmt.Sprintf("Hasil Pencarian  %d Tugas Ketemu", jumlah))
		if jumlah == 0 {
			tidakKetemu()
		} else {
			tampilkanArrayTugas(hasil, jumlah)
		}
	} else {
		tanggal := inputStr("  Ketik Tanggalnya (contoh: 20-08-2007): ", reader)
		hasil, jumlah := cariBinary(db, tanggal)
		judulKecil(fmt.Sprintf("Hasil Pencarian — %d Tugas Ketemu", jumlah))
		if jumlah == 0 {
			tidakKetemu()
		} else {
			tampilkanArrayTugas(hasil, jumlah)
		}
	}
}

// Spesifikasi: Menu pengurutan dengan pilihan kriteria dan arah (asc/desc)
func menuUrut(db ArrayTugas, reader *bufio.Reader) {
	judulBesar("URUTKAN TUGAS")

	if db.Jumlah == 0 {
		pesanPeringatan("Belum Ada Tugas Nih.")
		return
	}

	fmt.Println("  Mau Diurutkan Berdasarkan Apa?")
	fmt.Println()
	fmt.Println("  1. Prioritas ")
	fmt.Println("  2. Mood/Perasaan ")
	fmt.Println()
	pilihan := inputAngkaRentang(" Pilih (1 atau 2) : ", 1, 2, reader)

	fmt.Println()
	fmt.Println(" Urutan Tampilan : ")
	fmt.Println(" 1. Santai  →  Penting (1-10)")
	fmt.Println(" 2. Penting → Santai (10-1)")
	fmt.Println()
	arahPilih := inputAngkaRentang("  Pilih arah (1 atau 2): ", 1, 2, reader)
	ascending := arahPilih == 1

	if pilihan == 1 {
		arahLabel := "(10 - 1)"
		if ascending {
			arahLabel = "(1 - 10)"
		}
		judulKecil(fmt.Sprintf("Urutan Prioritas %s", arahLabel))
		hasil := selectionSortPrioritas(db.Data, db.Jumlah, ascending)
		i := 0
		for i < db.Jumlah {
			fmt.Printf("  %s#%d%s ", warnaBold, i+1, warnaNormal)
			tampilkanTugas(hasil[i])
			i++
		}
	} else {
		arahLabel := "(10 - 1)"
		if ascending {
			arahLabel = "(1 - 10)"
		}
		judulKecil(fmt.Sprintf("Urutan Mood %s", arahLabel))
		hasil := insertionSortMood(db.Data, db.Jumlah, ascending)
		i := 0
		for i < db.Jumlah {
			fmt.Printf("  %s#%d%s ", warnaBold, i+1, warnaNormal)
			tampilkanTugas(hasil[i])
			i++
		}
	}
}

// Spesifikasi: Menghitung dan menampilkan statistik tugas minggu ini ##meliputi rata-rata stres dan persentase penyelesaian
func statistik(db ArrayTugas) {
	judulBesar("LAPORAN MINGGUAN ")

	if db.Jumlah == 0 {
		pesanPeringatan("Belum Ada Tugas Nih. Tambah Dulu Yuk!")
		return
	}

	sekarang := time.Now()
	hariIni := int(sekarang.Weekday())
	if hariIni == 0 {
		hariIni = 7
	}
	senin := sekarang.AddDate(0, 0, -(hariIni - 1))
	minggu := senin.AddDate(0, 0, 6)

	format := "20-08-2007"
	tglSenin := senin.Format(format)
	tglMinggu := minggu.Format(format)

	totalStres := 0
	jumlahMingguIni := 0
	jumlahSelesai := 0
	totalSemua := db.Jumlah

	var tugasMingguIni [MAXDATA]Milestone
	jumlahMinggu := 0

	i := 0
	for i < db.Jumlah {
		if db.Data[i].Selesai {
			jumlahSelesai++
		}
		if db.Data[i].TanggalSelesai >= tglSenin && db.Data[i].TanggalSelesai <= tglMinggu {
			totalStres += db.Data[i].SkorStres
			jumlahMingguIni++
			tugasMingguIni[jumlahMinggu] = db.Data[i]
			jumlahMinggu++
		}
		i++
	}

	rataStres := 0.0
	if jumlahMingguIni > 0 {
		rataStres = float64(totalStres) / float64(jumlahMingguIni)
	}

	persen := 0.0
	if totalSemua > 0 {
		persen = float64(jumlahSelesai) / float64(totalSemua) * 100
	}

	fmt.Printf("\n  %sMinggu Ini%s : %s sampai %s\n", warnaBold, warnaNormal, tglSenin, tglMinggu)

	judulKecil("Rata-Rata Tingkat Stres Minggu Ini")
	if jumlahMingguIni == 0 {
		pesanPeringatan("Kamu Nggak Punya Tugas Yang Selesai Minggu Ini.")
	} else {
		nilaiBar := int(math.Round(rataStres))
		fmt.Printf("  Tugas Yang Selesai Minggu Ini : %d tugas\n", jumlahMingguIni)
		fmt.Printf("  Rata-Rata Stresnya            : %s\n", labelStres(nilaiBar))
		fmt.Printf("  Visualisasi Stres             : %s %.2f/10\n", barMood(nilaiBar), rataStres)
		fmt.Printf("\n  %sDetail Tugas Minggu Ini   :%s\n", warnaAbu, warnaNormal)
		j := 0
		for j < jumlahMinggu {
			fmt.Printf("   • %-32s | Stres: %s\n", tugasMingguIni[j].NamaTugas, labelStres(tugasMingguIni[j].SkorStres))
			j++
		}
	}

	judulKecil("Seberapa Banyak Target Yang Berhasil?")
	barProgress := int(persen / 10)
	fmt.Printf("  Total Semua Tugas : %d tugas\n", totalSemua)
	fmt.Printf("  Sudah Selesai     : %d tugas (%s%.0f%%%s)\n", jumlahSelesai, warnaHijau, persen, warnaNormal)
	fmt.Printf("  Belum Selesai     : %d tugas\n", totalSemua-jumlahSelesai)
	fmt.Printf("  Progress Kamu     : %s %.0f%%\n", barMood(barProgress), persen)

	judulKecil("Saran Buat Kamu")
	if rataStres >= 8 {
		fmt.Printf("  %s  Waduh, Stresnya Tinggi Banget! Istirahat Dulu Ya Sayang, Jangan Dipaksain Sayang🥰😘🥰.%s\n", warnaMerah, warnaNormal)
	} else if rataStres >= 5 {
		fmt.Printf("  %s Uemm Ciee Kamu Udah Stres Nyak%s\n", warnaKuning, warnaNormal)
	} else {
		fmt.Printf("  %s  Unch Keren Banget! Stresnya Terkontrol. Terus Jaga Kondisi Seperti Ini Nyak!%s\n", warnaHijau, warnaNormal)
	}

	if persen == 100 {
		fmt.Printf("  %s  Yeay, Semua Tugas Selesai! Kamu Luar Biasa, Keren Banget Beneran Deh Ngga Boong😁✌️!%s\n", warnaHijau, warnaNormal)
	} else if persen >= 70 {
		fmt.Printf("  %s  Produktivitasnya Udah Oke Banget Enih! Dikit Lagi, Semangattt!%s\n", warnaHijau, warnaNormal)
	} else if persen >= 40 {
		fmt.Printf("  %s  Lumayan Progresnya. Masih Bisa Lebih Baik Lagi, Kamu Pasti Bisa!%s\n", warnaKuning, warnaNormal)
	} else {
		fmt.Printf("  %s  Masih Banyak Yang Belum Selesai Enih. Mulai Dari Yang Mudah Dulu Nyak!%s\n", warnaMerah, warnaNormal)
	}
}

// Spesifikasi: Mengisi array dengan data contoh awal saat program dijalankan
func isiDataContoh(db *ArrayTugas) {
	hari := time.Now()
	mingguIni := hari.Format("20-08-2007")
	mingguLalu := hari.AddDate(0, 0, -7).Format("20-08-2007")

	contoh := [5]Milestone{
		{ID: 1, NamaTugas: "Tugas Laporan Praktikum", Deskripsi: "Laporan Praktikum Algoritma Pemrograman 2", Prioritas: 3, SkorStres: 6, SkorMood: 7, CatatanRasa: "Menantang Deadline", TanggalSelesai: mingguIni, Selesai: true},
		{ID: 2, NamaTugas: "Kuis Mingguan", Deskripsi: "Pengerjaan Soal-Soal", Prioritas: 2, SkorStres: 4, SkorMood: 8, CatatanRasa: "Menjemput Deadline", TanggalSelesai: mingguIni, Selesai: true},
		{ID: 3, NamaTugas: "Interview Kepanitiaan", Deskripsi: "Interview TODAYS dan WPI", Prioritas: 3, SkorStres: 9, SkorMood: 5, CatatanRasa: "Deg-Degan Banget, Mencoba Keberuntungan", TanggalSelesai: mingguIni, Selesai: false},
		{ID: 4, NamaTugas: "Belajar Ujian", Deskripsi: "Review Materi", Prioritas: 2, SkorStres: 5, SkorMood: 6, CatatanRasa: "3M, Mantap, Mumet, Mbuh", TanggalSelesai: mingguLalu, Selesai: true},
		{ID: 5, NamaTugas: "Olahraga Rutin", Deskripsi: "Jogging", Prioritas: 1, SkorStres: 2, SkorMood: 9, CatatanRasa: "Jujur Malas:)", TanggalSelesai: mingguLalu, Selesai: true},
	}

	i := 0
	for i < 5 {
		db.Data[i] = contoh[i]
		i++
	}
	db.Jumlah = 5
	db.Counter = 6
}

// Spesifikasi: Mencetak logo MindStone saat program dijalankan
func tampilanAwal() {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════╗%s\n", warnaCyan, warnaNormal)
	fmt.Printf("%s║%s   MindStone Catatan Tugas & Kesehatan Mental         %s║%s\n", warnaCyan, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s   Tugas Besar Algoritma Pemrograman 2 — TUP.         %s║%s\n", warnaCyan, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s╚══════════════════════════════════════════════════════╝%s\n", warnaCyan, warnaNormal)
}

// SUBPROGRAM: MENU UTAMA
// Spesifikasi: Loop utama program, menampilkan menu dan mendelegasikan aksi ke subprogram terkait berdasarkan pilihan pengguna
func menuUtama(db *ArrayTugas, reader *bufio.Reader) {
	tampilanAwal()
	selesai := false
	for !selesai {
		fmt.Printf("\n%s┌─ MAU NGAPAIN? ───────────────────────────┐%s\n", warnaBold, warnaNormal)
		fmt.Printf("%s│%s  1. Lihat semua tugas                    %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  2. Tambah tugas baru                    %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  3. Ubah data tugas                      %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  4. Hapus tugas                          %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  5. Cari tugas                           %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  6. Urutkan tugas                        %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  7. Lihat laporan & statistik mingguan   %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  0. Keluar dari aplikasi                 %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s└──────────────────────────────────────────┘%s\n", warnaBold, warnaNormal)

		pilihan := inputAngka("  Pilih menu : ", reader)

		if pilihan == 0 {
			selesai = true
		} else if pilihan == 1 {
			lihatSemuaTugas(*db)
		} else if pilihan == 2 {
			tambahTugas(db, reader)
		} else if pilihan == 3 {
			ubahTugas(db, reader)
		} else if pilihan == 4 {
			hapusTugas(db, reader)
		} else if pilihan == 5 {
			menuCari(*db, reader)
		} else if pilihan == 6 {
			menuUrut(*db, reader)
		} else if pilihan == 7 {
			statistik(*db)
		} else {
			pesanPeringatan("Pilihannya Tidak Ada Dalam Data")
		}
	}
	fmt.Printf("\n%s Terimakasih Sudah Memakai MindStone Cemangat Nyak Kamuh Cayang %s\n\n", warnaCyan, warnaNormal)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	isiDataContoh(&tugasDB)
	menuUtama(&tugasDB, reader)
}

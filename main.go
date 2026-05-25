package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ============================================================
//  STRUCT — Cetak biru data setiap tugas
// ============================================================

type Milestone struct {
	ID             int
	NamaTugas      string
	Deskripsi      string
	Prioritas      int    // 1=Santai, 2=Sedang, 3=Penting Banget
	SkorStres      int    // 1-10 (1=tenang, 10=panik)
	SkorMood       int    // 1-10 (1=sedih, 10=bahagia)
	CatatanRasa    string // Ceritain perasaanmu waktu ngerjain tugas ini
	TanggalSelesai string // format: TAHUN-BULAN-TANGGAL, contoh: 2026-05-18
	Selesai        bool
}

// ============================================================
//  DATA GLOBAL — Tempat nyimpen semua tugas
// ============================================================

var (
	data      []Milestone
	idCounter int = 1
	reader        = bufio.NewReader(os.Stdin)
)

// ============================================================
//  FUNGSI BANTU — Buat baca input dari keyboard
// ============================================================

func inputStr(pertanyaan string) string {
	fmt.Print(pertanyaan)
	teks, _ := reader.ReadString('\n')
	return strings.TrimSpace(teks)
}

func inputAngka(pertanyaan string) int {
	for {
		s := inputStr(pertanyaan)
		n, err := strconv.Atoi(s)
		if err == nil {
			return n
		}
		fmt.Println("  ⚠  Hei, itu bukan angka! Coba lagi ya.")
	}
}

func inputAngkaRentang(pertanyaan string, min, max int) int {
	for {
		n := inputAngka(pertanyaan)
		if n >= min && n <= max {
			return n
		}
		fmt.Printf("  ⚠  Angkanya harus antara %d sampai %d dong!\n", min, max)
	}
}

// ============================================================
//  WARNA — Biar tampilannya cantik di terminal
// ============================================================

const (
	warnaNormal = "\033[0m"
	warnaBold   = "\033[1m"
	warnaCyan   = "\033[36m"
	warnaHijau  = "\033[32m"
	warnaKuning = "\033[33m"
	warnaMerah  = "\033[31m"
	warnaAbu    = "\033[90m"
)

// ============================================================
//  FUNGSI TAMPILAN — Biar enak dibaca
// ============================================================

func judulBesar(teks string) {
	garis := strings.Repeat("─", 60)
	fmt.Printf("\n%s%s%s\n", warnaCyan, garis, warnaNormal)
	fmt.Printf("%s  ◆ %s%s\n", warnaBold, teks, warnaNormal)
	fmt.Printf("%s%s%s\n", warnaCyan, garis, warnaNormal)
}

func judulKecil(teks string) {
	fmt.Printf("\n%s── %s%s\n", warnaAbu, teks, warnaNormal)
}

func pesanOke(pesan string)        { fmt.Printf("%s  ✔  %s%s\n", warnaHijau, pesan, warnaNormal) }
func pesanPeringatan(pesan string) { fmt.Printf("%s  ⚠  %s%s\n", warnaKuning, pesan, warnaNormal) }
func pesanInfo(pesan string)       { fmt.Printf("%s  ℹ  %s%s\n", "\033[34m", pesan, warnaNormal) }
func tidakKetemu()                 { pesanPeringatan("Hmm, datanya tidak ketemu nih.") }

func labelPrioritas(p int) string {
	switch p {
	case 1:
		return warnaAbu + "Santai aja" + warnaNormal
	case 2:
		return warnaKuning + "Lumayan penting" + warnaNormal
	case 3:
		return warnaMerah + "Penting banget!" + warnaNormal
	}
	return "?"
}

func barMood(skor int) string {
	isi := skor
	kosong := 10 - skor
	return warnaHijau + strings.Repeat("█", isi) + warnaAbu + strings.Repeat("░", kosong) + warnaNormal
}

func labelStres(s int) string {
	if s <= 3 {
		return warnaHijau + strconv.Itoa(s) + "/10 (Tenang)" + warnaNormal
	} else if s <= 6 {
		return warnaKuning + strconv.Itoa(s) + "/10 (Lumayan stres)" + warnaNormal
	}
	return warnaMerah + strconv.Itoa(s) + "/10 (Stres tinggi!)" + warnaNormal
}

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

// ============================================================
//  SPESIFIKASI A — CRUD (Tambah, Lihat, Ubah, Hapus)
// ============================================================

func tambahTugas() {
	judulBesar("TAMBAH TUGAS BARU")
	m := Milestone{}
	m.ID = idCounter
	idCounter++

	fmt.Println("  Yuk isi info tugasnya! (Tekan Enter setelah nulis)")
	fmt.Println()
	m.NamaTugas = inputStr("  Nama tugasnya apa?           : ")
	m.Deskripsi = inputStr("  Ceritain dikit tugasnya      : ")
	m.Prioritas = inputAngkaRentang("  Seberapa penting? (1=Santai / 2=Lumayan / 3=Penting Banget): ", 1, 3)
	m.SkorStres = inputAngkaRentang("  Seberapa stres ngerjainnya? (1=Tenang ... 10=Panik): ", 1, 10)
	m.SkorMood = inputAngkaRentang("  Mood kamu waktu ngerjain? (1=Sedih ... 10=Happy banget): ", 1, 10)
	m.CatatanRasa = inputStr("  Mau curhat soal tugas ini?   : ")
	m.TanggalSelesai = inputStr("  Kapan selesainya? (YYYY-MM-DD, contoh: 2026-05-18): ")

	jwb := strings.ToLower(inputStr("  Udah selesai belum? (y = sudah / n = belum): "))
	m.Selesai = jwb == "y"

	data = append(data, m)
	pesanOke(fmt.Sprintf("Yeay! Tugas '%s' berhasil disimpan (ID: %d).", m.NamaTugas, m.ID))
}

func ubahTugas() {
	judulBesar("UBAH DATA TUGAS")
	if len(data) == 0 {
		tidakKetemu()
		return
	}
	id := inputAngka("  Masukkan ID tugas yang mau diubah: ")
	for i, m := range data {
		if m.ID == id {
			fmt.Printf("  Ketemu! Ini tugasnya: %s%s%s\n", warnaBold, m.NamaTugas, warnaNormal)
			fmt.Println("  Tekan Enter aja kalau nggak mau ubah bagian itu.")
			fmt.Println()

			if v := inputStr("  Nama baru (Enter = skip): "); v != "" {
				data[i].NamaTugas = v
			}
			if v := inputStr("  Keterangan baru (Enter = skip): "); v != "" {
				data[i].Deskripsi = v
			}
			if v := inputStr("  Prioritas baru 1/2/3 (Enter = skip): "); v != "" {
				if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 3 {
					data[i].Prioritas = n
				}
			}
			if v := inputStr("  Skor stres baru 1-10 (Enter = skip): "); v != "" {
				if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 10 {
					data[i].SkorStres = n
				}
			}
			if v := inputStr("  Skor mood baru 1-10 (Enter = skip): "); v != "" {
				if n, err := strconv.Atoi(v); err == nil && n >= 1 && n <= 10 {
					data[i].SkorMood = n
				}
			}
			if v := inputStr("  Curhatan baru (Enter = skip): "); v != "" {
				data[i].CatatanRasa = v
			}
			if v := inputStr("  Tanggal selesai baru YYYY-MM-DD (Enter = skip): "); v != "" {
				data[i].TanggalSelesai = v
			}
			jwb := strings.ToLower(inputStr("  Ubah status selesai? (y=sudah / n=belum / Enter=skip): "))
			if jwb == "y" {
				data[i].Selesai = true
			} else if jwb == "n" {
				data[i].Selesai = false
			}

			pesanOke("Data tugasnya berhasil diperbarui!")
			return
		}
	}
	tidakKetemu()
}

func hapusTugas() {
	judulBesar("HAPUS TUGAS")
	if len(data) == 0 {
		tidakKetemu()
		return
	}
	id := inputAngka("  Masukkan ID tugas yang mau dihapus: ")
	for i, m := range data {
		if m.ID == id {
			konfirmasi := strings.ToLower(inputStr(fmt.Sprintf("  Yakin mau hapus tugas '%s'? Nggak bisa balik lho! (y = yakin / n = batal): ", m.NamaTugas)))
			if konfirmasi == "y" {
				data = append(data[:i], data[i+1:]...)
				pesanOke("Oke, tugasnya udah dihapus!")
			} else {
				pesanInfo("Oke, nggak jadi hapus ya.")
			}
			return
		}
	}
	tidakKetemu()
}

func lihatSemuaTugas() {
	judulBesar("SEMUA TUGAS KAMU")
	if len(data) == 0 {
		pesanPeringatan("Belum ada tugas nih. Yuk tambah dulu!")
		return
	}
	fmt.Printf("  Total tugas: %d\n", len(data))
	for _, m := range data {
		tampilkanTugas(m)
	}
	fmt.Println()
}

// ============================================================
//  SPESIFIKASI C — PENCARIAN
// ============================================================

// Sequential Search: cek satu-satu dari awal sampai akhir
func cariNama(katakunci string) []Milestone {
	var hasil []Milestone
	katakunci = strings.ToLower(katakunci)
	for _, m := range data {
		if strings.Contains(strings.ToLower(m.NamaTugas), katakunci) {
			hasil = append(hasil, m)
		}
	}
	return hasil
}

// Binary Search: cari dengan cara belah tengah (lebih cepat, tapi data harus diurutkan dulu)
func cariTanggal(tanggal string) []Milestone {
	salinan := make([]Milestone, len(data))
	copy(salinan, data)

	sort.Slice(salinan, func(i, j int) bool {
		return salinan[i].TanggalSelesai < salinan[j].TanggalSelesai
	})

	var hasil []Milestone
	kiri, kanan := 0, len(salinan)-1

	for kiri <= kanan {
		tengah := (kiri + kanan) / 2
		if salinan[tengah].TanggalSelesai == tanggal {
			for k := tengah; k >= 0 && salinan[k].TanggalSelesai == tanggal; k-- {
				hasil = append([]Milestone{salinan[k]}, hasil...)
			}
			for k := tengah + 1; k < len(salinan) && salinan[k].TanggalSelesai == tanggal; k++ {
				hasil = append(hasil, salinan[k])
			}
			break
		} else if salinan[tengah].TanggalSelesai < tanggal {
			kiri = tengah + 1
		} else {
			kanan = tengah - 1
		}
	}
	return hasil
}

func menuCari() {
	judulBesar("CARI TUGAS")
	fmt.Println("  Mau cari pakai cara apa?")
	fmt.Println()
	fmt.Println("  1. Cari berdasarkan Nama Tugas  → pakai Sequential Search (cek satu-satu)")
	fmt.Println("  2. Cari berdasarkan Tanggal     → pakai Binary Search (cari dari tengah)")
	fmt.Println()
	pilihan := inputAngkaRentang("  Pilih (1 atau 2): ", 1, 2)

	if pilihan == 1 {
		katakunci := inputStr("  Ketik kata kunci nama tugas: ")
		hasil := cariNama(katakunci)
		judulKecil(fmt.Sprintf("Hasil pencarian — %d tugas ketemu", len(hasil)))
		if len(hasil) == 0 {
			tidakKetemu()
		} else {
			for _, m := range hasil {
				tampilkanTugas(m)
			}
		}
	} else {
		tanggal := inputStr("  Ketik tanggalnya (YYYY-MM-DD, contoh: 2026-05-18): ")
		hasil := cariTanggal(tanggal)
		judulKecil(fmt.Sprintf("Hasil pencarian — %d tugas ketemu", len(hasil)))
		if len(hasil) == 0 {
			tidakKetemu()
		} else {
			for _, m := range hasil {
				tampilkanTugas(m)
			}
		}
	}
}

// ============================================================
//  SPESIFIKASI D — PENGURUTAN
// ============================================================

// Selection Sort: tiap putaran, cari yang prioritasnya paling tinggi lalu taruh di depan
func urutPrioritas() []Milestone {
	salinan := make([]Milestone, len(data))
	copy(salinan, data)
	n := len(salinan)

	for i := 0; i < n-1; i++ {
		idxTertinggi := i
		for j := i + 1; j < n; j++ {
			if salinan[j].Prioritas > salinan[idxTertinggi].Prioritas {
				idxTertinggi = j
			}
		}
		salinan[i], salinan[idxTertinggi] = salinan[idxTertinggi], salinan[i]
	}
	return salinan
}

// Insertion Sort: ambil satu-satu, sisipkan di posisi yang tepat
func urutMood() []Milestone {
	salinan := make([]Milestone, len(data))
	copy(salinan, data)
	n := len(salinan)

	for i := 1; i < n; i++ {
		tugas := salinan[i]
		j := i - 1
		for j >= 0 && salinan[j].SkorMood < tugas.SkorMood {
			salinan[j+1] = salinan[j]
			j--
		}
		salinan[j+1] = tugas
	}
	return salinan
}

func menuUrut() {
	judulBesar("URUTKAN TUGAS")
	fmt.Println("  Mau diurutkan berdasarkan apa?")
	fmt.Println()
	fmt.Println("  1. Prioritas      → pakai Selection Sort (dari yang paling penting dulu)")
	fmt.Println("  2. Mood/Perasaan  → pakai Insertion Sort (dari yang paling happy dulu)")
	fmt.Println()
	pilihan := inputAngkaRentang("  Pilih (1 atau 2): ", 1, 2)

	if pilihan == 1 {
		judulKecil("Hasil urutan berdasarkan Prioritas (dari yang paling penting)")
		hasil := urutPrioritas()
		for urutan, m := range hasil {
			fmt.Printf("  %s#%d%s ", warnaBold, urutan+1, warnaNormal)
			tampilkanTugas(m)
		}
	} else {
		judulKecil("Hasil urutan berdasarkan Mood (dari yang paling happy)")
		hasil := urutMood()
		for urutan, m := range hasil {
			fmt.Printf("  %s#%d%s ", warnaBold, urutan+1, warnaNormal)
			tampilkanTugas(m)
		}
	}
}

// ============================================================
//  SPESIFIKASI E — STATISTIK MINGGUAN
// ============================================================

func statistik() {
	judulBesar("LAPORAN MINGGUAN KAMU")

	if len(data) == 0 {
		pesanPeringatan("Belum ada tugas nih. Tambah dulu yuk!")
		return
	}

	sekarang := time.Now()
	hariIni := int(sekarang.Weekday())
	if hariIni == 0 {
		hariIni = 7
	}
	senin := sekarang.AddDate(0, 0, -(hariIni - 1))
	minggu := senin.AddDate(0, 0, 6)

	format := "2006-01-02"
	tglSenin := senin.Format(format)
	tglMinggu := minggu.Format(format)

	var totalStres, jumlahMingguIni int
	var jumlahSelesai, totalSemua int
	var tugasMingguIni []Milestone

	for _, m := range data {
		totalSemua++
		if m.Selesai {
			jumlahSelesai++
		}
		if m.TanggalSelesai >= tglSenin && m.TanggalSelesai <= tglMinggu {
			totalStres += m.SkorStres
			jumlahMingguIni++
			tugasMingguIni = append(tugasMingguIni, m)
		}
	}

	var rataStres float64
	if jumlahMingguIni > 0 {
		rataStres = float64(totalStres) / float64(jumlahMingguIni)
	}

	persen := 0.0
	if totalSemua > 0 {
		persen = float64(jumlahSelesai) / float64(totalSemua) * 100
	}

	fmt.Printf("\n  %sMinggu ini%s : %s sampai %s\n", warnaBold, warnaNormal, tglSenin, tglMinggu)

	judulKecil("Rata-rata Tingkat Stres Minggu Ini")
	if jumlahMingguIni == 0 {
		pesanPeringatan("Kamu nggak punya tugas yang selesai minggu ini.")
	} else {
		nilaiBar := int(math.Round(rataStres))
		fmt.Printf("  Tugas yang kelar minggu ini : %d tugas\n", jumlahMingguIni)
		fmt.Printf("  Rata-rata stresnya          : %s\n", labelStres(nilaiBar))
		fmt.Printf("  Visualisasi stres           : %s %.2f/10\n", barMood(nilaiBar), rataStres)

		fmt.Printf("\n  %sDetail tugas minggu ini:%s\n", warnaAbu, warnaNormal)
		for _, m := range tugasMingguIni {
			fmt.Printf("   • %-32s | Stres: %s\n", m.NamaTugas, labelStres(m.SkorStres))
		}
	}

	judulKecil("Seberapa Banyak Target yang Berhasil?")
	barProgress := int(persen / 10)
	fmt.Printf("  Total semua tugas  : %d tugas\n", totalSemua)
	fmt.Printf("  Yang sudah selesai : %d tugas (%s%.0f%%%s)\n",
		jumlahSelesai, warnaHijau, persen, warnaNormal)
	fmt.Printf("  Yang belum selesai : %d tugas\n", totalSemua-jumlahSelesai)
	fmt.Printf("  Progress kamu      : %s %.0f%%\n", barMood(barProgress), persen)

	judulKecil("Saran Buat Kamu")
	switch {
	case rataStres >= 8:
		fmt.Printf("  %s⚠  Wah, stresnya tinggi banget! Istirahat dulu ya, jangan dipaksain.%s\n", warnaMerah, warnaNormal)
	case rataStres >= 5:
		fmt.Printf("  %s●  Stresnya lumayan nih. Jangan lupa rehat dan jaga keseimbangan ya!%s\n", warnaKuning, warnaNormal)
	default:
		fmt.Printf("  %s✔  Keren! Stresnya terkontrol. Terus jaga kondisi seperti ini ya!%s\n", warnaHijau, warnaNormal)
	}

	switch {
	case persen == 100:
		fmt.Printf("  %s🏆  Semua tugas beres! Kamu luar biasa, beneran!%s\n", warnaHijau, warnaNormal)
	case persen >= 70:
		fmt.Printf("  %s✔  Produktivitasnya oke banget! Dikit lagi, semangat!%s\n", warnaHijau, warnaNormal)
	case persen >= 40:
		fmt.Printf("  %s●  Lumayan progresnya. Masih bisa lebih baik lagi, kamu pasti bisa!%s\n", warnaKuning, warnaNormal)
	default:
		fmt.Printf("  %s⚠  Masih banyak yang belum kelar nih. Mulai dari yang kecil dulu ya!%s\n", warnaMerah, warnaNormal)
	}
}

// ============================================================
//  DATA CONTOH — Biar langsung ada isinya waktu pertama buka
// ============================================================

func isiDataContoh() {
	hari := time.Now()
	mingguIni := hari.Format("2006-01-02")
	mingguLalu := hari.AddDate(0, 0, -7).Format("2006-01-02")

	data = []Milestone{
		{ID: 1, NamaTugas: "Buat Laporan Mingguan", Deskripsi: "Laporan progres ke manajer", Prioritas: 3, SkorStres: 6, SkorMood: 7, CatatanRasa: "Sedikit tertekan tapi berhasil", TanggalSelesai: mingguIni, Selesai: true},
		{ID: 2, NamaTugas: "Review Kode Teman", Deskripsi: "Cek pull request teman satu tim", Prioritas: 2, SkorStres: 4, SkorMood: 8, CatatanRasa: "Seru, banyak hal baru yang dipelajari", TanggalSelesai: mingguIni, Selesai: true},
		{ID: 3, NamaTugas: "Presentasi ke Klien", Deskripsi: "Demo fitur terbaru ke klien", Prioritas: 3, SkorStres: 9, SkorMood: 5, CatatanRasa: "Deg-degan banget, tapi alhamdulillah lancar", TanggalSelesai: mingguIni, Selesai: false},
		{ID: 4, NamaTugas: "Belajar Algoritma", Deskripsi: "Tugas besar algoritma pemrograman", Prioritas: 2, SkorStres: 5, SkorMood: 6, CatatanRasa: "Susah tapi seru!", TanggalSelesai: mingguLalu, Selesai: true},
		{ID: 5, NamaTugas: "Olahraga Rutin", Deskripsi: "Jogging 30 menit tiap pagi", Prioritas: 1, SkorStres: 2, SkorMood: 9, CatatanRasa: "Bikin seger banget!", TanggalSelesai: mingguLalu, Selesai: true},
	}
	idCounter = 6
}

// ============================================================
//  TAMPILAN AWAL — Logo MindStone
// ============================================================

func tampilanAwal() {
	fmt.Println()
	fmt.Printf("%s╔══════════════════════════════════════════════════════════╗%s\n", warnaCyan, warnaNormal)
	fmt.Printf("%s║%s  %s ██▓▄▄▄█████▓ ███▄    █ ▓█████ %s                      %s║%s\n", warnaCyan, warnaNormal, warnaBold, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s  %s▓██▒▓  ██▒ ▓▒ ██ ▀█   █ ▓█   ▀ %s                      %s║%s\n", warnaCyan, warnaNormal, warnaBold, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s  %s▒██▒▒ ▓██░ ▒░▓██  ▀█ ██▒▒███   %s                      %s║%s\n", warnaCyan, warnaNormal, warnaBold, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s  %s░██░░ ▓██▓ ░ ▓██▒  ▐▌██▒▒▓█  ▄ %s                      %s║%s\n", warnaCyan, warnaNormal, warnaBold, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s  %s░██░  ▒██▒ ░ ▒██░   ▓██░░▒████▒%s                      %s║%s\n", warnaCyan, warnaNormal, warnaBold, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s  %s░▓    ▒ ░░   ░ ▒░   ▒ ▒ ░░ ▒░ ░%s  S T O N E           %s║%s\n", warnaCyan, warnaNormal, warnaBold, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s                                                          %s║%s\n", warnaCyan, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s   Catatan Tugas & Pemantau Kesehatan Mental Kamu         %s║%s\n", warnaCyan, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s║%s   Tugas Besar Algoritma Pemrograman 2 — Telkom Univ.     %s║%s\n", warnaCyan, warnaNormal, warnaCyan, warnaNormal)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════╝%s\n", warnaCyan, warnaNormal)
}

// ============================================================
//  MENU UTAMA
// ============================================================

func menuUtama() {
	tampilanAwal()
	for {
		fmt.Printf("\n%s┌─ MAU NGAPAIN? ────────────────────────────────────────┐%s\n", warnaBold, warnaNormal)
		fmt.Printf("%s│%s  1. Lihat semua tugas                                  %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  2. Tambah tugas baru                                  %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  3. Ubah data tugas                                    %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  4. Hapus tugas                                        %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  5. Cari tugas (Sequential / Binary Search)            %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  6. Urutkan tugas (Selection / Insertion Sort)         %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  7. Lihat laporan & statistik mingguan                 %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s│%s  0. Keluar dari aplikasi                               %s│%s\n", warnaBold, warnaNormal, warnaBold, warnaNormal)
		fmt.Printf("%s└───────────────────────────────────────────────────────┘%s\n", warnaBold, warnaNormal)

		pilihan := inputAngka("  Pilih menu (ketik angkanya): ")
		switch pilihan {
		case 1:
			lihatSemuaTugas()
		case 2:
			tambahTugas()
		case 3:
			ubahTugas()
		case 4:
			hapusTugas()
		case 5:
			menuCari()
		case 6:
			menuUrut()
		case 7:
			statistik()
		case 0:
			fmt.Printf("\n%s  ◆ Makasih udah pakai MindStone! Jaga kesehatan ya, semangat terus! ◆%s\n\n", warnaCyan, warnaNormal)
			os.Exit(0)
		default:
			pesanPeringatan("Pilihannya nggak ada tuh! Ketik angka 0 sampai 7 ya.")
		}
	}
}

// ============================================================
//  TITIK MULAI PROGRAM
// ============================================================

func main() {
	isiDataContoh()
	menuUtama()
}

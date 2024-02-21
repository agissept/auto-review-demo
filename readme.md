# Sebuah Auto Review
## Kriteria Submission
- Memiliki package.json
- Memiliki berkas file main.js
- Pada file main.js harus terdapat komentar yang berisi id siswa
- Root pada aplikasi main.js harus menampilkan html
- Port harus berada di 5000
- Pada file html harus terdapat element h1 dengan isi id siswa

# Requirement Untuk Menjalankan Program
- Linux
- NodeJS minimal Versi 14
- Golang(Opsional)

# Cara Menjalankan Program Tanpa Golang
- Clone project ini
- Jalankan binary yang ada pada folder `bin`
- Tentukan path submission dan path report, contoh

  `./bin/bin -submission ./submissions/submission-1 -report ./bin`
- Cek `report.json` pada folder bin

# Cara Menjalankan Program Dengan Golang
- Clone project ini
- Jalankan command ` go run *.go` pada root project beserta path submission dan path report, contoh
  
  `go run *.go -report submissions/submission-5 -submission submissions/submission-5`
- Cek `report.json` pada folder `submissions/submission-5`
Instal pustaka Firebase Admin SDK for Go menggunakan perintah go get:

Bash

go get firebase.google.com/go/v4
Mendapatkan Kunci dan Kredensial (JSON File)
Untuk mengautentikasi server Anda ke Firebase, Anda memerlukan kunci privat dalam bentuk file JSON yang dikenal sebagai Service Account Key.

1. Masuk ke Firebase Console
Masuk ke Firebase console dan pilih proyek Anda.

2. Unduh Kunci Privat
Navigasi ke Project settings > Service accounts.

Di bagian bawah halaman, klik tombol Generate new private key. Konfirmasi dengan mengklik Generate Key.

File JSON akan otomatis terunduh. Beri nama file ini, misalnya firebase-adminsdk.json.

⚠️ Penting: File ini adalah kunci rahasia. Jangan pernah mengunggahnya ke repositori publik seperti GitHub. Tambahkan file ini ke .gitignore Anda.

Konfigurasi di Proyek Golang
Setelah mendapatkan file JSON, Anda perlu menyimpannya di proyek dan menggunakannya.

1. Simpan File di Proyek
Letakkan file firebase-adminsdk.json di direktori yang aman di dalam proyek Anda. Direktori root proyek adalah lokasi yang umum.

my-goravel-app/
├── .env
├── .gitignore
├── firebase-adminsdk.json  
├── ...
2. Tambahkan Path ke .env
Tambahkan path ke file ini sebagai variabel lingkungan di file .env.

FIREBASE_ADMIN_SDK_PATH=./firebase-adminsdk.json
Hal ini memungkinkan Anda untuk dengan mudah mengubah path di lingkungan yang berbeda (misalnya, produksi).

3. Inisialisasi SDK di Kode
Gunakan variabel lingkungan ini untuk menginisialisasi SDK di kode Anda, biasanya di dalam sebuah service provider atau langsung di controller autentikasi Anda.

Go

package controllers

import (
    "context"
    "os"
    
    "firebase.google.com/go/v4"
    "google.golang.org/api/option"
    
    "github.com/goravel/framework/facades"
)

// Inisialisasi Firebase App
func initFirebaseApp() (*firebase.App, error) {
    serviceAccountKeyPath := facades.Config().GetString("FIREBASE_ADMIN_SDK_PATH")
    
    // Periksa apakah path file tersedia
    if _, err := os.Stat(serviceAccountKeyPath); os.IsNotExist(err) {
        return nil, err
    }
    
    // Inisialisasi Firebase dengan kredensial dari file JSON
    opt := option.WithCredentialsFile(serviceAccountKeyPath)
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        return nil, err
    }
    return app, nil
}

Dokumentasi ini mencakup
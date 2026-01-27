# Kasir API

REST API sederhana untuk sistem kasir yang dibangun dengan Go. API ini menyediakan operasi CRUD untuk produk dan kategori.

## 🚀 Live Demo

**Railway Deployment:** https://kasir-api-production-1404.up.railway.app
**Zeabur Deployment:** https://kasir-api-zeabur.zeabur.app (akan tersedia setelah deploy)

**Test Endpoints:**
- [GET /categories](https://kasir-api-production-1404.up.railway.app/categories) - Lihat semua kategori
- [GET /api/produk](https://kasir-api-production-1404.up.railway.app/api/produk) - Lihat semua produk
- [GET /health](https://kasir-api-production-1404.up.railway.app/health) - Health check

## 🚀 Fitur

- **Manajemen Produk**: Create, Read, Update, Delete produk
- **Manajemen Kategori**: Create, Read, Update, Delete kategori
- **Health Check**: Endpoint untuk monitoring status API
- **JSON Response**: Semua response dalam format JSON

## 📋 Persyaratan

- Go 1.25.6 atau lebih tinggi
- Git

## 🛠️ Instalasi

1. Clone repository
```bash
git clone https://github.com/ronroom/kasir-api.git
cd kasir-api
```

2. Jalankan aplikasi
```bash
go run .
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Endpoints

### Health Check
- `GET /health` - Cek status API

### Produk
- `GET /api/produk` - Ambil semua produk
- `POST /api/produk` - Tambah produk baru
- `GET /api/produk/{id}` - Ambil produk berdasarkan ID
- `PUT /api/produk/{id}` - Update produk
- `DELETE /api/produk/{id}` - Hapus produk

### Kategori
- `GET /categories` - Ambil semua kategori
- `POST /categories` - Tambah kategori baru
- `GET /categories/{id}` - Ambil kategori berdasarkan ID
- `PUT /categories/{id}` - Update kategori
- `DELETE /categories/{id}` - Hapus kategori

## 📝 Contoh Penggunaan

### Tambah Produk Baru
```bash
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Teh Botol",
    "harga": 4000,
    "stok": 25
  }'
```

### Tambah Kategori Baru
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Elektronik",
    "description": "Kategori untuk produk elektronik"
  }'
```

### Ambil Semua Produk
```bash
curl http://localhost:8080/api/produk
```

### Update Produk
```bash
curl -X PUT http://localhost:8080/api/produk/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Indomie Goreng",
    "harga": 3800,
    "stok": 15
  }'
```

## 📊 Model Data

### Produk
```json
{
  "id": 1,
  "nama": "Indomie",
  "harga": 3500,
  "stok": 10
}
```

### Kategori
```json
{
  "id": 1,
  "name": "Makanan",
  "description": "Kategori untuk produk makanan"
}
```

## 🏗️ Struktur Project

```
kasir-api/
├── main.go              # File utama aplikasi dengan semua handler
├── go.mod              # Go module
└── README.md           # Dokumentasi project
```

## 📊 Data Sample

### Produk
- Indomie (Rp 3.500)
- Vit 1000ml (Rp 3.000) 
- Kecap (Rp 12.000)

### Kategori
- Makanan
- Minuman
- Elektronik
- Pakaian
- Kesehatan
- Olahraga

## 🔧 Development

### Menjalankan dalam mode development
```bash
go run .
```

### Build aplikasi
```bash
go build -o kasir-api
./kasir-api
```

## 📄 Response Format

### Success Response
```json
{
  "id": 1,
  "nama": "Indomie",
  "harga": 3500,
  "stok": 10
}
```

### Error Response
```json
{
  "error": "Produk not found"
}
```

## 🤝 Kontribusi

1. Fork repository
2. Buat branch fitur (`git checkout -b feature/fitur-baru`)
3. Commit perubahan (`git commit -am 'Tambah fitur baru'`)
4. Push ke branch (`git push origin feature/fitur-baru`)
5. Buat Pull Request

## 📝 License

Project ini menggunakan MIT License.

## 👨‍💻 Author

**Ronny Romal**
- GitHub: [@ronroom](https://github.com/ronroom)

---

⭐ Jangan lupa berikan star jika project ini membantu!
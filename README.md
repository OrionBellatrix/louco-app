# Louco Event Backend

Bu proje Clean Architecture / Hexagonal Architecture prensiplerine göre geliştirilmiş bir Go backend uygulamasıdır.

## 🏗️ Mimari

Proje aşağıdaki katmanlardan oluşmaktadır:

- **Domain Layer**: İş kuralları ve entity'ler
- **Repository Layer**: Veri erişim katmanı
- **Service Layer**: İş mantığı (use-cases)
- **Transport Layer**: HTTP handlers (Gin)
- **Infrastructure**: Database, logger, config vb.

## 📁 Proje Yapısı

```
/cmd/app/main.go          # Ana uygulama giriş noktası
/internal/
    domain/               # Entity'ler ve iş kuralları
    repository/           # Repository interface'leri
        postgres/         # PostgreSQL implementasyonları
    service/              # İş mantığı katmanı
    transport/http/       # HTTP transport katmanı
        handler/          # Gin handler'ları
        router/           # Route tanımları
    middleware/           # Middleware bileşenleri
    dto/                  # Request/Response DTO'ları
    i18n/                 # Çoklu dil desteği
    factory/              # Dependency injection
    config/               # Konfigürasyon yönetimi
/pkg/
    utils/                # Yardımcı fonksiyonlar
    logger/               # Logger yapılandırması
    database/             # Database bağlantısı
    validator/            # Validation işlemleri
```

## 🚀 Teknolojiler

- **HTTP Framework**: Gin
- **Database**: PostgreSQL + GORM
- **Logger**: Zerolog
- **Validation**: go-playground/validator
- **JWT**: golang-jwt/jwt
- **File Storage**: AWS S3 Compatible
- **Config**: Environment variables + .env

## 🔧 Kurulum

1. **Gereksinimler**
   - Go 1.21+
   - PostgreSQL
   - AWS S3 Compatible Storage

2. **Projeyi klonlayın**
   ```bash
   git clone <repository-url>
   cd louco-event
   ```

3. **Bağımlılıkları yükleyin**
   ```bash
   go mod download
   ```

4. **Environment değişkenlerini ayarlayın**
   ```bash
   cp .env.example .env
   # .env dosyasını düzenleyin
   ```

5. **Uygulamayı çalıştırın**
   ```bash
   go run cmd/app/main.go
   ```

## ⚙️ Konfigürasyon

Aşağıdaki environment değişkenleri kullanılabilir:

### Server
- `SERVER_PORT`: HTTP server portu (varsayılan: 8080)
- `SERVER_MODE`: Çalışma modu (development/production)

### Database
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database kullanıcı adı
- `DB_PASSWORD`: Database şifresi
- `DB_NAME`: Database adı
- `DB_SSL_MODE`: SSL modu

### JWT
- `JWT_SECRET`: JWT secret key
- `JWT_EXPIRATION`: Token geçerlilik süresi (örn: 24h)

### AWS S3
- `AWS_ENDPOINT`: S3 endpoint URL'i
- `AWS_ACCESS_KEY_ID`: Access key
- `AWS_SECRET_ACCESS_KEY`: Secret key
- `AWS_DEFAULT_REGION`: Region
- `AWS_BUCKET`: Bucket adı

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register/step1` - Hesap oluşturma
- `POST /api/v1/auth/login` - Giriş yapma
- `POST /api/v1/auth/social-login` - Sosyal medya girişi

### User Management
- `GET /api/v1/users/profile` - Profil bilgilerini getir
- `PUT /api/v1/users/profile` - Profil güncelle
- `POST /api/v1/username/set` - Username belirleme
- `POST /api/v1/username/check` - Username kontrolü

### Media Upload
- `POST /api/v1/media/upload` - Dosya yükleme
- `GET /api/v1/media/:id` - Medya detayı
- `DELETE /api/v1/media/:id` - Medya silme

### Health Check
- `GET /health` - Sistem durumu

## 🌍 Çoklu Dil Desteği

Uygulama Türkçe ve İngilizce dillerini destekler. Dil seçimi `Accept-Language` header'ı ile yapılır.

Desteklenen diller:
- `tr`: Türkçe
- `en`: İngilizce (varsayılan)

## 🔐 Authentication

JWT tabanlı authentication kullanılır. Token'lar `Authorization: Bearer <token>` header'ı ile gönderilir.

## 📝 Validation

Tüm input'lar go-playground/validator ile doğrulanır. Custom validator'lar:
- `alphanum_underscore_dot`: Username formatı
- `e164`: Telefon numarası formatı

## 🗃️ Database

GORM AutoMigrate kullanılarak database şeması otomatik oluşturulur.

Ana tablolar:
- `users`: Kullanıcı bilgileri
- `media`: Medya dosyaları

## 🚦 Middleware

- **Logger**: Request/response loglama
- **Recovery**: Panic recovery
- **CORS**: Cross-origin resource sharing
- **I18n**: Çoklu dil desteği
- **Rate Limit**: İstek sınırlama
- **JWT Auth**: JWT doğrulama

## 🧪 Test

```bash
go test ./...
```

## 📦 Build

```bash
go build -o bin/app cmd/app/main.go
```

## 🐳 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/internal/i18n/locales ./internal/i18n/locales
CMD ["./main"]
```

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır.
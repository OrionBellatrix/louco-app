# Louco Event API Documentation

## Genel Bakış

Louco Event API, Clean Architecture prensiplerine göre geliştirilmiş, etkinlik yönetimi için tasarlanmış RESTful bir API'dir. Go dilinde Gin framework kullanılarak geliştirilmiştir.

## Temel Bilgiler

- **Base URL:** `http://localhost:8080`
- **API Version:** v1
- **Content-Type:** `application/json`
- **Authentication:** JWT Bearer Token

## Kimlik Doğrulama

API'nin çoğu endpoint'i JWT token gerektirir. Token'ı Authorization header'ında gönderin:

```
Authorization: Bearer <your_jwt_token>
```

## Dil Desteği

API Türkçe ve İngilizce dillerini destekler. Dil seçimi için `Accept-Language` header'ını kullanın:

- `tr` - Türkçe
- `en` - İngilizce (varsayılan)

## Standart Response Formatı

Tüm API yanıtları aşağıdaki standart formatı takip eder:

```json
{
  "success": true,
  "message": "i18n.key",
  "data": {},
  "errors": null
}
```

## Endpoint'ler

### 1. Health Check

#### GET /health
Sistem durumunu kontrol eder.

**Request:**
```bash
curl -X GET http://localhost:8080/health
```

**Response:**
```json
{
  "success": true,
  "message": "common.success",
  "data": {
    "status": "healthy",
    "timestamp": "2025-12-09T12:00:00Z",
    "version": "1.0.0",
    "database": "connected"
  }
}
```

---

### 2. Authentication Endpoints

#### POST /api/v1/auth/register/step1
Kullanıcı kaydının ilk adımı - hesap oluşturma ve **otomatik doğrulama kodu gönderimi**.

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "password": "SecurePass123!",
  "user_type": "user"
}
```

**Zorunlu Alanlar:**
- `identifier`: Geçerli email adresi veya ülke kodu ile telefon numarası
- `password`: Minimum 8 karakter
- `user_type`: "user" veya "creator"

**Otomatik Doğrulama Özelliği:**
- Sistem otomatik olarak identifier'ın email mi telefon mu olduğunu algılar
- 6 haneli doğrulama kodu üretir ve gönderir
- Email ise HTML template ile, telefon ise Twilio SMS ile gönderir
- Kod 10 dakika geçerlidir
- **Ayrı bir `/send` endpoint'ine gerek yoktur!**

**Identifier Örnekleri:**
- Email: `user@example.com`
- Telefon: `+905551234567`

**Response (Email ile kayıt):**
```json
{
  "success": true,
  "message": "auth.registration_success",
  "data": {
    "user_id": 1,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "verification": {
      "sent": true,
      "message": "verification.email_sent",
      "expires_in_minutes": 10
    }
  }
}
```

**Response (Telefon ile kayıt):**
```json
{
  "success": true,
  "message": "auth.registration_success",
  "data": {
    "user_id": 1,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "verification": {
      "sent": true,
      "message": "verification.sms_sent",
      "expires_in_minutes": 10
    }
  }
}
```

**Sonraki Adımlar:**
1. Email/SMS'inizden 6 haneli kodu alın
2. `/api/v1/verification/verify` endpoint'ini kullanarak kodu doğrulayın
3. Profil tamamlama adımlarına devam edin

#### POST /api/v1/auth/login
Kullanıcı girişi.

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "password": "SecurePass123!"
}
```

**Alanlar:**
- `identifier`: Email, telefon numarası veya kullanıcı adı
- `password`: Kullanıcı şifresi

**Response:**
```json
{
  "success": true,
  "message": "auth.login_success",
  "data": {
    "user": {
      "id": 1,
      "full_name": "John Doe",
      "username": "johndoe",
      "email": "user@example.com",
      "user_type": "user",
      "is_active": true
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### POST /api/v1/auth/social-login
Sosyal medya hesabı ile giriş.

**Request Body:**
```json
{
  "provider": "google",
  "social_id": "google_user_id_123",
  "email": "user@gmail.com",
  "full_name": "John Doe",
  "user_type": "user"
}
```

**Zorunlu Alanlar:**
- `provider`: "apple" veya "google"
- `social_id`: Sosyal medya sağlayıcısından gelen benzersiz ID
- `user_type`: "user" veya "creator"

#### POST /api/v1/auth/forgot-password
Şifre sıfırlama token'ı talep etme.

**Request Body:**
```json
{
  "identifier": "user@example.com"
}
```

#### POST /api/v1/auth/reset-password
Token ile şifre sıfırlama.

**Request Body:**
```json
{
  "token": "reset_token_from_email",
  "new_password": "NewSecurePass123!"
}
```

---

### 3. Verification System

Email ve telefon doğrulama sistemi. Kullanıcılar email ve telefon numaralarını 6 haneli OTP kodları ile doğrulayabilirler. Sistem otomatik olarak identifier'ın email mi telefon mu olduğunu algılar.

#### POST /api/v1/verification/send
Doğrulama kodu gönderme (Email/Phone Otomatik Algılama).

**Authentication:** Gerekli

**Request Body:**
```json
{
  "identifier": "user@example.com"
}
```

**Zorunlu Alanlar:**
- `identifier`: Geçerli email adresi veya E.164 formatında telefon numarası

**Davranış:**
- Identifier'ı otomatik olarak email veya telefon olarak algılar
- 6 haneli doğrulama kodu üretir
- Email ise HTML template ile, telefon ise Twilio SMS ile kod gönderir
- Kod 10 dakika geçerlidir
- Spam önleme için rate limit uygulanır

**Request (Email):**
```bash
curl -X POST http://localhost:8080/api/v1/verification/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "user@example.com"
  }'
```

**Request (Phone):**
```bash
curl -X POST http://localhost:8080/api/v1/verification/send \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "+905551234567"
  }'
```

**Response (Email):**
```json
{
  "success": true,
  "message": "verification.email_sent",
  "data": {
    "expires_in_minutes": 10
  }
}
```

**Response (Phone):**
```json
{
  "success": true,
  "message": "verification.sms_sent",
  "data": {
    "expires_in_minutes": 10
  }
}
```

#### POST /api/v1/verification/verify
Doğrulama kodu ile doğrulama (Email/Phone Otomatik Algılama).

**Authentication:** Gerekli

**Request Body:**
```json
{
  "identifier": "user@example.com",
  "code": "123456"
}
```

**Zorunlu Alanlar:**
- `identifier`: Doğrulama kodu gönderilen email adresi veya telefon numarası
- `code`: 6 haneli doğrulama kodu

**Davranış:**
- Identifier'ı otomatik olarak email veya telefon olarak algılar
- Doğrulama kodunu kontrol eder
- Email/telefonu doğrulanmış olarak işaretler
- `email_verified_at` veya `phone_verified_at` timestamp'ini ayarlar
- Kullanılan kodu siler
- Kullanıcının email/telefon bilgisini günceller (farklıysa)

**Request (Email):**
```bash
curl -X POST http://localhost:8080/api/v1/verification/verify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "user@example.com",
    "code": "123456"
  }'
```

**Request (Phone):**
```bash
curl -X POST http://localhost:8080/api/v1/verification/verify \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "+905551234567",
    "code": "123456"
  }'
```

**Response (Email Başarılı):**
```json
{
  "success": true,
  "message": "verification.email_verified",
  "data": null
}
```

**Response (Phone Başarılı):**
```json
{
  "success": true,
  "message": "verification.phone_verified",
  "data": null
}
```

**Response (Hatalı Kod):**
```json
{
  "success": false,
  "message": "verification.invalid_code",
  "data": null,
  "errors": ["Invalid or expired verification code"]
}
```

#### POST /api/v1/verification/resend
Doğrulama kodunu yeniden gönderme (Email/Phone Otomatik Algılama).

**Authentication:** Gerekli

**Request Body:**
```json
{
  "identifier": "user@example.com"
}
```

**Zorunlu Alanlar:**
- `identifier`: Doğrulama kodu yeniden gönderilecek email adresi veya telefon numarası

**Davranış:**
- Identifier'ı otomatik olarak email veya telefon olarak algılar
- Mevcut doğrulama kodlarını geçersiz kılar
- Yeni 6 haneli doğrulama kodu üretir
- Email ise HTML template ile, telefon ise SMS ile yeni kod gönderir
- Yeni kod 10 dakika geçerlidir
- Rate limit uygulanır (saatte max 3 yeniden gönderim)

**Request (Email):**
```bash
curl -X POST http://localhost:8080/api/v1/verification/resend \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "user@example.com"
  }'
```

**Request (Phone):**
```bash
curl -X POST http://localhost:8080/api/v1/verification/resend \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "+905551234567"
  }'
```

**Response (Email):**
```json
{
  "success": true,
  "message": "verification.email_resent",
  "data": {
    "expires_in_minutes": 10
  }
}
```

**Response (Phone):**
```json
{
  "success": true,
  "message": "verification.sms_resent",
  "data": {
    "expires_in_minutes": 10
  }
}
```

### Verification System Özellikleri

- **Otomatik Algılama:** Email ve telefon numarası otomatik olarak algılanır
- **6 Haneli OTP Kodları:** Güvenli ve kullanıcı dostu
- **HTML Email Templates:** Profesyonel görünümlü email'ler
- **Twilio SMS Integration:** Güvenilir SMS teslimatı
- **Rate Limiting:** Spam ve kötüye kullanım önleme
- **Code Expiration:** 10 dakika geçerlilik süresi
- **Multi-language Support:** Türkçe ve İngilizce mesajlar
- **Test Environment Support:** Test telefon numaraları
- **Automatic Cleanup:** Süresi dolan kodların otomatik temizlenmesi

### Verification Flow (Adım Adım)

#### Unified Verification Flow:
1. **Kod Gönderme:** `POST /api/v1/verification/send`
   - Kullanıcı identifier (email veya telefon) gönderir
   - Sistem otomatik olarak email mi telefon mu algılar
   - 6 haneli kod üretir ve uygun kanal ile gönderir
   - Kod 10 dakika geçerlidir

2. **Kod Doğrulama:** `POST /api/v1/verification/verify`
   - Kullanıcı identifier ve kodu gönderir
   - Sistem kodu doğrular
   - Email/telefon doğrulanmış olarak işaretlenir

3. **Yeniden Gönderim (İsteğe Bağlı):** `POST /api/v1/verification/resend`
   - Kod süresi dolmuşsa veya kaybolmuşsa
   - Yeni kod üretilir ve uygun kanal ile gönderilir

### Identifier Format Örnekleri

**Email Formatları:**
- `user@example.com`
- `test.user@domain.co.uk`
- `user+tag@example.org`

**Telefon Formatları (E.164):**
- `+905551234567` (Türkiye)
- `+15005550006` (Test numarası)
- `+447700900123` (İngiltere)

### Test Telefon Numaraları

Twilio test ortamı için kullanılabilir telefon numaraları:
- `+15005550006` - Başarılı SMS teslimi
- `+15005550001` - Geçersiz telefon numarası hatası
- `+15005550007` - SMS gönderim hatası

### Verification Middleware

Bazı endpoint'ler doğrulanmış email/telefon gerektirir. Bu endpoint'ler verification middleware ile korunur:

```json
{
  "success": false,
  "message": "verification.email_required",
  "data": null,
  "errors": ["Email verification required"]
}
```

**Korumalı Endpoint Örnekleri:**
- Profil güncelleme işlemleri
- Hassas bilgi değişiklikleri
- Premium özellikler

---

### 4. Username Management

#### POST /api/v1/username/check
Kullanıcı adı müsaitlik kontrolü.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "username": "johndoe123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "common.success",
  "data": {
    "exists": false
  }
}
```

#### POST /api/v1/username/set
Kullanıcı adı belirleme.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "username": "johndoe123"
}
```

---

### 4. User Management

#### GET /api/v1/users/profile
Kullanıcı profil bilgilerini getirme.

**Authentication:** Gerekli

**Response:**
```json
{
  "success": true,
  "message": "common.success",
  "data": {
    "user": {
      "id": 1,
      "full_name": "John Doe",
      "username": "johndoe123",
      "email": "user@example.com",
      "phone": "+905551234567",
      "user_type": "user",
      "address": "Istanbul, Turkey",
      "company_name": "Tech Corp",
      "biography": "Software developer",
      "birth_date": "1990-01-01T00:00:00Z",
      "profile_pic": "https://s3.example.com/profile.jpg",
      "is_active": true,
      "created_at": "2025-12-09T12:00:00Z",
      "updated_at": "2025-12-09T12:00:00Z"
    },
    "media_count": 5,
    "recent_media": []
  }
}
```

#### PUT /api/v1/users/profile
Profil bilgilerini güncelleme.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "address": "Ankara, Turkey",
  "company_name": "New Tech Corp",
  "biography": "Senior Software Developer",
  "birth_date": "1990-01-01T00:00:00Z"
}
```

#### PUT /api/v1/users/contact
İletişim bilgilerini güncelleme.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "phone": "+905559876543"
}
```

#### PUT /api/v1/users/profile-pic
Profil resmi ayarlama.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "media_id": 1
}
```

**Alanlar:**
- `media_id`: Profil resmi olarak kullanılacak medya ID'si (zorunlu)

**Davranış:**
- Medya ID'sinin var olduğunu doğrular
- Kullanıcının profil resmi referansını günceller
- Medya dosyası kimlik doğrulaması yapılan kullanıcıya ait olmalıdır

#### PUT /api/v1/users/cover-pic
Kapak resmi ayarlama.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "media_id": 2
}
```

**Alanlar:**
- `media_id`: Kapak resmi olarak kullanılacak medya ID'si (zorunlu)

**Davranış:**
- Medya ID'sinin var olduğunu doğrular
- Kullanıcının kapak resmi referansını günceller
- Medya dosyası kimlik doğrulaması yapılan kullanıcıya ait olmalıdır

#### POST /api/v1/users/register/step4
Kayıt işleminin 4. adımı - profil detaylarını tamamlama.

**Authentication:** Gerekli (1. adımdan gelen token)

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "user@example.com",
  "phone": "+905551234567",
  "address": "Istanbul, Turkey",
  "company_name": "Tech Corp",
  "biography": "Software developer passionate about technology",
  "birth_date": "1990-01-01T00:00:00Z"
}
```

#### POST /api/v1/users/change-password
Şifre değiştirme.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "current_password": "SecurePass123!",
  "new_password": "NewSecurePass123!"
}
```

#### DELETE /api/v1/users/deactivate
Hesap deaktivasyonu.

**Authentication:** Gerekli

---

### 5. Media Management

#### POST /api/v1/media/upload
Dosya yükleme.

**Authentication:** Gerekli

**Request:** Multipart form data
```
file: [binary file data]
```

**Desteklenen Formatlar:**
- **Resimler:** JPEG, JPG, PNG, WebP, HEIC, HEIF (max 10MB)
- **Videolar:** MP4, MOV, WebM, AVI, QuickTime (max 100MB)

**Özellikler:**
- Otomatik dosya türü tespiti
- Resim boyutlandırma (800px genişlik)
- Video format dönüştürme
- AWS S3 uyumlu depolama
- Metadata çıkarma (boyutlar, süre)

**Response:**
```json
{
  "success": true,
  "message": "media.file_uploaded",
  "data": {
    "media_id": 1,
    "file_type": "image",
    "original_name": "photo.jpg",
    "mime_type": "image/jpeg",
    "file_size": 1024000,
    "file_url": "https://s3.example.com/uploads/photo.jpg",
    "width": 800,
    "height": 600,
    "is_converted": true
  }
}
```

#### GET /api/v1/media/:id
ID ile medya detaylarını getirme.

**Authentication:** Gerekli

**Path Parameters:**
- `id`: Medya ID'si

#### GET /api/v1/media/user/:user_id
Kullanıcının medya dosyalarını getirme.

**Authentication:** Gerekli

**Path Parameters:**
- `user_id`: Kullanıcı ID'si

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe sayısı (varsayılan: 10, max: 100)
- `media_type`: Tür filtresi ("image" veya "video")

#### PUT /api/v1/media/:id
Medya metadata güncelleme.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "original_name": "Updated filename.jpg",
  "width": 1920,
  "height": 1080,
  "duration": 120
}
```

#### DELETE /api/v1/media/:id
Medya dosyası silme.

**Authentication:** Gerekli

---

### 6. Admin Operations

#### GET /api/v1/admin/users
Tüm kullanıcıları listeleme (Admin).

**Authentication:** Admin token gerekli

**Query Parameters:**
- `page`: Sayfa numarası
- `page_size`: Sayfa başına öğe sayısı
- `user_type`: Kullanıcı türü filtresi

#### GET /api/v1/admin/media
Tüm medya dosyalarını listeleme (Admin).

**Authentication:** Admin token gerekli

**Query Parameters:**
- `page`: Sayfa numarası
- `page_size`: Sayfa başına öğe sayısı
- `media_type`: Medya türü filtresi
- `user_id`: Kullanıcı ID filtresi

---

### 7. Industries (Sektörler)

#### GET /api/v1/industries
Tüm sektörleri listeleme.

**Authentication:** Gerekli değil (Public endpoint)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/industries \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "industry.get_all.success",
  "data": {
    "industries": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      },
      {
        "id": 2,
        "name": "Healthcare",
        "slug": "healthcare"
      },
      {
        "id": 3,
        "name": "Finance",
        "slug": "finance"
      }
    ],
    "total": 50
  }
}
```

**Özellikler:**
- Kimlik doğrulama gerektirmez
- Tüm sektörleri alfabetik sıraya göre döndürür
- Kayıt formlarında sektör dropdown'ı için kullanılabilir
- Creator kullanıcılar için sektör seçimi

#### GET /api/v1/industries/:id
ID ile sektör detaylarını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `id`: Sektör ID'si (integer)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/industries/1 \
  -H "Accept-Language: tr"
```

**Response (Başarılı):**
```json
{
  "success": true,
  "message": "industry.get_by_id.success",
  "data": {
    "id": 1,
    "name": "Technology",
    "slug": "technology"
  }
}
```

**Response (Bulunamadı):**
```json
{
  "success": false,
  "message": "industry.not_found",
  "data": null,
  "errors": ["Industry not found"]
}
```

#### GET /api/v1/industries/slug/:slug
Slug ile sektör detaylarını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `slug`: Sektör slug'ı (string, URL-friendly)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/industries/slug/technology \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "industry.get_by_slug.success",
  "data": {
    "id": 1,
    "name": "Technology",
    "slug": "technology"
  }
}
```

**Örnek Slug'lar:**
- `technology` - Teknoloji
- `healthcare` - Sağlık
- `finance` - Finans
- `food-beverage` - Gıda ve İçecek
- `real-estate` - Gayrimenkul
- `entertainment` - Eğlence
- `education` - Eğitim

**Kullanım Alanları:**
- SEO dostu URL'ler
- URL'lerde sektör filtreleme
- Sektör bazlı routing

---

### 8. Creator Management

#### POST /api/v1/creators
Creator profili oluşturma.

**Authentication:** Gerekli (Creator tipinde kullanıcı)

**Request Body:**
```json
{
  "weeztix_token": "{\"api_key\": \"your_api_key\", \"secret\": \"your_secret\"}",
  "company_name": "Tech Events Corp",
  "address": "İstanbul, Türkiye",
  "estimated_tickets": 1000,
  "estimated_events": 12,
  "industry_ids": [1, 3, 5]
}
```

**Zorunlu Alanlar:**
- `company_name`: Şirket adı (2-200 karakter)
- `address`: Adres (5-500 karakter)
- `estimated_tickets`: Tahmini bilet sayısı (minimum 1)
- `estimated_events`: Tahmini etkinlik sayısı (minimum 1)
- `industry_ids`: Sektör ID'leri dizisi (en az 1 sektör)

**Opsiyonel Alanlar:**
- `weeztix_token`: Weeztix entegrasyon token'ı (JSON formatında)

**Response:**
```json
{
  "success": true,
  "message": "creator.created",
  "data": {
    "id": 1,
    "user_id": 5,
    "weeztix_token": "{\"api_key\": \"your_api_key\", \"secret\": \"your_secret\"}",
    "company_name": "Tech Events Corp",
    "address": "İstanbul, Türkiye",
    "estimated_tickets": 1000,
    "estimated_events": 12,
    "industries": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      },
      {
        "id": 3,
        "name": "Finance",
        "slug": "finance"
      }
    ],
    "created_at": "2025-12-09T12:00:00Z",
    "updated_at": "2025-12-09T12:00:00Z"
  }
}
```

#### GET /api/v1/creators
Creator listesi getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe sayısı (varsayılan: 20, max: 100)
- `industry_id`: Sektör ID'si ile filtreleme

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/creators?page=1&page_size=10&industry_id=1" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "common.success",
  "data": {
    "creators": [
      {
        "id": 1,
        "user_id": 5,
        "weeztix_token": null,
        "company_name": "Tech Events Corp",
        "address": "İstanbul, Türkiye",
        "estimated_tickets": 1000,
        "estimated_events": 12,
        "industries": [
          {
            "id": 1,
            "name": "Technology",
            "slug": "technology"
          }
        ],
        "created_at": "2025-12-09T12:00:00Z",
        "updated_at": "2025-12-09T12:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

#### GET /api/v1/creators/:id
ID ile creator detaylarını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `id`: Creator ID'si

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/creators/1 \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "common.success",
  "data": {
    "id": 1,
    "user_id": 5,
    "weeztix_token": null,
    "company_name": "Tech Events Corp",
    "address": "İstanbul, Türkiye",
    "estimated_tickets": 1000,
    "estimated_events": 12,
    "industries": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      }
    ],
    "created_at": "2025-12-09T12:00:00Z",
    "updated_at": "2025-12-09T12:00:00Z"
  }
}
```

#### GET /api/v1/creators/me
Kendi creator profilimi getirme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Response:**
```json
{
  "success": true,
  "message": "common.success",
  "data": {
    "user": {
      "id": 5,
      "full_name": "John Doe",
      "username": "johndoe",
      "email": "creator@example.com",
      "phone": "+905551234567",
      "user_type": "creator",
      "biography": "Event organizer",
      "birth_date": "1985-01-01T00:00:00Z",
      "profile_pic_id": 1,
      "cover_pic_id": 2,
      "is_active": true,
      "created_at": "2025-12-09T11:00:00Z",
      "updated_at": "2025-12-09T12:00:00Z"
    },
    "creator": {
      "id": 1,
      "user_id": 5,
      "weeztix_token": "{\"api_key\": \"your_api_key\"}",
      "company_name": "Tech Events Corp",
      "address": "İstanbul, Türkiye",
      "estimated_tickets": 1000,
      "estimated_events": 12,
      "industries": [
        {
          "id": 1,
          "name": "Technology",
          "slug": "technology"
        }
      ],
      "created_at": "2025-12-09T12:00:00Z",
      "updated_at": "2025-12-09T12:00:00Z"
    }
  }
}
```

#### PUT /api/v1/creators/me
Creator profili güncelleme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Request Body:**
```json
{
  "company_name": "Updated Tech Events Corp",
  "address": "Ankara, Türkiye",
  "estimated_tickets": 1500,
  "estimated_events": 15,
  "industry_ids": [1, 2, 3]
}
```

**Opsiyonel Alanlar:**
- `company_name`: Şirket adı
- `address`: Adres
- `estimated_tickets`: Tahmini bilet sayısı
- `estimated_events`: Tahmini etkinlik sayısı
- `industry_ids`: Sektör ID'leri dizisi

**Response:**
```json
{
  "success": true,
  "message": "creator.updated",
  "data": null
}
```

#### PUT /api/v1/creators/me/weeztix-token
Weeztix token güncelleme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Request Body:**
```json
{
  "weeztix_token": "{\"api_key\": \"new_api_key\", \"secret\": \"new_secret\", \"webhook_url\": \"https://example.com/webhook\"}"
}
```

**Zorunlu Alanlar:**
- `weeztix_token`: Weeztix entegrasyon token'ı (JSON string formatında, 10-255 karakter)

**Response:**
```json
{
  "success": true,
  "message": "creator.token_updated",
  "data": null
}
```

#### DELETE /api/v1/creators/me
Creator profili silme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Response:**
```json
{
  "success": true,
  "message": "creator.deleted",
  "data": null
}
```

**Not:** Creator profili silindiğinde:
- Creator-Industry ilişkileri otomatik olarak silinir
- User profili silinmez, sadece Creator profili silinir
- Kullanıcı normal "user" tipinde kalmaya devam eder

---

### 9. Categories (Kategoriler)

Sınırsız derinliğe sahip nested (çok seviyeli) kategori yapısı. Etkinlik kategorilerini hiyerarşik olarak organize eder.

#### Kategori Tipleri
- `concerts_&_festivals` - Konserler ve Festivaller
- `party` - Parti
- `culture` - Kültür
- `shows` - Gösteriler
- `sports` - Spor
- `freetime_activities` - Boş Zaman Aktiviteleri
- `business` - İş
- `ethnic` - Etnik
- `other` - Diğer

#### GET /api/v1/categories/tree
Tam kategori ağacını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/tree \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.tree_success",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "Concerts & Festivals",
        "slug": "concerts-festivals",
        "type": "concerts_&_festivals",
        "parent_id": null,
        "depth": 0,
        "icon": {
          "id": 10,
          "file_url": "https://s3.example.com/icons/concerts.png",
          "file_type": "image"
        },
        "children": [
          {
            "id": 2,
            "name": "Rock Music",
            "slug": "rock-music",
            "type": "concerts_&_festivals",
            "parent_id": 1,
            "depth": 1,
            "icon": null,
            "children": [
              {
                "id": 3,
                "name": "Heavy Metal",
                "slug": "heavy-metal",
                "type": "concerts_&_festivals",
                "parent_id": 2,
                "depth": 2,
                "icon": null,
                "children": []
              }
            ]
          }
        ]
      }
    ]
  }
}
```

**Özellikler:**
- Sınırsız derinlikte hiyerarşik yapı
- 10 dakika Redis cache ile optimize edilmiş performans
- Her kategori için icon media desteği
- Nested Set Model ile verimli tree operasyonları

#### GET /api/v1/categories/tree/type/:type
Tipe göre kategori ağacını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `type`: Kategori tipi (concerts_&_festivals, party, culture, shows, sports, freetime_activities, business, ethnic, other)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/tree/type/concerts_&_festivals \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.tree_by_type_success",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "Concerts & Festivals",
        "slug": "concerts-festivals",
        "type": "concerts_&_festivals",
        "parent_id": null,
        "depth": 0,
        "children": [
          {
            "id": 2,
            "name": "Rock Music",
            "slug": "rock-music",
            "type": "concerts_&_festivals",
            "parent_id": 1,
            "depth": 1,
            "children": []
          }
        ]
      }
    ]
  }
}
```

#### GET /api/v1/categories/:id
ID ile kategori detaylarını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `id`: Kategori ID'si

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/1 \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.get_success",
  "data": {
    "id": 1,
    "name": "Concerts & Festivals",
    "slug": "concerts-festivals",
    "type": "concerts_&_festivals",
    "parent_id": null,
    "depth": 0,
    "lft": 1,
    "rgt": 20,
    "icon": {
      "id": 10,
      "file_url": "https://s3.example.com/icons/concerts.png",
      "file_type": "image"
    },
    "created_at": "2025-12-09T12:00:00Z",
    "updated_at": "2025-12-09T12:00:00Z"
  }
}
```

#### GET /api/v1/categories/slug/:slug
Slug ile kategori detaylarını getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `slug`: Kategori slug'ı (URL-friendly)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/slug/rock-music \
  -H "Accept-Language: tr"
```

**Örnek Slug'lar:**
- `rock-music` - Rock Müzik
- `electronic-music` - Elektronik Müzik
- `food-festivals` - Yemek Festivalleri
- `art-exhibitions` - Sanat Sergileri
- `football-matches` - Futbol Maçları

#### GET /api/v1/categories/:id/children
Kategori alt kategorilerini getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `id`: Ana kategori ID'si

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/1/children \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.children_success",
  "data": {
    "children": [
      {
        "id": 2,
        "name": "Rock Music",
        "slug": "rock-music",
        "type": "concerts_&_festivals",
        "parent_id": 1,
        "depth": 1,
        "icon": null
      },
      {
        "id": 3,
        "name": "Electronic Music",
        "slug": "electronic-music",
        "type": "concerts_&_festivals",
        "parent_id": 1,
        "depth": 1,
        "icon": null
      }
    ]
  }
}
```

#### GET /api/v1/categories/:id/parents
Kategori üst kategorilerini getirme (breadcrumb).

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `id`: Kategori ID'si

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/5/parents \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.parents_success",
  "data": {
    "parents": [
      {
        "id": 1,
        "name": "Concerts & Festivals",
        "slug": "concerts-festivals",
        "type": "concerts_&_festivals",
        "parent_id": null,
        "depth": 0
      },
      {
        "id": 2,
        "name": "Rock Music",
        "slug": "rock-music",
        "type": "concerts_&_festivals",
        "parent_id": 1,
        "depth": 1
      }
    ]
  }
}
```

#### GET /api/v1/categories/roots
Ana kategorileri getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/roots \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.roots_success",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "Concerts & Festivals",
        "slug": "concerts-festivals",
        "type": "concerts_&_festivals",
        "parent_id": null,
        "depth": 0,
        "icon": {
          "id": 10,
          "file_url": "https://s3.example.com/icons/concerts.png"
        }
      },
      {
        "id": 10,
        "name": "Party",
        "slug": "party",
        "type": "party",
        "parent_id": null,
        "depth": 0,
        "icon": null
      }
    ]
  }
}
```

#### GET /api/v1/categories/leaves
Son seviye kategorileri getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/leaves \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.leaves_success",
  "data": {
    "categories": [
      {
        "id": 5,
        "name": "Heavy Metal",
        "slug": "heavy-metal",
        "type": "concerts_&_festivals",
        "parent_id": 2,
        "depth": 2,
        "icon": null
      },
      {
        "id": 8,
        "name": "House Music",
        "slug": "house-music",
        "type": "concerts_&_festivals",
        "parent_id": 3,
        "depth": 2,
        "icon": null
      }
    ]
  }
}
```

#### GET /api/v1/categories/type/:type
Tipe göre kategorileri getirme (düz liste).

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `type`: Kategori tipi

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/categories/type/concerts_&_festivals \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.by_type_success",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "Concerts & Festivals",
        "slug": "concerts-festivals",
        "type": "concerts_&_festivals",
        "parent_id": null,
        "depth": 0
      },
      {
        "id": 2,
        "name": "Rock Music",
        "slug": "rock-music",
        "type": "concerts_&_festivals",
        "parent_id": 1,
        "depth": 1
      }
    ]
  }
}
```

#### GET /api/v1/categories/search
Kategori arama.

**Authentication:** Gerekli değil (Public endpoint)

**Query Parameters:**
- `q`: Arama sorgusu (minimum 2 karakter)
- `type`: Opsiyonel kategori tipi filtresi

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/categories/search?q=music&type=concerts_&_festivals" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.search_success",
  "data": {
    "categories": [
      {
        "id": 2,
        "name": "Rock Music",
        "slug": "rock-music",
        "type": "concerts_&_festivals",
        "parent_id": 1,
        "depth": 1,
        "icon": null
      },
      {
        "id": 3,
        "name": "Electronic Music",
        "slug": "electronic-music",
        "type": "concerts_&_festivals",
        "parent_id": 1,
        "depth": 1,
        "icon": null
      }
    ]
  }
}
```

#### POST /api/v1/admin/categories/cache/refresh
Kategori cache'ini yenileme (Admin).

**Authentication:** Gerekli (Admin - creator user type)

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/categories/cache/refresh \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.cache_refresh_success",
  "data": null
}
```

#### DELETE /api/v1/admin/categories/cache/clear
Kategori cache'ini temizleme (Admin).

**Authentication:** Gerekli (Admin - creator user type)

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/categories/cache/clear \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "category.cache_clear_success",
  "data": null
}
```

**Kategori Sistemi Özellikleri:**
- **Nested Set Model:** Sınırsız derinlikte hiyerarşik yapı
- **Redis Cache:** 10 dakika cache ile optimize edilmiş performans
- **Readonly System:** Sadece okuma operasyonları (CRUD yok)
- **Icon Support:** Her kategori için medya icon desteği
- **Multi-language:** i18n desteği ile çoklu dil
- **SEO Friendly:** URL-friendly slug'lar
- **Type-based Filtering:** Kategori tipine göre filtreleme
- **Search Functionality:** İsim ve slug'da arama
- **Tree Operations:** Parent/child, breadcrumb, leaf/root operations

---

### 10. Creator Registration Flow

Creator tipindeki kullanıcılar için kayıt akışı:

#### 1. Register Step 1
```bash
curl -X POST http://localhost:8080/api/v1/auth/register/step1 \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "creator@example.com",
    "password": "SecurePass123!",
    "user_type": "creator"
  }'
```

#### 2. Register Step 4 (Creator için ek alanlar)
```bash
curl -X POST http://localhost:8080/api/v1/users/register/step4 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "full_name": "John Doe",
    "email": "creator@example.com",
    "phone": "+905551234567",
    "biography": "Professional event organizer",
    "birth_date": "1985-01-01T00:00:00Z",
    "company_name": "Tech Events Corp",
    "address": "İstanbul, Türkiye",
    "estimated_tickets": 1000,
    "estimated_events": 12,
    "industry_ids": [1, 3, 5]
  }'
```

**Creator için Zorunlu Alanlar (Step 4):**
- `company_name`: Şirket adı
- `address`: Adres
- `estimated_tickets`: Tahmini bilet sayısı
- `estimated_events`: Tahmini etkinlik sayısı
- `industry_ids`: En az bir sektör seçimi

#### 3. Creator Profili Otomatik Oluşturma
RegisterStep4 tamamlandığında, creator tipindeki kullanıcılar için otomatik olarak Creator profili oluşturulur.

---

## Hata Kodları

| HTTP Status | Error Code | Açıklama |
|-------------|------------|----------|
| 400 | BAD_REQUEST | Geçersiz istek |
| 401 | UNAUTHORIZED | Kimlik doğrulama gerekli |
| 403 | FORBIDDEN | Erişim izni yok |
| 404 | NOT_FOUND | Kaynak bulunamadı |
| 409 | CONFLICT | Çakışma (örn: email zaten mevcut) |
| 422 | VALIDATION_ERROR | Doğrulama hatası |
| 500 | INTERNAL_SERVER_ERROR | Sunucu hatası |

## Validation Kuralları

### Email
- RFC 5322 standardına uygun olmalı
- Örnek: `user@example.com`

### Telefon
- E.164 formatında olmalı
- Örnek: `+905551234567`

### Şifre
- Minimum 8 karakter
- En az 1 büyük harf, 1 küçük harf, 1 rakam içermeli

### Kullanıcı Adı
- 3-30 karakter arası
- Sadece alfanumerik karakterler, alt çizgi ve nokta
- Örnek: `johndoe123`, `user.name`, `user_123`

## Rate Limiting

API rate limiting uygulanmıştır:
- Genel endpoint'ler: 100 istek/dakika
- Upload endpoint'leri: 10 istek/dakika
- Auth endpoint'leri: 20 istek/dakika

## CORS

API CORS desteği sunar ve aşağıdaki origin'lere izin verir:
- `http://localhost:3000` (Development)
- `https://yourdomain.com` (Production)

## Güvenlik

- JWT token'lar 24 saat geçerlidir
- Şifreler bcrypt ile hash'lenir
- Rate limiting uygulanır
- CORS koruması aktiftir
- Request logging yapılır
- Panic recovery middleware'i aktiftir

## Örnek Kullanım Senaryosu

### 1. Kullanıcı Kaydı
```bash
# 1. Adım: Hesap oluşturma
curl -X POST http://localhost:8080/api/v1/auth/register/step1 \
  -H "Content-Type: application/json" \
  -H "Accept-Language: tr" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "user_type": "user"
  }'

# Response'dan token'ı al ve sonraki isteklerde kullan
```

### 2. Kullanıcı Girişi
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "Accept-Language: tr" \
  -d '{
    "identifier": "user@example.com",
    "password": "SecurePass123!"
  }'
```

### 3. Dosya Yükleme
```bash
curl -X POST http://localhost:8080/api/v1/media/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -F "file=@/path/to/your/image.jpg"
```

### 4. Profil Güncelleme
```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "full_name": "Updated Name",
    "biography": "Updated biography"
  }'
```

## Postman Collection

API'yi test etmek için hazırlanmış Postman collection'ı `postman/Louco-Event-API.postman_collection.json` dosyasında bulabilirsiniz.

### Collection'ı İçe Aktarma:
1. Postman'i açın
2. Import butonuna tıklayın
3. `postman/Louco-Event-API.postman_collection.json` dosyasını seçin
4. Collection içindeki environment variables'ları ayarlayın:
   - `base_url`: `http://localhost:8080`
   - `jwt_token`: Login'den aldığınız token
   - `admin_jwt_token`: Admin token (gerekirse)

## Geliştirme Ortamı

### Gereksinimler
- Go 1.21+
- PostgreSQL 13+
- AWS S3 uyumlu storage (MinIO, AWS S3, vb.)

### Çalıştırma
```bash
# Bağımlılıkları yükle
go mod tidy

# Uygulamayı çalıştır
go run cmd/app/main.go
```

### Environment Variables
`.env` dosyasında aşağıdaki değişkenleri ayarlayın:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=louco_event_db
JWT_SECRET=your_jwt_secret
AWS_ENDPOINT=your_s3_endpoint
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_BUCKET=your_bucket_name
```

---

### 11. Follow System (Takip Sistemi)

Instagram benzeri tek yönlü takip sistemi. Kullanıcılar birbirlerini takip edebilir, takipten çıkabilir ve takipçi/takip edilen listelerini görüntüleyebilirler.

#### Sistem Özellikleri
- **Tek Yönlü Takip:** Instagram gibi karşılıklı onay gerektirmez
- **Kendini Takip Etme Engeli:** Kullanıcılar kendilerini takip edemez
- **Duplicate Takip Engeli:** Aynı kullanıcıyı birden fazla takip etme engellenir
- **Otomatik Sayaç Güncelleme:** Takipçi/takip edilen sayıları otomatik güncellenir
- **Sayfalama Desteği:** Takipçi/takip edilen listeleri sayfalama ile sunulur
- **Karşılıklı Takip Kontrolü:** İki kullanıcı arasındaki takip durumu kontrol edilebilir

#### POST /api/v1/follows
Bir kullanıcıyı takip etme.

**Authentication:** Gerekli

**Request Body:**
```json
{
  "user_id": 2
}
```

**Zorunlu Alanlar:**
- `user_id`: Takip edilecek kullanıcının ID'si

**Davranış:**
- Mevcut kullanıcı ile hedef kullanıcı arasında takip ilişkisi oluşturur
- Kendini takip etmeyi engeller
- Duplicate takip ilişkisini engeller
- Takipçi/takip edilen sayılarını otomatik günceller
- Tek yönlü ilişki (karşılıklı takip gerektirmez)

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/follows \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "user_id": 2
  }'
```

**Response (Başarılı):**
```json
{
  "success": true,
  "message": "follow.success",
  "data": {
    "follower_id": 1,
    "following_id": 2,
    "created_at": "2025-12-10T12:00:00Z"
  }
}
```

**Response (Kendini Takip Etme Hatası):**
```json
{
  "success": false,
  "message": "follow.cannot_follow_self",
  "data": null,
  "errors": ["Cannot follow yourself"]
}
```

**Response (Zaten Takip Ediliyor):**
```json
{
  "success": false,
  "message": "follow.already_following",
  "data": null,
  "errors": ["Already following this user"]
}
```

#### DELETE /api/v1/follows/:user_id
Bir kullanıcıyı takipten çıkarma.

**Authentication:** Gerekli

**Path Parameters:**
- `user_id`: Takipten çıkarılacak kullanıcının ID'si

**Davranış:**
- Mevcut takip ilişkisini kaldırır
- Takipçi/takip edilen sayılarını otomatik günceller
- İlişki yoksa bile başarılı yanıt döner
- Sadece kendi takip ilişkilerini etkiler

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/follows/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "follow.unfollow_success",
  "data": null
}
```

#### GET /api/v1/follows/status/:user_id
İki kullanıcı arasındaki takip durumunu kontrol etme.

**Authentication:** Gerekli

**Path Parameters:**
- `user_id`: Takip durumu kontrol edilecek kullanıcının ID'si

**Response:**
- `user_id`: Hedef kullanıcının ID'si
- `is_following`: Bu kullanıcıyı takip ediyor musunuz?
- `is_followed_by`: Bu kullanıcı sizi takip ediyor mu?

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/follows/status/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "follow.status_retrieved",
  "data": {
    "user_id": 2,
    "is_following": true,
    "is_followed_by": false
  }
}
```

**Kullanım Alanları:**
- Takip et/takipten çık buton durumlarını gösterme
- Karşılıklı takip göstergelerini gösterme
- UI için ilişki durumunu belirleme

#### GET /api/v1/follows/followers
Takipçi listesini getirme.

**Authentication:** Gerekli

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe sayısı (varsayılan: 20, max: 100)

**Response:**
- Sayfalanmış takipçi listesi
- Her takipçi için kullanıcı detayları ve takip timestamp'i
- Profil resimleri ve temel kullanıcı bilgileri
- Sayfalama metadata'sı (toplam, sayfa, vb.)

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/follows/followers?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "follow.followers_retrieved",
  "data": {
    "followers": [
      {
        "user": {
          "id": 2,
          "full_name": "Jane Smith",
          "username": "janesmith",
          "profile_pic": {
            "id": 5,
            "file_url": "https://s3.example.com/profile2.jpg",
            "file_type": "image"
          },
          "followers_count": 150,
          "following_count": 200
        },
        "followed_at": "2025-12-10T10:30:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

**Kullanım Alanları:**
- Profilde takipçi listesi gösterme
- Sizi kimin takip ettiğini görme
- Takipçi yönetim arayüzü

#### GET /api/v1/follows/following
Takip edilen kullanıcı listesini getirme.

**Authentication:** Gerekli

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe sayısı (varsayılan: 20, max: 100)

**Response:**
- Sayfalanmış takip edilen kullanıcı listesi
- Her kullanıcı için detaylar ve takip timestamp'i
- Profil resimleri ve temel kullanıcı bilgileri
- Sayfalama metadata'sı

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/follows/following?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "follow.following_retrieved",
  "data": {
    "following": [
      {
        "user": {
          "id": 3,
          "full_name": "Bob Johnson",
          "username": "bobjohnson",
          "profile_pic": {
            "id": 7,
            "file_url": "https://s3.example.com/profile3.jpg",
            "file_type": "image"
          },
          "followers_count": 300,
          "following_count": 100
        },
        "followed_at": "2025-12-10T11:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 10,
    "total_pages": 2
  }
}
```

**Kullanım Alanları:**
- Profilde takip edilen listesi gösterme
- Kimi takip ettiğinizi görme
- Takip yönetim arayüzü
- Takipten çıkma işlevselliği

#### GET /api/v1/follows/mutual/:user_id
Karşılıklı takip edilen kullanıcıları getirme.

**Authentication:** Gerekli

**Path Parameters:**
- `user_id`: Karşılıklı takip edilenleri bulunacak kullanıcının ID'si

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe sayısı (varsayılan: 20, max: 100)

**Response:**
- Hem sizin hem de hedef kullanıcının takip ettiği kullanıcılar
- Ortak bağlantıları gösterir
- Sayfalanmış liste formatında

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/follows/mutual/2?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "follow.mutual_follows_retrieved",
  "data": {
    "mutual_follows": [
      {
        "user": {
          "id": 4,
          "full_name": "Alice Brown",
          "username": "alicebrown",
          "profile_pic": {
            "id": 9,
            "file_url": "https://s3.example.com/profile4.jpg",
            "file_type": "image"
          },
          "followers_count": 500,
          "following_count": 250
        }
      }
    ],
    "total": 8,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

**Kullanım Alanları:**
- Kullanıcılar arası ortak bağlantıları gösterme
- Ortak ilgi alanlarını keşfetme
- Sosyal ağ analizi
- Ortak bağlantılara dayalı arkadaş önerileri

### Follow System Özellikleri

- **Instagram-like Follow System:** Tek yönlü takip sistemi
- **Automatic Count Management:** Takipçi/takip edilen sayıları otomatik güncelleme
- **Duplicate Prevention:** Aynı kullanıcıyı birden fazla takip etme engeli
- **Self-Follow Prevention:** Kendini takip etme engeli
- **Pagination Support:** Tüm listelerde sayfalama desteği
- **Mutual Follow Detection:** Karşılıklı takip durumu tespiti
- **Profile Integration:** Kullanıcı profillerinde takip sayıları gösterimi
- **Multi-language Support:** Türkçe ve İngilizce mesaj desteği
- **Efficient Queries:** Optimize edilmiş veritabanı sorguları
- **User Info Preloading:** Takip listelerinde kullanıcı bilgileri önceden yükleme

### Follow System Database Schema

```sql
-- Follows table
CREATE TABLE follows (
    id SERIAL PRIMARY KEY,
    follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

-- Users table (updated with follow counts)
ALTER TABLE users ADD COLUMN followers_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN following_count INTEGER DEFAULT 0;

-- Indexes for performance
CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
CREATE INDEX idx_follows_created_at ON follows(created_at);
```

### Follow System API Flow

#### 1. Kullanıcı Takip Etme
```bash
# 1. Kullanıcıyı takip et
curl -X POST http://localhost:8080/api/v1/follows \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"user_id": 2}'

# 2. Takip durumunu kontrol et
curl -X GET http://localhost:8080/api/v1/follows/status/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 2. Takipçi/Takip Edilen Listeleri
```bash
# Takipçilerimi getir
curl -X GET http://localhost:8080/api/v1/follows/followers \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Takip ettiklerimi getir
curl -X GET http://localhost:8080/api/v1/follows/following \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 3. Karşılıklı Takip ve Takipten Çıkma
```bash
# Karşılıklı takip edilenleri getir
curl -X GET http://localhost:8080/api/v1/follows/mutual/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Kullanıcıyı takipten çıkar
curl -X DELETE http://localhost:8080/api/v1/follows/2 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 12. Address Management System

Google Places API entegrasyonu ile adres yönetim sistemi. Location-based etkinlikler için adres oluşturma, arama ve yönetim işlevleri sunar.

#### Sistem Özellikleri
- **Google Places API Integration:** Place ID ile adres doğrulama
- **Geographic Data Support:** Latitude/longitude koordinatları
- **Advanced Search & Filtering:** Şehir, ülke, bölge bazlı filtreleme
- **Location-based Events:** Etkinlikler için adres referansı
- **Creator-only Access:** Sadece Creator kullanıcılar adres oluşturabilir
- **Comprehensive Address Data:** Tam adres, posta kodu, kapı numarası desteği

#### POST /api/v1/addresses
Yeni adres oluşturma.

**Authentication:** Gerekli (Creator kullanıcısı)

**Request Body:**
```json
{
  "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
  "full_address": "Istanbul Convention Center, Harbiye Mahallesi, Şişli/Istanbul, Turkey",
  "country": "Turkey",
  "city": "Istanbul",
  "district": "Şişli",
  "street": "Harbiye Mahallesi",
  "postal_code": "34367",
  "latitude": 41.0082,
  "longitude": 28.9784,
  "door_number": "1"
}
```

**Zorunlu Alanlar:**
- `place_id`: Google Places API'den gelen benzersiz ID
- `full_address`: Tam adres metni (5-500 karakter)
- `country`: Ülke adı (2-100 karakter)
- `city`: Şehir adı (2-100 karakter)
- `latitude`: Enlem (-90 ile 90 arası)
- `longitude`: Boylam (-180 ile 180 arası)

**Opsiyonel Alanlar:**
- `district`: İlçe/bölge adı
- `street`: Sokak adı
- `postal_code`: Posta kodu
- `door_number`: Kapı numarası

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr" \
  -d '{
    "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
    "full_address": "Istanbul Convention Center, Harbiye Mahallesi, Şişli/Istanbul, Turkey",
    "country": "Turkey",
    "city": "Istanbul",
    "district": "Şişli",
    "street": "Harbiye Mahallesi",
    "postal_code": "34367",
    "latitude": 41.0082,
    "longitude": 28.9784,
    "door_number": "1"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "address.create.success",
  "data": {
    "id": 1,
    "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
    "full_address": "Istanbul Convention Center, Harbiye Mahallesi, Şişli/Istanbul, Turkey",
    "country": "Turkey",
    "city": "Istanbul",
    "district": "Şişli",
    "street": "Harbiye Mahallesi",
    "postal_code": "34367",
    "latitude": 41.0082,
    "longitude": 28.9784,
    "door_number": "1",
    "created_at": "2025-12-11T12:00:00Z",
    "updated_at": "2025-12-11T12:00:00Z"
  }
}
```

#### GET /api/v1/addresses/:id
ID ile adres detaylarını getirme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Path Parameters:**
- `id`: Adres ID'si

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/addresses/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "address.get.success",
  "data": {
    "id": 1,
    "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
    "full_address": "Istanbul Convention Center, Harbiye Mahallesi, Şişli/Istanbul, Turkey",
    "country": "Turkey",
    "city": "Istanbul",
    "district": "Şişli",
    "street": "Harbiye Mahallesi",
    "postal_code": "34367",
    "latitude": 41.0082,
    "longitude": 28.9784,
    "door_number": "1",
    "created_at": "2025-12-11T12:00:00Z",
    "updated_at": "2025-12-11T12:00:00Z"
  }
}
```

#### GET /api/v1/addresses/search
Filtreleme ile adres arama.

**Authentication:** Gerekli (Creator kullanıcısı)

**Query Parameters:**
- `country`: Ülke filtresi
- `city`: Şehir filtresi
- `district`: İlçe filtresi
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/addresses/search?country=Turkey&city=Istanbul&page=1&page_size=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "address.search.success",
  "data": {
    "items": [
      {
        "id": 1,
        "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
        "full_address": "Istanbul Convention Center...",
        "country": "Turkey",
        "city": "Istanbul",
        "district": "Şişli",
        "latitude": 41.0082,
        "longitude": 28.9784,
        "created_at": "2025-12-11T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 25,
      "total_pages": 2
    }
  }
}
```

#### GET /api/v1/addresses/city/:city
Şehir bazlı adres listeleme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Path Parameters:**
- `city`: Şehir adı

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)

**Request:**
```bash
curl -X GET "http://localhost:8080/api/v1/addresses/city/Istanbul?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Accept-Language: tr"
```

**Response:**
```json
{
  "success": true,
  "message": "address.search.success",
  "data": {
    "items": [
      {
        "id": 1,
        "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
        "full_address": "Istanbul Convention Center...",
        "country": "Turkey",
        "city": "Istanbul",
        "district": "Şişli",
        "latitude": 41.0082,
        "longitude": 28.9784
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 15,
      "total_pages": 2
    }
  }
}
```

### Address + Event Workflow

Address Management System, Event Management System ile entegre çalışır. Location-based etkinlikler için iki aşamalı süreç:

#### 1. Adres Oluşturma
```bash
# Önce adres oluştur
curl -X POST http://localhost:8080/api/v1/addresses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
    "full_address": "Istanbul Convention Center...",
    "country": "Turkey",
    "city": "Istanbul",
    "latitude": 41.0082,
    "longitude": 28.9784
  }'

# Response'dan address_id'yi al (örn: 1)
```

#### 2. Etkinlik Oluşturma
```bash
# Sonra etkinliği address_id ile oluştur
curl -X POST http://localhost:8080/api/v1/events/manage \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Tech Conference 2025",
    "description": "Annual technology conference...",
    "type": "public",
    "location_type": "location",
    "address_id": 1,
    "start_date": "2025-03-15",
    "start_time": "09:00",
    "category_ids": [1, 5, 8]
  }'
```

### Address Validation Rules

**Google Places Integration:**
- `place_id`: Google Places API'den gelen geçerli ID
- Duplicate place_id kontrolü (aynı place_id ile birden fazla adres oluşturulamaz)
- Koordinat doğrulama (latitude: -90/+90, longitude: -180/+180)

**Address Data:**
- `full_address`: 5-500 karakter arası
- `country`: 2-100 karakter arası
- `city`: 2-100 karakter arası
- `district`, `street`: Opsiyonel, 2-100 karakter arası
- `postal_code`: Opsiyonel, 2-20 karakter arası
- `door_number`: Opsiyonel, 1-50 karakter arası

**Geographic Constraints:**
- Latitude: -90.0 ile +90.0 arası (Kuzey/Güney kutupları dahil)
- Longitude: -180.0 ile +180.0 arası (Doğu/Batı meridyenleri dahil)
- Koordinat hassasiyeti: 6 ondalık basamak desteklenir

---

### 13. Event Management System

Kapsamlı etkinlik yönetim sistemi. Creator kullanıcılar etkinlik oluşturabilir, yönetebilir ve bilet satabilirler. Sistem public/private etkinlikler, location/online/announcement türleri, davetiye sistemi ve bilet satış yönetimi sunar.

#### Event Types & Features
- **Public Events**: Herkese açık etkinlikler
- **Private Events**: Davetiye ile katılım
- **Location Types**:
  - `location`: Fiziksel mekan etkinlikleri
  - `online`: Online etkinlikler
  - `announcement`: Duyuru türü etkinlikler
- **Status Workflow**: `draft` → `pending` → `published`/`rejected` → `cancelled`
- **Ticket System**: Bilet oluşturma, satış takibi, gelir yönetimi
- **Invitation System**: Kullanıcı ve email bazlı davetler

#### POST /api/v1/events/manage
Yeni etkinlik oluşturma.

**Authentication:** Gerekli (Creator kullanıcısı)

**Request Body:**
```json
{
  "name": "Tech Conference 2025",
  "description": "Annual technology conference with industry leaders",
  "type": "public",
  "location_type": "location",
  "start_date": "2025-03-15",
  "start_time": "09:00",
  "end_date": "2025-03-15",
  "end_time": "18:00",
  "address_id": 1,
  "category_ids": [1, 5, 8],
  "image_id": 10,
  "video_id": 15
}
```

**Zorunlu Alanlar:**
- `name`: Etkinlik adı (3-200 karakter)
- `description`: Etkinlik açıklaması (10-2000 karakter)
- `type`: Etkinlik türü (`public` veya `private`)
- `location_type`: Konum türü (`location`, `online`, `announcement`)
- `start_date`: Başlangıç tarihi (YYYY-MM-DD)
- `start_time`: Başlangıç saati (HH:MM)
- `category_ids`: Kategori ID'leri dizisi (en az 1)

**Koşullu Zorunlu Alanlar:**
- `address_id`: `location_type` = `location` ise zorunlu (önceden oluşturulmuş adres ID'si)
- `end_date`, `end_time`: `location_type` ≠ `announcement` ise zorunlu

**Not:** Location-based etkinlikler için önce `/api/v1/addresses` endpoint'i ile adres oluşturulmalı, sonra dönen `address_id` kullanılmalıdır.

**Response:**
```json
{
  "success": true,
  "message": "event.created",
  "data": {
    "id": 1,
    "name": "Tech Conference 2025",
    "description": "Annual technology conference...",
    "type": "public",
    "location_type": "location",
    "status": "draft",
    "creator_id": 5,
    "start_date": "2025-03-15",
    "start_time": "09:00",
    "end_date": "2025-03-15",
    "end_time": "18:00",
    "address": {
      "id": 1,
      "place_id": "ChIJVTPokywsQkARmtVhquZmAQM",
      "full_address": "Istanbul Convention Center...",
      "latitude": 41.0082,
      "longitude": 28.9784
    },
    "categories": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      }
    ],
    "image": {
      "id": 10,
      "file_url": "https://s3.example.com/event-image.jpg"
    },
    "created_at": "2025-12-11T12:00:00Z"
  }
}
```

#### GET /api/v1/events/manage/:id
Etkinlik detaylarını getirme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Response:**
```json
{
  "success": true,
  "message": "event.retrieved",
  "data": {
    "id": 1,
    "name": "Tech Conference 2025",
    "description": "Annual technology conference...",
    "type": "public",
    "location_type": "location",
    "status": "published",
    "creator": {
      "id": 5,
      "company_name": "Tech Events Corp",
      "user": {
        "full_name": "John Doe",
        "username": "johndoe"
      }
    },
    "address": {
      "full_address": "Istanbul Convention Center...",
      "latitude": 41.0082,
      "longitude": 28.9784
    },
    "categories": [
      {
        "id": 1,
        "name": "Technology",
        "slug": "technology"
      }
    ],
    "tickets": [
      {
        "id": 1,
        "title": "Early Bird",
        "price": 150.00,
        "total_quantity": 100,
        "sold_quantity": 25
      }
    ],
    "stats": {
      "total_tickets": 100,
      "sold_tickets": 25,
      "total_revenue": 3750.00,
      "invitation_count": 50
    }
  }
}
```

#### PUT /api/v1/events/manage/:id
Etkinlik güncelleme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Request Body:** (Tüm alanlar opsiyonel)
```json
{
  "name": "Updated Tech Conference 2025",
  "description": "Updated description...",
  "start_date": "2025-03-20",
  "start_time": "10:00",
  "category_ids": [1, 2, 3],
  "image_id": 12
}
```

#### DELETE /api/v1/events/manage/:id
Etkinlik silme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Davranış:**
- Sadece `draft` durumundaki etkinlikler silinebilir
- Cascade delete: tickets, invitations, categories ilişkileri
- Address kaydı da silinir

#### GET /api/v1/events/manage/my
Kendi etkinliklerimi getirme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)
- `status`: Durum filtresi (`draft`, `pending`, `published`, `rejected`, `cancelled`)
- `type`: Tür filtresi (`public`, `private`)
- `location_type`: Konum türü filtresi (`location`, `online`, `announcement`)

#### GET /api/v1/events/manage/my/drafts
Taslak etkinliklerimi getirme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)

#### GET /api/v1/events/manage/my/published
Yayınlanan etkinliklerimi getirme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)

#### PUT /api/v1/events/manage/:id/status
Etkinlik durumu güncelleme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Request Body:**
```json
{
  "status": "pending"
}
```

**Geçerli Durum Geçişleri:**
- `draft` → `pending` (İncelemeye gönder)
- `pending` → `published` (Yayınla - Admin)
- `pending` → `rejected` (Reddet - Admin)
- `published` → `cancelled` (İptal et)

#### POST /api/v1/events/manage/:id/submit
Etkinliği incelemeye gönderme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Davranış:**
- Etkinlik durumunu `draft`'tan `pending`'e değiştirir
- Etkinlik tamamlanmış olmalı (gerekli alanlar dolu)

#### POST /api/v1/events/manage/:id/publish
Etkinliği yayınlama (Admin).

**Authentication:** Gerekli (Admin - Creator kullanıcısı)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Davranış:**
- Etkinlik durumunu `pending`'den `published`'a değiştirir
- Sadece admin yetkisi olan kullanıcılar yapabilir

#### POST /api/v1/events/manage/:id/cancel
Etkinliği iptal etme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Etkinlik ID'si

**Davranış:**
- Etkinlik durumunu `published`'dan `cancelled`'a değiştirir
- İptal edilen etkinlikler geri alınamaz

#### GET /api/v1/events/manage/stats
Etkinlik istatistiklerimi getirme.

**Authentication:** Gerekli (Creator kullanıcısı)

**Response:**
```json
{
  "success": true,
  "message": "event.stats_retrieved",
  "data": {
    "total_events": 15,
    "draft_events": 3,
    "pending_events": 2,
    "published_events": 8,
    "cancelled_events": 2,
    "total_tickets_sold": 1250,
    "total_revenue": 125000.00,
    "total_invitations_sent": 500,
    "upcoming_events": 5
  }
}
```

#### GET /api/v1/events
Herkese açık etkinlikleri getirme.

**Authentication:** Gerekli değil (Public endpoint)

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)
- `category_id`: Kategori ID'si ile filtreleme
- `location_type`: Konum türü filtresi
- `city`: Şehir filtresi
- `start_date_from`: Başlangıç tarihi filtresi (YYYY-MM-DD)
- `start_date_to`: Bitiş tarihi filtresi (YYYY-MM-DD)

#### GET /api/v1/events/search
Etkinlik arama.

**Authentication:** Gerekli değil (Public endpoint)

**Query Parameters:**
- `q`: Arama sorgusu (minimum 2 karakter)
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)
- `category_id`: Kategori filtresi
- `city`: Şehir filtresi

**Response:**
```json
{
  "success": true,
  "message": "event.search_results",
  "data": {
    "events": [
      {
        "id": 1,
        "name": "Tech Conference 2025",
        "description": "Annual technology conference...",
        "type": "public",
        "location_type": "location",
        "status": "published",
        "start_date": "2025-03-15",
        "start_time": "09:00",
        "address": {
          "city": "Istanbul",
          "full_address": "Istanbul Convention Center..."
        },
        "creator": {
          "company_name": "Tech Events Corp"
        },
        "image": {
          "file_url": "https://s3.example.com/event-image.jpg"
        }
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 20,
    "total_pages": 2
  }
}
```

#### GET /api/v1/events/location/:city
Şehir bazlı etkinlikler.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `city`: Şehir adı

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)

#### GET /api/v1/events/upcoming
Yaklaşan etkinlikler.

**Authentication:** Gerekli değil (Public endpoint)

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)
- `days`: Kaç gün sonrasına kadar (varsayılan: 30)

#### GET /api/v1/events/category/:category_id
Kategori bazlı etkinlikler.

**Authentication:** Gerekli değil (Public endpoint)

**Path Parameters:**
- `category_id`: Kategori ID'si

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)

---

### 13. Ticket Management

#### POST /api/v1/tickets/event/:event_id
Etkinlik için bilet oluşturma.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `event_id`: Etkinlik ID'si

**Request Body:**
```json
{
  "title": "Early Bird Ticket",
  "price": 150.00,
  "total_quantity": 100
}
```

**Zorunlu Alanlar:**
- `title`: Bilet başlığı (2-100 karakter)
- `price`: Bilet fiyatı (0 veya pozitif)
- `total_quantity`: Toplam bilet sayısı (minimum 1)

#### GET /api/v1/tickets/event/:event_id
Etkinlik biletlerini getirme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `event_id`: Etkinlik ID'si

**Response:**
```json
{
  "success": true,
  "message": "ticket.list_retrieved",
  "data": {
    "tickets": [
      {
        "id": 1,
        "title": "Early Bird Ticket",
        "price": 150.00,
        "total_quantity": 100,
        "sold_quantity": 25,
        "remaining_quantity": 75,
        "total_revenue": 3750.00,
        "created_at": "2025-12-11T12:00:00Z"
      }
    ]
  }
}
```

#### PUT /api/v1/tickets/:id
Bilet güncelleme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Bilet ID'si

**Request Body:**
```json
{
  "title": "Updated Early Bird",
  "price": 175.00,
  "total_quantity": 120
}
```

#### DELETE /api/v1/tickets/:id
Bilet silme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Bilet ID'si

**Davranış:**
- Sadece satılmamış biletler silinebilir
- `sold_quantity` > 0 ise silme işlemi reddedilir

---

### 14. Invitation Management

#### POST /api/v1/invitations/event/:event_id
Etkinlik daveti gönderme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `event_id`: Etkinlik ID'si

**Request Body:**
```json
{
  "invited_user_id": 10,
  "invited_email": "guest@example.com"
}
```

**Alanlar:**
- `invited_user_id`: Davet edilecek kullanıcı ID'si (opsiyonel)
- `invited_email`: Davet edilecek email adresi (opsiyonel)

**Not:** `invited_user_id` veya `invited_email`'den en az biri gerekli

#### GET /api/v1/invitations/event/:event_id
Etkinlik davetlerini getirme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `event_id`: Etkinlik ID'si

**Query Parameters:**
- `page`: Sayfa numarası (varsayılan: 1)
- `page_size`: Sayfa başına öğe (varsayılan: 20, max: 100)
- `status`: Durum filtresi (`pending`, `accepted`, `declined`)

**Response:**
```json
{
  "success": true,
  "message": "invitation.list_retrieved",
  "data": {
    "invitations": [
      {
        "id": 1,
        "event_id": 1,
        "invited_user": {
          "id": 10,
          "full_name": "Jane Doe",
          "email": "jane@example.com"
        },
        "invited_email": null,
        "status": "pending",
        "sent_at": "2025-12-11T12:00:00Z",
        "responded_at": null
      }
    ],
    "total": 15,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

#### PUT /api/v1/invitations/:id/respond
Davete yanıt verme.

**Authentication:** Gerekli (Davet edilen kullanıcı)

**Path Parameters:**
- `id`: Davet ID'si

**Request Body:**
```json
{
  "status": "accepted"
}
```

**Geçerli Durumlar:**
- `accepted`: Daveti kabul et
- `declined`: Daveti reddet

#### DELETE /api/v1/invitations/:id
Daveti iptal etme.

**Authentication:** Gerekli (Etkinlik sahibi)

**Path Parameters:**
- `id`: Davet ID'si

**Davranış:**
- Sadece `pending` durumundaki davetler iptal edilebilir
- İptal edilen davet geri alınamaz

---

### Event Management System Özellikleri

- **Multi-type Event Support:** Public/private etkinlikler, location/online/announcement türleri
- **Advanced Status Workflow:** Draft → Pending → Published/Rejected → Cancelled
- **Google Places Integration:** Adres yönetimi ve konum bazlı arama
- **Comprehensive Ticket System:** Bilet oluşturma, satış takibi, gelir yönetimi
- **Flexible Invitation System:** Kullanıcı ve email bazlı davetler
- **Category Integration:** Çok seviyeli kategori sistemi ile etkinlik sınıflandırma
- **Advanced Search & Filtering:** İsim, kategori, konum, tarih bazlı filtreleme
- **Creator Analytics:** Detaylı istatistikler ve performans takibi
- **Permission-based Access Control:** Public/private etkinlik erişim kontrolü
- **Multi-language Support:** Türkçe ve İngilizce i18n desteği

### Event Validation Rules

**Event Creation:**
- Name: 3-200 karakter
- Description: 10-2000 karakter
- Start date: Gelecek tarih olmalı
- End date: Start date'den sonra olmalı (announcement hariç)
- Address: Location type = location ise zorunlu
- Categories: En az 1 kategori seçilmeli

**Status Transitions:**
- Draft → Pending: Tüm zorunlu alanlar dolu olmalı
- Pending → Published: Admin yetkisi gerekli
- Published → Cancelled: Sadece etkinlik sahibi yapabilir

**Ticket Management:**
- Price: 0 veya pozitif değer
- Quantity: Minimum 1
- Deletion: Sadece satılmamış biletler silinebilir

**Invitation System:**
- Private events: Sadece davet edilenler katılabilir
- Email invitations: Sistem dışı kullanıcılar için
- User invitations: Kayıtlı kullanıcılar için

---

## Environment Variables
`.env` dosyasında aşağıdaki değişkenleri ayarlayın:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=louco_event_db
JWT_SECRET=your_jwt_secret
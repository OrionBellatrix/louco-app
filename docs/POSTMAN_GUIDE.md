# Postman Kullanım Kılavuzu - Louco Event API

Bu kılavuz, Louco Event API'sini Postman ile test etmek için gerekli adımları açıklar.

## 📁 Dosyalar

Proje içinde aşağıdaki Postman dosyaları bulunmaktadır:

- `postman/Louco-Event-API.postman_collection.json` - API endpoint'lerinin tamamı
- `postman/Louco-Event-Environment.postman_environment.json` - Environment değişkenleri

## 🚀 Kurulum

### 1. Collection'ı İçe Aktarma

1. Postman'i açın
2. Sol üst köşedeki **Import** butonuna tıklayın
3. **Upload Files** sekmesini seçin
4. `postman/Louco-Event-API.postman_collection.json` dosyasını sürükleyin veya seçin
5. **Import** butonuna tıklayın

### 2. Environment'ı İçe Aktarma

1. Postman'de **Import** butonuna tekrar tıklayın
2. `postman/Louco-Event-Environment.postman_environment.json` dosyasını seçin
3. **Import** butonuna tıklayın
4. Sağ üst köşeden **Louco Event - Development** environment'ını seçin

## ⚙️ Environment Değişkenleri

Environment dosyasında aşağıdaki değişkenler tanımlıdır:

| Değişken | Açıklama | Örnek Değer |
|----------|----------|-------------|
| `base_url` | API'nin temel URL'i | `http://localhost:8080` |
| `jwt_token` | Login'den alınan JWT token | (otomatik doldurulur) |
| `admin_jwt_token` | Admin JWT token | (manuel girilir) |
| `user_id` | Test kullanıcısının ID'si | (otomatik doldurulur) |
| `media_id` | Test medya dosyasının ID'si | (otomatik doldurulur) |
| `username` | Test kullanıcı adı | `testuser123` |
| `test_email` | Test email adresi | `test@example.com` |
| `test_phone` | Test telefon numarası | `+905551234567` |
| `test_password` | Test şifresi | `TestPass123!` |

## 📋 Collection Yapısı

Collection aşağıdaki klasörlere ayrılmıştır:

### 1. Health Check
- **Health Check** - API durumunu kontrol eder

### 2. Authentication
- **Register Step 1** - Hesap oluşturma
- **Login** - Kullanıcı girişi
- **Social Login** - Sosyal medya girişi
- **Forgot Password** - Şifre sıfırlama talebi
- **Reset Password** - Şifre sıfırlama
- **Change Password** - Şifre değiştirme

### 3. Username Management
- **Check Username Availability** - Kullanıcı adı müsaitlik kontrolü
- **Set Username** - Kullanıcı adı belirleme

### 4. User Management
- **Get User Profile** - Profil bilgilerini getirme
- **Update User Profile** - Profil güncelleme
- **Update Contact Information** - İletişim bilgileri güncelleme
- **Set Profile Picture** - Profil resmi ayarlama
- **Register Step 4** - Profil detaylarını tamamlama
- **Deactivate Account** - Hesap deaktivasyonu

### 5. Media Management
- **Upload File** - Dosya yükleme
- **Get Media by ID** - Medya detaylarını getirme
- **Get User Media** - Kullanıcının medya dosyaları
- **Update Media** - Medya güncelleme
- **Delete Media** - Medya silme

### 6. Admin Operations
- **Get All Users (Admin)** - Tüm kullanıcıları listeleme
- **Get All Media (Admin)** - Tüm medya dosyalarını listeleme

## 🔄 Test Senaryosu

### Adım 1: API Durumunu Kontrol Etme
1. **Health Check** → **Health Check** endpoint'ini çalıştırın
2. Response'da `"status": "healthy"` olduğunu kontrol edin

### Adım 2: Kullanıcı Kaydı
1. **Authentication** → **Register Step 1** endpoint'ini çalıştırın
2. Response'dan `user_id` ve `token` değerlerini not alın
3. `jwt_token` environment değişkeni otomatik olarak güncellenecek

### Adım 3: Kullanıcı Adı Belirleme
1. **Username Management** → **Check Username Availability** ile kullanıcı adını kontrol edin
2. **Username Management** → **Set Username** ile kullanıcı adını belirleyin

### Adım 4: Profil Tamamlama
1. **User Management** → **Register Step 4** ile profil detaylarını tamamlayın

### Adım 5: Login Test
1. **Authentication** → **Login** endpoint'ini test edin
2. Yeni token ile environment güncellenecek

### Adım 6: Dosya Yükleme
1. **Media Management** → **Upload File** endpoint'ini kullanın
2. Bir resim veya video dosyası seçin
3. Response'dan `media_id` değerini not alın

### Adım 7: Profil Resmi Ayarlama
1. **Set Profile Picture** endpoint'ini kullanın
2. Yüklediğiniz medya dosyasının ID'sini girin
3. Profil resmi başarıyla güncellendiğini kontrol edin

### Adım 8: Medya İşlemleri
1. **Get Media by ID** ile yüklenen dosyayı görüntüleyin
2. **Update Media** ile metadata güncelleyin
3. **Get User Media** ile kullanıcının tüm dosyalarını listeleyin

## 🔧 Pre-request Scripts

Bazı endpoint'lerde otomatik token yönetimi için pre-request script'ler eklenmiştir:

```javascript
// Login endpoint'inde token'ı environment'a kaydetme
pm.test("Save JWT Token", function () {
    var jsonData = pm.response.json();
    if (jsonData.success && jsonData.data.token) {
        pm.environment.set("jwt_token", jsonData.data.token);
        pm.environment.set("user_id", jsonData.data.user.id);
    }
});
```

## 🧪 Test Scripts

Her endpoint için otomatik test script'leri eklenmiştir:

```javascript
// Temel response kontrolü
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has success field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('success');
    pm.expect(jsonData.success).to.be.true;
});

pm.test("Response has message field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('message');
});
```

## 🌐 Dil Desteği

API çoklu dil desteği sunar. Test etmek için:

1. Request header'larında `Accept-Language` değerini değiştirin:
   - `tr` - Türkçe
   - `en` - İngilizce

2. Response mesajlarının seçilen dile göre geldiğini kontrol edin

## 🔐 Authentication

### JWT Token Yönetimi

1. **Login** veya **Register Step 1** endpoint'lerini çalıştırdıktan sonra token otomatik olarak environment'a kaydedilir
2. Diğer endpoint'ler bu token'ı otomatik olarak kullanır
3. Token süresi dolduğunda yeniden login yapmanız gerekir

### Admin İşlemleri

Admin endpoint'lerini test etmek için:

1. Admin yetkisine sahip bir kullanıcı ile login olun
2. Alınan token'ı `admin_jwt_token` environment değişkenine manuel olarak kaydedin
3. Admin endpoint'lerini çalıştırın

## 📊 Response Örnekleri

### Başarılı Response
```json
{
  "success": true,
  "message": "auth.login_success",
  "data": {
    "user": {
      "id": 1,
      "full_name": "John Doe",
      "email": "user@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "errors": null
}
```

### Hata Response
```json
{
  "success": false,
  "message": "common.validation_failed",
  "data": null,
  "errors": [
    {
      "field": "email",
      "message": "Email format is invalid",
      "value": "invalid-email"
    }
  ]
}
```

## 🐛 Hata Ayıklama

### Yaygın Sorunlar

1. **401 Unauthorized**
   - JWT token'ın süresi dolmuş olabilir
   - Yeniden login yapın

2. **422 Validation Error**
   - Request body'deki alanları kontrol edin
   - Zorunlu alanların eksik olmadığından emin olun

3. **500 Internal Server Error**
   - API sunucusunun çalıştığından emin olun
   - Database bağlantısını kontrol edin

### Debug İpuçları

1. **Console** sekmesinde hata mesajlarını kontrol edin
2. **Headers** sekmesinde request/response header'larını inceleyin
3. **Body** sekmesinde gönderilen/alınan veriyi kontrol edin

## 📝 Notlar

- Tüm endpoint'ler `Content-Type: application/json` header'ı gerektirir (dosya yükleme hariç)
- Dosya yükleme endpoint'i `multipart/form-data` kullanır
- Rate limiting uygulandığı için çok hızlı istek göndermeyin
- Test verilerini gerçek production verilerinden ayırın

## 🔄 Collection Güncelleme

API'de değişiklik olduğunda:

1. Yeni collection dosyasını indirin
2. Postman'de mevcut collection'ı silin
3. Yeni collection'ı import edin
4. Environment değişkenlerini yeniden ayarlayın

## 📞 Destek

Sorun yaşadığınızda:

1. API dokümantasyonunu kontrol edin (`docs/API_DOCUMENTATION.md`)
2. Server loglarını inceleyin
3. Postman console'unda hata mesajlarını kontrol edin
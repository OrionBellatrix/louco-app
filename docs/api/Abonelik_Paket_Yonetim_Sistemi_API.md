# Abonelik & Paket Yönetim Sistemi API Dokümantasyonu

## Genel Bakış

Abonelik & Paket Yönetim Sistemi, abonelik planları ve tek seferlik paketler aracılığıyla kapsamlı etkinlik yayınlama hakları yönetimi sağlar. Ödeme işlemleri için Stripe entegrasyonu ve otomatik kullanım takibi ve doğrulama içerir.

## İçindekiler

1. [Sistem Mimarisi](#sistem-mimarisi)
2. [Kimlik Doğrulama](#kimlik-doğrulama)
3. [Abonelik Planları](#abonelik-planları)
4. [API Endpointleri](#api-endpointleri)
5. [Etkinlik Yayınlama Entegrasyonu](#etkinlik-yayınlama-entegrasyonu)
6. [Stripe Entegrasyonu](#stripe-entegrasyonu)
7. [Hata Yönetimi](#hata-yönetimi)
8. [Kullanım Örnekleri](#kullanım-örnekleri)
9. [Test Senaryoları](#test-senaryoları)

## Sistem Mimarisi

### Abonelik Türleri

#### 1. Yinelenen Abonelikler
- **Temel Plan**: 5 etkinlik/hafta veya 20 etkinlik/ay
- **Plus Plan**: 15 etkinlik/hafta veya 60 etkinlik/ay  
- **Pro Plan**: Sınırsız etkinlik
- **Faturalama Döngüleri**: Haftalık (Pazartesi sıfırlama) veya Aylık (1. gün sıfırlama)

#### 2. Tek Seferlik Paketler
- **1 Etkinlik Paketi**: Tek etkinlik kredisi
- **10 Etkinlik Paketi**: 10 etkinlik kredisi
- **25 Etkinlik Paketi**: 25 etkinlik kredisi
- **Süre Sınırı Yok**: Krediler kullanılana kadar hiç bitmez

### Yayınlama Hakları Doğrulaması

Sistem, etkinlik gönderimlerine izin vermeden önce yayınlama haklarını doğrular:

1. **Aktif Abonelik/Paket Kontrolü**: Kullanıcının en az bir aktif aboneliği veya paketi olmalı
2. **Limit/Kredi Doğrulaması**: Yeterli kalan limit veya kredi olduğundan emin olun
3. **Kullanım Tüketimi**: Başarılı gönderim sonrası otomatik olarak limit/kredilerden düş
4. **Hata Yönetimi**: Yetersiz haklar için kullanıcı dostu hata mesajları döndür

## Kimlik Doğrulama

Tüm endpointler (webhook'lar hariç) Bearer token kimlik doğrulaması gerektirir:

```
Authorization: Bearer {access_token}
```

Dil tercihi header ile ayarlanabilir:

```
Accept-Language: en|tr
```

## Abonelik Planları

### Plan Yapısı

```json
{
  "id": 1,
  "name": "Temel",
  "description": "Küçük etkinlikler için mükemmel",
  "plan_type": "subscription",
  "billing_cycle": "monthly",
  "price": 29.99,
  "currency": "USD",
  "weekly_limit": 5,
  "monthly_limit": 20,
  "credits": null,
  "is_active": true,
  "stripe_price_id": "price_basic_monthly"
}
```

### Mevcut Planlar

| ID | İsim | Tür | Haftalık Limit | Aylık Limit | Kredi | Fiyat |
|----|------|-----|----------------|-------------|-------|-------|
| 1 | Temel | abonelik | 5 | 20 | - | $29.99/ay |
| 2 | Plus | abonelik | 15 | 60 | - | $79.99/ay |
| 3 | Pro | abonelik | sınırsız | sınırsız | - | $199.99/ay |
| 4 | 1 Etkinlik | paket | - | - | 1 | $9.99 |
| 5 | 10 Etkinlik | paket | - | - | 10 | $79.99 |
| 6 | 25 Etkinlik | paket | - | - | 25 | $179.99 |

## API Endpointleri

### Abonelik Planları

#### GET /api/v1/subscriptions/plans
Tüm mevcut abonelik planlarını getir.

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.plan.list.success",
  "data": [
    {
      "id": 1,
      "name": "Temel",
      "description": "Küçük etkinlikler için mükemmel",
      "plan_type": "subscription",
      "billing_cycle": "monthly",
      "price": 29.99,
      "currency": "USD",
      "weekly_limit": 5,
      "monthly_limit": 20,
      "credits": null,
      "is_active": true
    }
  ]
}
```

#### GET /api/v1/subscriptions/plans/{id}
Belirli abonelik planı detaylarını getir.

**Parametreler:**
- `id` (path): Plan ID

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.plan.get.success",
  "data": {
    "id": 1,
    "name": "Temel",
    "description": "Küçük etkinlikler için mükemmel",
    "plan_type": "subscription",
    "billing_cycle": "monthly",
    "price": 29.99,
    "currency": "USD",
    "weekly_limit": 5,
    "monthly_limit": 20,
    "credits": null,
    "is_active": true
  }
}
```

### Kullanıcı Abonelikleri

#### GET /api/v1/subscriptions/my
Mevcut kullanıcının aktif abonelik ve paketlerini getir.

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.user.list.success",
  "data": [
    {
      "id": 1,
      "user_id": 123,
      "plan_id": 1,
      "plan_name": "Temel",
      "plan_type": "subscription",
      "billing_cycle": "monthly",
      "status": "active",
      "current_period_start": "2024-01-01T00:00:00Z",
      "current_period_end": "2024-02-01T00:00:00Z",
      "weekly_limit": 5,
      "monthly_limit": 20,
      "weekly_usage": 2,
      "monthly_usage": 8,
      "credits": null,
      "credits_used": null,
      "stripe_subscription_id": "sub_stripe_id",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### GET /api/v1/subscriptions/publishing-rights
Mevcut kullanıcının etkinlik yayınlama haklarını kontrol et.

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.rights.success",
  "data": {
    "can_publish": true,
    "active_subscriptions": 1,
    "active_packages": 0,
    "weekly_remaining": 3,
    "monthly_remaining": 12,
    "total_credits": 0,
    "used_credits": 0,
    "remaining_credits": 0,
    "next_reset": "2024-02-01T00:00:00Z"
  }
}
```

#### GET /api/v1/subscriptions/usage-stats
Detaylı kullanım istatistiklerini getir.

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.stats.success",
  "data": {
    "current_period": {
      "start": "2024-01-01T00:00:00Z",
      "end": "2024-02-01T00:00:00Z",
      "weekly_usage": 2,
      "monthly_usage": 8,
      "events_published": 8
    },
    "subscriptions": [
      {
        "id": 1,
        "plan_name": "Temel",
        "status": "active",
        "weekly_limit": 5,
        "monthly_limit": 20,
        "weekly_usage": 2,
        "monthly_usage": 8
      }
    ],
    "packages": [],
    "total_events_published": 8,
    "total_credits_purchased": 0,
    "total_credits_used": 0
  }
}
```

### Abonelik Satın Alma

#### POST /api/v1/subscriptions/purchase
Abonelik planı satın al.

**İstek Gövdesi:**
```json
{
  "plan_id": 1,
  "billing_cycle": "monthly",
  "payment_method_id": "pm_card_visa"
}
```

**Parametreler:**
- `plan_id`: Abonelik planı ID (1-3)
- `billing_cycle`: "weekly" veya "monthly"
- `payment_method_id`: Stripe ödeme yöntemi ID

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.purchase.success",
  "data": {
    "id": 1,
    "user_id": 123,
    "plan_id": 1,
    "status": "active",
    "stripe_subscription_id": "sub_stripe_id",
    "stripe_customer_id": "cus_stripe_id",
    "current_period_start": "2024-01-01T00:00:00Z",
    "current_period_end": "2024-02-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### POST /api/v1/subscriptions/purchase-package
Tek seferlik etkinlik paketi satın al.

**İstek Gövdesi:**
```json
{
  "plan_id": 5,
  "payment_method_id": "pm_card_visa"
}
```

**Parametreler:**
- `plan_id`: Paket planı ID (4-6)
- `payment_method_id`: Stripe ödeme yöntemi ID

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.purchase.success",
  "data": {
    "id": 2,
    "user_id": 123,
    "plan_id": 5,
    "status": "active",
    "credits": 10,
    "credits_used": 0,
    "stripe_payment_intent_id": "pi_stripe_id",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Abonelik Yönetimi

#### POST /api/v1/subscriptions/cancel
Aktif aboneliği iptal et.

**İstek Gövdesi:**
```json
{
  "subscription_id": 1,
  "reason": "Artık gerekli değil"
}
```

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.user.cancel.success",
  "data": {
    "id": 1,
    "status": "cancelled",
    "cancelled_at": "2024-01-15T00:00:00Z",
    "cancellation_reason": "Artık gerekli değil"
  }
}
```

#### PUT /api/v1/subscriptions/update
Mevcut aboneliği güncelle (yükseltme/düşürme).

**İstek Gövdesi:**
```json
{
  "subscription_id": 1,
  "new_plan_id": 2,
  "billing_cycle": "monthly"
}
```

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.user.update.success",
  "data": {
    "id": 1,
    "plan_id": 2,
    "plan_name": "Plus",
    "status": "active",
    "updated_at": "2024-01-15T00:00:00Z"
  }
}
```

## Etkinlik Yayınlama Entegrasyonu

### Yayınlama Hakları Doğrulaması

Sistem, etkinlikler inceleme için gönderildiğinde otomatik olarak yayınlama haklarını doğrular:

#### POST /api/v1/events/{id}/submit
Etkinliği inceleme için gönder (hak doğrulaması ile).

**Gelişmiş Davranış:**
1. **Ön Doğrulama**: Kullanıcının aktif abonelik/paketi olup olmadığını kontrol et
2. **Limit Kontrolü**: Kalan haftalık/aylık limit veya kredileri doğrula
3. **Durum Güncelleme**: Etkinlik durumunu "draft"tan "pending"e değiştir
4. **Kullanım Tüketimi**: Otomatik olarak limit/kredilerden düş
5. **Hata Yönetimi**: Yetersiz haklar için detaylı hata mesajları döndür

**Başarı Yanıtı:**
```json
{
  "success": true,
  "message": "event.submit.success",
  "data": {
    "id": 1,
    "status": "pending",
    "submitted_at": "2024-01-15T00:00:00Z",
    "usage_consumed": {
      "type": "monthly_limit",
      "amount": 1,
      "remaining": 19
    }
  }
}
```

**Hata Yanıtı (Yetersiz Haklar):**
```json
{
  "success": false,
  "message": "subscription.insufficient_publishing_rights",
  "errors": [
    {
      "field": "subscription",
      "message": "Yeterli yayınlama hakkınız bulunmamaktadır. Etkinlik yayınlamak için lütfen bir abonelik veya paket satın alın."
    }
  ]
}
```

#### GET /api/v1/subscriptions/can-publish
Etkinlik göndermeden yayınlama haklarını test et.

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.rights.success",
  "data": {
    "can_publish": true,
    "reason": "Kalan limitli aktif abonelik",
    "details": {
      "active_subscription": true,
      "weekly_remaining": 3,
      "monthly_remaining": 12,
      "credits_remaining": 0
    }
  }
}
```

## Stripe Entegrasyonu

### Webhook Yönetimi

#### POST /api/v1/webhooks/stripe
Stripe webhook olaylarını işle.

**Desteklenen Olaylar:**
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`
- `invoice.payment_succeeded`
- `invoice.payment_failed`
- `payment_intent.succeeded`
- `payment_intent.payment_failed`

**Header'lar:**
```
Content-Type: application/json
Stripe-Signature: {webhook_signature}
```

**Yanıt:**
```json
{
  "success": true,
  "message": "subscription.webhook.success",
  "data": {
    "event_id": "evt_stripe_id",
    "event_type": "customer.subscription.created",
    "processed_at": "2024-01-15T00:00:00Z"
  }
}
```

### Ödeme İşlemi Akışı

1. **Frontend**: Stripe.js kullanarak ödeme yöntemi topla
2. **Backend**: Stripe API aracılığıyla abonelik/ödeme intent oluştur
3. **Webhook**: Ödeme onayını işle
4. **Veritabanı**: Abonelik durumunu güncelle
5. **Kullanıcı**: Onay al ve yayınlama haklarına erişim sağla

## Hata Yönetimi

### Yaygın Hata Kodları

| Kod | Mesaj Anahtarı | Açıklama |
|-----|----------------|----------|
| 400 | `subscription.purchase.invalid_plan` | Geçersiz abonelik planı ID |
| 400 | `subscription.validation.invalid_billing_cycle` | Geçersiz faturalama döngüsü |
| 402 | `subscription.insufficient_publishing_rights` | Yetersiz yayınlama hakkı |
| 404 | `subscription.plan.not_found` | Abonelik planı bulunamadı |
| 409 | `subscription.purchase.already_active` | Kullanıcının zaten aktif aboneliği var |
| 409 | `subscription.rights.limit_exceeded` | Yayınlama limiti aşıldı |
| 409 | `subscription.rights.credits_exhausted` | Paket kredileri tükendi |
| 500 | `subscription.purchase.payment_failed` | Ödeme işlemi başarısız |

### Hata Yanıt Formatı

```json
{
  "success": false,
  "message": "subscription.insufficient_publishing_rights",
  "errors": [
    {
      "field": "subscription",
      "message": "Yeterli yayınlama hakkınız bulunmamaktadır. Etkinlik yayınlamak için lütfen bir abonelik veya paket satın alın."
    }
  ],
  "data": null
}
```

## Kullanım Örnekleri

### Örnek 1: Temel Abonelik Satın Alma

```bash
# 1. Mevcut planları getir
curl -X GET "http://localhost:8080/api/v1/subscriptions/plans" \
  -H "Accept-Language: tr"

# 2. Temel aylık abonelik satın al
curl -X POST "http://localhost:8080/api/v1/subscriptions/purchase" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: tr" \
  -d '{
    "plan_id": 1,
    "billing_cycle": "monthly",
    "payment_method_id": "pm_card_visa"
  }'

# 3. Yayınlama haklarını kontrol et
curl -X GET "http://localhost:8080/api/v1/subscriptions/publishing-rights" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: tr"
```

### Örnek 2: Etkinlik Paketi Satın Alma

```bash
# 1. 10 etkinlik paketi satın al
curl -X POST "http://localhost:8080/api/v1/subscriptions/purchase-package" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -H "Accept-Language: tr" \
  -d '{
    "plan_id": 5,
    "payment_method_id": "pm_card_visa"
  }'

# 2. Kullanım istatistiklerini kontrol et
curl -X GET "http://localhost:8080/api/v1/subscriptions/usage-stats" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: tr"
```

### Örnek 3: Hak Doğrulaması ile Etkinlik Gönderimi

```bash
# 1. Abonelik olmadan göndermeyi dene (başarısız olmalı)
curl -X POST "http://localhost:8080/api/v1/events/1/submit" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: tr"

# Beklenen yanıt: 402 Yetersiz yayınlama hakkı

# 2. Paket satın al
curl -X POST "http://localhost:8080/api/v1/subscriptions/purchase-package" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_id": 4,
    "payment_method_id": "pm_card_visa"
  }'

# 3. Etkinlik gönder (başarılı olmalı)
curl -X POST "http://localhost:8080/api/v1/events/1/submit" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: tr"

# 4. Kalan kredileri kontrol et
curl -X GET "http://localhost:8080/api/v1/subscriptions/usage-stats" \
  -H "Authorization: Bearer {token}" \
  -H "Accept-Language: tr"
```

## Test Senaryoları

### Senaryo 1: Abonelik İş Akışı
1. **Creator kullanıcı olarak kayıt ol**
2. **Taslak etkinlik oluştur**
3. **Abonelik olmadan göndermeyi dene** (hata bekle)
4. **Temel aylık abonelik satın al**
5. **Etkinliği başarıyla gönder**
6. **Kullanım istatistiklerini kontrol et**
7. **Limite ulaşana kadar daha fazla etkinlik gönder**
8. **Limiti aşmaya çalış** (hata bekle)

### Senaryo 2: Paket İş Akışı
1. **1 Etkinlik paketi satın al**
2. **Etkinliği başarıyla gönder** (1 kredi tüket)
3. **Başka etkinlik göndermeyi dene** (hata bekle - kredi yok)
4. **10 Etkinlik paketi satın al**
5. **Birden fazla etkinlik gönder**
6. **Her gönderim sonrası kalan kredileri kontrol et**

### Senaryo 3: Karma Kullanım
1. **Temel haftalık abonelik satın al**
2. **Haftalık limiti kullan**
3. **Ek etkinlikler için 10 Etkinlik paketi satın al**
4. **Paket kredilerini kullanarak etkinlik gönder**
5. **Haftalık sıfırlamayı bekle**
6. **Yenilenen haftalık limiti kullanarak etkinlik gönder**

### Senaryo 4: Abonelik Yönetimi
1. **Temel abonelik satın al**
2. **Plus aboneliğe yükselt**
3. **Güncellenmiş limitleri kontrol et**
4. **Aboneliği iptal et**
5. **İptal sonrası göndermeyi dene** (hata bekle)

## Hız Sınırlaması

- **Abonelik endpointleri**: Kullanıcı başına dakikada 100 istek
- **Webhook endpointleri**: Dakikada 1000 istek (kullanıcı sınırı yok)
- **Yayınlama hakları kontrolü**: Kullanıcı başına dakikada 200 istek

## Güvenlik Hususları

1. **Webhook İmza Doğrulaması**: Tüm Stripe webhook'ları imza kullanılarak doğrulanır
2. **Ödeme Yöntemi Güvenliği**: Ödeme yöntemleri tamamen Stripe tarafından işlenir
3. **Kullanıcı Yetkilendirmesi**: Tüm abonelik işlemleri geçerli kullanıcı kimlik doğrulaması gerektirir
4. **Kullanım Takibi**: Tüm kullanım tüketimi denetim amaçlı loglanır
5. **İdempotency**: Ödeme işlemleri duplikasyonu önlemek için idempotency anahtarları içerir

## İzleme ve Loglama

Sistem şunlar için kapsamlı loglama sağlar:
- **Abonelik satın almaları** ve iptalleri
- **Kullanım tüketimi** ve limit sıfırlamaları
- **Yayınlama hakları doğrulama** denemeleri
- **Ödeme işlemi** olayları
- **Webhook işleme** durumu
- **Hata koşulları** ve başarısızlıklar

Tüm loglar kolay izleme ve uyarı için yapılandırılmış veri içerir.
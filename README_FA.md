# SchemaGit
[English](README.md) | [فارسی](README_FA.md)

SchemaGit یک سیستم کامل برای **نسخه‌بندی اسکیمای دیتابیس**، **Diff**، **Plan**، **Pipeline مهاجرت**، **CLI کامل**، **توکن‌های API**، **وبهوک‌ها**، **نوتیفیکیشن‌ها**، **Audit Logs** و یک **داشبورد Cloud** است.

نسخه **۱.۰** کاملاً تکمیل شده و آماده انتشار است.

---

## 🚀 قابلیت‌ها

- نسخه‌بندی اسکیمای دیتابیس  
- Diff و Plan خودکار  
- پایپلاین مهاجرت ابری  
- ابزار CLI کامل  
- API Tokens  
- Webhooks  
- Notifications (ایمیل + اسلک)  
- Audit Logs  
- Production Hardening (Rate limit، CORS، Security headers)  
- سازمان‌ها / تیم‌ها / پروژه‌ها  
- داشبورد Cloud

---

## 📥 نصب CLI

```bash
go install github.com/sepanta/schemagit/cli/schemagit@latest
```

---

## 🔐 ورود به Cloud

```bash
schemagit login
```

---

## 📁 ساخت پروژه

```bash
schemagit project create mydb
```

---

## 🔄 Diff

```bash
schemagit diff mydb
```

---

## 📜 Plan

```bash
schemagit plan mydb
```

---

## 🚚 Migration

```bash
schemagit migrate mydb
```

---

## ⬆️ Push اسکیمای دیتابیس

```bash
schemagit push mydb schema.sql
```

---

## ⬇️ Pull اسکیمای دیتابیس

```bash
schemagit pull mydb schema.sql
```

---

## 🔐 API Tokens

توکن را در Cloud بسازید و سپس:

```bash
export SCHEMAGIT_TOKEN=xxxx
```

---

## 🔔 نوتیفیکیشن‌ها

تنظیمات ایمیل و اسلک:

```
/notifications
```

---

## 🔗 وبهوک‌ها

وبهوک برای رویدادهای مهاجرت:

```
/webhooks
```

---

## 📊 Audit Logs

تمام رویدادهای مهم ثبت می‌شوند:

- ورود  
- CLI  
- API  
- مهاجرت‌ها  
- وبهوک‌ها  
- نوتیفیکیشن‌ها  
- ساخت توکن  

مشاهده در:

```
/audit
```

---

## 🛡 Production Hardening

نسخه ۱.۰ شامل:

- Rate limiting  
- CORS  
- Security headers  
- Panic recovery  
- Structured logging  

---

## 🧱 ساختار پروژه

```
internal/cloud/
    api/
    store/
    pipeline/
    server.go
cli/
    cmd/
ui/
    cloud/
```

---

## 🏁 استقرار

### Docker

```bash
docker build -t schemagit .
docker run -p 8080:8080 schemagit
```

### اجرا روی سرور

```bash
go build -o schemagit-cloud
./schemagit-cloud
```

---

## 📄 License

MIT

# 📦 SchemaGit

[English](README.md) | [فارسی](README_FA.md)

SchemaGit is a complete system for **database schema versioning**, **diff**, **plan**, **migration pipeline**, **CLI**, **API tokens**, **webhooks**, **notifications**, **audit logs**, and a full **Cloud dashboard**.

Version **1.0** is now complete and production-ready.

---

## 🚀 Features

- Schema versioning
- Automatic diff & plan
- Cloud migration pipeline
- Full CLI tool
- API Tokens
- Webhooks
- Notifications (Email + Slack)
- Audit Logs
- Production Hardening (Rate limit, CORS, Security headers)
- Organizations / Teams / Projects
- Cloud Dashboard UI

---

## 📥 Install CLI

```bash
go install github.com/sepanta/schemagit/cli/schemagit@latest
```

---

## 🔐 Login to Cloud

```bash
schemagit login
```

---

## 📁 Create Project

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

## 🚚 Migrate

```bash
schemagit migrate mydb
```

---

## ⬆️ Push Schema

```bash
schemagit push mydb schema.sql
```

---

## ⬇️ Pull Schema

```bash
schemagit pull mydb schema.sql
```

---

## 🔐 API Tokens

Create a token in the Cloud dashboard, then:

```bash
export SCHEMAGIT_TOKEN=xxxx
```

---

## 🔔 Notifications

Configure Email / Slack notifications in:

```
/notifications
```

---

## 🔗 Webhooks

Configure webhooks for migration events:

```
/webhooks
```

---

## 📊 Audit Logs

All important events are recorded:

- Login
- CLI actions
- API usage
- Migrations
- Webhooks
- Notifications
- Token creation

View logs at:

```
/audit
```

---

## 🛡 Production Hardening

Version 1.0 includes:

- Rate limiting
- CORS
- Security headers
- Panic recovery
- Structured logging

---

## 🧱 Project Structure

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

## 🏁 Deployment

### Docker

```bash
docker build -t schemagit .
docker run -p 8080:8080 schemagit
```

### Manual

```bash
go build -o schemagit-cloud
./schemagit-cloud
```

---

## 📄 License

MIT

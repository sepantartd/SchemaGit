FROM golang:1.22-alpine AS builder

WORKDIR /app

# کپی سورس پروژه
COPY . .

# ساخت باینری SchemaGit
RUN go build -o schemagit ./cmd/schemagit

# مرحله دوم: ایمیج سبک برای اجرا
FROM alpine:3.19

WORKDIR /app

# کپی باینری
COPY --from=builder /app/schemagit /usr/local/bin/schemagit

# کپی فایل‌های لازم
COPY schema.sql /app/schema.sql
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

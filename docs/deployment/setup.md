# راهنمای راه‌اندازی محیط توسعه (Onboarding Guide)

این مستند به شما کمک می‌کند تا محیط توسعه کامل پروژه "معماری روایت" را روی سیستم خود راه‌اندازی کنید.

## 1. پیش‌نیازها

قبل از شروع، مطمئن شوید ابزارهای زیر روی سیستم شما نصب و در `PATH` موجود است:

-   **Git:** برای مدیریت سورس‌کد.
-   **Node.js:** نسخه `18.17.0` یا بالاتر.
-   **pnpm:** (`npm install -g pnpm`) - مدیر پکیج Monorepo ما.
-   **Docker & Docker Compose:** برای اجرای سرویس‌های جانبی (دیتابیس‌ها).
-   **Go:** نسخه `1.21` یا بالاتر.
-   **`golang-migrate`:** ابزار مدیریت مایگریشن‌های PostgreSQL.
    ```bash
    # نصب با Go
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    # یا با Homebrew (macOS)
    brew install golang-migrate
    ```
-   **cURL** یا **Postman:** برای تست APIها.

## 2. راه‌اندازی اولیه

این مراحل فقط یک‌بار در ابتدا انجام می‌شود.

### قدم 1: دریافت کد و نصب وابستگی‌ها

```bash
# ریپازیتوری را از GitHub کلون کنید
git clone https://github.com/your-org/narrative-architecture-platform.git
cd narrative-architecture-platform

# تمام وابستگی‌های Monorepo را با pnpm نصب کنید
pnpm install
```

### قدم 2: تنظیم متغیرهای محیطی

برای هر سرویس، یک کپی از فایل `.env.example` ایجاد کنید و مقادیر لازم را (در صورت نیاز) تغییر دهید. برای توسعه محلی، مقادیر پیش‌فرض معمولاً کافی هستند.

```bash
cp apps/platform/.env.example apps/platform/.env.local
cp apps/backend/.env.example apps/backend/.env
cp apps/ai-service/.env.example apps/ai-service/.env
cp apps/workers/search-sync/.env.example apps/workers/search-sync/.env
```
**مهم:** کلیدهای API و توکن‌های محرمانه را از مدیر پروژه دریافت کرده و در این فایل‌ها قرار دهید. هرگز این فایل‌ها را در Git کامیت نکنید.

## 3. اجرای زیرساخت با Docker

این دستور تمام سرویس‌های جانبی (PostgreSQL, MongoDB, Redis, Elasticsearch) را در کانتینرهای Docker اجرا می‌کند.

```bash
# کانتینرها را در پس‌زمینه اجرا کنید
docker-compose up -d

# برای مشاهده وضعیت سرویس‌ها
docker-compose ps
```
✅ **نقطه اطمینان:** تمام سرویس‌ها باید وضعیت `running` یا `healthy` داشته باشند.

## 4. آماده‌سازی دیتابیس‌ها

حالا باید شِماها و داده‌های اولیه را روی دیتابیس‌های در حال اجرا اعمال کنیم.

```bash
# 1. اجرای مایگریشن‌های PostgreSQL
# مطمئن شوید در ریشه پروژه هستید
# این دستور از DATABASE_URL در apps/backend/.env استفاده می‌کند
pnpm run db:migrate:up

# 2. اجرای اسکریپت MongoDB
# این دستور از MONGO_URI و MONGODB_DB در .env استفاده می‌کند
pnpm run db:mongo:init

# 3. اعمال تمپلیت‌ها و Aliasهای Elasticsearch
# این اسکریپت از ELASTIC_URL در .env استفاده می‌کند
pnpm run db:es:init
```
✅ **نقطه اطمینان:** اگر تمام این دستورات بدون خطا اجرا شدند، دیتابیس‌های شما آماده استفاده هستند.

**نکته:** ما برای راحتی، این دستورات را در `package.json` ریشه پروژه تعریف کرده‌ایم.

## 5. اجرای کامل پروژه

حالا می‌توانید تمام اپلیکیشن‌ها (Frontend, Backend, Worker) را با یک دستور اجرا کنید.

```bash
# این دستور تمام اسکریپت‌های 'dev' را به‌صورت موازی اجرا می‌کند
pnpm run dev
```

پس از اجرای موفقیت‌آمیز، سرویس‌ها در آدرس‌های زیر در دسترس خواهند بود:

-   **Frontend (Next.js):** `http://localhost:3000`
-   **Backend (Go):** `http://localhost:8080`
-   **AI Service (Python):** `http://localhost:8000`
-   **Search Sync Worker:** در پس‌زمینه در حال اجراست و لاگ‌های آن در ترمینال نمایش داده می‌شود.

**تبریک! محیط توسعه شما کاملاً آماده است.**

## عیب‌یابی (Troubleshooting)

-   **خطای `Port Conflict`:** مطمئن شوید پورت‌های `3000`, `8080`, `8000`, `5432`, `27017`, `6379`, `9200` آزاد هستند.
-   **خطای `Connection Refused`:** ابتدا `docker-compose ps` را اجرا کرده و از `healthy` بودن تمام سرویس‌های Docker مطمئن شوید.
-   **خطای مایگریشن:** مطمئن شوید `golang-migrate` نصب است و `DATABASE_URL` در `apps/backend/.env` صحیح است.
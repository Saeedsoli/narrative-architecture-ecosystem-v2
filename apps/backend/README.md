# Backend Service

این سرویس، هسته اصلی پلتفرم "معماری روایت" است که با زبان Go و معماری Clean Architecture پیاده‌سازی شده است.

## راه‌اندازی

1.  مطمئن شوید که در ریشه پروژه، `docker-compose up -d` را اجرا کرده‌اید.
2.  متغیرهای محیطی را در فایل `.env` این پوشه تنظیم کنید.
3.  برای اجرای سرویس در حالت توسعه، از ریشه پروژه دستور `pnpm run dev` را اجرا کنید.

## افزودن یک Use Case جدید (مثال: `GetMyProfile`)

1.  **Domain:**
    *   در `domain/user/repository.go`، متد `FindByID` را (اگر وجود ندارد) اضافه کنید.
2.  **Infrastructure:**
    *   در `infrastructure/database/postgres/user_repository.go`، متد `FindByID` را پیاده‌سازی کنید.
3.  **Application:**
    *   فایل `application/user/get_profile.go` را ایجاد کنید.
    *   `GetProfileUseCase` را با منطق لازم (فراخوانی `userRepo.FindByID`) بنویسید.
4.  **Presentation:**
    *   فایل `interfaces/http/handlers/user_handler.go` را ایجاد کنید.
    *   `GetMyProfile` handler را بنویسید که `userID` را از context گرفته و Use Case را فراخوانی می‌کند.
5.  **DI & Routing:**
    *   در `cmd/api/main.go`، `UserRepository`، `GetProfileUseCase`، و `UserHandler` را ایجاد و تزریق کنید.
    *   روت `protected.GET("/users/me", userHandler.GetMyProfile)` را اضافه کنید.
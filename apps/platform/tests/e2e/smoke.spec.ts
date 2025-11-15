// apps/platform/tests/e2e/smoke.spec.ts

import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:3000';
const UNIQUE_EMAIL = `test-user-${Date.now()}@example.com`;
const PASSWORD = 'Password123!';

test.describe('Smoke Test: Full User Journey', () => {

  test('should allow a new user to register, login, view an article, and logout', async ({ page }) => {
    // 1. ثبت‌نام
    await page.goto(`${BASE_URL}/auth/register`);
    await page.fill('input[type="email"]', UNIQUE_EMAIL);
    await page.fill('input[type="password"]', PASSWORD);
    await page.fill('input[id="fullName"]', 'کاربر تستی');
    await page.fill('input[id="username"]', `testuser${Date.now()}`);
    await page.click('button[type="submit"]');
    
    // باید به داشبورد هدایت شود
    await expect(page).toHaveURL(`${BASE_URL}/dashboard`);
    await expect(page.locator('h1')).toContainText('داشبورد');
    await expect(page.locator('h2')).toContainText('سلام، کاربر تستی!');

    // 2. خروج از سیستم
    await page.click('button:has-text("خروج")');
    await expect(page).toHaveURL(BASE_URL);

    // 3. ورود مجدد
    await page.goto(`${BASE_URL}/auth/login`);
    await page.fill('input[type="email"]', UNIQUE_EMAIL);
    await page.fill('input[type="password"]', PASSWORD);
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(`${BASE_URL}/dashboard`);
    
    // 4. مشاهده لیست مقالات
    // (فرض می‌کنیم مقاله‌ای با اسلاگ 'is-our-brain-shrinking' وجود دارد)
    await page.goto(`${BASE_URL}/articles`);
    await page.click('a[href*="is-our-brain-shrinking"]');
    
    // 5. بررسی محتوای مقاله
    await expect(page).toHaveURL(`${BASE_URL}/articles/is-our-brain-shrinking`);
    await expect(page.locator('h1')).toContainText('آیا مغز ما دارد آب می رود!؟');

    // 6. بوکمارک کردن مقاله
    const bookmarkButton = page.locator('button:has-text("بوکمارک کردن")');
    await bookmarkButton.click();
    
    // (اختیاری) بررسی اینکه دکمه تغییر حالت می‌دهد
    // await expect(bookmarkButton).toContainText('بوکمارک شد');

    // 7. مشاهده لیست بوکمارک‌ها
    await page.goto(`${BASE_URL}/dashboard/bookmarks`);
    await expect(page.locator('h2')).toContainText('آیا مغز ما دارد آب می رود!؟');
  });
});
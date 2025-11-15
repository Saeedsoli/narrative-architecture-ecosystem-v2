// apps/backend/lib/utils/slug.go

package utils

import (
	"regexp"
	"strings"
)

var (
	nonLatinRegex = regexp.MustCompile(`[^a-z0-9\s-]`)
	spaceRegex    = regexp.MustCompile(`\s+`)
	dashRegex     = regexp.MustCompile(`-+`)
)

// CreateSlug یک اسلاگ URL-friendly از یک رشته ایجاد می‌کند.
func CreateSlug(title string) string {
	// 1. تبدیل به حروف کوچک
	slug := strings.ToLower(title)

	// 2. Transliteration برای کاراکترهای فارسی (نمونه)
	// در عمل، استفاده از یک کتابخانه کامل‌تر توصیه می‌شود.
	slug = strings.ReplaceAll(slug, "آ", "a")
	slug = strings.ReplaceAll(slug, "ا", "a")
	slug = strings.ReplaceAll(slug, "ب", "b")
	// ... (سایر حروف)

	// 3. حذف کاراکترهای غیرمجاز
	slug = nonLatinRegex.ReplaceAllString(slug, "")

	// 4. جایگزینی فاصله‌ها با خط تیره
	slug = spaceRegex.ReplaceAllString(slug, "-")

	// 5. حذف خط تیره‌های تکراری
	slug = dashRegex.ReplaceAllString(slug, "-")

	// 6. حذف خط تیره از ابتدا و انتها
	slug = strings.Trim(slug, "-")

	return slug
}
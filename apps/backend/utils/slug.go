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
// این تابع برای هر دو زبان فارسی و انگلیسی کار می‌کند.
func CreateSlug(title string) string {
	// 1. تبدیل به حروف کوچک
	slug := strings.ToLower(title)

	// 2. Transliteration برای کاراکترهای فارسی (نمونه ساده)
	// در عمل، از یک کتابخانه کامل‌تر استفاده کنید.
	slug = strings.ReplaceAll(slug, "آ", "a")
	slug = strings.ReplaceAll(slug, "ا", "a")
	slug = strings.ReplaceAll(slug, "ب", "b")
	slug = strings.ReplaceAll(slug, "پ", "p")
	slug = strings.ReplaceAll(slug, "ت", "t")
	slug = strings.ReplaceAll(slug, "ث", "s")
	slug = strings.ReplaceAll(slug, "ج", "j")
	slug = strings.ReplaceAll(slug, "چ", "ch")
	slug = strings.ReplaceAll(slug, "ح", "h")
	slug = strings.ReplaceAll(slug, "خ", "kh")
	slug = strings.ReplaceAll(slug, "د", "d")
	slug = strings.ReplaceAll(slug, "ذ", "z")
	slug = strings.ReplaceAll(slug, "ر", "r")
	slug = strings.ReplaceAll(slug, "ز", "z")
	slug = strings.ReplaceAll(slug, "ژ", "zh")
	slug = strings.ReplaceAll(slug, "س", "s")
	slug = strings.ReplaceAll(slug, "ش", "sh")
	slug = strings.ReplaceAll(slug, "ص", "s")
	slug = strings.ReplaceAll(slug, "ض", "z")
	slug = strings.ReplaceAll(slug, "ط", "t")
	slug = strings.ReplaceAll(slug, "ظ", "z")
	slug = strings.ReplaceAll(slug, "ع", "a")
	slug = strings.ReplaceAll(slug, "غ", "gh")
	slug = strings.ReplaceAll(slug, "ف", "f")
	slug = strings.ReplaceAll(slug, "ق", "gh")
	slug = strings.ReplaceAll(slug, "ک", "k")
	slug = strings.ReplaceAll(slug, "گ", "g")
	slug = strings.ReplaceAll(slug, "ل", "l")
	slug = strings.ReplaceAll(slug, "م", "m")
	slug = strings.ReplaceAll(slug, "ن", "n")
	slug = strings.ReplaceAll(slug, "و", "v")
	slug = strings.ReplaceAll(slug, "ه", "h")
	slug = strings.ReplaceAll(slug, "ی", "y")
	
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
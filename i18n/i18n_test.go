package i18n

import "testing"

func TestLocale_DefaultEnglish(t *testing.T) {
	t.Setenv("FILEX_LANG", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")

	if got := Locale(); got != "en" {
		t.Fatalf("Locale() = %q, want %q", got, "en")
	}
}

func TestLocale_FilexLangPrecedence(t *testing.T) {
	t.Setenv("FILEX_LANG", "zh_CN.UTF-8")
	t.Setenv("LC_ALL", "en_US.UTF-8")
	t.Setenv("LC_MESSAGES", "en_US.UTF-8")
	t.Setenv("LANG", "en_US.UTF-8")

	if got := Locale(); got != "zh" {
		t.Fatalf("Locale() = %q, want %q", got, "zh")
	}
}

func TestLocale_FallbackToLang(t *testing.T) {
	t.Setenv("FILEX_LANG", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "zh_TW.UTF-8")

	if got := Locale(); got != "zh" {
		t.Fatalf("Locale() = %q, want %q", got, "zh")
	}
}

func TestTranslateChinese(t *testing.T) {
	t.Setenv("FILEX_LANG", "zh")

	if got := T("Files"); got != "文件" {
		t.Fatalf("T(Files) = %q, want %q", got, "文件")
	}
	if got := T("%d items", 3); got != "3 项" {
		t.Fatalf("T(%%d items, 3) = %q, want %q", got, "3 项")
	}
}

func TestTranslateFallbackEnglish(t *testing.T) {
	t.Setenv("FILEX_LANG", "en_US")

	if got := T("Files"); got != "Files" {
		t.Fatalf("T(Files) = %q, want %q", got, "Files")
	}
	if got := T("Unknown Key"); got != "Unknown Key" {
		t.Fatalf("T(Unknown Key) = %q, want key itself", got)
	}
}

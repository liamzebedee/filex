package i18n

import (
	"fmt"
	"os"
	"strings"
)

var zhCN = map[string]string{
	"Files":                       "文件",
	"Back":                        "后退",
	"Forward":                     "前进",
	"Search…":                     "搜索…",
	"List view":                   "列表视图",
	"Icon view":                   "图标视图",
	"Name":                        "名称",
	"Size":                        "大小",
	"Modified":                    "修改时间",
	"Open":                        "打开",
	"Open in Terminal":            "在终端中打开",
	"New Folder…":                 "新建文件夹…",
	"Cut":                         "剪切",
	"Copy":                        "复制",
	"Paste":                       "粘贴",
	"Copy Path":                   "复制路径",
	"Rename…":                     "重命名…",
	"Extract Here":                "解压到此处",
	"Move to Trash":               "移到回收站",
	"Properties":                  "属性",
	"Rename":                      "重命名",
	"Cancel":                      "取消",
	"Enter new name:":             "输入新名称：",
	"Rename Failed":               "重命名失败",
	"New Folder":                  "新建文件夹",
	"Create":                      "创建",
	"Folder name:":                "文件夹名称：",
	"Create Folder Failed":        "创建文件夹失败",
	"Move \"%s\" to the Trash?":   "要将“%s”移到回收站吗？",
	"Move %d items to the Trash?": "要将 %d 个项目移到回收站吗？",
	"Error":                       "错误",
	"Close":                       "关闭",
	"Type:":                       "类型：",
	"Folder":                      "文件夹",
	"Location:":                   "位置：",
	"Permissions:":                "权限：",
	"Places":                      "位置",
	"Remove Bookmark":             "移除书签",
	"1 item":                      "1 项",
	"%d items":                    "%d 项",
	"Free space: %s":              "可用空间：%s",
	"Home":                        "主目录",
	"Desktop":                     "桌面",
	"Documents":                   "文档",
	"Downloads":                   "下载",
	"Music":                       "音乐",
	"Pictures":                    "图片",
	"Videos":                      "视频",
	"Trash":                       "回收站",
	"File System":                 "文件系统",
	"TodayPrefix":                 "今天",
	"Date.ThisYear":               "01-02 15:04",
	"Date.PastYear":               "2006-01-02",
	"Name:":                       "名称：",
	"Size:":                       "大小：",
	"Modified:":                   "修改时间：",
}

// Locale returns a normalized locale token used by this app.
func Locale() string {
	for _, key := range []string{"FILEX_LANG", "LC_ALL", "LC_MESSAGES", "LANG"} {
		if raw := strings.TrimSpace(os.Getenv(key)); raw != "" {
			return normalizeLocale(raw)
		}
	}
	return "en"
}

func normalizeLocale(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "en"
	}
	if idx := strings.IndexAny(raw, ".@"); idx >= 0 {
		raw = raw[:idx]
	}
	raw = strings.ReplaceAll(raw, "-", "_")
	if strings.HasPrefix(raw, "zh") {
		return "zh"
	}
	return "en"
}

func IsChinese() bool {
	return Locale() == "zh"
}

// T translates a message key and applies fmt.Sprintf formatting when args are provided.
func T(key string, args ...interface{}) string {
	msg := key
	if Locale() == "zh" {
		if translated, ok := zhCN[key]; ok {
			msg = translated
		}
	}
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

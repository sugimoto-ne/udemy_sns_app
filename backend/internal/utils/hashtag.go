package utils

import (
	"regexp"
	"strings"
)

// ExtractHashtags コンテンツからハッシュタグを抽出する
// @param content 投稿内容
// @return ハッシュタグ名のスライス（#を除く、重複削除済み、最大10個）
func ExtractHashtags(content string) []string {
	// 正規表現: #の後に1文字以上の英数字、アンダースコア、Unicode文字
	// \p{L}はUnicodeのLetterカテゴリ（日本語など対応）
	re := regexp.MustCompile(`#([a-zA-Z0-9_\p{L}]+)`)

	matches := re.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		return []string{}
	}

	// 重複削除のためmapを使用
	hashtagSet := make(map[string]bool)
	var hashtags []string

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// #を除いた部分を取得
		tag := strings.ToLower(match[1]) // 小文字に統一

		// 既に存在しない場合のみ追加
		if !hashtagSet[tag] {
			hashtagSet[tag] = true
			hashtags = append(hashtags, tag)

			// 最大10個に制限
			if len(hashtags) >= 10 {
				break
			}
		}
	}

	return hashtags
}

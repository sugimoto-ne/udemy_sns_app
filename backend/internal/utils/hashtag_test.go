package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractHashtags(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "単一のハッシュタグ",
			content:  "これは #テスト です",
			expected: []string{"テスト"},
		},
		{
			name:     "複数のハッシュタグ",
			content:  "#Go言語 と #プログラミング について",
			expected: []string{"go言語", "プログラミング"},
		},
		{
			name:     "英数字のハッシュタグ",
			content:  "Learning #golang and #programming",
			expected: []string{"golang", "programming"},
		},
		{
			name:     "アンダースコア付きハッシュタグ",
			content:  "#tech_news #web_development",
			expected: []string{"tech_news", "web_development"},
		},
		{
			name:     "重複するハッシュタグ",
			content:  "#テスト #test #テスト",
			expected: []string{"テスト", "test"},
		},
		{
			name:     "10個を超えるハッシュタグ（最大10個に制限）",
			content:  "#1 #2 #3 #4 #5 #6 #7 #8 #9 #10 #11 #12",
			expected: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
		{
			name:     "ハッシュタグなし",
			content:  "これは普通のテキストです",
			expected: []string{},
		},
		{
			name:     "数字を含むハッシュタグ",
			content:  "価格は#1000です",
			expected: []string{"1000です"},
		},
		{
			name:     "大文字小文字を統一（小文字に変換）",
			content:  "#GoLang #GOLANG #golang",
			expected: []string{"golang"},
		},
		{
			name:     "日本語と英語混在",
			content:  "#プログラミング学習 と #100DaysOfCode チャレンジ",
			expected: []string{"プログラミング学習", "100daysofcode"},
		},
		{
			name:     "空文字列",
			content:  "",
			expected: []string{},
		},
		{
			name:     "ハッシュタグが文末",
			content:  "今日は良い天気 #晴れ",
			expected: []string{"晴れ"},
		},
		{
			name:     "連続するハッシュタグ",
			content:  "#tag1#tag2#tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractHashtags(tt.content)
			assert.Equal(t, tt.expected, result, "ハッシュタグの抽出結果が期待値と一致しません")
		})
	}
}

func TestExtractHashtagsPerformance(t *testing.T) {
	// 大量のテキストからハッシュタグを抽出するパフォーマンステスト
	longContent := ""
	for i := 0; i < 100; i++ {
		longContent += "#tag" + string(rune(i)) + " "
	}

	result := ExtractHashtags(longContent)
	assert.LessOrEqual(t, len(result), 10, "ハッシュタグは最大10個まで")
}

func TestExtractHashtagsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "特殊文字の後のハッシュタグ",
			content:  "Hello!#world",
			expected: []string{"world"},
		},
		{
			name:     "改行を含むテキスト",
			content:  "First line #tag1\nSecond line #tag2",
			expected: []string{"tag1", "tag2"},
		},
		{
			name:     "タブを含むテキスト",
			content:  "Tab\t#tag1\tseparated",
			expected: []string{"tag1"},
		},
		{
			name:     "絵文字の後のハッシュタグ",
			content:  "😀 #happy #emoji",
			expected: []string{"happy", "emoji"},
		},
		{
			name:     "URLに含まれる#",
			content:  "Visit https://example.com#section and use #hashtag",
			expected: []string{"section", "hashtag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractHashtags(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

package renderer

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render - テンプレートをレンダリング
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// NewTemplateRenderer - テンプレートレンダラーを初期化
func NewTemplateRenderer(templatesDir string) (*TemplateRenderer, error) {
	// テンプレートディレクトリの存在確認（デバッグ用）
	absPath, _ := filepath.Abs(templatesDir)
	println("Loading templates from:", absPath)

	// 空のテンプレートを初期化
	templates := template.New("root")

	// すべてのテンプレートを読み込み
	pattern := filepath.Join(templatesDir, "*.html")
	println("Pattern:", pattern)
	matches, _ := filepath.Glob(pattern)
	println("Found", len(matches), "files in root")

	// ルートディレクトリのテンプレートがあれば読み込み
	if len(matches) > 0 {
		var err error
		templates, err = templates.ParseGlob(pattern)
		if err != nil {
			println("Error loading templates:", err.Error())
			return nil, err
		}
		println("Successfully loaded", len(templates.Templates()), "templates from root")
	}

	// サブディレクトリのテンプレートも読み込み
	patterns := []string{
		filepath.Join(templatesDir, "users", "*.html"),
		filepath.Join(templatesDir, "password_resets", "*.html"),
		filepath.Join(templatesDir, "logs", "*.html"),
	}

	for _, pattern := range patterns {
		matches, _ := filepath.Glob(pattern)
		println("Pattern:", pattern, "- Found", len(matches), "files")
		if len(matches) > 0 {
			var err error
			templates, err = templates.ParseGlob(pattern)
			if err != nil {
				println("Error loading templates from", pattern, ":", err.Error())
				continue
			}
			println("Successfully loaded", len(matches), "templates from", pattern)
		}
	}

	// テンプレートが1つも読み込まれていないかチェック
	totalTemplates := len(templates.Templates())
	println("Total templates loaded:", totalTemplates)
	if totalTemplates == 0 {
		return nil, fmt.Errorf("no templates found in %s", templatesDir)
	}

	return &TemplateRenderer{
		templates: templates,
	}, nil
}

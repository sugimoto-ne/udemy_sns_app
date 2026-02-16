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

	// すべてのテンプレートを読み込み
	pattern := filepath.Join(templatesDir, "*.html")
	println("Pattern:", pattern)
	templates, err := template.ParseGlob(pattern)
	if err != nil {
		println("Error loading templates:", err.Error())
		return nil, err
	}

	if templates == nil {
		println("WARNING: templates is nil after ParseGlob")
		return nil, fmt.Errorf("no templates found in %s", templatesDir)
	}

	println("Successfully loaded", len(templates.Templates()), "templates from root")

	// サブディレクトリのテンプレートも読み込み
	patterns := []string{
		filepath.Join(templatesDir, "users", "*.html"),
		filepath.Join(templatesDir, "password_resets", "*.html"),
		filepath.Join(templatesDir, "logs", "*.html"),
	}

	for _, pattern := range patterns {
		templates, err = templates.ParseGlob(pattern)
		if err != nil {
			// ファイルが存在しない場合はスキップ
			continue
		}
	}

	return &TemplateRenderer{
		templates: templates,
	}, nil
}

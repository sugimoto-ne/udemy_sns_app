package renderer

import (
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
	// すべてのテンプレートを読み込み
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

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

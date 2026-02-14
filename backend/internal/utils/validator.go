package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct - 構造体をバリデーション
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// GetValidator - バリデータインスタンスを取得
func GetValidator() *validator.Validate {
	return validate
}

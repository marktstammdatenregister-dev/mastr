package generator

import "marktstammdatenregister.dev/internal/spec"

func MapGoType(field spec.Field) any {
	switch field.Xsd {
	case "nonNegativeInteger":
		return "uint"
	case "decimal":
		return "float32"
	default:
		return "string"
	}
}

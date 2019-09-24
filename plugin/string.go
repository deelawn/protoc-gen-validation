package plugin

import (
	"fmt"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	pb "github.com/neophenix/protoc-gen-validation"
)

func (p *Plugin) generateStringValidationCode(fieldName string, fieldValue string, v *pb.FieldValidation, mv *pb.MessageValidation, field *descriptor.FieldDescriptorProto) {
	if (v.Trim != nil && *v.Trim) || (mv != nil && mv.TrimStrings != nil && *mv.TrimStrings) {
		p.P(`%s = %s.Trim(%s, " ")`, fieldValue, p.stringsPkg.Use(), fieldValue)
	}
	if v.Lc != nil && *v.Lc {
		p.P(`%s = %s.ToLower(%s)`, fieldValue, p.stringsPkg.Use(), fieldValue)
	}
	if v.Uc != nil && *v.Uc {
		p.P(`%s = %s.ToUpper(%s)`, fieldValue, p.stringsPkg.Use(), fieldValue)
	}
	if v.NotEmptyString != nil {
		// start this now, we will close it out later, if we've set not_empty_string there is no sense doing extra validation on it it if is one
		// so this begins the wrap of the other possible validations
		p.P(`if %s != "" {`, fieldValue)
	}
	if v.Matches != nil {
		p.P(`if %s != "%s" {`, fieldValue, v.GetMatches())
		p.generateErrorCode(fieldName, v.GetMatches(), "{field} must equal {value}", v, mv, field, "")
		p.P(`}`)
	}
	if v.Contains != nil {
		p.P(`if !%s.Contains(%s, "%s") {`, p.stringsPkg.Use(), fieldValue, v.GetContains())
		p.generateErrorCode(fieldName, v.GetContains(), "{field} must contain {value}", v, mv, field, "")
		p.P(`}`)
	}
	if v.Regex != nil {
		p.P(`if !%s.MustCompile("%s").MatchString(%s) {`, p.regexPkg.Use(), v.GetRegex(), fieldValue)
		p.generateErrorCode(fieldName, v.GetRegex(), "{field} must match regex {value}", v, mv, field, "")
		p.P(`}`)
	}
	if v.MinLen != nil {
		p.P(`if len(%s) < %d {`, fieldValue, v.GetMinLen())
		p.generateErrorCode(fieldName, fmt.Sprintf("%d", v.GetMinLen()), "{field} must be at least {value} characters long", v, mv, field, "")
		p.P(`}`)
	}
	if v.MaxLen != nil {
		p.P(`if len(%s) > %d {`, fieldValue, v.GetMaxLen())
		p.generateErrorCode(fieldName, fmt.Sprintf("%d", v.GetMaxLen()), "{field} must be no more than {value} characters long", v, mv, field, "")
		p.P(`}`)
	}
	if v.EqLen != nil {
		p.P(`if len(%s) != %d {`, fieldValue, v.GetEqLen())
		p.generateErrorCode(fieldName, fmt.Sprintf("%d", v.GetEqLen()), "{field} must be exactly {value} characters long", v, mv, field, "")
		p.P(`}`)
	}
	if v.IsUuid != nil && *v.IsUuid {
		p.P(`if !isValidUUID(%s) {`, fieldValue)
		p.generateErrorCode(fieldName, "", "{field} must be a valid UUID", v, mv, field, "")
		p.P(`}`)
	}
	if v.IsEmail != nil && *v.IsEmail {
		p.P(`if !isValidEmail(%s) {`, fieldValue)
		p.generateErrorCode(fieldName, "", "{field} must be a valid email address", v, mv, field, "")
		p.P(`}`)
	}

	// so if we set not_empty_string, we only did validation if we had some value, now we need to decide if we need to output an error that it was empty
	if v.NotEmptyString != nil {
		if *v.NotEmptyString {
			p.P(`} else {`)
			p.generateErrorCode(fieldName, "", "{field} can not be an empty string", v, mv, field, "")
		}
		// we opened these up at the start, so close them now
		p.P(`}`)
	}
}

func isString(field *descriptor.FieldDescriptorProto) bool {
	if field.GetType() == descriptor.FieldDescriptorProto_TYPE_STRING {
		return true
	}
	if isWKTString(field.GetTypeName()) {
		return true
	}
	return false
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodingError := json.NewEncoder(w).Encode(data); encodingError != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func RespondError(w http.ResponseWriter, status int, message string, err *error) {
	RespondJSON(w, status, map[string]any{"error": message})
	if err != nil {
		fmt.Println(err)
	}
}

func ToJSON(value any) string {
	jsonBytes, marshalError := json.Marshal(value)
	if marshalError != nil {
		return "{}"
	}
	return string(jsonBytes)
}

func ToJSONPretty(value any) string {
	jsonBytes, marshalError := json.MarshalIndent(value, "", "  ")
	if marshalError != nil {
		return "{}"
	}
	return string(jsonBytes)
}

var validate = validator.New()

func init() {
	_ = validate.RegisterValidation("rfc3339", func(fieldLevel validator.FieldLevel) bool {
		if fieldLevel.Field().Kind() != reflect.String {
			return false
		}

		value := fieldLevel.Field().String()
		if value == "" {
			return true
		}

		_, parseError := time.Parse(time.RFC3339, value)
		return parseError == nil
	})
}

type RequestOptions struct {
	Body    any
	Params  any
	Query   any
	Headers any
}

func ParseRequest(r *http.Request, options RequestOptions) error {
	if options.Body != nil {
		if r.Body == nil {
			return fmt.Errorf("Missing request body")
		}
		if decodingError := json.NewDecoder(r.Body).Decode(options.Body); decodingError != nil {
			return fmt.Errorf("Invalid JSON body: %w", decodingError)
		}
		defer r.Body.Close()
		if validationError := validate.Struct(options.Body); validationError != nil {
			return fmt.Errorf("Body validation failed: %w", validationError)
		}
	}

	if options.Params != nil {
		if parsingError := parsePathParams(r, options.Params); parsingError != nil {
			return fmt.Errorf("Invalid path params: %w", parsingError)
		}
		if validationError := validate.Struct(options.Params); validationError != nil {
			return fmt.Errorf("Path params validation failed: %w", validationError)
		}
	}

	if options.Query != nil {
		if parsingError := parseQueryParams(r, options.Query); parsingError != nil {
			return fmt.Errorf("Invalid query params: %w", parsingError)
		}
		if validationError := validate.Struct(options.Query); validationError != nil {
			return fmt.Errorf("Query params validation failed: %w", validationError)
		}
	}

	if options.Headers != nil {
		if parsingError := parseHeaders(r, options.Headers); parsingError != nil {
			return fmt.Errorf("Invalid headers: %w", parsingError)
		}
		if validationError := validate.Struct(options.Headers); validationError != nil {
			return fmt.Errorf("Headers validation failed: %w", validationError)
		}
	}

	return nil
}

func parsePathParams(r *http.Request, targetStruct any) error {
	reflectValue := reflect.ValueOf(targetStruct)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("ParamsScheme must be a pointer to struct")
	}

	reflectValue = reflectValue.Elem()
	reflectType := reflectValue.Type()

	for fieldIndex := 0; fieldIndex < reflectValue.NumField(); fieldIndex++ {
		field := reflectValue.Field(fieldIndex)
		fieldType := reflectType.Field(fieldIndex)

		tagName := fieldType.Tag.Get("param")
		if tagName == "" {
			tagName = strings.ToLower(fieldType.Name)
		}

		paramValue := r.PathValue(tagName)
		if paramValue == "" {
			continue
		}

		if settingError := setFieldValue(field, paramValue); settingError != nil {
			return fmt.Errorf("field %s: %w", fieldType.Name, settingError)
		}
	}

	return nil
}

func parseQueryParams(r *http.Request, targetStruct any) error {
	reflectValue := reflect.ValueOf(targetStruct)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("QueryScheme must be a pointer to struct")
	}

	reflectValue = reflectValue.Elem()
	reflectType := reflectValue.Type()
	queryParams := r.URL.Query()

	for fieldIndex := 0; fieldIndex < reflectValue.NumField(); fieldIndex++ {
		field := reflectValue.Field(fieldIndex)
		fieldType := reflectType.Field(fieldIndex)

		tagName := fieldType.Tag.Get("query")
		if tagName == "" {
			tagName = strings.ToLower(fieldType.Name)
		}

		queryValue := queryParams.Get(tagName)
		if queryValue == "" {
			continue
		}

		if settingError := setFieldValue(field, queryValue); settingError != nil {
			return fmt.Errorf("Field %s: %w", fieldType.Name, settingError)
		}
	}

	return nil
}

func parseHeaders(r *http.Request, targetStruct any) error {
	reflectValue := reflect.ValueOf(targetStruct)
	if reflectValue.Kind() != reflect.Ptr || reflectValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("HeadersScheme must be a pointer to struct")
	}

	reflectValue = reflectValue.Elem()
	reflectType := reflectValue.Type()

	for fieldIndex := 0; fieldIndex < reflectValue.NumField(); fieldIndex++ {
		field := reflectValue.Field(fieldIndex)
		fieldType := reflectType.Field(fieldIndex)

		tagName := fieldType.Tag.Get("header")
		if tagName == "-" {
			continue
		}
		if tagName == "" {
			tagName = http.CanonicalHeaderKey(fieldType.Name)
		}

		headerValue := r.Header.Get(tagName)
		if headerValue == "" {
			continue
		}

		if settingError := setFieldValue(field, headerValue); settingError != nil {
			return fmt.Errorf("Field %s: %w", fieldType.Name, settingError)
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, stringValue string) error {
	if !field.CanSet() {
		return fmt.Errorf("Field cannot be set")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(stringValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		integerValue, parsingError := strconv.ParseInt(stringValue, 10, 64)
		if parsingError != nil {
			return fmt.Errorf("Cannot parse as int: %w", parsingError)
		}
		field.SetInt(integerValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		unsignedIntegerValue, parsingError := strconv.ParseUint(stringValue, 10, 64)
		if parsingError != nil {
			return fmt.Errorf("Cannot parse as uint: %w", parsingError)
		}
		field.SetUint(unsignedIntegerValue)
	case reflect.Float32, reflect.Float64:
		floatValue, parsingError := strconv.ParseFloat(stringValue, 64)
		if parsingError != nil {
			return fmt.Errorf("Cannot parse as float: %w", parsingError)
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		booleanValue, parsingError := strconv.ParseBool(stringValue)
		if parsingError != nil {
			return fmt.Errorf("Cannot parse as bool: %w", parsingError)
		}
		field.SetBool(booleanValue)
	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setFieldValue(field.Elem(), stringValue)
	default:
		return fmt.Errorf("Unsupported field type: %s", field.Kind())
	}

	return nil
}

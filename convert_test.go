package protobaggins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestConvertProtoValueToInterface(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()
		result := ConvertProtoValueToInterface(nil)
		assert.Nil(t, result)
	})

	t.Run("null value", func(t *testing.T) {
		t.Parallel()
		nullValue, err := structpb.NewValue(nil)
		assert.NoError(t, err)
		result := ConvertProtoValueToInterface(nullValue)
		assert.Nil(t, result)
	})

	t.Run("number value", func(t *testing.T) {
		t.Parallel()
		numberValue, err := structpb.NewValue(42.5)
		assert.NoError(t, err)
		result := ConvertProtoValueToInterface(numberValue)
		assert.Equal(t, 42.5, result)
	})

	t.Run("string value", func(t *testing.T) {
		t.Parallel()
		stringValue, err := structpb.NewValue("test string")
		assert.NoError(t, err)
		result := ConvertProtoValueToInterface(stringValue)
		assert.Equal(t, "test string", result)
	})

	t.Run("bool value", func(t *testing.T) {
		t.Parallel()
		boolValue, err := structpb.NewValue(true)
		assert.NoError(t, err)
		result := ConvertProtoValueToInterface(boolValue)
		assert.Equal(t, true, result)
	})

	t.Run("list value", func(t *testing.T) {
		t.Parallel()
		listValue, err := structpb.NewValue([]interface{}{1, "two", true})
		assert.NoError(t, err)
		result := ConvertProtoValueToInterface(listValue)
		expected := []interface{}{float64(1), "two", true}
		assert.Equal(t, expected, result)
	})

	t.Run("struct value", func(t *testing.T) {
		t.Parallel()
		structValue, err := structpb.NewValue(map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		})
		assert.NoError(t, err)
		result := ConvertProtoValueToInterface(structValue)
		expected := map[string]interface{}{
			"key1": "value1",
			"key2": float64(42),
			"key3": true,
		}
		assert.Equal(t, expected, result)
	})
}

func TestMapToStructValues(t *testing.T) {
	t.Parallel()

	t.Run("nil map", func(t *testing.T) {
		t.Parallel()
		result := MapToStructValues(nil)
		assert.Nil(t, result)
	})

	t.Run("empty map", func(t *testing.T) {
		t.Parallel()
		result := MapToStructValues(map[string]any{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("map with primitive values", func(t *testing.T) {
		t.Parallel()
		input := map[string]any{
			"string":  "value",
			"number":  42.5,
			"bool":    true,
			"null":    nil,
			"integer": 10,
		}

		result := MapToStructValues(input)

		assert.Len(t, result, 5)
		assert.Equal(t, "value", result["string"].GetStringValue())
		assert.Equal(t, 42.5, result["number"].GetNumberValue())
		assert.Equal(t, true, result["bool"].GetBoolValue())
		assert.NotNil(t, result["null"])
		assert.Equal(t, float64(10), result["integer"].GetNumberValue())
	})

	t.Run("map with complex values", func(t *testing.T) {
		t.Parallel()
		input := map[string]any{
			"list": []any{1, "two", true},
			"map": map[string]any{
				"nested": "value",
			},
		}

		result := MapToStructValues(input)

		assert.Len(t, result, 2)

		// Check list value
		listVal := result["list"].GetListValue().GetValues()
		assert.Len(t, listVal, 3)
		assert.Equal(t, float64(1), listVal[0].GetNumberValue())
		assert.Equal(t, "two", listVal[1].GetStringValue())
		assert.Equal(t, true, listVal[2].GetBoolValue())

		// Check map value
		mapVal := result["map"].GetStructValue().GetFields()
		assert.Len(t, mapVal, 1)
		assert.Equal(t, "value", mapVal["nested"].GetStringValue())
	})

	t.Run("map with unconvertible values", func(t *testing.T) {
		t.Parallel()

		// Create a value that cannot be represented in protobuf
		type unconvertible struct {
			Field string
		}

		input := map[string]any{
			"valid":   "value",
			"invalid": unconvertible{Field: "test"},
		}

		result := MapToStructValues(input)

		// Only the valid field should be present
		assert.Len(t, result, 1)
		assert.Equal(t, "value", result["valid"].GetStringValue())
		assert.NotContains(t, result, "invalid")
	})
}

func TestStructValuesToMap(t *testing.T) {
	t.Parallel()

	t.Run("nil map", func(t *testing.T) {
		t.Parallel()
		result := StructValuesToMap(nil)
		assert.Nil(t, result)
	})

	t.Run("empty map", func(t *testing.T) {
		t.Parallel()
		result := StructValuesToMap(map[string]*structpb.Value{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("map with primitive values", func(t *testing.T) {
		t.Parallel()
		input := map[string]*structpb.Value{}

		strVal, _ := structpb.NewValue("value")
		numVal, _ := structpb.NewValue(42.5)
		boolVal, _ := structpb.NewValue(true)
		nullVal, _ := structpb.NewValue(nil)

		input["string"] = strVal
		input["number"] = numVal
		input["bool"] = boolVal
		input["null"] = nullVal

		result := StructValuesToMap(input)

		assert.Len(t, result, 4)
		assert.Equal(t, "value", result["string"])
		assert.Equal(t, 42.5, result["number"])
		assert.Equal(t, true, result["bool"])
		assert.Nil(t, result["null"])
	})

	t.Run("map with complex values", func(t *testing.T) {
		t.Parallel()
		input := map[string]*structpb.Value{}

		listVal, _ := structpb.NewValue([]any{1, "two", true})
		mapVal, _ := structpb.NewValue(map[string]any{"nested": "value"})

		input["list"] = listVal
		input["map"] = mapVal

		result := StructValuesToMap(input)

		assert.Len(t, result, 2)

		// Check list value
		listResult, ok := result["list"].([]any)
		assert.True(t, ok)
		assert.Len(t, listResult, 3)
		assert.Equal(t, float64(1), listResult[0])
		assert.Equal(t, "two", listResult[1])
		assert.Equal(t, true, listResult[2])

		// Check map value
		mapResult, ok := result["map"].(map[string]any)
		assert.True(t, ok)
		assert.Len(t, mapResult, 1)
		assert.Equal(t, "value", mapResult["nested"])
	})
}

func TestTryNewStructValue(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()
		result := TryNewStructValue(nil)
		assert.NotNil(t, result)
		// Should create a null value
		_, ok := result.Kind.(*structpb.Value_NullValue)
		assert.True(t, ok)
	})

	t.Run("primitive values", func(t *testing.T) {
		t.Parallel()

		stringResult := TryNewStructValue("test")
		assert.Equal(t, "test", stringResult.GetStringValue())

		numberResult := TryNewStructValue(42.5)
		assert.Equal(t, 42.5, numberResult.GetNumberValue())

		boolResult := TryNewStructValue(true)
		assert.Equal(t, true, boolResult.GetBoolValue())
	})

	t.Run("complex values", func(t *testing.T) {
		t.Parallel()

		listResult := TryNewStructValue([]any{1, "two", true})
		assert.NotNil(t, listResult.GetListValue())
		assert.Len(t, listResult.GetListValue().GetValues(), 3)

		mapResult := TryNewStructValue(map[string]any{"key": "value"})
		assert.NotNil(t, mapResult.GetStructValue())
		assert.Len(t, mapResult.GetStructValue().GetFields(), 1)
		assert.Equal(t, "value", mapResult.GetStructValue().GetFields()["key"].GetStringValue())
	})

	t.Run("unconvertible value", func(t *testing.T) {
		t.Parallel()

		type unconvertible struct {
			Field string
		}

		result := TryNewStructValue(unconvertible{Field: "test"})
		assert.Nil(t, result)
	})
}

func TestSliceToStructValues(t *testing.T) {
	t.Parallel()

	t.Run("nil slice", func(t *testing.T) {
		t.Parallel()
		result := SliceToStructValues(nil)
		assert.Nil(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		result := SliceToStructValues([]any{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("slice with primitive values", func(t *testing.T) {
		t.Parallel()
		input := []any{
			"string",
			42.5,
			true,
			nil,
		}

		result := SliceToStructValues(input)

		assert.Len(t, result, 4)
		assert.Equal(t, "string", result[0].GetStringValue())
		assert.Equal(t, 42.5, result[1].GetNumberValue())
		assert.Equal(t, true, result[2].GetBoolValue())
		_, ok := result[3].Kind.(*structpb.Value_NullValue)
		assert.True(t, ok)
	})

	t.Run("slice with complex values", func(t *testing.T) {
		t.Parallel()
		input := []any{
			[]any{1, 2, 3},
			map[string]any{"nested": "value"},
		}

		result := SliceToStructValues(input)

		assert.Len(t, result, 2)

		// Check nested list
		listVal := result[0].GetListValue().GetValues()
		assert.Len(t, listVal, 3)
		assert.Equal(t, float64(1), listVal[0].GetNumberValue())
		assert.Equal(t, float64(2), listVal[1].GetNumberValue())
		assert.Equal(t, float64(3), listVal[2].GetNumberValue())

		// Check map
		mapVal := result[1].GetStructValue().GetFields()
		assert.Len(t, mapVal, 1)
		assert.Equal(t, "value", mapVal["nested"].GetStringValue())
	})

	t.Run("slice with unconvertible values", func(t *testing.T) {
		t.Parallel()

		type unconvertible struct {
			Field string
		}

		input := []any{
			"valid",
			unconvertible{Field: "test"},
			42,
		}

		result := SliceToStructValues(input)

		// Only valid values should be present
		assert.Len(t, result, 2)
		assert.Equal(t, "valid", result[0].GetStringValue())
		assert.Equal(t, float64(42), result[1].GetNumberValue())
	})
}

func TestStructValuesToSlice(t *testing.T) {
	t.Parallel()

	t.Run("nil slice", func(t *testing.T) {
		t.Parallel()
		result := StructValuesToSlice(nil)
		assert.Nil(t, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		result := StructValuesToSlice([]*structpb.Value{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("slice with primitive values", func(t *testing.T) {
		t.Parallel()
		var input []*structpb.Value

		strVal, err := structpb.NewValue("value")
		assert.NoError(t, err)
		numVal, err := structpb.NewValue(42.5)
		assert.NoError(t, err)
		boolVal, err := structpb.NewValue(true)
		assert.NoError(t, err)
		nullVal, err := structpb.NewValue(nil)
		assert.NoError(t, err)

		input = append(input, strVal, numVal, boolVal, nullVal)

		result := StructValuesToSlice(input)

		assert.Len(t, result, 4)
		assert.Equal(t, "value", result[0])
		assert.Equal(t, 42.5, result[1])
		assert.Equal(t, true, result[2])
		assert.Nil(t, result[3])
	})

	t.Run("slice with complex values", func(t *testing.T) {
		t.Parallel()
		var input []*structpb.Value

		listVal, err := structpb.NewValue([]any{1, "two", true})
		assert.NoError(t, err)
		mapVal, err := structpb.NewValue(map[string]any{"nested": "value"})
		assert.NoError(t, err)

		input = append(input, listVal, mapVal)

		result := StructValuesToSlice(input)

		assert.Len(t, result, 2)

		// Check list value
		listResult, ok := result[0].([]any)
		assert.True(t, ok)
		assert.Len(t, listResult, 3)
		assert.Equal(t, float64(1), listResult[0])
		assert.Equal(t, "two", listResult[1])
		assert.Equal(t, true, listResult[2])

		// Check map value
		mapResult, ok := result[1].(map[string]any)
		assert.True(t, ok)
		assert.Len(t, mapResult, 1)
		assert.Equal(t, "value", mapResult["nested"])
	})
}

func TestStringFromProto(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()
		result := StringFromProto(nil)
		assert.Equal(t, "", result, "should return empty string for nil input")
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		emptyStr := ""
		result := StringFromProto(&emptyStr)
		assert.Equal(t, "", result, "should return empty string for empty string input")
	})

	t.Run("non-empty string", func(t *testing.T) {
		t.Parallel()
		testStr := "test string"
		result := StringFromProto(&testStr)
		assert.Equal(t, testStr, result, "should return string value for non-empty string input")
	})
}

func TestStringToProto(t *testing.T) {
	t.Parallel()

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		result := StringToProto("")
		assert.NotNil(t, result, "should not return nil for empty string")
		assert.Equal(t, "", *result, "should properly store empty string")
	})

	t.Run("non-empty string", func(t *testing.T) {
		t.Parallel()
		testStr := "test string"
		result := StringToProto(testStr)
		assert.NotNil(t, result, "should not return nil")
		assert.Equal(t, testStr, *result, "should properly store string value")
	})
}

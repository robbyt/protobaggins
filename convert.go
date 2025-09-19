package protobaggins

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// ConvertProtoValueToInterface converts a protobuf structpb.Value to a Go any
// This should be replaced with v.AsInterface() for new code, but is kept for compatibility
func ConvertProtoValueToInterface(v *structpb.Value) any {
	if v == nil {
		return nil
	}
	return v.AsInterface()
}

// MapToStructValues converts a Go map[string]any to a map[string]*structpb.Value
// Silently skips values that cannot be converted to protobuf values
func MapToStructValues(m map[string]any) map[string]*structpb.Value {
	if m == nil {
		return nil
	}

	result := make(map[string]*structpb.Value, len(m))
	for k, v := range m {
		pbValue, err := structpb.NewValue(v)
		if err == nil {
			result[k] = pbValue
		}
	}
	return result
}

// SliceToStructValues converts a slice of any Go values to a slice of protocol buffer values
// Silently skips values that cannot be converted to protobuf values
func SliceToStructValues(values []any) []*structpb.Value {
	if values == nil {
		return nil
	}

	result := make([]*structpb.Value, 0, len(values))
	for _, v := range values {
		pbValue, err := structpb.NewValue(v)
		if err == nil {
			result = append(result, pbValue)
		}
	}
	return result
}

// StringFromProto safely converts a protocol buffer string pointer to a Go string
// Returns an empty string if the pointer is nil
func StringFromProto(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringToProto converts a Go string to a protocol buffer string pointer
// This is a convenience wrapper around proto.String
func StringToProto(s string) *string {
	return proto.String(s)
}

// StructValuesToMap converts a map[string]*structpb.Value to a Go map[string]any
func StructValuesToMap(m map[string]*structpb.Value) map[string]any {
	if m == nil {
		return nil
	}

	result := make(map[string]any, len(m))
	for k, v := range m {
		result[k] = v.AsInterface()
	}
	return result
}

// StructValuesToSlice converts a list of protocol buffer values to a slice of Go values
func StructValuesToSlice(values []*structpb.Value) []any {
	if values == nil {
		return nil
	}

	result := make([]any, len(values))
	for i, v := range values {
		result[i] = v.AsInterface()
	}
	return result
}

// TryNewStructValue creates a new *structpb.Value from a Go value
// Returns nil if the value cannot be converted to a protocol buffer value
func TryNewStructValue(v any) *structpb.Value {
	pbValue, err := structpb.NewValue(v)
	if err != nil {
		return nil
	}
	return pbValue
}

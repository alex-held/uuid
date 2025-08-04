package uuid

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatedUUID_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "nil UUID",
			input:   "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "not-a-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, result.IsZero())
			} else {
				assert.NoError(t, err)
				assert.False(t, result.IsZero())
				assert.Equal(t, tt.input, result.String())
			}
		})
	}
}

func TestValidatedUUID_FromGoogleUUID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		googleUUID := uuid.New()
		result, err := FromGoogleUUID(googleUUID)
		require.NoError(t, err)
		assert.Equal(t, googleUUID.String(), result.String())
	})

	t.Run("nil UUID", func(t *testing.T) {
		result, err := FromGoogleUUID(uuid.Nil)
		assert.Error(t, err)
		assert.True(t, result.IsZero())
	})
}

func TestValidatedUUID_JSON(t *testing.T) {
	t.Run("marshal valid UUID", func(t *testing.T) {
		u := New()
		data, err := json.Marshal(u)
		require.NoError(t, err)

		var unmarshaled ValidatedUUID
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)
		assert.Equal(t, u.String(), unmarshaled.String())
	})

	t.Run("marshal zero UUID fails", func(t *testing.T) {
		var u ValidatedUUID
		_, err := json.Marshal(u)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("unmarshal invalid UUID fails", func(t *testing.T) {
		data := []byte(`"invalid-uuid"`)
		var u ValidatedUUID
		err := json.Unmarshal(data, &u)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})
}

func TestValidatedUUID_Proto(t *testing.T) {
	t.Run("to proto valid UUID", func(t *testing.T) {
		u := New()
		pb, err := u.ToProto()
		require.NoError(t, err)
		assert.Equal(t, u.String(), pb.GetVal())
	})

	t.Run("to proto zero UUID fails", func(t *testing.T) {
		var u ValidatedUUID
		_, err := u.ToProto()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("from proto valid UUID", func(t *testing.T) {
		original := New()
		pb, err := original.ToProto()
		require.NoError(t, err)

		result, err := FromProto(pb)
		require.NoError(t, err)
		assert.Equal(t, original.String(), result.String())
	})

	t.Run("from proto nil fails", func(t *testing.T) {
		_, err := FromProto(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be nil")
	})

	t.Run("from proto invalid UUID fails", func(t *testing.T) {
		pb := &UUID{Val: "invalid-uuid"}
		_, err := FromProto(pb)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid UUID format")
	})
}

func TestHelpers(t *testing.T) {
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("StringToProto", func(t *testing.T) {
		pb, err := StringToProto(validUUIDStr)
		require.NoError(t, err)
		assert.Equal(t, validUUIDStr, pb.GetVal())

		_, err = StringToProto("invalid")
		assert.Error(t, err)
	})

	t.Run("ProtoToString", func(t *testing.T) {
		pb := &UUID{Val: validUUIDStr}
		result, err := ProtoToString(pb)
		require.NoError(t, err)
		assert.Equal(t, validUUIDStr, result)

		_, err = ProtoToString(nil)
		assert.Error(t, err)
	})

	t.Run("ValidateProtoUUID", func(t *testing.T) {
		pb := &UUID{Val: validUUIDStr}
		err := ValidateProtoUUID(pb)
		assert.NoError(t, err)

		err = ValidateProtoUUID(nil)
		assert.Error(t, err)

		pb = &UUID{Val: "invalid"}
		err = ValidateProtoUUID(pb)
		assert.Error(t, err)
	})

	t.Run("ValidateStringUUID", func(t *testing.T) {
		err := ValidateStringUUID(validUUIDStr)
		assert.NoError(t, err)

		err = ValidateStringUUID("invalid")
		assert.Error(t, err)
	})
}

package uuid

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ValidatedUUID wraps google/uuid.UUID with protobuf marshalling validation
type ValidatedUUID struct {
	uuid.UUID
}

// New creates a new ValidatedUUID
func New() ValidatedUUID {
	return ValidatedUUID{UUID: uuid.New()}
}

// Parse parses a string into a ValidatedUUID with validation
func Parse(s string) (ValidatedUUID, error) {
	if s == "" {
		return ValidatedUUID{}, fmt.Errorf("UUID cannot be empty")
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return ValidatedUUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	if parsed == uuid.Nil {
		return ValidatedUUID{}, fmt.Errorf("UUID cannot be nil/zero value")
	}

	return ValidatedUUID{UUID: parsed}, nil
}

// MustParse parses a string into a ValidatedUUID, panicking on error
func MustParse(s string) ValidatedUUID {
	u, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

// FromGoogleUUID converts a google/uuid.UUID to our ValidatedUUID type
func FromGoogleUUID(u uuid.UUID) (ValidatedUUID, error) {
	if u == uuid.Nil {
		return ValidatedUUID{}, fmt.Errorf("UUID cannot be nil/zero value")
	}
	return ValidatedUUID{UUID: u}, nil
}

// MustFromGoogleUUID converts a google/uuid.UUID to our ValidatedUUID type, panicking on error
func MustFromGoogleUUID(u uuid.UUID) ValidatedUUID {
	result, err := FromGoogleUUID(u)
	if err != nil {
		panic(err)
	}
	return result
}

// IsZero returns true if the UUID is the zero value
func (u ValidatedUUID) IsZero() bool {
	return u.UUID == uuid.Nil
}

// Validate ensures the UUID is not zero and is properly formatted
func (u ValidatedUUID) Validate() error {
	if u.UUID == uuid.Nil {
		return fmt.Errorf("UUID cannot be nil/zero value")
	}
	return nil
}

// String returns the string representation of the UUID
func (u ValidatedUUID) String() string {
	return u.UUID.String()
}

// MarshalJSON implements json.Marshaler with validation
func (u ValidatedUUID) MarshalJSON() ([]byte, error) {
	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("UUID validation failed during JSON marshalling: %w", err)
	}
	return json.Marshal(u.UUID.String())
}

// UnmarshalJSON implements json.Unmarshaler with validation
func (u *ValidatedUUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := Parse(s)
	if err != nil {
		return fmt.Errorf("UUID validation failed during JSON unmarshalling: %w", err)
	}

	*u = parsed
	return nil
}

// Value implements driver.Valuer for database operations
func (u ValidatedUUID) Value() (driver.Value, error) {
	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("UUID validation failed during database write: %w", err)
	}
	return u.UUID.String(), nil
}

// Scan implements sql.Scanner for database operations
func (u *ValidatedUUID) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("UUID cannot be nil")
	}

	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("cannot scan %T into UUID", value)
	}

	parsed, err := Parse(s)
	if err != nil {
		return fmt.Errorf("UUID validation failed during database scan: %w", err)
	}

	*u = parsed
	return nil
}

// ToProto converts the ValidatedUUID to a protobuf UUID message with validation
func (u ValidatedUUID) ToProto() (*UUID, error) {
	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("UUID validation failed during protobuf marshalling: %w", err)
	}
	return &UUID{
		Val: u.String(),
	}, nil
}

// MustToProto converts the ValidatedUUID to a protobuf UUID message, panicking on validation error
func (u ValidatedUUID) MustToProto() *UUID {
	pb, err := u.ToProto()
	if err != nil {
		panic(err)
	}
	return pb
}

// FromProto creates a ValidatedUUID from a protobuf UUID message
func FromProto(pb *UUID) (ValidatedUUID, error) {
	if pb == nil {
		return ValidatedUUID{}, fmt.Errorf("protobuf UUID cannot be nil")
	}
	return Parse(pb.GetVal())
}

// MustFromProto creates a ValidatedUUID from a protobuf UUID message, panicking on error
func MustFromProto(pb *UUID) ValidatedUUID {
	u, err := FromProto(pb)
	if err != nil {
		panic(err)
	}
	return u
}

// ToStringValue converts ValidatedUUID to protobuf StringValue with validation
func (u ValidatedUUID) ToStringValue() (*wrapperspb.StringValue, error) {
	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("UUID validation failed: %w", err)
	}
	return wrapperspb.String(u.String()), nil
}

// FromStringValue creates ValidatedUUID from protobuf StringValue with validation
func FromStringValue(sv *wrapperspb.StringValue) (ValidatedUUID, error) {
	if sv == nil {
		return ValidatedUUID{}, fmt.Errorf("StringValue cannot be nil")
	}
	return Parse(sv.Value)
}

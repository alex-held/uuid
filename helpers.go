package uuid

import (
	"fmt"
)

// StringToProto converts a string UUID to protobuf UUID with validation
func StringToProto(s string) (*UUID, error) {
	validated, err := Parse(s)
	if err != nil {
		return nil, err
	}
	return validated.ToProto()
}

// MustStringToProto converts a string UUID to protobuf UUID, panicking on error
func MustStringToProto(s string) *UUID {
	pb, err := StringToProto(s)
	if err != nil {
		panic(err)
	}
	return pb
}

// ProtoToString converts a protobuf UUID to string with validation
func ProtoToString(pb *UUID) (string, error) {
	validated, err := FromProto(pb)
	if err != nil {
		return "", err
	}
	return validated.String(), nil
}

// MustProtoToString converts a protobuf UUID to string, panicking on error
func MustProtoToString(pb *UUID) string {
	s, err := ProtoToString(pb)
	if err != nil {
		panic(err)
	}
	return s
}

// ValidateProtoUUID validates a protobuf UUID without conversion
func ValidateProtoUUID(pb *UUID) error {
	if pb == nil {
		return fmt.Errorf("protobuf UUID cannot be nil")
	}
	_, err := Parse(pb.GetVal())
	return err
}

// ValidateStringUUID validates a string UUID without conversion
func ValidateStringUUID(s string) error {
	_, err := Parse(s)
	return err
}

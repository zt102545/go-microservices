package utils

import (
	"github.com/zeromicro/go-zero/core/jsonx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func JsonUnmarshalClosure(jsonStr string, v any, callback func() bool) bool {
	err := jsonx.UnmarshalFromString(jsonStr, v)
	if err == nil {
		ret := callback()

		return ret
	}

	return false
}

func JsonMarshalClosure(v any, callback func(string) bool) bool {
	jsonStr, err := jsonx.MarshalToString(v)
	if err == nil {
		return callback(jsonStr)
	}

	return false
}

func TransformMessageToStruct(message proto.Message) *structpb.Struct {
	protoJson, marshalErr := jsonx.MarshalToString(message)

	if marshalErr != nil {
		return nil
	}

	pbStruct := &structpb.Struct{}

	unmarshalErr := jsonx.UnmarshalFromString(protoJson, pbStruct)

	if unmarshalErr != nil {
		return nil
	}

	return pbStruct
}

package binder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/getfrontierhq/buf-plugins/pkg/gohttp/errors"
	"github.com/getfrontierhq/buf-plugins/pkg/gohttp/option"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (d *RequestDecoder) BindBody(v interface{}) error {
	// check content type of request
	switch d.Request.Header.Get("Content-Type") {
	case option.ContentTypeApplicationJson.String():
		body, err := io.ReadAll(d.Request.Body)
		if err != nil {
			return err
		}

		return protojson.Unmarshal(body, v.(protoreflect.ProtoMessage))
	case "":
		return nil
	default:
		return fmt.Errorf("content-type is not supported, %w", errors.ErrGeneralUnsupportedMediaType)
	}
}

func (d *ResponseDecoder) BindBody(v interface{}) error {
	// check content type of request
	switch d.Response.Header.Get("Content-Type") {
	case option.ContentTypeApplicationJson.String():
		body, err := io.ReadAll(d.Response.Body)
		if err != nil {
			return err
		}

		if protoMessage, ok := v.(protoreflect.ProtoMessage); ok {
			return protojson.Unmarshal(body, protoMessage)
		} else {
			return json.Unmarshal(body, v)
		}
	case "":
		return nil
	default:
		return fmt.Errorf("content-type is not supported, %w", errors.ErrGeneralUnsupportedMediaType)
	}
}

func (e *RequestEncoder) BindBody(v interface{}) error {
	var content []byte
	switch e.Opts.ContentType {
	case option.ContentTypeApplicationJson:
		var err error
		content, err = protojson.Marshal(v.(protoreflect.ProtoMessage))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("content-type is not supported, %w", errors.ErrGeneralUnsupportedMediaType)
	}

	e.Request.Body = io.NopCloser(bytes.NewBuffer(content))
	return nil
}

func (e *ResponseEncoder) BindBody(v interface{}) error {
	var content []byte
	var err error

	if protoMessage, ok := v.(protoreflect.ProtoMessage); ok {
		content, err = protojson.Marshal(protoMessage)
	} else {
		content, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	e.ResponseWriter.Header().Set("Content-Type", option.ContentTypeApplicationJson.String())
	_, err = e.ResponseWriter.Write(content)
	return err
}

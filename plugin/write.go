// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.starlark.net/starlark"
)

func write(out *bytes.Buffer, v starlark.Value) error {
	if marshaler, ok := v.(json.Marshaler); ok {
		jsonData, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}
		out.Write(jsonData)
		return nil
	}

	switch v := v.(type) {
	case starlark.NoneType:
		out.WriteString("null")
	case starlark.Bool:
		fmt.Fprintf(out, "%t", v)
	case starlark.Int:
		out.WriteString(v.String())
	case starlark.Float:
		fmt.Fprintf(out, "%g", v)
	case starlark.String:
		s := string(v)
		if isQuoteSafe(s) {
			fmt.Fprintf(out, "%q", s)
		} else {
			data, _ := json.Marshal(s)
			out.Write(data)
		}
	case starlark.Indexable:
		out.WriteByte('[')
		for i, n := 0, starlark.Len(v); i < n; i++ {
			if i > 0 {
				out.WriteString(", ")
			}
			if err := write(out, v.Index(i)); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case *starlark.Dict:
		out.WriteByte('{')
		for i, itemPair := range v.Items() {
			key := itemPair[0]
			value := itemPair[1]
			if i > 0 {
				out.WriteString(", ")
			}
			if err := write(out, key); err != nil {
				return err
			}
			out.WriteString(": ")
			if err := write(out, value); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		return fmt.Errorf("value %s (type `%s') can't be converted to JSON", v.String(), v.Type())
	}
	return nil
}

func isQuoteSafe(s string) bool {
	for _, r := range s {
		if r < 0x20 || r >= 0x10000 {
			return false
		}
	}
	return true
}

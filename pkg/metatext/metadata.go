package metatext

import (
	"google.golang.org/grpc/metadata"
	"strings"
)

type MetadataTextMap struct {
	metadata.MD
}

func (m MetadataTextMap) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range m.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m MetadataTextMap) Set(key, val string) {
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key], val)
}


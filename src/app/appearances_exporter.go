package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	appearances "github.com/simivar/tibia-sprites-exporter/src/app/pb"
)

func ExportAppearancesJSON(catalogDir, appearancesFileName, outPath string) error {
	datPath := filepath.Join(catalogDir, appearancesFileName)
	if _, err := os.Stat(datPath); err != nil {
		return fmt.Errorf("appearances dat file not found: %w", err)
	}

	data, err := os.ReadFile(datPath)
	if err != nil {
		return fmt.Errorf("failed to read appearances dat file: %w", err)
	}
	log.Debug().Int("bytes", len(data)).Str("file", datPath).Msg("read appearances.dat")

	var root appearances.Appearances
	if err := proto.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("protobuf unmarshal failed: %w", err)
	}

	sanitizeProtoStrings(root.ProtoReflect())

	mo := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: false,
	}

	objectsJSON, err := marshalRepeatedAppearance(mo, root.GetObject())
	if err != nil {
		return err
	}
	outfitsJSON, err := marshalRepeatedAppearance(mo, root.GetOutfit())
	if err != nil {
		return err
	}
	effectsJSON, err := marshalRepeatedAppearance(mo, root.GetEffect())
	if err != nil {
		return err
	}
	missilesJSON, err := marshalRepeatedAppearance(mo, root.GetMissile())
	if err != nil {
		return err
	}

	payload := map[string]any{
		"meta": map[string]any{
			"schema_version":     1,
			"generated_at":       time.Now().UTC().Format(time.RFC3339),
			"appearances_file":   appearancesFileName,
			"appearances_sha256": filepath.Base(datPath),
		},
		"categories": map[string]any{
			"objects":  objectsJSON,
			"outfits":  outfitsJSON,
			"effects":  effectsJSON,
			"missiles": missilesJSON,
		},
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	log.Info().
		Int("objects", len(root.GetObject())).
		Int("outfits", len(root.GetOutfit())).
		Int("effects", len(root.GetEffect())).
		Int("missiles", len(root.GetMissile())).
		Str("out", outPath).
		Msg("appearances export ok")

	return nil
}

func marshalRepeatedAppearance(mo protojson.MarshalOptions, list []*appearances.Appearance) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(list))

	for _, a := range list {
		sanitizeProtoStrings(a.ProtoReflect())

		b, err := mo.Marshal(a)
		if err != nil {
			return nil, fmt.Errorf("protojson marshal failed: %w", err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, fmt.Errorf("json unmarshal failed: %w", err)
		}
		out = append(out, m)
	}
	return out, nil
}

func sanitizeProtoStrings(m protoreflect.Message) {
	if !m.IsValid() {
		return
	}

	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.Kind() == protoreflect.StringKind && !fd.IsList() && !fd.IsMap() {
			s := v.String()
			if !utf8.ValidString(s) {
				m.Set(fd, protoreflect.ValueOfString(strings.ToValidUTF8(s, "�")))
			}
			return true
		}

		switch {
		case fd.IsList():
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				elem := list.Get(i)
				if fd.Kind() == protoreflect.MessageKind {
					sanitizeProtoStrings(elem.Message())
				} else if fd.Kind() == protoreflect.StringKind {
					s := elem.String()
					if !utf8.ValidString(s) {
						list.Set(i, protoreflect.ValueOfString(strings.ToValidUTF8(s, "�")))
					}
				}
			}

		case fd.IsMap():
			mp := v.Map()
			mp.Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
				if fd.MapValue().Kind() == protoreflect.MessageKind {
					sanitizeProtoStrings(mv.Message())
				} else if fd.MapValue().Kind() == protoreflect.StringKind {
					s := mv.String()
					if !utf8.ValidString(s) {
						mp.Set(k, protoreflect.ValueOfString(strings.ToValidUTF8(s, "�")))
					}
				}
				return true
			})

		case fd.Kind() == protoreflect.MessageKind:
			sanitizeProtoStrings(v.Message())
		}

		return true
	})
}
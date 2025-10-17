// Package godynamo extracts DynamoDB annotations from proto field options
// and converts them into Go struct tag strings.
package godynamo

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	dynamopb "buf.build/gen/go/getfrontierhq/public-apis/protocolbuffers/go/dynamo"
)

// DynamoTags maps message names to field names to tag strings.
// Example: {"UserAction": {"UserId": `dynamo:"ID,hash" index:"Seq-ID-index,range"`}}
type DynamoTags map[string]map[string]string

type tagExtractor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	pgsgo.Context

	tags DynamoTags
}

func newTagExtractor(d pgs.DebuggerCommon, ctx pgsgo.Context) *tagExtractor {
	v := &tagExtractor{DebuggerCommon: d, Context: ctx}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

// VisitField extracts dynamo annotations from a proto field and builds the tag string.
func (v *tagExtractor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	msgName := v.Context.Name(f.Message()).String()

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]string{}
	}

	tagStr := buildTagsFromField(f)
	if tagStr != "" {
		fieldName := v.Context.Name(f).String()
		v.tags[msgName][fieldName] = tagStr
	}

	return v, nil
}

// Extract walks the proto file and returns all dynamo tags.
func (v *tagExtractor) Extract(f pgs.File) DynamoTags {
	v.tags = DynamoTags{}
	v.CheckErr(pgs.Walk(v, f))
	return v.tags
}

// buildTagsFromField reads field options and constructs the complete tag string.
// Example: `dynamo:"id,hash" index:"username-index,hash"`
func buildTagsFromField(f pgs.Field) string {
	var parts []string

	// Check for primary key
	if keyCfg, err := getKeyConfig(f); err == nil && keyCfg != nil {
		tagStr := buildKeyTag(keyCfg)
		if tagStr != "" {
			parts = append(parts, tagStr)
		}
	}

	// Check for GSIs (now repeated)
	if gsis, err := getGSIs(f); err == nil && len(gsis) > 0 {
		for _, gsi := range gsis {
			parts = append(parts, buildGSITag(gsi))
		}
	}

	// Check for LSIs (now repeated)
	if lsis, err := getLSIs(f); err == nil && len(lsis) > 0 {
		for _, lsi := range lsis {
			parts = append(parts, buildLSITag(lsi))
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return joinParts(parts)
}

func getKeyConfig(f pgs.Field) (*dynamopb.KeyConfig, error) {
	var cfg dynamopb.KeyConfig
	ok, err := f.Extension(dynamopb.E_Key, &cfg)
	if err != nil || !ok {
		return nil, err
	}
	return &cfg, nil
}

func getGSIs(f pgs.Field) ([]*dynamopb.IndexConfig, error) {
	var cfgs []*dynamopb.IndexConfig
	ok, err := f.Extension(dynamopb.E_Gsi, &cfgs)
	if err != nil || !ok {
		return nil, err
	}
	return cfgs, nil
}

func getLSIs(f pgs.Field) ([]*dynamopb.IndexConfig, error) {
	var cfgs []*dynamopb.IndexConfig
	ok, err := f.Extension(dynamopb.E_Lsi, &cfgs)
	if err != nil || !ok {
		return nil, err
	}
	return cfgs, nil
}

func buildKeyTag(cfg *dynamopb.KeyConfig) string {
	columnName := cfg.ColumnName

	switch cfg.Type {
	case dynamopb.KeyType_KEY_TYPE_HASH:
		return fmt.Sprintf(`dynamo:"%s,hash"`, columnName)
	case dynamopb.KeyType_KEY_TYPE_RANGE:
		return fmt.Sprintf(`dynamo:"%s,range"`, columnName)
	case dynamopb.KeyType_KEY_TYPE_UNSPECIFIED:
		// Type not specified, just output column name
		if columnName == "" {
			return ""
		}
		return fmt.Sprintf(`dynamo:"%s"`, columnName)
	default:
		return ""
	}
}

func buildGSITag(cfg *dynamopb.IndexConfig) string {
	if cfg.Name == "" {
		return ""
	}
	keyType := keyTypeString(cfg.Key)
	if keyType == "" {
		return ""
	}
	return fmt.Sprintf(`index:"%s,%s"`, cfg.Name, keyType)
}

func buildLSITag(cfg *dynamopb.IndexConfig) string {
	if cfg.Name == "" {
		return ""
	}
	keyType := keyTypeString(cfg.Key)
	if keyType == "" {
		return ""
	}
	return fmt.Sprintf(`localIndex:"%s,%s"`, cfg.Name, keyType)
}

func keyTypeString(kt dynamopb.KeyType) string {
	switch kt {
	case dynamopb.KeyType_KEY_TYPE_HASH:
		return "hash"
	case dynamopb.KeyType_KEY_TYPE_RANGE:
		return "range"
	default:
		return ""
	}
}

func joinParts(parts []string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}
	return result
}

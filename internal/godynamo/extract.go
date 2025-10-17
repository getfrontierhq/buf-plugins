// Package dynamo extracts DynamoDB annotations from proto field options
// and converts them into Go struct tag strings.
package godynamo

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	dynamopb "github.com/getfrontierhq/buf-plugins/gen/go/dynamo"
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
// Example: `dynamo:"id,hash" index:"username-index,hash" index:"email-index,range"`
func buildTagsFromField(f pgs.Field) string {
	var parts []string

	// Check for primary key annotation
	var keyStr string
	if ok, _ := f.Extension(dynamopb.E_Key, &keyStr); ok && keyStr != "" {
		tagStr := buildKeyTag(keyStr, f)
		if tagStr != "" {
			parts = append(parts, tagStr)
		}
	}

	// Check for GSI annotations (now repeated)
	var gsiStrs []string
	if ok, _ := f.Extension(dynamopb.E_Gsi, &gsiStrs); ok {
		for _, gsiStr := range gsiStrs {
			if gsiStr != "" {
				tagStr := buildIndexTag("index", gsiStr)
				if tagStr != "" {
					parts = append(parts, tagStr)
				}
			}
		}
	}

	// Check for LSI annotations (now repeated)
	var lsiStrs []string
	if ok, _ := f.Extension(dynamopb.E_Lsi, &lsiStrs); ok {
		for _, lsiStr := range lsiStrs {
			if lsiStr != "" {
				tagStr := buildIndexTag("localIndex", lsiStr)
				if tagStr != "" {
					parts = append(parts, tagStr)
				}
			}
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return joinParts(parts)
}

// buildKeyTag parses the key string annotation and builds a dynamo tag.
// Format: "column_name,key_type" or "key_type" or "column_name"
// Examples: "ID,hash", "range", "hash", "id"
func buildKeyTag(keyStr string, f pgs.Field) string {
	parts := strings.Split(keyStr, ",")

	if len(parts) == 1 {
		value := strings.TrimSpace(parts[0])

		// If it's exactly "hash" or "range", treat as key type only
		if value == "hash" || value == "range" {
			return fmt.Sprintf(`dynamo:",%s"`, value)
		}

		// Otherwise treat as column name only
		return fmt.Sprintf(`dynamo:"%s"`, value)
	}

	if len(parts) == 2 {
		// Column name and key type: "ID,hash"
		columnName := strings.TrimSpace(parts[0])
		keyType := strings.TrimSpace(parts[1])
		if keyType == "hash" || keyType == "range" {
			return fmt.Sprintf(`dynamo:"%s,%s"`, columnName, keyType)
		}
		// If no valid key type, just use column name
		return fmt.Sprintf(`dynamo:"%s"`, columnName)
	}

	return ""
}

// buildIndexTag parses an index string annotation and builds a tag.
// Format: "index_name,key_type"
// Example: "UserID-index,hash"
func buildIndexTag(tagName, indexStr string) string {
	parts := strings.Split(indexStr, ",")
	if len(parts) != 2 {
		return ""
	}

	indexName := strings.TrimSpace(parts[0])
	keyType := strings.TrimSpace(parts[1])

	if indexName == "" || (keyType != "hash" && keyType != "range") {
		return ""
	}

	return fmt.Sprintf(`%s:"%s,%s"`, tagName, indexName, keyType)
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

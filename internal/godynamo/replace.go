// Package dynamo modifies Go AST to inject DynamoDB struct tags.
// It walks the AST, finds struct type definitions, and appends
// dynamo tags to existing struct field tags.
package godynamo

import (
	"go/ast"
	"go/token"
	"strings"
)

// Retag walks the AST and injects dynamo tags into matching struct fields.
// It finds struct definitions by name, then updates field tags by appending
// the dynamo tag strings to existing tags.
func Retag(n ast.Node, tags DynamoTags) error {
	r := &retagVisitor{}

	walker := &structVisitor{
		visitor: func(n ast.Node) ast.Visitor {
			if r.err != nil {
				return nil
			}

			if ts, ok := n.(*ast.TypeSpec); ok {
				r.tags = tags[ts.Name.String()]
				return r
			}

			return nil
		},
	}

	ast.Walk(walker, n)
	return r.err
}

// structVisitor finds struct type definitions
type structVisitor struct {
	visitor func(n ast.Node) ast.Visitor
}

func (v *structVisitor) Visit(n ast.Node) ast.Visitor {
	ts, ok := n.(*ast.TypeSpec)
	if !ok {
		return v
	}

	_, ok = ts.Type.(*ast.StructType)
	if !ok {
		return v
	}

	// Found a struct, visit it with our visitor function
	ast.Walk(v.visitor(n), n)
	return nil // Don't traverse nested structs
}

// retagVisitor modifies struct field tags
type retagVisitor struct {
	err  error
	tags map[string]string
}

func (v *retagVisitor) Visit(n ast.Node) ast.Visitor {
	if v.err != nil {
		return nil
	}

	field, ok := n.(*ast.Field)
	if !ok {
		return v
	}

	// Skip fields without names (embedded fields)
	if len(field.Names) == 0 {
		return nil
	}

	fieldName := field.Names[0].String()
	newTag := v.tags[fieldName]

	// No dynamo tags for this field
	if newTag == "" {
		return nil
	}

	// Get existing tag value
	existingTag := ""
	if field.Tag != nil {
		existingTag = strings.Trim(field.Tag.Value, "`")
	}

	// Append dynamo tags to existing tags
	combinedTag := existingTag
	if combinedTag != "" {
		combinedTag += " " + newTag
	} else {
		combinedTag = newTag
	}

	// Update the tag
	field.Tag = &ast.BasicLit{
		Kind:  token.STRING,
		Value: "`" + combinedTag + "`",
	}

	return nil
}

package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed httpTemplate.tpl
var httpTemplate string

type serviceDescriptor struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	Methods     []*methodDescriptor
	MethodSets  map[string]*methodDescriptor
}

type methodDescriptor struct {
	// method
	Name         string
	OriginalName string // The parsed original name
	Num          int
	Request      string
	Reply        string
	Comment      string

	// http_rule
	Path         string
	Method       string
	Body         string
	ResponseBody string
}

func (s *serviceDescriptor) execute() string {
	s.MethodSets = make(map[string]*methodDescriptor)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return strings.Trim(buf.String(), "\r\n")
}

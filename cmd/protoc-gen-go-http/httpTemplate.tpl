{{$svcType := .ServiceType}}
{{$svcName := .ServiceName}}

const (
{{- range .MethodSets}}
  Operation_{{$svcType}}_{{.OriginalName}} = "/{{$svcName}}/{{.OriginalName}}"
{{- end}}

{{- range .MethodSets}}
  {{$svcType}}_{{.OriginalName}}_Method = "{{.Method}}"
  {{$svcType}}_{{.OriginalName}}_Path = "{{.Path}}"
{{- end}}
)

type {{$svcType}}HTTPServer interface {
{{- range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{.Name}}(ctx context.Context, in *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{$svcType}}HTTPServer(srv {{$svcType}}HTTPServer) http.Handler {
  return pot.RegisterService(&_{{$svcType}}_HTTP_ServiceDesc, srv)
}

func Register{{$svcType}}HTTPServerWithChi(srv {{$svcType}}HTTPServer, router v5.Router) http.Handler {
  return pot.RegisterServiceWithChi(&_{{$svcType}}_HTTP_ServiceDesc, srv, router)
}

{{range .Methods}}
func _{{$svcType}}_{{.Name}}{{.Num}}_HTTP_Handler(ctx context.Context, srv interface{}, dec pot.DecoderFunc, middleware pot.MiddlewareFunc) (interface{}, error) {
  in := new({{.Request}})
  if err := dec(in); err != nil {
    return nil, err
  }
  if middleware == nil {
    return srv.({{$svcType}}HTTPServer).{{.Name}}(ctx, in)
  }
  h := middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
    return srv.({{$svcType}}HTTPServer).{{.Name}}(ctx, req.(*{{.Request}}))
  })
  return h(ctx, in)
}
{{end}}

var _{{$svcType}}_HTTP_ServiceDesc = pot.ServiceDescriptor{
  ServiceName: "{{$svcName}}",
  HandlerType: (*{{$svcType}}HTTPServer)(nil),
  Methods: []pot.MethodDescriptor{
    {{- range .Methods}}
    {
      MethodName: "{{.Name}}",
      HttpMethod: "{{.Method}}",
      HttpPath: "{{.Path}}",
      Handler: _{{$svcType}}_{{.Name}}{{.Num}}_HTTP_Handler,
    },
    {{- end}}
  },
}

type {{$svcType}}HTTPClient interface {
{{- range .MethodSets}}
	{{.Name}}(ctx context.Context, in *{{.Request}}, opts ...option.BinderOption) (*{{.Reply}}, error)
{{- end}}
}

type {{$svcType}}HTTPClientImpl struct{
  baseUrl string
	client  *http.Client
}

func New{{$svcType}}HTTPClient (opts ...option.ClientOption) {{$svcType}}HTTPClient {
  options := option.NewClientOptions(opts...)
	return &{{$svcType}}HTTPClientImpl{
    baseUrl: options.BaseURL,
    client: &http.Client{
      Timeout: options.Timeout,
    },
  }
}

{{range .MethodSets}}
func (c *{{$svcType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...option.BinderOption) (*{{.Reply}}, error) {
	out := new({{.Reply}})
  url := fmt.Sprintf("%s%s", c.baseUrl, {{$svcType}}_{{.OriginalName}}_Path)
  req, err := http.NewRequest({{$svcType}}_{{.OriginalName}}_Method, url, nil)
  if err != nil {
    return nil, err
  }
  opts = append(opts, option.WithOperation(Operation_{{$svcType}}_{{.OriginalName}}))
  if err = binder.NewRequestEncoder(req, opts...).Bind(in); err != nil {
      return nil, err
  }
  req = req.WithContext(ctx)
  res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
  defer res.Body.Close()
	dec := &binder.ResponseDecoder{Response: res}
	if err := errors.ErrorMap[res.StatusCode]; err != nil {
		customErr := new(errors.Error)
    if err := dec.BindBody(customErr); err != nil {
      return nil, err
    }
		return nil, customErr
	}
	if err := dec.BindBody(out); err != nil {
    return nil, err
  }
	return out, nil
}
{{end}}

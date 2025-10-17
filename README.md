# Buf Plugins

A collection of protobuf plugins and annotations for Frontier Technologies.

## Plugins

### protoc-gen-dynamo

A protoc plugin that adds DynamoDB struct tags to generated Go protobuf code.

#### Installation

```bash
go install github.com/getfrontierhq/buf-plugins/cmd/protoc-gen-dynamo@latest
```

#### Usage

1. Add the proto dependency to your `buf.yaml`:

```yaml
deps:
  - buf.build/getfrontierhq/public-apis
```

2. Import and use the annotations in your proto files:

```protobuf
syntax = "proto3";

import "dynamo/annotations.proto";

message User {
  string id = 1 [(dynamo.key) = "ID,hash"];
  string email = 2 [
    (dynamo.key) = "email",
    (dynamo.gsi) = "email-index,hash"
  ];
  int64 created_at = 3 [(dynamo.key) = "range"];
}
```

3. Configure the plugin in your `buf.gen.yaml`:

```yaml
version: v2
plugins:
  - local: protoc-gen-dynamo
    out: gen/go
    opt:
      - paths=source_relative
      - outdir=gen/go
```

4. Generate:

```bash
buf generate
```

#### Annotations

##### (dynamo.key)

Primary table key annotation.

Format: `"column_name,key_type"` or `"key_type"` or `"column_name"`

- `key_type`: `"hash"` (partition key) or `"range"` (sort key)
- `column_name`: DynamoDB column name (optional)

Examples:
- `[(dynamo.key) = "ID,hash"]` → `` `dynamo:"ID,hash"` ``
- `[(dynamo.key) = "range"]` → `` `dynamo:",range"` ``
- `[(dynamo.key) = "email"]` → `` `dynamo:"email"` ``

##### (dynamo.gsi)

Global Secondary Index annotation (repeatable).

Format: `"index_name,key_type"`

- `index_name`: Name of the GSI (required)
- `key_type`: `"hash"` or `"range"`

Example:
- `[(dynamo.gsi) = "email-index,hash"]` → `` `index:"email-index,hash"` ``

Multiple GSIs on one field:
```protobuf
string email = 1 [
  (dynamo.gsi) = "email-index,hash",
  (dynamo.gsi) = "secondary-email-index,range"
];
```

##### (dynamo.lsi)

Local Secondary Index annotation (repeatable).

Format: Same as GSI

Example:
- `[(dynamo.lsi) = "timestamp-index,range"]` → `` `localIndex:"timestamp-index,range"` ``

## License

Apache 2.0

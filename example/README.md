# depot example

This repository contains a working example using depot's code generation facility
to generate a repository for a simple persistent struct.

# Code Elements
## Model

The file [`models/models.go`](./models/models.go) contains a single type definition
for a struct name `Message`:

```go
type Message struct {
    ID         string    `depot:"id,id"`
    Text       string    `depot:"text"`
    OrderIndex int       `depot:"order_index"`
    Length     float32   `depot:"len"`
    Attachment []byte    `depot:"attachment"`
    Created    time.Time `depot:"created"`
}
```

The fields contain tags that define the column names as well as the id field.

Based on that definition, `depot` generates a repository type. The generation
is configured by placing a  `//go:generate` comment in [`models.go`](./models/models.go).

To run the generation, execute

```
$ go generate ./models
```

The result is also part of this git repo: [`repo/gen-messagerepo.go`](./repo/gen-messagerepo.go).
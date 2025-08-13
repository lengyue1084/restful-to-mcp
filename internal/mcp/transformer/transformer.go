package transformer

type Transformer interface {
	Transform(ctx context.Context, src []*schema.Document, opts ...TransformerOption) ([]*schema.Document, error)
}
type Document struct {
	// ID 是文档的唯一标识符
	ID string
	// Content 是文档的内容
	Content string
	// MetaData 用于存储文档的元数据信息
	MetaData map[string]any
}

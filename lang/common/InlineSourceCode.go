package common

// InlineSourceCode represents source code, which may come from anywhere (stdin, file etc.)
type InlineSourceCode struct {
	AbstractSourceCode
}

func NewInlineSourceCode(text string) *InlineSourceCode {
	return &InlineSourceCode{
		AbstractSourceCode: AbstractSourceCode{
			text: text,
		},
	}
}

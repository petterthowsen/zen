package common

// FileSourceCode represents source code from a .zen file
type FileSourceCode struct {
	AbstractSourceCode
	Path string
}

func NewFileSourceCode(path string, code string) *FileSourceCode {
	return &FileSourceCode{
		AbstractSourceCode: AbstractSourceCode{text: code},
		Path:               path,
	}
}

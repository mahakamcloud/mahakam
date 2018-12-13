package writers

const (
	basePath      = "/tmp/mahakam/terraform"
	templatesPath = "templates"
)

type Writer interface {
	writeFile()
}

type AbstractWriter struct {
	Writer Writer
}

func (aw *AbstractWriter) Start() {
	aw.Writer.writeFile()
}

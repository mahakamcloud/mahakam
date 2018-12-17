package parsers

type Parser interface {
	parseTemplate() string
}

type AbstractParser struct {
	Parser Parser
}

func (self *AbstractParser) Parse() string {
	return self.Parser.parseTemplate()
}

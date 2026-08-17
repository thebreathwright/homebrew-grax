package parsers

import "github.com/ollama/ollama/api"

type GraxParser struct {
	LagunaParser
}

func (p *GraxParser) Init(tools []api.Tool, lastMessage *api.Message, thinkValue *api.ThinkValue) []api.Tool {
	return p.LagunaParser.Init(tools, lastMessage, thinkValue)
}

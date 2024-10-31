package parsing

import "zen/lang/lexing"

type Parser struct {
	tokens []lexing.Token
}

func NewParser(tokens []lexing.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() {

}

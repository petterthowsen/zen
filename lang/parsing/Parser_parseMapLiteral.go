package parsing

import (
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

// parseMapLiteral parses map literals like { "name": "john", "volume": 0.5 }
func (p *Parser) parseMapLiteral() ast.Expression {
	openBrace := p.advance() // consume {
	entries := make([]expression.MapEntry, 0)

	// Handle empty map
	if p.match(lexing.RIGHT_BRACE) {
		return expression.NewMapLiteralExpression(entries, openBrace.Location)
	}

	// Parse first entry
	entry := p.parseMapEntry()
	if entry == nil {
		return nil
	}
	entries = append(entries, *entry)

	// Parse remaining entries
	for p.match(lexing.COMMA) {
		// Allow trailing comma
		if p.check(lexing.RIGHT_BRACE) {
			break
		}

		entry = p.parseMapEntry()
		if entry == nil {
			return nil
		}
		entries = append(entries, *entry)
	}

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after map entries")
		return nil
	}

	return expression.NewMapLiteralExpression(entries, openBrace.Location)
}

// parseMapEntry parses a single key-value pair in a map literal
func (p *Parser) parseMapEntry() *expression.MapEntry {
	key := p.parseExpression()
	if key == nil {
		return nil
	}

	if !p.match(lexing.COLON) {
		p.error("Expected ':' after map key")
		return nil
	}

	value := p.parseExpression()
	if value == nil {
		return nil
	}

	return &expression.MapEntry{
		Key:   key,
		Value: value,
	}
}

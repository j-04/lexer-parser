package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func (lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex *regexp.Regexp
	handler regexHandler
}

type Lexer struct {
	patterns []regexPattern
	Tokens []Token
	source string
	pos int
}

func (lex *Lexer) advanceN(n int) {
	lex.pos += n
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) at() byte {
	return lex.source[lex.pos]
}

func (lex *Lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *Lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

func CreateLexer(source string) *Lexer {
	return &Lexer {
		pos: 0,
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`\/\/.*`), skipHandler},
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMICOLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		// advance the lexer's position past the value we just reached
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value))
	}
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.push(NewToken(NUMBER, match))
	lex.advanceN(len(match))
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	stringLiteral := lex.remainder()[match[0]:match[1]]
	lex.push(NewToken(STRING, stringLiteral))
	lex.advanceN(len(stringLiteral))
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	value := regex.FindString(lex.remainder())

	if kind, exists := reservedWords[value]; exists {
		lex.push(NewToken(kind, value))
	} else {
		lex.push(NewToken(IDENTIFIER, value))
	}

	lex.advanceN(len(value))
}

func Tokenize(source string) []Token {
	lex := CreateLexer(source)

	for !lex.at_eof() {
		matched := false

		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindStringIndex(lex.remainder())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			} 
		}

		if !matched {
			panic(fmt.Sprintf("lexer::Error -> unrecognized token near %s\n", lex.remainder()))
		}
	}
	
	lex.push(NewToken(EOF, "EOF"))
	return lex.Tokens
}
package icfg

import (
	"fmt"
	"unicode/utf8"
)

const (
	defaultTabLen = 8 // Number of spaces to make a tab.  Kinda a hack.
)

type itemType int

type item struct {
	typ itemType
	pos int
	val string
}

func (i *item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	}
	return fmt.Sprintf("%q", i.val)
}

const (
	itemError itemType = iota
	itemEOF
	itemIndent
	itemDedent
	itemStatement
	itemComment
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input, name                string
	state                      stateFn
	items                      chan item
	tabLen                     int
	start, pos, lastPos, width int
	indentStack                []int // store the last indents. Needed so we know how far we dedented
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		tabLen: defaultTabLen,
		items:  make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += w
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) run() {
	for l.state = lexSection; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func lexSection(l *lexer) stateFn {
	switch r := l.next(); {
	case isSpace(r):
		return lexIndent
	case r == eof:
		break
	}
	l.emit(itemEOF)
	return nil
}

func lexIndent(l *lexer) stateFn {
	indent := 0
Loop:
	for {
		switch l.next() {
		case ' ':
			indent += 1
		case '\t':
			indent += l.tabLen
		default:
			l.backup()
			break Loop
		}
	}

	n := l.peek()

	if isEndOfLine(n) {
		// Ignore any indent we got if we just have a newline
		return lexSection
	}
	return lexSection
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

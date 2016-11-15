package icfg

import "testing"

type lexTest struct {
	name  string
	input string
	items []item
}

var (
	tEOF = item{itemEOF, 0, ""}
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF}},
	{"spaces", "    ", []item{}},
	{"empty comment", "!", []item{}},
	{"comment", "! comment", []item{}},
	{"statement", "service unsupported-transciver", []item{}},
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		l := lex(test.name, test.input)
		var items []item
		for {
			item := l.nextItem()
			items = append(items, item)
			if item.typ == itemEOF || item.typ == itemError {
				break
			}
		}
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%+v", test.name, items, test.items)
		}
	}
}

// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"miller/parsing/token"
)

const (
	NoState    = -1
	NumStates  = 238
	NumSymbols = 374
)

type Lexer struct {
	src    []byte
	pos    int
	line   int
	column int
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:    src,
		pos:    0,
		line:   1,
		column: 1,
	}
	return lexer
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return NewLexer(src), nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = new(token.Token)
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: '"'
1: '"'
2: '0'
3: 'x'
4: '0'
5: 'b'
6: '.'
7: '-'
8: '.'
9: '.'
10: '-'
11: '.'
12: '.'
13: '-'
14: 'M'
15: '_'
16: 'P'
17: 'I'
18: 'M'
19: '_'
20: 'E'
21: 'I'
22: 'P'
23: 'S'
24: 'I'
25: 'F'
26: 'S'
27: 'I'
28: 'R'
29: 'S'
30: 'O'
31: 'P'
32: 'S'
33: 'O'
34: 'F'
35: 'S'
36: 'O'
37: 'R'
38: 'S'
39: 'N'
40: 'F'
41: 'N'
42: 'R'
43: 'F'
44: 'N'
45: 'R'
46: 'F'
47: 'I'
48: 'L'
49: 'E'
50: 'N'
51: 'A'
52: 'M'
53: 'E'
54: 'F'
55: 'I'
56: 'L'
57: 'E'
58: 'N'
59: 'U'
60: 'M'
61: 'b'
62: 'e'
63: 'g'
64: 'i'
65: 'n'
66: 'd'
67: 'o'
68: 'd'
69: 'u'
70: 'm'
71: 'p'
72: 'e'
73: 'd'
74: 'u'
75: 'm'
76: 'p'
77: 'e'
78: 'l'
79: 'i'
80: 'f'
81: 'e'
82: 'l'
83: 's'
84: 'e'
85: 'e'
86: 'm'
87: 'i'
88: 't'
89: 'e'
90: 'n'
91: 'd'
92: 'f'
93: 'i'
94: 'l'
95: 't'
96: 'e'
97: 'r'
98: 'f'
99: 'o'
100: 'r'
101: 'i'
102: 'f'
103: 'i'
104: 'n'
105: 'w'
106: 'h'
107: 'i'
108: 'l'
109: 'e'
110: 'b'
111: 'r'
112: 'e'
113: 'a'
114: 'k'
115: 'c'
116: 'o'
117: 'n'
118: 't'
119: 'i'
120: 'n'
121: 'u'
122: 'e'
123: 'f'
124: 'u'
125: 'n'
126: 'c'
127: 'r'
128: 'e'
129: 't'
130: 'u'
131: 'r'
132: 'n'
133: 'i'
134: 'n'
135: 't'
136: 'f'
137: 'l'
138: 'o'
139: 'a'
140: 't'
141: '$'
142: '$'
143: '{'
144: '}'
145: '$'
146: '*'
147: '@'
148: '@'
149: '{'
150: '}'
151: '@'
152: '*'
153: '%'
154: '%'
155: '%'
156: 'p'
157: 'a'
158: 'n'
159: 'i'
160: 'c'
161: '%'
162: '%'
163: '%'
164: ';'
165: '{'
166: '}'
167: '='
168: '$'
169: '['
170: ']'
171: '@'
172: '['
173: '|'
174: '|'
175: '='
176: '^'
177: '^'
178: '='
179: '&'
180: '&'
181: '='
182: '|'
183: '='
184: '^'
185: '='
186: '<'
187: '<'
188: '='
189: '>'
190: '>'
191: '='
192: '>'
193: '>'
194: '>'
195: '='
196: '+'
197: '='
198: '.'
199: '='
200: '-'
201: '='
202: '*'
203: '='
204: '/'
205: '='
206: '/'
207: '/'
208: '='
209: '%'
210: '='
211: '*'
212: '*'
213: '='
214: '?'
215: ':'
216: '|'
217: '|'
218: '^'
219: '^'
220: '&'
221: '&'
222: '='
223: '~'
224: '!'
225: '='
226: '~'
227: '='
228: '='
229: '!'
230: '='
231: '>'
232: '>'
233: '='
234: '<'
235: '<'
236: '='
237: '|'
238: '^'
239: '&'
240: '<'
241: '<'
242: '>'
243: '>'
244: '>'
245: '>'
246: '>'
247: '+'
248: '-'
249: '.'
250: '+'
251: '.'
252: '-'
253: '.'
254: '*'
255: '/'
256: '/'
257: '/'
258: '%'
259: '.'
260: '*'
261: '.'
262: '/'
263: '.'
264: '/'
265: '/'
266: '!'
267: '~'
268: '*'
269: '*'
270: '('
271: ')'
272: '['
273: ','
274: '_'
275: ' '
276: '!'
277: '#'
278: '$'
279: '%'
280: '&'
281: '''
282: '\'
283: '('
284: ')'
285: '*'
286: '+'
287: ','
288: '-'
289: '.'
290: '/'
291: ':'
292: ';'
293: '<'
294: '='
295: '>'
296: '?'
297: '@'
298: '['
299: ']'
300: '^'
301: '_'
302: '`'
303: '{'
304: '|'
305: '}'
306: '~'
307: '\'
308: '"'
309: 'e'
310: 'E'
311: 't'
312: 'r'
313: 'u'
314: 'e'
315: 'f'
316: 'a'
317: 'l'
318: 's'
319: 'e'
320: ' '
321: '!'
322: '#'
323: '$'
324: '%'
325: '&'
326: '''
327: '\'
328: '('
329: ')'
330: '*'
331: '+'
332: ','
333: '-'
334: '.'
335: '/'
336: ':'
337: ';'
338: '<'
339: '='
340: '>'
341: '?'
342: '@'
343: '['
344: ']'
345: '^'
346: '_'
347: '`'
348: '|'
349: '~'
350: '\'
351: '{'
352: '\'
353: '}'
354: ' '
355: '\t'
356: '\n'
357: '\r'
358: 'a'-'z'
359: 'A'-'Z'
360: '0'-'9'
361: '0'-'9'
362: 'a'-'f'
363: 'A'-'F'
364: '0'-'1'
365: 'A'-'Z'
366: 'a'-'z'
367: '0'-'9'
368: \u0100-\U0010ffff
369: 'A'-'Z'
370: 'a'-'z'
371: '0'-'9'
372: \u0100-\U0010ffff
373: .
*/

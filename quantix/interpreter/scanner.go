package interpreter

import (
	"Quantos/quantix/token"
)

func (s *Scanner) scanTokens() []token.Token {

	isAtEnd := false
	for !isAtEnd {
		_, err := s.reader.Read([]byte(s._source))
		if err != nil {
			isAtEnd = true
		}

	}

}

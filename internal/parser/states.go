package parser

import (
    "markee/internal/lexer"
)

type parseFunc func(*Parser) parseFunc



package block

// See: https://spec.commonmark.org/0.31.2/#characters-and-lines
const reNewline = regexp.MustCompile(`\n|\r|\r\n`)

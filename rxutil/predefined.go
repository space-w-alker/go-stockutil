package rxutil

import "regexp"

var Numbers = regexp.MustCompile(`-?\d+(?:\.\d+)?`)
var Integers = regexp.MustCompile(`-?\d+`)
var NonNegativeIntegers = regexp.MustCompile(`\d+`)
var Whitespace = regexp.MustCompile(`\s+`)
var LineBreak = regexp.MustCompile(`(?:\n|\r|\n\r)`)
var LineBreaks = regexp.MustCompile(`[\n\r]+`)
var LowerCaseLetter = regexp.MustCompile(`\p{Ll}`)
var UpperCaseLetter = regexp.MustCompile(`\p{Lu}`)
var Letter = regexp.MustCompile(`[\p{Lu}\p{Ll}]`)
var LowerCaseLetters = regexp.MustCompile(`\p{Ll}+`)
var UpperCaseLetters = regexp.MustCompile(`\p{Lu}+`)
var Letters = regexp.MustCompile(`[\p{Lu}\p{Ll}]+`)

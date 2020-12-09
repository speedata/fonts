package type1

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func parseCharMetrics(t *Type1, scanner *bufio.Scanner) error {
	wsSemicolonWs := regexp.MustCompile(`\s*;\s*`)
	for scanner.Scan() {
		tokens := wsSemicolonWs.Split(scanner.Text(), -1)
		if tokens[0] == "EndCharMetrics" {
			return nil
		}
		var err error
		var codepoint rune
		char := new(Char)
		char.Kernx = make(map[rune]int)
		for _, k := range tokens {
			v := strings.Split(k, " ")
			switch v[0] {
			case "C":
				char.OrigCodepoint, err = strconv.Atoi(v[1])
				if err != nil {
					return err
				}

			case "WX":
				char.Wx, err = strconv.Atoi(v[1])
				if err != nil {
					return err
				}
			case "N":
				char.Name = v[1]
				codepoint = adobeToUnicodeCodepoint[char.Name]
				char.Codepoint = codepoint
			case "B":
				a, err := strconv.Atoi(v[1])
				if err != nil {
					return err
				}
				b, err := strconv.Atoi(v[2])
				if err != nil {
					return err
				}
				c, err := strconv.Atoi(v[3])
				if err != nil {
					return err
				}
				d, err := strconv.Atoi(v[4])
				if err != nil {
					return err
				}
				char.BBox = []int{a, b, c, d}
			}
		}
		t.CharsName[char.Name] = *char
		if codepoint > 0 {
			t.CharsCodepoint[codepoint] = *char
		}
	}
	return nil
}

func parseKernData(t *Type1, scanner *bufio.Scanner) error {
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		switch tokens[0] {
		case "StartKernPairs":
			return parseKernPairs(t, scanner)
		}
	}
	return nil
}

func parseKernPairs(t *Type1, scanner *bufio.Scanner) error {

	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		switch tokens[0] {
		case "KPX":
			codpoint_dest := adobeToUnicodeCodepoint[tokens[2]]
			char := t.CharsName[tokens[1]]
			amount, err := strconv.Atoi(tokens[3])
			if err != nil {
				return err
			}
			char.Kernx[codpoint_dest] = amount
		}
	}
	return nil
}

func (t *Type1) ParseAFM(f io.Reader) error {
	var err error
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		switch tokens[0] {
		case "FontName":
			t.FontName = tokens[1]
		case "FullName":
			t.FullName = tokens[1]
		case "FamilyName":
			t.FamilyName = tokens[1]
		case "Weight":
			t.Weight = tokens[1]
		case "ItalicAngle":
			t.ItalicAngle, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "IsFixedPitch":
			t.IsFixedPitch = tokens[1] == "true"
		case "UnderlinePosition":
			t.UnderlinePosition, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "UnderlineThickness":
			t.UnderlineThickness, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "Version":
			t.Version = tokens[1]
		case "EncodingScheme":
			t.EncodingScheme = tokens[1]
		case "FontBBox":
			a, err := strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
			b, err := strconv.Atoi(tokens[2])
			if err != nil {
				return err
			}
			c, err := strconv.Atoi(tokens[3])
			if err != nil {
				return err
			}
			d, err := strconv.Atoi(tokens[4])
			if err != nil {
				return err
			}
			t.FontBBox = []int{a, b, c, d}
		case "CapHeight":
			t.UnderlineThickness, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "XHeight":
			t.XHeight, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "Descender":
			t.Descender, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "Ascender":
			t.Ascender, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
		case "StartCharMetrics":
			t.NumChars, err = strconv.Atoi(tokens[1])
			if err != nil {
				return err
			}
			t.CharsName = make(map[string]Char, t.NumChars)
			t.CharsCodepoint = make(map[rune]Char, t.NumChars)
			err := parseCharMetrics(t, scanner)
			if err != nil {
				return err
			}
		case "StartKernData":
			err = parseKernData(t, scanner)
			if err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

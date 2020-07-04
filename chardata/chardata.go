package chardata

import (
	"fmt"
	"regexp"

	"golang.org/x/text/encoding/charmap"
)

// CharmapEncodings :
// These are the encodings we will try to fix in ftfy, in the
// order that they should be tried.
var CharmapEncodings = []string{
	"latin-1",
	// "sloppy-windows-1252",
	// "sloppy-windows-1250",
	"iso-8859-2",
	// "sloppy-windows-1251",
	// "macroman",
	// "cp437",
}

// CharmapEncodingsLookup is a substitute for Python encoding lookups.
var CharmapEncodingsLookup = map[string]*charmap.Charmap{
	"latin-1":    charmap.ISO8859_1,
	"iso-8859-1": charmap.ISO8859_1,
	"iso-8859-2": charmap.ISO8859_2,
}

/* ENCODING_REGEXES contain reasonably fast ways to detect if we
could represent a given string in a given encoding. The simplest one is
the 'ascii' detector, which of course just determines if all characters
are between U+0000 and U+007F.
*/
func buildRegexes() map[string]*regexp.Regexp {
	// Define a regex that matches ASCII text.
	encodingRegexes := make(map[string]*regexp.Regexp)

	encodingRegexes["ascii"] = regexp.MustCompile(`^[\x00-\x7f]*$`)

	toDecode := make([]byte, 129)
	for i := 0; i < 128; i++ {
		toDecode[i] = byte(128 + i)
	}
	toDecode[128] = '\x1a'

	fmt.Printf("toDecode=%#v\n", toDecode)

	for _, encoding := range CharmapEncodings {
		// Make a sequence of characters that bytes \x80 to \xFF decode to
		// in each encoding, as well as byte \x1A, which is used to represent
		// the replacement character � in the sloppy-* encodings.
		// byte_range = bytes(list(range(0x80, 0x100)) + [0x1a])
		// charlist = byte_range.decode(encoding)

		cMap := CharmapEncodingsLookup[encoding]
		cDec := cMap.NewDecoder()
		decoded, err := cDec.Bytes(toDecode)
		if err != nil {
			panic(err)
		}
		regexString := fmt.Sprintf(`^[\x00-\x19\x1b-\x7f%s]*$`, string(decoded))
		encodingRegexes[encoding] = regexp.MustCompile(regexString)
	}

	return encodingRegexes
}

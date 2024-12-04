package onig

/*
#include <oniguruma.h>
*/
import "C"

// Syntax is a wrapper for Onig Syntax
//
// Each syntax defines a flavour of regex syntax.
// This type allows interaction with the built-in syntaxes through the static accessor functions
// (Syntax::emacs(), Syntax::default() etc.) and the creation of custom syntaxes.
type Syntax struct {
	ReplacerFactory RegexReplacerFactory
	raw             *C.OnigSyntaxType
}

// SyntaxPython is the Python syntax.
var SyntaxPython = &Syntax{
	ReplacerFactory: NewPythonRegexReplacer,
	raw:             C.ONIG_SYNTAX_PYTHON,
}

// SyntaxAsis is the plain text syntax.
var SyntaxAsis = &Syntax{
	raw: C.ONIG_SYNTAX_ASIS,
}

// SyntaxPosixBasic is the POSIX Basic regular expression syntax.
var SyntaxPosixBasic = &Syntax{
	raw: C.ONIG_SYNTAX_POSIX_BASIC,
}

// SyntaxPosixExtended is the POSIX Extended regular expression syntax.
var SyntaxPosixExtended = &Syntax{
	raw: C.ONIG_SYNTAX_POSIX_EXTENDED,
}

// SyntaxEmacs is the Emacs regular expression syntax.
var SyntaxEmacs = &Syntax{
	raw: C.ONIG_SYNTAX_EMACS,
}

// SyntaxGrep is the grep regular expression syntax.
var SyntaxGrep = &Syntax{
	raw: C.ONIG_SYNTAX_GREP,
}

// SyntaxGnuRegex is the GNU regex regular expression syntax.
var SyntaxGnuRegex = &Syntax{
	raw: C.ONIG_SYNTAX_GNU_REGEX,
}

// SyntaxJava is the Java (Sun java.util.regex) regular expression syntax.
var SyntaxJava = &Syntax{
	raw: C.ONIG_SYNTAX_JAVA,
}

// SyntaxPerl is the Perl regular expression syntax.
var SyntaxPerl = &Syntax{
	raw: C.ONIG_SYNTAX_PERL,
}

// SyntaxPerlNG is the Perl + named group regular expression syntax.
var SyntaxPerlNG = &Syntax{
	raw: C.ONIG_SYNTAX_PERL_NG,
}

// SyntaxRuby is the Ruby regular expression syntax.
var SyntaxRuby = &Syntax{
	ReplacerFactory: NewRubyRegexReplacer,
	raw:             C.ONIG_SYNTAX_RUBY,
}

// SyntaxOniguruma is the Oniguruma regular expression syntax.
var SyntaxOniguruma = &Syntax{
	raw: C.ONIG_SYNTAX_ONIGURUMA,
}

// SyntaxDefault is the default syntax (Ruby syntax).
var SyntaxDefault = &Syntax{
	ReplacerFactory: NewRubyRegexReplacer,
	raw:             C.ONIG_SYNTAX_RUBY,
}

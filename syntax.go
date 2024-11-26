package onig

/*
#include <oniguruma.h>
*/
import "C"

type Syntax struct {
	raw *C.OnigSyntaxType
}

var SyntaxPython = &Syntax{
	raw: C.ONIG_SYNTAX_PYTHON,
}

var SyntaxAsis = &Syntax{
	raw: C.ONIG_SYNTAX_ASIS,
}

var SyntaxPosixBasic = &Syntax{
	raw: C.ONIG_SYNTAX_POSIX_BASIC,
}

var SyntaxPosixExtended = &Syntax{
	raw: C.ONIG_SYNTAX_POSIX_EXTENDED,
}

var SyntaxEmacs = &Syntax{
	raw: C.ONIG_SYNTAX_EMACS,
}

var SyntaxGrep = &Syntax{
	raw: C.ONIG_SYNTAX_GREP,
}

var SyntaxGnuRegex = &Syntax{
	raw: C.ONIG_SYNTAX_GNU_REGEX,
}

var SyntaxJava = &Syntax{
	raw: C.ONIG_SYNTAX_JAVA,
}

var SyntaxPerl = &Syntax{
	raw: C.ONIG_SYNTAX_PERL,
}

var SyntaxPerl_NG = &Syntax{
	raw: C.ONIG_SYNTAX_PERL_NG,
}

var SyntaxRuby = &Syntax{
	raw: C.ONIG_SYNTAX_RUBY,
}

var SyntaxOniguruma = &Syntax{
	raw: C.ONIG_SYNTAX_ONIGURUMA,
}

var SyntaxDefault = &Syntax{
	raw: C.ONIG_SYNTAX_DEFAULT,
}

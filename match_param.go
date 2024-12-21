package onig

/*
#include <oniguruma.h>
*/
import "C"
import "runtime"

// MatchParam contains parameters for a Match or Search.
type MatchParam struct {
	raw *C.OnigMatchParam
}

// NewMatchParam creates a new MatchParam.
func NewMatchParam() *MatchParam {
	matchParam := &MatchParam{
		raw: C.onig_new_match_param(),
	}
	C.onig_initialize_match_param(matchParam.raw)
	runtime.SetFinalizer(matchParam, func(matchParam *MatchParam) {
		C.onig_free_match_param(matchParam.raw)
	})
	return matchParam
}

// SetMatchStackLimit sets the match stack limit.
func (p *MatchParam) SetMatchStackLimit(limit uint32) {
	C.onig_set_match_stack_limit_size_of_match_param(p.raw, C.uint(limit))
}

// SetRetryLimitInMatch sets the retry limit in match.
func (p *MatchParam) SetRetryLimitInMatch(limit uint32) {
	C.onig_set_retry_limit_in_match_of_match_param(p.raw, C.ulong(limit))
}

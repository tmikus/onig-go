#include "regex.h"

int goOnigForeachNameCallbackWrapper(
	const UChar* name,
	const UChar* nameEnd,
	int nGroupNum,
	int* groupNums,
	OnigRegex regex,
	void* arg
) {
    return goOnigForeachNameCallback((UChar*)name, (UChar*)nameEnd, nGroupNum, groupNums, regex, arg);
}

// Helper function to call onig_foreach_name
int callOnigForeachName(OnigRegex regex, void* arg) {
    return onig_foreach_name(regex, goOnigForeachNameCallbackWrapper, arg);
}

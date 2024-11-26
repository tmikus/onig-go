#include <stdlib.h>
#include <string.h>
#include <oniguruma.h>

int goOnigForeachNameCallback(
	UChar* name,
	UChar* nameEnd,
	int nGroupNum,
	int* groupNums,
	OnigRegex regex,
	void* arg
);

// Helper function to call onig_foreach_name
int callOnigForeachName(OnigRegex regex, void* arg);

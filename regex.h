#include <stdlib.h>
#include "oniguruma.h"

#ifdef __cplusplus
extern "C" {
#endif
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

    typedef struct {
        int* groupStartIndices;
        int* groupEndIndices;
        unsigned int groupCount;
    } region;

    typedef struct {
        unsigned int count;
        region** regions;
    } regionsArray;

    typedef struct {
        int result;
        region* region;
    } searchFirstResult;

    searchFirstResult searchFirstWithParam(
        regex_t* reg,
        const char* text,
        unsigned int textLen,
        unsigned int from,
        unsigned int to,
        OnigOptionType option,
        unsigned int maxStackSize,
        unsigned int retryLimitInMath
    );

    typedef struct {
        int result;
        regionsArray* array;
    } searchAllResult;

    searchAllResult searchAllWithParam(
        regex_t* reg,
        const char* text,
        unsigned int textLen,
        unsigned int from,
        unsigned int to,
        OnigOptionType option,
        unsigned int maxStackSize,
        unsigned int retryLimitInMath
    );

    void freeRegion(region* region);
    void freeRegionsArray(regionsArray* array);
    void freeRegionsArrayWithRegions(regionsArray* array);
#ifdef __cplusplus
}
#endif

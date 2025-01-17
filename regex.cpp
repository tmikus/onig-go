#include "regex.h"
#include "oniguruma.h"
#include <cstdlib>
#include <vector>

class OptionalInt {
public:
    bool hasValue;
    int value;

    OptionalInt() : hasValue(false), value(0) {}

    void setValue(int value) {
        this->hasValue = true;
        this->value = value;
    }
};

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

region* regionPtrFromOnigRegion(OnigRegion* onigRegion) {
    region* result = (region*)malloc(sizeof(region));
    result->groupCount = onigRegion->num_regs;
    result->groupStartIndices = (int*)malloc(sizeof(int) * result->groupCount);
    result->groupEndIndices = (int*)malloc(sizeof(int) * result->groupCount);
    for (int i = 0; i < result->groupCount; i++) {
        result->groupStartIndices[i] = onigRegion->beg[i];
        result->groupEndIndices[i] = onigRegion->end[i];
    }
    return result;
}

searchFirstResult searchFirstWithParam(
    regex_t* reg,
    const char* text,
    unsigned int textLen,
    unsigned int from,
    unsigned int to,
    OnigOptionType option,
    unsigned int maxStackSize,
    unsigned int retryLimitInMath
) {
    const UChar* textStart = (const UChar*)text;
    const UChar* end = textStart + textLen;
    const UChar* start = textStart + from;
    const UChar* range = textStart + to;
    OnigMatchParam* match_param = onig_new_match_param();
    onig_initialize_match_param(match_param);
    if (maxStackSize != 0) {
        onig_set_match_stack_limit_size_of_match_param(match_param, maxStackSize);
    }
    if (retryLimitInMath != 0) {
        onig_set_retry_limit_in_match_of_match_param(match_param, retryLimitInMath);
    }
    OnigRegion* onigRegion = onig_region_new();
    searchFirstResult result;
    result.result = onig_search_with_param(
        reg,
        textStart,
        end,
        start,
        range,
        onigRegion,
        option,
        match_param
    );
    result.region = regionPtrFromOnigRegion(onigRegion);
    onig_region_free(onigRegion, 1);
    onig_free_match_param(match_param);
    free((void*)text);
    return result;
}

regionsArray* regionsArrayFromVector(std::vector<region*>& regions) {
    regionsArray* result = (regionsArray*)malloc(sizeof(regionsArray));
    result->count = regions.size();
    result->regions = (region**)malloc(sizeof(region*) * result->count);
    for (int i = 0; i < result->count; i++) {
        result->regions[i] = regions[i];
    }
    return result;
}

searchAllResult searchAllWithParam(
    regex_t* reg,
    const char* text,
    unsigned int textLen,
    unsigned int from,
    unsigned int to,
    OnigOptionType option,
    unsigned int maxStackSize,
    unsigned int retryLimitInMath
) {
    std::vector<region*> regions;
    OnigMatchParam* match_param = onig_new_match_param();
    onig_initialize_match_param(match_param);
    if (maxStackSize != 0) {
        onig_set_match_stack_limit_size_of_match_param(match_param, maxStackSize);
    }
    if (retryLimitInMath != 0) {
        onig_set_retry_limit_in_match_of_match_param(match_param, retryLimitInMath);
    }

    OnigRegion* onigRegion = onig_region_new();
    int onigSearchResult = 0;
    int lastEnd = 0;
    OptionalInt lastMatchEnd;
    while (lastEnd <= textLen) {
        const UChar* begin = (const UChar*)text;
        const UChar* end = begin + textLen;
        const UChar* limitStart = begin + lastEnd;
        const UChar* limitRange = begin + textLen;
        onigSearchResult = onig_search_with_param(
            reg,
            begin,
            end,
            limitStart,
            limitRange,
            onigRegion,
            option,
            match_param
        );
        if (onigSearchResult < 0) {
            break;
        }
        int posFrom = 0;
        int posTo = 0;
        if (onigRegion->beg != NULL && onigRegion->end != NULL && onigRegion->num_regs > 0) {
            posFrom = *(onigRegion->beg);
            posTo = *(onigRegion->end);
        } else {
            break;
        }
        // Don't accept empty matches immediately following the last match.
        // i.e., no infinite loops please.
        if (posFrom == posTo && lastMatchEnd.hasValue && lastMatchEnd.value == posTo) {
            // In Go we used to have all this stuff ... is it relevant to C++?
            // offset := 1
            // if lastEnd < textLength-1 {
            //     offset = len(c.text[lastEnd : lastEnd+1])
            // }
            lastEnd += 1;
            continue;
        } else {
            lastEnd = posTo;
            lastMatchEnd.setValue(posTo);
        }
        regions.push_back(regionPtrFromOnigRegion(onigRegion));
    }
    if (onigSearchResult == ONIG_MISMATCH) {
        onigSearchResult = 0;
    }
    onig_region_free(onigRegion, 1);
    onig_free_match_param(match_param);
    free((void*)text);
    searchAllResult result;
    result.result = onigSearchResult;
    result.array = regionsArrayFromVector(regions);
    if (onigSearchResult < 0) {
        freeRegionsArray(result.array);
        result.array = NULL;
    }
    return result;
}

void freeRegion(region* region) {
    if (region->groupStartIndices != NULL) {
        free(region->groupStartIndices);
    }
    if (region->groupEndIndices != NULL) {
        free(region->groupEndIndices);
    }
}

void freeRegionsArray(regionsArray* array) {
    free(array->regions);
    free(array);
}

void freeRegionsArrayWithRegions(regionsArray* array) {
    for (int i = 0; i < array->count; i++) {
        freeRegion(array->regions[i]);
    }
    freeRegionsArray(array);
}

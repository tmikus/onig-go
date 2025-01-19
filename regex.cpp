#include "regex.h"
#include "oniguruma.h"
#include <cstdlib>
#include <cstring>
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

typedef struct {
    int currentIndex;
    groupNamesArray* result;
} appendGroupNameToArrayState;

int appendGroupNameToArray(
	const UChar* name,
	const UChar* nameEnd,
	int nGroupNum,
	int* groupNums,
	OnigRegex regex,
	void* arg
) {
    appendGroupNameToArrayState* state = (appendGroupNameToArrayState*)arg;
    groupName* current = &state->result->names[state->currentIndex];
    current->nameLength = nameEnd - name;
    if (current->nameLength > 0) {
        current->name = (char*)calloc(current->nameLength, sizeof(char));
        strncpy(current->name, (const char*)name, current->nameLength);
    } else {
        current->name = NULL;
    }
    current->indicesCount = nGroupNum;
    if (current->indicesCount > 0) {
        current->indices = (int*)calloc(current->indicesCount, sizeof(int));
        memcpy(current->indices, groupNums, current->indicesCount);
    } else {
        current->indices = NULL;
    }
    state->currentIndex++;
    return 0;
}

groupNamesArray* readGroupNames(OnigRegex regex) {
    groupNamesArray* result = (groupNamesArray*)calloc(1, sizeof(groupNamesArray));
    result->count = onig_number_of_names(regex);
    if (result->count > 0) {
        result->names = (groupName*)calloc(result->count, sizeof(groupName));
    } else {
        result->names = NULL;
    }
    appendGroupNameToArrayState state;
    state.currentIndex = 0;
    state.result = result;
    onig_foreach_name(regex, appendGroupNameToArray, &state);
    return result;
}

newRegexResult newRegex(
    const char* text,
    unsigned int textLen,
    OnigOptionType options,
    OnigSyntaxType* syntax
) {
    newRegexResult result;
    UChar* patternStart = (UChar*)text;
    UChar* patternEnd = patternStart + textLen;
    OnigErrorInfo errorInfo;
    result.result = onig_new(
        &result.regex,
        patternStart,
        patternEnd,
        options,
        ONIG_ENCODING_UTF8,
        syntax,
        &errorInfo
    );
    if (result.result == ONIG_NORMAL) {
        result.groupNames = readGroupNames(result.regex);
    }
    free((void*)text);
    return result;
}

region* regionPtrFromOnigRegion(OnigRegion* onigRegion) {
    region* result = (region*)calloc(1, sizeof(region));
    result->groupCount = onigRegion->num_regs;
    result->groupStartIndices = (int*)calloc(result->groupCount, sizeof(int));
    result->groupEndIndices = (int*)calloc(result->groupCount, sizeof(int));
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
    regionsArray* result = (regionsArray*)calloc(1, sizeof(regionsArray));
    result->count = regions.size();
    result->regions = (region**)calloc(result->count, sizeof(region*));
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

void freeGroupNamesArray(groupNamesArray* array) {
    if (array == NULL) {
        return;
    }
    if (array->names != NULL) {
        for (int i = 0; i < array->count; i++) {
            groupName* name = &array->names[i];
            if (name->indices != NULL) {
                free(name->indices);
            }
            if (name->name != NULL) {
                free(name->name);
            }
        }
        free(array->names);
    }
    free(array);
}

void freeRegion(region* region) {
    if (region == NULL) {
        return;
    }
    if (region->groupStartIndices != NULL) {
        free(region->groupStartIndices);
    }
    if (region->groupEndIndices != NULL) {
        free(region->groupEndIndices);
    }
    free(region);
}

void freeRegionsArray(regionsArray* array) {
    if (array == NULL) {
        return;
    }
    if (array->regions != NULL) {
        free(array->regions);
    }
    free(array);
}

void freeRegionsArrayWithRegions(regionsArray* array) {
    if (array == NULL) {
        return;
    }
    if (array->regions != NULL) {
        for (int i = 0; i < array->count; i++) {
            freeRegion(array->regions[i]);
        }
    }
    freeRegionsArray(array);
}

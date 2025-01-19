#include <stdlib.h>
#include "oniguruma.h"

#ifdef __cplusplus
extern "C" {
#endif
    typedef struct {
        char* name;
        int nameLength;
        int* indices;
        unsigned int indicesCount;
    } groupName;

    typedef struct {
        unsigned int count;
        groupName* names;
    } groupNamesArray;

    typedef struct {
        OnigRegex regex;
        int result;
        groupNamesArray* groupNames;
    } newRegexResult;

    newRegexResult newRegex(
        const char* text,
        unsigned int textLen,
        OnigOptionType type,
        OnigSyntaxType* syntax
    );

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

    void freeGroupNamesArray(groupNamesArray* array);
    void freeRegion(region* region);
    void freeRegionsArray(regionsArray* array);
    void freeRegionsArrayWithRegions(regionsArray* array);
#ifdef __cplusplus
}
#endif

#include <stdio.h>
#include <string.h>

#include "memory.h"
#include "object.h"
#include "table.h"
#include "value.h"
#include "vm.h"

#define ALLOCATE_OBJ(type, objectType) \
    (type*) allocateObject(sizeof(type), objectType)

static Obj* allocateObject(size_t size, ObjType type) {
    Obj* object = (Obj*) reallocate(NULL, 0, size);
    object->type = type;

    object->next = vm.objects;
    vm.objects = object;

    return object;
}

ObjFunction* newFunction() {
    ObjFunction* function = ALLOCATE_OBJ(ObjFunction, OBJ_FUNCTION);

    function->arity = 0;
    function->name = NULL;
    initChunk(&function->chunk);
    return function;
}

ObjNative* newNative(NativeFn function) {
    ObjNative* native = ALLOCATE_OBJ(ObjNative, OBJ_NATIVE);
    native->function = function;
    return native;
}

static ObjString* allocateString(char* chars, int length, uint32_t hash) {
    ObjString* string = ALLOCATE_OBJ(ObjString, OBJ_STRING);
    string->length = length;
    string->chars = chars;
    string->hash = hash;

    tableSet(&vm.strings, string, NIL_VAL);

    return string;
}

ObjArray* allocateArray() {
    ObjArray* array = ALLOCATE_OBJ(ObjArray, OBJ_ARRAY);
    initValueArray(&array->elements);

    return array;
}

static uint32_t hashString(const char* key, int length) {
    uint32_t hash = 2166136261u;

    for (int i = 0; i < length; i++) {
        hash ^= key[i];
        hash *= 16777619;
    }

    return hash;
}

ObjString* takeString(char* chars, int length) {
    uint32_t hash = hashString(chars, length);
    ObjString* interned = tableFindString(&vm.strings, chars, length, hash);

    if (interned != NULL) {
        FREE_ARRAY(char, chars, length);
        return interned;
    }

    return allocateString(chars, length, hash);
}

ObjString* copyString(const char* chars, int length) {
    uint32_t hash = hashString(chars, length);
    ObjString* interned = tableFindString(&vm.strings, chars, length, hash);

    if (interned != NULL) return interned;

    char* headChars = ALLOCATE(char, length + 1);
    memcpy(headChars, chars, length);
    headChars[length] = '\0';

    return allocateString(headChars, length, hash);
}

ObjArray* writeArray(Value* objArrayValue, int index, Value value) {
    ObjArray* array = (IS_ARRAY(*objArrayValue))
            ? AS_ARRAY(*objArrayValue) : allocateArray();
    writeValueArray(&array->elements, index, value);

    return array;
}

void printObject(Value value) {
    switch (OBJ_TYPE(value)) {
        case OBJ_FUNCTION:
            if (AS_FUNCTION(value)->name == NULL) {
                printf("<script>");
                break;
            }
            printf("<fn %s>", AS_FUNCTION(value)->name->chars);
            break;
        case OBJ_NATIVE:
            printf("<native fn>");
            break;
        case OBJ_STRING:
            printf("%s", AS_CSTRING(value));
            break;
        case OBJ_ARRAY:
            printf("%d", AS_ARRAY_COUNT(value));
            break;
    }
}

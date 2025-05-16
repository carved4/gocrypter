#pragma once
#include <windows.h>

#ifdef __cplusplus
extern "C" {
#endif

int RunPE32(wchar_t* targetPath, unsigned char* payload, size_t payload_size);

#ifdef __cplusplus
}
#endif 
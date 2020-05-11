/**
 * Part of Wingows - Win32 API layer for Go
 * https://github.com/rodrigocfd/wingows
 * Copyright 2020-present Rodrigo Cesar de Freitas Dias
 * This library is released under the MIT license
 */

package proc

import (
	"syscall"
)

var (
	dllKernel32 = syscall.NewLazyDLL("kernel32.dll")

	MulDiv              = dllKernel32.NewProc("MulDiv")
	GetModuleHandle     = dllKernel32.NewProc("GetModuleHandleW")
	VerifyVersionInfo   = dllKernel32.NewProc("VerifyVersionInfoW")
	VerSetConditionMask = dllKernel32.NewProc("VerSetConditionMask")
)

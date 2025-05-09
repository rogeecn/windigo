//go:build windows

package shell

import (
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/internal/dll"
	"github.com/rodrigocfd/windigo/internal/vt"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/ole"
	"github.com/rodrigocfd/windigo/win/wstr"
)

// [SHCreateItemFromIDList] function.
//
// Return type is tipically [IShellItem] of [IShellItem2].
//
// # Example
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	item, _ := shell.SHCreateItemFromParsingName[shell.IShellItem](
//		rel, "C:\\Temp\\foo.txt")
//
//	pidl, _ := shell.SHGetIDListFromObject(item)
//	defer pidl.Free()
//
//	sameItem, _ := shell.SHCreateItemFromIDList[shell.IShellItem2](
//		rel, pidl)
//
// [SHCreateItemFromIDList]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-shcreateitemfromidlist
// [IShellItem]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitem
// [IShellItem2]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitem2
func SHCreateItemFromIDList[T any, P ole.ComCtor[T]](
	releaser *ole.Releaser,
	pidl ITEMIDLIST,
) (*T, error) {
	pObj := P(new(T)) // https://stackoverflow.com/a/69575720/6923555
	var ppvtQueried **vt.IUnknown
	guidIid := win.GuidFrom(pObj.IID())

	ret, _, _ := syscall.SyscallN(_SHCreateItemFromIDList.Addr(),
		uintptr(pidl), uintptr(unsafe.Pointer(&guidIid)),
		uintptr(unsafe.Pointer(&ppvtQueried)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		pObj.Set(ppvtQueried)
		releaser.Add(pObj)
		return pObj, nil
	} else {
		return nil, hr
	}
}

var _SHCreateItemFromIDList = dll.Shell32.NewProc("SHCreateItemFromIDList")

// [SHCreateItemFromParsingName] function.
//
// Return type is tipically [IShellItem] of [IShellItem2].
//
// # Example
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	ish, _ := shell.SHCreateItemFromParsingName[shell.IShellItem](
//		rel, "C:\\Temp\\foo.txt")
//
// [SHCreateItemFromParsingName]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-shcreateitemfromparsingname
// [IShellItem]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitem
// [IShellItem2]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nn-shobjidl_core-ishellitem2
func SHCreateItemFromParsingName[T any, P ole.ComCtor[T]](
	releaser *ole.Releaser,
	folderOrFilePath string,
) (*T, error) {
	pObj := P(new(T)) // https://stackoverflow.com/a/69575720/6923555
	var ppvtQueried **vt.IUnknown
	riidGuid := win.GuidFrom(pObj.IID())

	folderOrFilePath16 := wstr.NewBufWith[wstr.Stack20](folderOrFilePath, wstr.EMPTY_IS_NIL)

	ret, _, _ := syscall.SyscallN(_SHCreateItemFromParsingName.Addr(),
		uintptr(folderOrFilePath16.UnsafePtr()),
		0, uintptr(unsafe.Pointer(&riidGuid)),
		uintptr(unsafe.Pointer(&ppvtQueried)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		pObj.Set(ppvtQueried)
		releaser.Add(pObj)
		return pObj, nil
	} else {
		return nil, hr
	}
}

var _SHCreateItemFromParsingName = dll.Shell32.NewProc("SHCreateItemFromParsingName")

// [SHGetDesktopFolder] function.
//
// # Example
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	folder, _ := shell.SHGetDesktopFolder(rel)
//
// [SHGetDesktopFolder]: https://learn.microsoft.com/en-us/windows/win32/api/shlobj_core/nf-shlobj_core-shgetdesktopfolder
func SHGetDesktopFolder(releaser *ole.Releaser) (*IShellFolder, error) {
	var ppvtQueried **vt.IUnknown
	ret, _, _ := syscall.SyscallN(_SHGetDesktopFolder.Addr(),
		uintptr(unsafe.Pointer(&ppvtQueried)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		pObj := vt.NewObj[IShellFolder](ppvtQueried)
		releaser.Add(pObj)
		return pObj, nil
	} else {
		return nil, hr
	}
}

var _SHGetDesktopFolder = dll.Shell32.NewProc("SHGetDesktopFolder")

// [SHGetIDListFromObject] function.
//
// ⚠️ You must defer ITEMIDLIST.Free().
//
// # Example
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	item, _ := shell.SHCreateItemFromParsingName[shell.IShellItem](
//		rel, "C:\\Temp\\foo.txt")
//
//	pidl, _ := shell.SHGetIDListFromObject(item)
//	defer pidl.Free()
//
// [SHGetIDListFromObject]: https://learn.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-shgetidlistfromobject
func SHGetIDListFromObject(obj ole.ComPtr) (ITEMIDLIST, error) {
	var pidl ITEMIDLIST
	ret, _, _ := syscall.SyscallN(_SHGetIDListFromObject.Addr(),
		uintptr(unsafe.Pointer(obj.Ppvt())), uintptr(unsafe.Pointer(&pidl)))
	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		return pidl, nil
	} else {
		return ITEMIDLIST(0), hr
	}
}

var _SHGetIDListFromObject = dll.Shell32.NewProc("SHGetIDListFromObject")

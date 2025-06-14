//go:build windows

package oleaut

import (
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/internal/utl"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/ole"
	"github.com/rodrigocfd/windigo/win/wstr"
)

// [OleLoadPicture] function.
//
// Pass size = 0 to read all the bytes from the stream.
//
// The bytes are copied, so [ole.IStream] can be released after this function
// returns.
//
// # Example
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	data := []byte{0x10, 0x11, 0x12}
//	defer runtime.KeepAlive(data)
//
//	stream, _ := ole.SHCreateMemStream(rel, data)
//	pic, _ := oleaut.OleLoadPicture(rel, &stream, 0, true)
//
// [OleLoadPicture]: https://learn.microsoft.com/en-us/windows/win32/api/olectl/nf-olectl-oleloadpicture
func OleLoadPicture(
	releaser *ole.Releaser,
	stream *ole.IStream,
	size uint,
	keepOriginalFormat bool,
) (*IPicture, error) {
	var ppvtQueried **ole.IUnknownVt
	guid := win.GuidFrom(co.IID_IPicture)

	ret, _, _ := syscall.SyscallN(dllOleaut(_PROC_OleLoadPicture),
		uintptr(unsafe.Pointer(stream.Ppvt())),
		uintptr(int32(size)),
		utl.BoolToUintptr(!keepOriginalFormat), // note: reversed
		uintptr(unsafe.Pointer(&guid)),
		uintptr(unsafe.Pointer(&ppvtQueried)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		var pObj *IPicture
		utl.ComCreateObj(&pObj, unsafe.Pointer(ppvtQueried))
		releaser.Add(pObj)
		return pObj, nil
	} else {
		return nil, hr
	}
}

// [OleLoadPicturePath] function.
//
// The picture must be in the following formats:
//   - BMP (bitmap)
//   - JPEG
//   - WMF (metafile)
//   - ICO (icon)
//   - GIF
//
// [OleLoadPicturePath]: https://learn.microsoft.com/en-us/windows/win32/api/olectl/nf-olectl-oleloadpicturepath
func OleLoadPicturePath(
	releaser *ole.Releaser,
	path string,
	transparentColor win.COLORREF,
) (*IPicture, error) {
	path16 := wstr.NewBufWith[wstr.Stack20](path, wstr.EMPTY_IS_NIL)
	var ppvtQueried **ole.IUnknownVt
	guid := win.GuidFrom(co.IID_IPicture)

	ret, _, _ := syscall.SyscallN(dllOleaut(_PROC_OleLoadPicturePath),
		uintptr(path16.UnsafePtr()),
		0, 0,
		uintptr(transparentColor),
		uintptr(unsafe.Pointer(&guid)),
		uintptr(unsafe.Pointer(&ppvtQueried)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		var pObj *IPicture
		utl.ComCreateObj(&pObj, unsafe.Pointer(ppvtQueried))
		releaser.Add(pObj)
		return pObj, nil
	} else {
		return nil, hr
	}
}

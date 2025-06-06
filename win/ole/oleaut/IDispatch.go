//go:build windows

package oleaut

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/rodrigocfd/windigo/internal/utl"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
	"github.com/rodrigocfd/windigo/win/ole"
	"github.com/rodrigocfd/windigo/win/wstr"
)

// [IDispatch] COM interface.
//
// Implements [ole.ComObj] and [ole.ComResource].
//
// [IDispatch]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nn-oaidl-idispatch
type IDispatch struct{ ole.IUnknown }

// Returns the unique COM [interface ID].
//
// [interface ID]: https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/iid
func (*IDispatch) IID() co.IID {
	return co.IID_IDispatch
}

// [GetIDsOfNames] method.
//
// [GetIDsOfNames]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-getidsofnames
func (me *IDispatch) GetIDsOfNames(
	lcid win.LCID,
	member string,
	parameters ...string,
) ([]MEMBERID, error) {
	nParams := uint(1 + len(parameters)) // member + parameters
	nullGuid := win.GuidFrom(co.IID_NULL)
	memberIds := make([]MEMBERID, nParams) // will be returned

	allStrs16 := wstr.NewArray()
	allStrs16.Append(member)
	allStrs16.Append(parameters...)

	strPtrs := make([]*uint16, 0, nParams)
	for i := uint(0); i < nParams; i++ {
		strPtrs = append(strPtrs, allStrs16.PtrOf(i))
	}

	ret, _, _ := syscall.SyscallN(
		(*_IDispatchVt)(unsafe.Pointer(*me.Ppvt())).GetIDsOfNames,
		uintptr(unsafe.Pointer(me.Ppvt())),
		uintptr(unsafe.Pointer(&nullGuid)),
		uintptr(unsafe.Pointer(&strPtrs[0])),
		uintptr(uint32(nParams)),
		uintptr(lcid),
		uintptr(unsafe.Pointer(&memberIds[0])))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		return memberIds, nil
	} else {
		return nil, hr
	}
}

// [GetTypeInfo] method.
//
// # Example
//
//	var iDisp oleaut.IDispatch // initialized somewhere
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	nfo, _ := iDisp.GetTypeInfo(rel, win.LCID_USER_DEFAULT)
//
// [GetTypeInfo]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-gettypeinfo
func (me *IDispatch) GetTypeInfo(releaser *ole.Releaser, lcid win.LCID) (*ITypeInfo, error) {
	var ppvtQueried **ole.IUnknownVt
	ret, _, _ := syscall.SyscallN(
		(*_IDispatchVt)(unsafe.Pointer(*me.Ppvt())).GetTypeInfo,
		uintptr(unsafe.Pointer(me.Ppvt())),
		0,
		uintptr(lcid),
		uintptr(unsafe.Pointer(&ppvtQueried)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		var pObj *ITypeInfo
		utl.ComCreateObj(&pObj, unsafe.Pointer(ppvtQueried))
		releaser.Add(pObj)
		return pObj, nil
	} else {
		return nil, hr
	}
}

// [GetTypeInfoCount] method.
//
// If the object provides type information, this number is 1; otherwise the
// number is 0.
//
// [GetTypeInfoCount]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-gettypeinfocount
func (me *IDispatch) GetTypeInfoCount() (uint, error) {
	var pctInfo uint32
	ret, _, _ := syscall.SyscallN(
		(*_IDispatchVt)(unsafe.Pointer(*me.Ppvt())).GetTypeInfoCount,
		uintptr(unsafe.Pointer(me.Ppvt())),
		uintptr(unsafe.Pointer(&pctInfo)))

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		return uint(pctInfo), nil
	} else {
		return 0, hr
	}
}

// [Invoke] method.
//
// This is a low-level method, prefer using [IDispatch.InvokeGet],
// [IDispatch.InvokeMethod] or [IDispatch.InvokePut].
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
func (me *IDispatch) Invoke(
	releaser *ole.Releaser,
	dispIdMember MEMBERID,
	lcid win.LCID,
	flags co.DISPATCH,
	dispParams *DISPPARAMS,
) (*VARIANT, error) {
	var remoteErr _EXCEPINFO // in case of remote error, will be converted to *EXCEPINFO
	defer remoteErr.Free()

	remoteResult := NewVariantEmpty(releaser) // result returned from the remote call
	nullGuid := win.GuidFrom(co.IID_NULL)

	ret, _, _ := syscall.SyscallN(
		(*_IDispatchVt)(unsafe.Pointer(*me.Ppvt())).Invoke,
		uintptr(unsafe.Pointer(me.Ppvt())),
		uintptr(dispIdMember),
		uintptr(unsafe.Pointer(&nullGuid)),
		uintptr(lcid),
		uintptr(flags),
		uintptr(unsafe.Pointer(dispParams)),
		uintptr(unsafe.Pointer(remoteResult)),
		uintptr(unsafe.Pointer(&remoteErr)),
		0) // puArgErr is not retrieved

	if hr := co.HRESULT(ret); hr == co.HRESULT_S_OK {
		return remoteResult, nil
	} else if hr == co.HRESULT_DISP_E_EXCEPTION {
		return nil, remoteErr.Serialize()
	} else {
		return nil, hr
	}
}

// Calls [Invoke] with [co.DISPATCH_PROPERTYGET].
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// Parameters must be one of the valid [VARIANT] types.
//
// # Example
//
//	ole.CoInitializeEx(co.COINIT_APARTMENTTHREADED | co.COINIT_DISABLE_OLE1DDE)
//	defer ole.CoUninitialize()
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	clsId, _ := ole.CLSIDFromProgID("Excel.Application")
//
//	var dispExcel *oleaut.IDispatch
//	ole.CoCreateInstance(
//		rel, clsId, nil, co.CLSCTX_LOCAL_SERVER, &dispExcel)
//
//	varBooks, _ := dispExcel.InvokeGet(rel, "Workbooks")
//	dispBooks, _ := varBooks.IDispatch(rel)
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
// [EXCEPINFO]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-excepinfo
func (me *IDispatch) InvokeGet(
	releaser *ole.Releaser,
	propertyName string,
	params ...interface{},
) (*VARIANT, error) {
	return me.rawInvoke(releaser, co.DISPATCH_PROPERTYGET, propertyName, params...)
}

// Calls [Invoke] with [co.DISPATCH_PROPERTYGET], and tries to convert the
// [VARIANT] result to an [IDispatch] object.
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// Parameters must be one of the valid [VARIANT] types.
//
// # Example
//
//	ole.CoInitializeEx(co.COINIT_APARTMENTTHREADED | co.COINIT_DISABLE_OLE1DDE)
//	defer ole.CoUninitialize()
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	clsId, _ := ole.CLSIDFromProgID("Excel.Application")
//
//	var dispExcel *oleaut.IDispatch
//	ole.CoCreateInstance(rel, clsId, nil, co.CLSCTX_LOCAL_SERVER, &dispExcel)
//
//	books, _ := dispExcel.InvokeGetIDispatch(rel, "Workbooks")
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
// [EXCEPINFO]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-excepinfo
func (me *IDispatch) InvokeGetIDispatch(
	releaser *ole.Releaser,
	propertyName string,
	params ...interface{},
) (*IDispatch, error) {
	variant, err := me.InvokeGet(releaser, propertyName, params...)
	if err != nil {
		return nil, err
	}
	if idisp, ok := variant.IDispatch(releaser); ok {
		return idisp, nil
	} else {
		return nil, fmt.Errorf("InvokeGet \"%s\" didn't return an IDispatch object", propertyName)
	}
}

// Calls [Invoke] with [co.DISPATCH_METHOD].
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// Parameters must be one of the valid [VARIANT] types.
//
// # Example
//
//	ole.CoInitializeEx(co.COINIT_APARTMENTTHREADED | co.COINIT_DISABLE_OLE1DDE)
//	defer ole.CoUninitialize()
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	clsId, _ := ole.CLSIDFromProgID("Excel.Application")
//
//	var dispExcel *oleaut.IDispatch
//	ole.CoCreateInstance(rel, clsId, nil, co.CLSCTX_LOCAL_SERVER, &dispExcel)
//
//	varBooks, _ := dispExcel.InvokeGet(rel, "Workbooks")
//	dispBooks, _ := varBooks.IDispatch(rel)
//
//	varFile, _ := dispBooks.InvokeMethod(rel, "Open", "C:\\Temp\\file.xlsx")
//	dispFile, _ := varFile.IDispatch(rel)
//
//	dispFile.InvokeMethod(rel, "SaveAs", "C:\\Temp\\copy.xlsx")
//	dispFile.InvokeMethod(rel, "Close")
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
// [EXCEPINFO]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-excepinfo
func (me *IDispatch) InvokeMethod(
	releaser *ole.Releaser,
	methodName string,
	params ...interface{},
) (*VARIANT, error) {
	return me.rawInvoke(releaser, co.DISPATCH_METHOD, methodName, params...)
}

// Calls [Invoke] with [co.DISPATCH_METHOD], and tries to convert the
// [VARIANT] result to an [IDispatch] object.
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// Parameters must be one of the valid [VARIANT] types.
//
// # Example
//
//	ole.CoInitializeEx(co.COINIT_APARTMENTTHREADED | co.COINIT_DISABLE_OLE1DDE)
//	defer ole.CoUninitialize()
//
//	rel := ole.NewReleaser()
//	defer rel.Release()
//
//	clsId, _ := ole.CLSIDFromProgID("Excel.Application")
//
//	var excel *oleaut.IDispatch
//	ole.CoCreateInstance(rel, clsId, nil, co.CLSCTX_LOCAL_SERVER, &excel)
//
//	books, _ := excel.InvokeGetIDispatch(rel, "Workbooks")
//	file, _ := books.InvokeMethodIDispatch(rel, "Open", "C:\\Temp\\file.xlsx")
//	file.InvokeMethod(rel, "SaveAs", "C:\\Temp\\copy.xlsx")
//	file.InvokeMethod(rel, "Close")
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
// [EXCEPINFO]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-excepinfo
func (me *IDispatch) InvokeMethodIDispatch(
	releaser *ole.Releaser,
	methodName string,
	params ...interface{},
) (*IDispatch, error) {
	variant, err := me.InvokeMethod(releaser, methodName, params...)
	if err != nil {
		return nil, err
	}
	if idisp, ok := variant.IDispatch(releaser); ok {
		return idisp, nil
	} else {
		return nil, fmt.Errorf("InvokeMethod \"%s\" didn't return an IDispatch object", methodName)
	}
}

// Calls [Invoke] with [co.DISPATCH_PROPERTYPUT].
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// Parameter must be one of the valid [VARIANT] types.
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
// [EXCEPINFO]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-excepinfo
func (me *IDispatch) InvokePut(
	releaser *ole.Releaser,
	propertyName string,
	value interface{},
) (*VARIANT, error) {
	return me.rawInvoke(releaser, co.DISPATCH_PROPERTYPUT, propertyName, value)
}

// Calls [Invoke] with [co.DISPATCH_PROPERTYPUT], and tries to convert the
// [VARIANT] result to an [IDispatch] object.
//
// If the remote call raises an exception, the returned error will be an
// instance of *[EXCEPINFO].
//
// Parameter must be one of the valid [VARIANT] types.
//
// [Invoke]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/nf-oaidl-idispatch-invoke
// [EXCEPINFO]: https://learn.microsoft.com/en-us/windows/win32/api/oaidl/ns-oaidl-excepinfo
func (me *IDispatch) InvokePutIDispatch(
	releaser *ole.Releaser,
	propertyName string,
	value interface{},
) (*IDispatch, error) {
	variant, err := me.InvokePut(releaser, propertyName, value)
	if err != nil {
		return nil, err
	}
	if idisp, ok := variant.IDispatch(releaser); ok {
		return idisp, nil
	} else {
		return nil, fmt.Errorf("InvokePut \"%s\" didn't return an IDispatch object", propertyName)
	}
}

func (me *IDispatch) rawInvoke(
	releaser *ole.Releaser,
	method co.DISPATCH,
	methodName string,
	params ...interface{},
) (*VARIANT, error) {
	memberIds, err := me.GetIDsOfNames(win.LCID_USER_DEFAULT, methodName) // will return 1 element
	if err != nil {
		return nil, err
	}

	localRel := ole.NewReleaser()
	defer localRel.Release()

	arrVars := make([]VARIANT, 0, len(params))
	for i := len(params) - 1; i >= 0; i-- { // in reverse order
		arrVars = append(arrVars, *NewVariant(localRel, params[i])) // copy bytes, and trust they won't be changed
	}

	var dp DISPPARAMS
	if len(params) > 0 {
		dp.SetArgs(arrVars)
	}
	if method == co.DISPATCH_PROPERTYPUT {
		dp.SetNamedArgs(co.DISPID_PROPERTYPUT)
	}

	v, err := me.Invoke(releaser, memberIds[0], win.LCID_USER_DEFAULT, method, &dp)
	if err != nil {
		return nil, err
	}
	return v, nil
}

type _IDispatchVt struct {
	ole.IUnknownVt
	GetTypeInfoCount uintptr
	GetTypeInfo      uintptr
	GetIDsOfNames    uintptr
	Invoke           uintptr
}

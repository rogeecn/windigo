/**
 * Part of Wingows - Win32 API layer for Go
 * https://github.com/rodrigocfd/wingows
 * This library is released under the MIT license.
 */

package api

import (
	c "wingows/consts"
)

type ACCEL struct {
	FVirt uint8
	Key   uint16
	Cmd   uint16
}

type CREATESTRUCT struct {
	LpCreateParams uintptr
	HInstance      HINSTANCE
	HMenu          HMENU
	HwndParent     HWND
	Cy, Cx, Y, X   int32
	Style          c.WS
	LpszName       *uint16
	LpszClass      *uint16
	ExStyle        c.WS_EX
}

type HELPINFO struct {
	CbSize       uint32
	IContextType c.HELPINFO
	ICtrlId      int32
	HItemHandle  HANDLE
	DwContextId  uintptr
	MousePos     POINT
}

type MENUINFO struct {
	CbSize          uint32
	FMask           c.MIM
	DwStyle         c.MNS
	CyMax           uint32
	HbrBack         HBRUSH
	DwContextHelpID uint32
	DwMenuData      uintptr
}

type MENUITEMINFO struct {
	CbSize        uint32
	FMask         c.MIIM
	FType         c.MFT
	FState        c.MFS
	WId           uint32
	HSubMenu      HMENU
	HBmpChecked   HBITMAP
	HBmpUnchecked HBITMAP
	DwItemData    uintptr
	DwTypeData    *uint16
	Cch           uint32
	HBmpItem      HBITMAP
}

type MONITORINFOEX struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	Flags     uint32
	SzDevice  [32]uint16 // CCHDEVICENAME
}

type NMHDR struct {
	HWndFrom HWND
	IdFrom   uintptr
	Code     uint32 // in fact it should be int32
}

type NONCLIENTMETRICS struct {
	CbSize             uint32
	IBorderWidth       int32
	IScrollWidth       int32
	IScrollHeight      int32
	ICaptionWidth      int32
	ICaptionHeight     int32
	LfCaptionFont      LOGFONT
	ISmCaptionWidth    int32
	ISmCaptionHeight   int32
	LfSmCaptionFont    LOGFONT
	IMenuWidth         int32
	IMenuHeight        int32
	LfMenuFont         LOGFONT
	LfStatusFont       LOGFONT
	LfMessageFont      LOGFONT
	IPaddedBorderWidth int32
}

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

type SIZE struct {
	Cx, Cy int32
}

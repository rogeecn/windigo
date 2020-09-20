/**
 * Part of Wingows - Win32 API layer for Go
 * https://github.com/rodrigocfd/wingows
 * This library is released under the MIT license.
 */

package directshow

type (
	// https://docs.microsoft.com/en-us/windows/win32/api/evr/ns-evr-mfvideonormalizedrect
	MFVideoNormalizedRect struct {
		Left, Top, Right, Bottom float32
	}
)

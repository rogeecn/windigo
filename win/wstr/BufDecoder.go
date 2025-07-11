//go:build windows

package wstr

import (
	"unsafe"
)

// Buffer to receive Windows UTF-16 strings and convert them to Go strings.
//
// Tries to use the internal global buffer for faster operations.
type BufDecoder struct {
	gotGlobal bool // Do we have the lock to the global buffer?
	faceSize  uint // Declared size, which can be smaller than actual buffer size.
	localBuf  []uint16
}

// Constructs a new buffer to receive Windows UTF-16 strings, with the given
// initial length. Tries to use the global buffer; if already in use, allocates
// a dynamic buffer.
func NewBufDecoder(numChars uint) BufDecoder {
	me := BufDecoder{
		gotGlobal: globalBuf.tryGet(),
	}
	me.Resize(numChars)
	return me
}

// Releases the lock for the buffer. No further strings should be added: the
// object must be discarded.
func (me *BufDecoder) Free() {
	if me.gotGlobal {
		globalBuf.release()
	}
	me.localBuf = nil
}

// Returns the size of the receiving buffer.
func (me *BufDecoder) Len() uint {
	return me.faceSize
}

// Resizes the receiving buffer to the given number of chars. Always grows.
func (me *BufDecoder) Resize(numChars uint) {
	if me.gotGlobal {
		if !globalBuf.canFit(int(numChars)) {
			me.localBuf = make([]uint16, numChars)
			copy(me.localBuf, globalBuf.buf[:])
			globalBuf.release() // we won't use the global buffer anymore
			me.gotGlobal = false
		}
	} else {
		if numChars > uint(len(me.localBuf)) { // requesting a buffer even larger
			newLocalBuf := make([]uint16, numChars)
			copy(newLocalBuf, me.localBuf)
			me.localBuf = newLocalBuf
		}
	}
	me.faceSize = numChars
}

// Converts the receiving buffer content to a Go string.
func (me *BufDecoder) String() string {
	return DecodeSlice(me.HotSlice())
}

// Returns a slice over the block.
func (me *BufDecoder) HotSlice() []uint16 {
	if me.gotGlobal {
		return globalBuf.buf[:me.faceSize]
	} else {
		return me.localBuf[:me.faceSize]
	}
}

// Returns the *uint16 to the beginning of the block.
func (me *BufDecoder) UnsafePtr() unsafe.Pointer {
	if me.gotGlobal {
		return unsafe.Pointer(&globalBuf.buf[0])
	} else {
		return unsafe.Pointer(&me.localBuf[0])
	}
}

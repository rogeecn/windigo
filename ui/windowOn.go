package ui

import (
	c "gowinui/consts"
	"gowinui/parm"
)

// Custom hash for WM_NOTIFY messages.
type nfyHash struct {
	IdFrom c.ID
	Code   c.WM
}

// Keeps all user message handlers.
type windowOn struct {
	msgs map[c.WM]func(p parm.Raw) uintptr
	cmds map[c.ID]func(p parm.WmCommand)
	nfys map[nfyHash]func(p parm.WmNotify) uintptr
}

func makeWindowOn() windowOn {
	msgs := make(map[c.WM]func(p parm.Raw) uintptr)
	cmds := make(map[c.ID]func(p parm.WmCommand))
	nfys := make(map[nfyHash]func(p parm.WmNotify) uintptr)

	return windowOn{
		msgs: msgs,
		cmds: cmds,
		nfys: nfys,
	}
}

func (me *windowOn) processMessage(p parm.Raw) (uintptr, bool) {
	switch p.Msg {
	case c.WM_COMMAND:
		paramCmd := parm.WmCommand(p)
		if userFunc, hasCmd := me.cmds[paramCmd.ControlId()]; hasCmd {
			userFunc(paramCmd)
			return 0, true
		}
	case c.WM_NOTIFY:
		paramNfy := parm.WmNotify(p)
		hash := nfyHash{
			IdFrom: c.ID(paramNfy.NmHdr().IdFrom),
			Code:   c.WM(paramNfy.NmHdr().Code),
		}
		if userFunc, hasNfy := me.nfys[hash]; hasNfy {
			return userFunc(paramNfy), true
		}
	default:
		if userFunc, hasMsg := me.msgs[p.Msg]; hasMsg {
			return userFunc(p), true
		}
	}

	return 0, false // no user handler found
}

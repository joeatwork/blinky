package main

import (
	"strconv"
)

func parseColor(rs, gs, bs string) (color uint32, err error) {
	var r, g, b uint64
	color = 0
	err = nil
	r, err = strconv.ParseUint(rs, 0, 32)
	if err != nil {
		return
	}
	g, gerror := strconv.ParseUint(gs, 0, 32)
	if gerror != nil {
		return
	}
	b, berror := strconv.ParseUint(bs, 0, 32)
	if berror != nil {
		return
	}

	color64 :=
		((r & 0xFF) << 16) |
		((g & 0xFF) << 8) |
		(b & 0xFF);

	color = uint32(color64)
	return
}

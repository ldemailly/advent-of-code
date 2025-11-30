package main

import "testing"

func TestRotate(t *testing.T) {
	e := north.RotateRight().RotateRight()
	if e != east {
		t.Error("RotateRight N->E")
	}
	e = e.RotateRight().RotateRight()
	if e != south {
		t.Error("RotateRight E->S")
	}
	e = e.RotateRight().RotateRight()
	if e != west {
		t.Error("RotateRight S->W")
	}
	e = e.RotateRight().RotateRight()
	if e != north {
		t.Error("RotateRight W->N")
	}
}

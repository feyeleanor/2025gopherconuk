package main

import "golang.org/x/sys/unix"

func SemOpen(s string) (uintptr, error) {
	r, _, e := unix.Syscall6(
		268,
		uintptr(CString(s)),
		unix.O_CREAT,
		0644,
		1, 0, 0)
	if r == 0 {
		return 0, ErrorCode("open failed", e)
	}
	return r, nil
}

func SemClose(s uintptr) error {
	return Syscall1(269, s, "close failed")
}

func SemUnlink(s string) error {
	n := uintptr(CString(s))
	return Syscall1(270, n, "unlink failed")
}

func SemWait(s uintptr) error {
	return Syscall1(271, s, "wait failed")
}

func SemTryWait(s uintptr) error {
	return Syscall1(272, s, "trywait failed")
}

func SemPost(s uintptr) error {
	return Syscall1(273, s, "post failed")
}

func Syscall1(op, v uintptr, s string) error {
	r, _, e := unix.Syscall(op, v, 0, 0)
	return SyscallError(r, e, s)
}

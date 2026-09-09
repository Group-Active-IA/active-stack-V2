package system

// SetBinaryInstallDirWritableFn replaces the writable-directory probe used by
// BinaryInstallDir and BinaryInstallDirForHome.
//
// Without this seam a test that exercises binary-install-dir resolution from
// another package silently depends on whether the machine running it happens
// to have a writable /usr/local/bin: the same test resolves to /usr/local/bin
// on one machine and to ~/.local/bin on the next. Tests inject a fake so the
// resolved path is a property of the code, not of the host.
//
// Returns a restore function; always defer it.
func SetBinaryInstallDirWritableFn(fn func(string) bool) (restore func()) {
	old := binaryInstallDirWritableFn
	binaryInstallDirWritableFn = fn
	return func() { binaryInstallDirWritableFn = old }
}

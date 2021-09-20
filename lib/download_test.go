package lib_test

var downloadTests = []struct {
	name    string
	version string
	path    string
	want    interface{}
}{
	// {"invalid path", "1.12.3", "%asd/~pqs", error},
	{"valid path invalid file", "1.12.3", "./tests/", nil},
	{"valid path valid file", "1.12.3", "./", nil},
	// {"empty path", "1.12.3", "", errors.Err},
	// {"inaccessible path", "1.12.3", "/bin", errors.Err},
}

// func TestDownloadKctl(t *testing.T) {
//
// }

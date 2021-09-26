package lib_test

var downloadTests = []struct {
	name    string
	version string
	path    string
	want    interface{}
}{
	{"valid path invalid file", "1.12.3", "./tests/", nil},
	{"valid path valid file", "1.12.3", "./", nil},
	// check file permissions
}

// func TestDownloadKctl(t *testing.T) {
//
// }

/*MIT License Copyright (c) 2021 seaweed843

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
==============================================================================*/

package gozipper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func prepareTestDirTree(path string) {

	test_a := filepath.Join(path, "a")
	os.MkdirAll(path, 0755)
	a_b := filepath.Join(test_a, "b")
	os.MkdirAll(a_b, 0755)
	a_b_c := filepath.Join(a_b, "c")
	os.MkdirAll(a_b_c, 0755)

	ioutil.WriteFile(filepath.Join(test_a, "a0.txt"), []byte("a0\n"), 0644)
	ioutil.WriteFile(filepath.Join(test_a, "a1.txt"), []byte("a1\n"), 0644)
	ioutil.WriteFile(filepath.Join(a_b, "b.txt"), []byte("b\n"), 0644)
	ioutil.WriteFile(filepath.Join(a_b, ".DS_Store"), []byte("b\n"), 0644)
	ioutil.WriteFile(filepath.Join(a_b, "thumbs.db"), []byte("b\n"), 0644)
}

func TestZipPath(t *testing.T) {
	testPath := filepath.Join(".", "test")
	prepareTestDirTree(testPath)
	got := ZipPath(testPath, ".", "tname.tar.zip")
	eq := reflect.DeepEqual(nil, got)
	if !eq {
		t.Errorf("ZipPath() = %q", got)
	}
	got = ZipPath(filepath.Join(testPath, "a", "a0.txt"), "./", "a0_no_txt.zip")
	eq = reflect.DeepEqual(nil, got)
	if !eq {
		t.Errorf("ZipPath() = %q", got)
	}
	os.RemoveAll(testPath)
}

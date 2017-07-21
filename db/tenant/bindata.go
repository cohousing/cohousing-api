// Code generated by go-bindata.
// sources:
// db/tenant/1_base.sql
// DO NOT EDIT!

package tenant

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dbTenant1_baseSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x93\xc1\x4f\x83\x30\x14\xc6\xef\xfd\x2b\xde\x91\xc5\x2d\x51\x93\x9d\x76\xaa\xf0\x36\x1b\xa1\x2c\xa5\x98\xed\xd4\x54\x5b\x95\x38\x18\x01\x16\xfd\xf3\x4d\x8d\x1b\xb0\x0d\x2e\x72\x02\xbe\xf7\x1e\xef\xfb\x7d\x74\x36\x83\x9b\x3c\x7b\xaf\x74\x63\x21\x2d\x89\x2f\x90\x4a\x04\x49\x1f\x42\x04\x5d\xea\xaa\xc9\x6d\xd1\xd4\xc4\x23\x00\x99\x81\xe3\xc5\xb8\xf4\xee\x6e\x27\x90\xf2\x84\xad\x38\x06\x40\x53\x19\x2b\xc6\x7d\x81\x11\x72\x09\x6b\xc1\x22\x2a\xb6\xf0\x84\xdb\x29\x01\x78\xad\xac\x6e\xac\x51\xba\x01\xc9\x22\x4c\x24\x8d\xd6\x6e\x0c\x4f\xc3\xd0\xe9\x87\xd2\x8c\xea\xc6\xee\xec\x98\xae\x8d\xa9\x6c\x5d\xbb\x57\xcf\x54\xf8\x8f\x54\x78\xf7\xf3\xf9\xe4\x57\x27\x93\x05\x21\x34\x94\x28\x2e\x5d\x01\xd0\x20\x00\xc6\x03\xdc\x40\x66\xbe\x55\x2b\xa9\xce\x27\xbd\xf6\xde\xcd\xea\x21\xaa\x6c\x9d\x99\x6b\x84\xfe\x05\x09\xa0\xef\x73\x80\xd5\x58\x59\x67\xff\xb1\xb2\x42\xe7\xb6\xdd\xb8\x07\xaf\x5b\x56\x7e\xec\x0b\xab\x8a\x43\xfe\x62\xab\x91\x32\x9b\xeb\x6c\xa7\x8e\x69\x0c\x96\x9d\x30\x2b\x07\xec\x02\xd3\xd5\xd4\x5a\xd0\xe7\xa1\x9d\x94\xa1\xcc\xa6\x23\x2d\xbd\x55\xbc\xee\xd3\xa9\xcd\x8f\x79\x22\x05\x65\x5c\xc2\xdb\xe7\x50\xeb\x32\x16\xc8\x56\xdc\x65\x79\x36\x06\x04\x2e\x51\x20\xf7\x31\xe9\xfc\x7a\xe0\x65\xc6\x59\xec\x1e\xbf\x60\xff\x55\x90\x40\xc4\xeb\x3f\xcb\x6c\x09\xb8\x61\x89\x4c\x5a\xf3\x8b\xeb\x7a\x3b\x77\xf1\x13\x00\x00\xff\xff\xc1\x79\xf3\x0f\xd0\x03\x00\x00")

func dbTenant1_baseSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbTenant1_baseSql,
		"db/tenant/1_base.sql",
	)
}

func dbTenant1_baseSql() (*asset, error) {
	bytes, err := dbTenant1_baseSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/tenant/1_base.sql", size: 976, mode: os.FileMode(436), modTime: time.Unix(1500582053, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"db/tenant/1_base.sql": dbTenant1_baseSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"db": &bintree{nil, map[string]*bintree{
		"tenant": &bintree{nil, map[string]*bintree{
			"1_base.sql": &bintree{dbTenant1_baseSql, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
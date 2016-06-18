package io

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func Copy(src string, dst string, mode os.FileMode) error {
	//fmt.Println("src", src, "dst", dst)
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return CopyFile(src, dst, mode)
	} else {
		return CopyDir(src, dst, mode)
	}
}

func CopyFile(src string, dst string, mode os.FileMode) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, b, mode)
}

func CopyDir(src string, dst string, mode os.FileMode) error {
	return filepath.Walk(src, createCopyWalkFn(src, dst, mode))
}

func createCopyWalkFn(src_dir string, dst_dir string, mode os.FileMode) filepath.WalkFunc {
	counter := 0
	return func(path string, info os.FileInfo, err error) error {
		counter++
		if counter == 1 {
			// this is the root of the directory being copied
			return os.MkdirAll(dst_dir, mode)
		}

		dstFullPath := filepath.Join(dst_dir, path[len(src_dir):])

		if info.IsDir() {
			return os.MkdirAll(dstFullPath, mode)
		} else {
			return CopyFile(path, dstFullPath, mode)
		}
		return nil
	}
}

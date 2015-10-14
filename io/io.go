package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Copy(src string, dst string, mode os.FileMode) error {
	fmt.Println(src, dst)
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

		fmt.Println("===", path)
		dstFullPath := filepath.Join(dst_dir, path[len(src_dir):])
		fmt.Println("src", path, "dst", dstFullPath)

		if info.IsDir() {
			return os.MkdirAll(dstFullPath, mode)
		} else {
			return CopyFile(path, dstFullPath, mode)
		}
		return nil
	}
}

// CopyDir copies the directory 'src' into the folder 'dst_parent' such that 'dst_parent/basename(src)'
// contains the same contents as 'src'.
func CopyDirAuto(src string, dst_parent string, mode os.FileMode) error {
	return filepath.Walk(src, createCopyWalkFnAuto(src, dst_parent, mode))
}

func createCopyWalkFnAuto(src_dir string, dst_parent_dir string, mode os.FileMode) filepath.WalkFunc {
	counter := 0
	src_dir_parent := filepath.Dir(src_dir)
	return func(path string, info os.FileInfo, err error) error {
		counter++
		if counter == 1 {
			// this is the root of the directory being copied
			return os.MkdirAll(filepath.Join(dst_parent_dir, info.Name()), mode)
		}

		dstFullPath := filepath.Join(dst_parent_dir, strings.Replace(path, src_dir_parent, "", 1))
		fmt.Println("src", path, "dst", dstFullPath)

		if info.IsDir() {
			return os.MkdirAll(dstFullPath, mode)
		} else {
			// copy the file to the directory
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(dstFullPath, b, mode)
		}
		return nil
	}
}

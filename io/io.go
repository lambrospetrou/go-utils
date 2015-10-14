package io

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func compileSinglePost(src_post_dir string, info os.FileInfo, dst_posts_dir string, viewBuilder *view.Builder) error {
	postName := info.Name()[11 : len(info.Name())-3]
	postDir := filepath.Join(dst_posts_dir, postName)
	// create the directory of the post
	if err := os.MkdirAll(postDir, SITE_DST_PERM); err != nil {
		return err
	}

	// get the markdown filename
	postMarkdownPath := filepath.Join(src_post_dir, info.Name())

	// create the actual HTML file for the post
	bundle := &view.TemplateBundle{
		Footer: &view.FooterStruct{Year: time.Now().Year()},
		Header: &view.HeaderStruct{Title: "Single Post"},
		Post:   post.FromFile(postMarkdownPath),
	}
	f, err := os.Create(filepath.Join(postDir, postName+".html"))
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = viewBuilder.Render(w, view.LAYOUT_POST, bundle)
	w.Flush()

	// copy the markdown file to the directory
	markdown, err := ioutil.ReadFile(postMarkdownPath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(postDir, info.Name()), markdown, SITE_DST_PERM)
	return err
}

// CopyDir copies the directory 'src' into the folder 'dst_parent' such that 'dst_parent/basename(src)'
// contains the same contents as 'src'.
func CopyDir(src string, dst_parent string, mode os.FileMode) error {
	return filepath.Walk(src, createCopyWalkFn(filepath.Dir(src), dst_parent, mode))
}

func createCopyWalkFn(src_parent string, dst_parent string, mode os.FileMode) {
	return func(path string, info os.FileInfo, err error) error {
		srcBasePath := strings.Replace(path, src_parent, "", 1) + info.Name()
		dstFullPath := filepath.Join(dst_parent, srcBasePath)

		fmt.Println("src", path, "dst", dstFullPath)

		if info.IsDir() {
			if err := os.MkdirAll(dstFullPath, mode); err != nil {
				return nil
			}
		} else {
			// copy the markdown file to the directory
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(dstFullPath, b, mode)
			return err
		}
	}
}

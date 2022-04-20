package builder

import (
	"haisite/internal/model"
	"haisite/internal/zlog"
	"os"
	"path/filepath"
)

var (
	// ContentPostsDirectory is the directory of posts.
	ContentPostsDirectory = "./content/posts"
)

func (b *Builder) buildPosts() error {
	err := filepath.Walk(ContentPostsDirectory, func(path string, info os.FileInfo, err error) error {
		// skip directory
		if info.IsDir() {
			return nil
		}

		// only process markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		post, err := model.NewPostFromFile(path)
		if err != nil {
			zlog.Warn("failed to build post", "path", path, "err", err)
			return nil
		}

		// save post into parsed data
		b.parsedData.Posts = append(b.parsedData.Posts, post)
		zlog.Info("posts: parsed ok", "path", path, "title", post.Title)

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

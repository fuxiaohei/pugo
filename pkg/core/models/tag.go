package models

import "sort"

type (
	// TagLink returns the tag link of the post.
	TagLink struct {
		Name      string
		Link      string
		LocalFile string
		PostCount int
	}

	// TagPosts is posts categorized by tag.
	TagPosts struct {
		Tag   *TagLink
		Posts []*Post
	}
)

// BuildTagPosts returns the tag posts.
func BuildTagPosts(posts []*Post) []*TagPosts {
	tagData := make(map[string]*TagPosts)
	for _, p := range posts {
		/*if len(p.Tags) != len(p.TagLinks) {
			zlog.Warn("posts: invalid tag links", "title", p.Title, "tags", p.Tags)
			continue
		}*/
		for _, t := range p.TagLinks {
			t2 := t
			if _, ok := tagData[t.Name]; !ok {
				tagData[t.Name] = &TagPosts{
					Tag: t2,
				}
			}
			tagData[t.Name].Posts = append(tagData[t.Name].Posts, p)
		}
	}
	result := make([]*TagPosts, 0, len(tagData))
	for _, tag := range tagData {
		tag.Tag.PostCount = len(tag.Posts)
		result = append(result, tag)
	}

	// aplhabetical sort
	sort.Slice(result, func(i, j int) bool {
		return result[i].Tag.Name < result[j].Tag.Name
	})
	return result
}

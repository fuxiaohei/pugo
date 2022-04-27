package models

import "sort"

// Archive is the archive model.
type Archive struct {
	Year  int
	Posts []*Post
}

// NewArchives creates archives list from posts.
func NewArchives(posts []*Post) []*Archive {
	archivesMap := make(map[int]*Archive)
	for _, p := range posts {
		year := int(p.Date().Year())
		if _, ok := archivesMap[year]; !ok {
			archivesMap[year] = &Archive{
				Year:  year,
				Posts: []*Post{},
			}
		}
		archivesMap[year].Posts = append(archivesMap[year].Posts, p)
	}
	result := make([]*Archive, 0, len(archivesMap))
	for _, a := range archivesMap {
		result = append(result, a)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Year > result[j].Year
	})
	return result
}

package model

import (
	"strconv"
	"strings"
)

type (
	// Pager is pager model.
	Pager struct {
		all   int
		size  int
		pages int
	}
	// PagerItem is pager item model.
	PagerItem struct {
		Begin   int
		End     int
		Prev    int
		Next    int
		Current int
		Total   int
		AllSize int
		HasPrev bool
		HasNext bool
		Link    string
		IsFirst bool

		layout string
	}
)

// PageSize returns page size after pagination
func (p *Pager) PageSize() int {
	return p.pages
}

// NewPager with size and all count
func NewPager(size, all int) *Pager {
	pc := &Pager{
		all:  all,
		size: size,
	}
	if all%size == 0 {
		pc.pages = all / size
	} else {
		pc.pages = all/size + 1
	}
	return pc
}

// Page creates Pager on a page number
func (pg *Pager) Page(i int, layout string) *PagerItem {
	if i < 1 {
		return nil
	}
	begin := (i - 1) * pg.size
	if begin > pg.all {
		return nil // no pager when begin number over all
	}
	pager := &PagerItem{
		Begin:   begin,
		Prev:    i - 1,
		Next:    i + 1,
		Current: i,
		Total:   pg.pages,
		AllSize: pg.all,
		layout:  layout,
	}
	end := begin + pg.size
	if end >= pg.all {
		end = pg.all
		pager.Next = 0 // no next
	}
	pager.End = end
	pager.HasPrev = pager.Prev > 0
	pager.HasNext = pager.Next > 0
	pager.IsFirst = pager.Current == 1
	pager.Link = pager.linkFormat(pager.Current)
	return pager
}

func (pi *PagerItem) PrevLink(i ...int) string {
	if !pi.HasPrev {
		return pi.Link
	}
	prev := 1
	if len(i) > 0 {
		prev = i[0]
	}
	if prev < 1 {
		prev = 1
	}
	return pi.linkFormat(pi.Current - prev)
}

func (pi *PagerItem) NextLink(i ...int) string {
	if !pi.HasNext {
		return pi.Link
	}
	next := 1
	if len(i) > 0 {
		next = i[0]
	}
	if next < 1 {
		next = 1
	}
	return pi.linkFormat(pi.Current + next)
}

func (pi *PagerItem) linkFormat(i int) string {
	// FIXME: use template
	return strings.ReplaceAll(pi.layout, "{{.Page}}", strconv.Itoa(i))
}

package condmodel

import ()

type CondComposerItem struct {
	Where map[string]interface{}
	Rule  map[string]string
	Comp  string
}

type CondComposerLinker struct {
	Item     *CondComposerItem
	CompNext string

	Next *CondComposerLinker
}

func NewCondComposerLinker(compnext string) *CondComposerLinker {
	return &CondComposerLinker{
		Item:     nil,
		Next:     nil,
		CompNext: compnext,
	}
}

func (this *CondComposerLinker) SetItem(item *CondComposerItem) {
	this.Item = item
}

func (this *CondComposerLinker) SetNext(next *CondComposerLinker) {
	this.Next = next
}

type CondComposer struct {
	Items []*CondComposerItem
	Comp  string
}

func NewCondComposerItem(whereIn map[string]interface{}, ruleIn map[string]string, compIn string) *CondComposerItem {
	return &CondComposerItem{
		Where: whereIn,
		Rule:  ruleIn,
		Comp:  compIn,
	}
}

func NewCondComposer(compIn string) *CondComposer {
	return &CondComposer{
		Comp:  compIn,
		Items: make([]*CondComposerItem, 0),
	}
}

func (this *CondComposer) AddItem(item *CondComposerItem) {
	this.Items = append(this.Items, item)
}

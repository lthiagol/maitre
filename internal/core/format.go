package core

import (
	"reflect"
	"sort"
	"strings"
)

// planFormatFacts are the facts of the plan's format: what the writer puts on disk.
//
// The version stamp fences a stale tool off from a plan written under newer rules.
// The **format** is part of those rules: a tool that writes `R01` where the plan
// says `RCP01`, or that writes a record under a field name the reader does not
// know, produces a plan the reader cannot read — and the reader cannot tell that
// the plan is broken.
//
// The list is small and explicit on purpose. A fact that is not in it is a fact
// the stamp cannot see, and a format change that moves nothing here moves no
// version: a stale tool keeps writing.
func planFormatFacts() []string {
	out := []string{
		"id prefix recipe " + recipeIDPrefix,
		"id prefix backlog " + backlogIDPrefix,
	}
	for _, dir := range []string{DirRecipes, DirCourses, DirBacklog} {
		out = append(out, "directory "+dir)
	}
	for _, file := range []string{FileProject, FileActivity, FileDecisions} {
		out = append(out, "file "+file)
	}
	for _, field := range planRecordFields() {
		out = append(out, "field "+field)
	}
	return out
}

// planRecordFields are the JSON field names of the records a plan holds.
//
// The names are read from the struct tags, so they cannot drift from what is
// written: a field renamed on disk is a field renamed here.
//
// The walk follows the records the plan holds — a recipe and everything inside it,
// a backlog item, a journal entry and the journal's envelope — because a rename
// anywhere in them changes what a file on disk says. The envelope is in the walk
// because the journal's own fields are the format too: a journal that carries a
// different `schema` is a journal a reader cannot read.
func planRecordFields() []string {
	var out []string
	seen := map[reflect.Type]bool{}
	for _, record := range []any{Recipe{}, BacklogItem{}, JournalEntry{}, Journal{}, projectFile{}} {
		out = append(out, jsonFieldNames(reflect.TypeOf(record), seen)...)
	}
	sort.Strings(out)
	return out
}

// jsonFieldNames collects the JSON names of a struct's fields and of the structs it
// holds.
//
// The name is the tag's first part — the part before the comma — because that is
// what is on disk. A field with no tag, or one tagged `-`, is not written and so
// is not a fact about the format.
func jsonFieldNames(t reflect.Type, seen map[reflect.Type]bool) []string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct || seen[t] {
		return nil
	}
	seen[t] = true

	var out []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		if name != "" && name != "-" {
			out = append(out, t.Name()+"."+name)
		}
		ft := f.Type
		for ft.Kind() == reflect.Slice || ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct {
			out = append(out, jsonFieldNames(ft, seen)...)
		}
	}
	return out
}

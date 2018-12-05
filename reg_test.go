package gblocks

import (
	// "encoding/json"
	"github.com/kr/pretty"
	r "reflect"
	"testing"
)

type Possessives int

//go:generate stringer -type=Possessives
const (
	IsA Possessives = iota
	// for plural nouns
	IsAn
	AreAn
	// for multiple nouns
	Are
	AreA
)

// TheNouns are-an [ Adjectives ] Class [ NounRelation ] stop.
type NounAssertion struct {
	//*Assertion
	// noun list: must have at least one.
	//Nouns
	// enums generate FieldDropdown
	Possessives
	// optional: adjectives
	// Adjectives []*Adjectives
	// class name
	Class *Class
	// optional: relation.
	PreviousStatement, NextStatement *NounAssertion
}

type Class struct {
	Class string
}

type Noun struct {
	Noun string
}

func TestEnumLabels(t *testing.T) {
	var reg Registry
	reg.RegisterEnum(map[Possessives]string{
		IsA:   "is a",
		IsAn:  "is an",
		AreAn: "are an",
		Are:   "are",
		AreA:  "are a",
	})
	opt := make(map[string]interface{})
	reg.initJson(r.TypeOf((*NounAssertion)(nil)).Elem(), opt)
	// out, _ := json.MarshalIndent(opt, "", "    ")
	//t.Log(string(out)) //
	v := pretty.Diff(opt, possessives)
	if len(v) != 0 {
		t.Fatal(v)
		t.Log(v)
	}
}

var possessives = map[string]interface{}{
	"message0": "%1 %2",
	"args0": []Options{
		{
			"name": "Possessives",
			"type": "field_dropdown",
			"options": []stringPair{
				{"is a", "IsA"},
				{"is an", "IsAn"},
				{"are an", "AreAn"},
				{"are", "Are"},
				{"are a", "AreA"},
			},
		},
		{
			"check": "Class",
			"name":  "Class",
			"type":  "input_value",
		},
	},
	"previousStatement": "NounAssertion",
	"nextStatement":     "NounAssertion",
	"type":              "NounAssertion",
}

package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/dom"
)

// via Blockly.Xml
type BlocklyXml struct {
	*js.Object
}

var Xml BlocklyXml

func (x *BlocklyXml) init() (okay bool) {
	if bl := getBlockly(); bl != nil {
		if obj := bl.Get("Xml"); obj != nil && obj.Bool() {
			x.Object = obj
		}
	}
	return x.Object != nil
}

// func (x *BlocklyXml) workspaceToDom()  (workspace, OptNoId) {
// func (x *BlocklyXml) variablesToDom (variableList) {
// func (x *BlocklyXml) blockToDomWithXY (block, OptNoId) {
// func (x *BlocklyXml) fieldToDomVariable_ (field) {
// func (x *BlocklyXml) fieldToDom_ (field) {
// func (x *BlocklyXml) allFieldsToDom_ (block, element) {
// func (x *BlocklyXml) blockToDom (block, OptNoId) {
// func (x *BlocklyXml) cloneShadow_ (shadow) {
func (x *BlocklyXml) DomToText(xml *dom.Element) (ret string) {
	if x.init() {
		if obj := x.Call("domToText", xml.Object); obj != nil && obj.Bool() {
			ret = obj.String()
		}
	}
	return
}

// func (x *BlocklyXml) domToPrettyText (dom) {
func (x *BlocklyXml) TextToDom(text string) (ret *dom.Element) {
	if x.init() {
		if obj := x.Call("textToDom", text); obj != nil && obj.Bool() {
			ret = &dom.Element{Object: obj}
		}
	}
	return
}

// func (x *BlocklyXml) clearWorkspaceAndLoadFromXml (xml, workspace) {
// func (x *BlocklyXml) domToWorkspace (xml, workspace) {
// func (x *BlocklyXml) appendDomToWorkspace (xml, workspace) {
func (x *BlocklyXml) DomToBlock(xmlBlock *dom.Element, ws *Workspace) (ret *Block) {
	if x.init() {
		if obj := x.Call("domToBlock", xmlBlock.Object, ws.Object); obj != nil && obj.Bool() {
			ret = &Block{Object: obj}
		}
	}
	return
}

// func (x *BlocklyXml) domToVariables (xmlVariables, workspace) {
// func (x *BlocklyXml) domToBlockHeadless_ (xmlBlock, workspace) {
// func (x *BlocklyXml) domToFieldVariable_ (workspace, xml, text, field) {
// func (x *BlocklyXml) domToField_ (block, fieldName, xml) {
// func (x *BlocklyXml) deleteNext (xmlBlock) {

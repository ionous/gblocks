package blockly

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/ionous/gblocks/jsdom"
)

// via Blockly.Xml
type Xml struct {
	*js.Object
}

// func (x *Xml) workspaceToDom()  (workspace, OptNoId) {
// func (x *Xml) variablesToDom (variableList) {
// func (x *Xml) blockToDomWithXY (block, OptNoId) {
// func (x *Xml) fieldToDomVariable_ (field) {
// func (x *Xml) fieldToDom_ (field) {
// func (x *Xml) allFieldsToDom_ (block, element) {
// func (x *Xml) blockToDom (block, OptNoId) {
// func (x *Xml) cloneShadow_ (shadow) {
func (x *Xml) DomToText(xml *jsdom.Element) (ret string) {
	if obj := x.Call("domToText", xml.Object); obj != nil && obj.Bool() {
		ret = obj.String()
	}
	return
}

// func (x *Xml) domToPrettyText (jsdom) {
func (x *Xml) TextToDom(text string) (ret *jsdom.Element) {
	if obj := x.Call("textToDom", text); obj != nil && obj.Bool() {
		ret = &jsdom.Element{Object: obj}
	}
	return
}

// func (x *Xml) clearWorkspaceAndLoadFromXml (xml, workspace) {
// func (x *Xml) domToWorkspace (xml, workspace) {
// func (x *Xml) appendDomToWorkspace (xml, workspace) {
func (x *Xml) DomToBlock(xmlBlock *jsdom.Element, ws *Workspace) (ret *Block) {
	if obj := x.Call("domToBlock", xmlBlock.Object, ws.Object); obj != nil && obj.Bool() {
		ret = &Block{Object: obj}
	}
	return
}

// func (x *Xml) domToVariables (xmlVariables, workspace) {
// func (x *Xml) domToBlockHeadless_ (xmlBlock, workspace) {
// func (x *Xml) domToFieldVariable_ (workspace, xml, text, field) {
// func (x *Xml) domToField_ (block, fieldName, xml) {
// func (x *Xml) deleteNext (xmlBlock) {

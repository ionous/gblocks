package gblocks

import (
	"github.com/gopherjs/gopherjs/js"
)

// via Blockly.Xml
type Xml struct {
	*js.Object
}

// func (x *Xml) workspaceToDom()  (workspace, opt_noId) {
// func (x *Xml) variablesToDom (variableList) {
// func (x *Xml) blockToDomWithXY (block, opt_noId) {
// func (x *Xml) fieldToDomVariable_ (field) {
// func (x *Xml) fieldToDom_ (field) {
// func (x *Xml) allFieldsToDom_ (block, element) {
// func (x *Xml) blockToDom (block, opt_noId) {
// func (x *Xml) cloneShadow_ (shadow) {
func (x *Xml) DomToText(dom *XmlElement) string {
	obj := x.Call("domToText", dom.Object)
	return obj.String()
}

// func (x *Xml) domToPrettyText (dom) {
func (x *Xml) TextToDom(text string) (ret *XmlElement) {
	if obj := x.Call("textToDom", text); obj != nil && obj.Bool() {
		ret = &XmlElement{Object: obj}
	}
	return
}

// func (x *Xml) clearWorkspaceAndLoadFromXml (xml, workspace) {
// func (x *Xml) domToWorkspace (xml, workspace) {
// func (x *Xml) appendDomToWorkspace (xml, workspace) {
func (x *Xml) DomToBlock(xmlBlock *XmlElement, ws *Workspace) (ret *Block) {
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

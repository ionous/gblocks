package blockly

import "github.com/gopherjs/gopherjs/js"

type Extensions struct {
	*js.Object
	// all map[string] interface{}
}

type Callback func(*js.Object, []*js.Object) (ret interface{})

//  map of name to funcion object
type Mixin map[string]*js.Object

func (x *Extensions) Register(name string, cb Callback) {
	x.Call("register", name, cb)
}

func (x *Extensions) RegisterMixin(name string, mixin Mixin) {
	x.Call("registerMixin", name, mixin)
}

func (x *Extensions) RegisterMutator(
	name string,
	mixin Mixin,
	postFn *js.Object,
	quarks []string,
) {
	x.Call("registerMutator", name, mixin, postFn, quarks)
}

func (x *Extensions) Apply(name string, block *Block, isMutator bool) {
	x.Call("apply", name, block, isMutator)
}

// func (x *Extensions) BuildTooltipForDropdown(dropdownName string, table map[string]string) (ret *js.Object) {
// 	ret = x.Call("buildTooltipForDropdown", dropdownName, table)
// 	return
// }

// CheckHasFunction_ = function(errorPrefix, func,
// CheckNoMutatorProperties_ = function(name, block) {
// CheckMutatorDialog_ = function(object, errorPrefix) {
// GetMutatorProperties_ = function(block) {
// CheckBlockHasMutatorProperties_ = function(errorPrefix,
// MutatorPropertiesMatch_ = function(oldProperties, block) {
// CheckDropdownOptionsInTable_ = function(block, dropdownName,
// ExtensionParentTooltip_ = function() {

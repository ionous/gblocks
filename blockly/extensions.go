package blockly

import "github.com/gopherjs/gopherjs/js"

type BlocklyExtensions struct {
	*js.Object
	// all map[string] interface{}
}

var Extensions BlocklyExtensions

func (x *BlocklyExtensions) init() (okay bool) {
	if bl := getBlockly(); bl != nil {
		if obj := bl.Get("Extensions"); obj != nil && obj.Bool() {
			x.Object = obj
		}
	}
	return x.Object != nil
}

type Callback func(*js.Object, []*js.Object) (ret interface{})

//  map of name to funcion object
type Mixin map[string]*js.Object

type PostMixinCallback func(*Block)

// * @this Blockly.Block
func (x *BlocklyExtensions) Register(name string, cb Callback) {
	if x.init() {
		x.Call("register", name, cb)
	}
}

func (x *BlocklyExtensions) RegisterMixin(name string, mixin Mixin) {
	if x.init() {
		x.Call("registerMixin", name, mixin)
	}
}

func (x *BlocklyExtensions) RegisterMutator(name string, mixin Mixin,
	postFn PostMixinCallback, quarks []string) {
	if x.init() {
		cb := js.MakeFunc(func(this *js.Object, args []*js.Object) (ret interface{}) {
			b := &Block{Object: this}
			postFn(b)
			return
		})
		x.Call("registerMixin", name, mixin, cb, quarks)
	}
}

func (x *BlocklyExtensions) Apply(name string, block *Block, isMutator bool) {
	if x.init() {
		x.Call("apply", name, block, isMutator)
	}
}

func (x *BlocklyExtensions) BuildTooltipForDropdown(dropdownName string, table map[string]string) (ret *js.Object) {
	if x.init() {
		ret = x.Call("buildTooltipForDropdown", dropdownName, table)
	}
	return
}

// CheckHasFunction_ = function(errorPrefix, func,
// CheckNoMutatorProperties_ = function(mutationName, block) {
// CheckMutatorDialog_ = function(object, errorPrefix) {
// GetMutatorProperties_ = function(block) {
// CheckBlockHasMutatorProperties_ = function(errorPrefix,
// MutatorPropertiesMatch_ = function(oldProperties, block) {
// CheckDropdownOptionsInTable_ = function(block, dropdownName,
// ExtensionParentTooltip_ = function() {

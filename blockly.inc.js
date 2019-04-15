'use strict';

// the following is needed for testing gblocks --
// it interferes with using gblocks for real though.
// FIX: tbd.
if (typeof window === "undefined") { 
	const jsdom = require("jsdom");
	const { JSDOM } = jsdom;
	const dom  = new JSDOM(``);
	// apparently require('jsdom-global')() does this too?
	// https://www.npmjs.com/package/jsdom-global
	global.window= dom.window;
	global.document = dom.window.document;
	global.DOMParser= dom.window.DOMParser;
	global.XMLSerializer= dom.window.XMLSerializer;

	const jsblockly = require("blockly/blockly_uncompressed.js");
	var ns = jsblockly.Events;
	ns.fire= function(evt) {
		if (ns.isEnabled) {
			ns.FIRE_QUEUE_.push(evt);
			ns.fireNow_();
		}
	}
}
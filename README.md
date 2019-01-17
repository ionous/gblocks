# gblocks
A [gopherjs](https://github.com/gopherjs) wrapper for [Blockly](https://developers.google.com/blockly/guides/overview).

The impetuous for this library is to provide a visual editor for [iffy](https://github.com/ionous/iffy) - an interactive fiction engine.
Gblocks, the editor, and iffy are all works in progress.

## Goals

1. Define [blocks](https://developers.google.com/blockly/guides/create-custom-blocks/define-blocks) using Go-lang types.
2. Simplify [mutations](https://developers.google.com/blockly/guides/create-custom-blocks/web/mutators) to reduce the need for per-block custom code.
3. Build [toolboxes](https://developers.google.com/blockly/guides/configure/web/toolbox) from Go-lang instances.
4. Mirror Go-lang data to/from Blocky blocks ( to provide alternative serialization and processing. )


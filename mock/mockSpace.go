package mock

import (
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/gblocks/block"
)

// created by NewMockSpace
type MockSpace struct {
	*Registry
	Shapes map[string]block.Shape
	ids    map[string]int
	ondel  []block.OnDelete
}

func (ws *MockSpace) Delete(id string) {
	for _, ondel := range ws.ondel {
		ondel.OnDelete(id)
	}
}

func (ws *MockSpace) OnDelete(ondel block.OnDelete) {
	ws.ondel = append(ws.ondel, ondel)
}

func (ws *MockSpace) NewBlockWithId(blockId, blockType string) (ret block.Shape, err error) {
	if desc, ok := ws.Blocks[blockType]; !ok {
		err = errutil.New("unknown block", blockType)
	} else {
		b := CreateBlock(blockId, desc)
		b.Workspace = ws
		if ws.Shapes == nil {
			ws.Shapes = map[string]block.Shape{blockId: b}
		} else {
			ws.Shapes[blockId] = b
		}
		ret = b
	}
	return
}

func (ws *MockSpace) NewBlock(blockType string) (ret block.Shape, err error) {
	idc := ws.ids[blockType]
	zeroIndexed := strconv.Itoa(idc)
	blockId := blockType + "#" + zeroIndexed
	if b, e := ws.NewBlockWithId(blockId, blockType); e != nil {
		err = e
	} else {
		ws.ids[blockType] = idc + 1
		ret = b
	}
	return
}

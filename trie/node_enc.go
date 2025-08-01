// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"github.com/ethereum/go-ethereum/rlp"
)

func nodeToBytes(n node) []byte {
	w := rlp.NewEncoderBuffer(nil)
	n.encode(w)
	result := w.ToBytes()
	w.Flush()
	return result
}

func (n *fullNode) encode(w rlp.EncoderBuffer) {
	offset := w.List()
	for _, c := range n.Children {
		if c != nil {
			c.encode(w)
		} else {
			w.Write(rlp.EmptyString)
		}
	}
	w.ListEnd(offset)
}

func (n *fullnodeEncoder) encode(w rlp.EncoderBuffer) {
	offset := w.List()
	for i, c := range n.Children {
		if len(c) == 0 {
			w.Write(rlp.EmptyString)
		} else {
			// valueNode or hashNode
			if i == 16 || len(c) >= 32 {
				w.WriteBytes(c)
			} else {
				w.Write(c) // rawNode
			}
		}
	}
	w.ListEnd(offset)
}

func (n *fullnodeEncoder) reset() {
	for i, c := range n.Children {
		if len(c) != 0 {
			n.Children[i] = n.Children[i][:0]
		}
	}
}

func (n *shortNode) encode(w rlp.EncoderBuffer) {
	offset := w.List()
	w.WriteBytes(n.Key)
	if n.Val != nil {
		n.Val.encode(w)
	} else {
		w.Write(rlp.EmptyString)
	}
	w.ListEnd(offset)
}

func (n *extNodeEncoder) encode(w rlp.EncoderBuffer) {
	offset := w.List()
	w.WriteBytes(n.Key)

	if n.Val == nil {
		w.Write(rlp.EmptyString) // theoretically impossible to happen
	} else if len(n.Val) < 32 {
		w.Write(n.Val) // rawNode
	} else {
		w.WriteBytes(n.Val) // hashNode
	}
	w.ListEnd(offset)
}

func (n *leafNodeEncoder) encode(w rlp.EncoderBuffer) {
	offset := w.List()
	w.WriteBytes(n.Key) // Compact format key
	w.WriteBytes(n.Val) // Value node, must be non-nil
	w.ListEnd(offset)
}

func (n hashNode) encode(w rlp.EncoderBuffer) {
	w.WriteBytes(n)
}

func (n valueNode) encode(w rlp.EncoderBuffer) {
	w.WriteBytes(n)
}

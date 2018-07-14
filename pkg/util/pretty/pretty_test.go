// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package pretty_test

import (
	"fmt"

	"github.com/cockroachdb/cockroach/pkg/util/pretty"
)

// Example_align demonstrates alignment.
func Example_align() {
	testData := []pretty.Doc{
		pretty.JoinGroupAligned("SELECT", ",",
			pretty.Text("aaa"),
			pretty.Text("bbb"),
			pretty.Text("ccc")),
		pretty.RLTable(
			pretty.RLTableRow{Label: "SELECT",
				Doc: pretty.Join(",",
					pretty.Text("aaa"),
					pretty.Text("bbb"),
					pretty.Text("ccc")),
			},
			pretty.RLTableRow{Label: "FROM",
				Doc: pretty.Join(",",
					pretty.Text("t"),
					pretty.Text("u"),
					pretty.Text("v")),
			}),
	}
	for _, n := range []int{1, 15, 30, 80} {
		for _, doc := range testData {
			p := pretty.Pretty(doc, n, true /*useTabs*/, 4 /*tabWidth*/)
			fmt.Printf("%d:\n%s\n\n", n, p)
		}
	}

	// Output:
	// 1:
	// SELECT
	// 	aaa,
	// 	bbb,
	// 	ccc
	//
	// 1:
	// SELECT
	// 	aaa,
	// 	bbb,
	// 	ccc
	// FROM
	// 	t,
	// 	u,
	// 	v
	//
	// 15:
	// SELECT aaa,
	//        bbb,
	//        ccc
	//
	// 15:
	// SELECT aaa,
	//        bbb,
	//        ccc
	//   FROM t, u, v
	//
	// 30:
	// SELECT aaa, bbb, ccc
	//
	// 30:
	// SELECT aaa, bbb, ccc
	//   FROM t, u, v
	//
	// 80:
	// SELECT aaa, bbb, ccc
	//
	// 80:
	// SELECT aaa, bbb, ccc FROM t, u, v
}

// ExampleTree demonstrates the Tree example from the paper.
func Example_tree() {
	type Tree struct {
		s  string
		n  []Tree
		op string
	}
	tree := Tree{
		s: "aaa",
		n: []Tree{
			{
				s: "bbbbb",
				n: []Tree{
					{s: "ccc"},
					{s: "dd"},
					{s: "ee", op: "*", n: []Tree{
						{s: "some"},
						{s: "another", n: []Tree{{s: "2a"}, {s: "2b"}}},
						{s: "final"},
					}},
				},
			},
			{s: "eee"},
			{
				s: "ffff",
				n: []Tree{
					{s: "gg"},
					{s: "hhh"},
					{s: "ii"},
				},
			},
		},
	}
	var (
		showTree    func(Tree) pretty.Doc
		showTrees   func([]Tree) pretty.Doc
		showBracket func([]Tree) pretty.Doc
	)
	showTrees = func(ts []Tree) pretty.Doc {
		if len(ts) == 1 {
			return showTree(ts[0])
		}
		return pretty.Fold(pretty.Concat,
			showTree(ts[0]),
			pretty.Text(","),
			pretty.Line,
			showTrees(ts[1:]),
		)
	}
	showBracket = func(ts []Tree) pretty.Doc {
		if len(ts) == 0 {
			return pretty.Nil
		}
		return pretty.Fold(pretty.Concat,
			pretty.Text("["),
			pretty.NestT(showTrees(ts)),
			pretty.Text("]"),
		)
	}
	showTree = func(t Tree) pretty.Doc {
		var doc pretty.Doc
		if t.op != "" {
			var operands []pretty.Doc
			for _, o := range t.n {
				operands = append(operands, showTree(o))
			}
			doc = pretty.Fold(pretty.Concat,
				pretty.Text("("),
				pretty.JoinNestedRight(
					pretty.Text(t.op), operands...),
				pretty.Text(")"),
			)
		} else {
			doc = showBracket(t.n)
		}
		return pretty.Group(pretty.Concat(
			pretty.Text(t.s),
			pretty.NestS(int16(len(t.s)), doc),
		))
	}
	for _, n := range []int{1, 30, 80} {
		p := pretty.Pretty(showTree(tree), n, false /*useTabs*/, 4 /*tabWidth*/)
		fmt.Printf("%d:\n%s\n\n", n, p)
	}
	// Output:
	// 1:
	// aaa[bbbbb[ccc,
	//             dd,
	//             ee(some
	//               * another[2a,
	//                         2b]
	//               * final)],
	//     eee,
	//     ffff[gg,
	//             hhh,
	//             ii]]
	//
	// 30:
	// aaa[bbbbb[ccc,
	//             dd,
	//             ee(some
	//               * another[2a,
	//                         2b]
	//               * final)],
	//     eee,
	//     ffff[gg, hhh, ii]]
	//
	// 80:
	// aaa[bbbbb[ccc, dd, ee(some * another[2a, 2b] * final)], eee, ffff[gg, hhh, ii]]
}

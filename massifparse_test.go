package massifparse

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/martinlindhe/unit"
)

// A truncated "detailed" heap.
const detailedHeap = `n6: 1123587 (heap allocation functions) malloc/new/new[], --alloc-fns, etc.
 n2: 639061 0x10F42A: void std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> >::_M_construct<char*>(char*, char*, std::forward_iterator_tag) (basic_string.tcc:219)
  n1: 637595 0x119232: _M_construct_aux<char*> (basic_string.h:247)
   n1: 637595 0x119232: _M_construct<char*> (basic_string.h:266)
    n1: 637595 0x119232: basic_string (basic_string.h:451)
 n0: 12558 in 28 places, all below massif's threshold (1.00%)`

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    *Massif
	}{
		{
			name: "preamble only",
			content: strings.Join([]string{
				`desc: --massif-out-file=/outs/blah.massif`,
				`cmd: /some/foo/command --first-option --second-option`,
				`time_unit: i`,
			}, "\n"),
			want: &Massif{
				Description: "--massif-out-file=/outs/blah.massif",
				Binary:      "/some/foo/command",
				Args:        []string{"--first-option", "--second-option"},
				TimeUnit:    "i",
			},
		},
		{
			name: "single snapshot, no header",
			content: strings.Join([]string{
				`#-----------`,
				`snapshot=0`,
				`#-----------`,
				`time=0`,
				`mem_heap_B=0`,
				`mem_heap_extra_B=0`,
				`mem_stacks_B=0`,
				`heap_tree=empty`,
			}, "\n"),
			want: &Massif{
				Snapshots: []Snapshot{{
					Index:           0,
					Time:            "0",
					MemoryHeap:      0 * unit.Byte,
					MemoryHeapExtra: 0 * unit.Byte,
					MemoryStack:     0 * unit.Byte,
				}},
			},
		},
		{
			name: "single snapshot, no header, detailed heap",
			content: strings.Join([]string{
				`#-----------`,
				`snapshot=2`,
				`#-----------`,
				`time=1`,
				`mem_heap_B=2`,
				`mem_heap_extra_B=3`,
				`mem_stacks_B=4`,
				`heap_tree=detailed`,
				detailedHeap,
			}, "\n"),
			want: &Massif{
				Snapshots: []Snapshot{{
					Index:           2,
					Time:            "1",
					MemoryHeap:      2 * unit.Byte,
					MemoryHeapExtra: 3 * unit.Byte,
					MemoryStack:     4 * unit.Byte,
					HeapTree:        detailedHeap,
				}},
			},
		},
		{
			name: "three snapshots",
			content: strings.Join([]string{
				`#-----------`,
				`snapshot=0`,
				`#-----------`,
				`time=0`,
				`mem_heap_B=0`,
				`mem_heap_extra_B=0`,
				`mem_stacks_B=0`,
				`heap_tree=empty`,
				`#-----------`,
				`snapshot=1`,
				`#-----------`,
				`time=103242501`,
				`mem_heap_B=1123587`,
				`mem_heap_extra_B=89197`,
				`mem_stacks_B=0`,
				`heap_tree=detailed`,
				detailedHeap + "1",
				`#-----------`,
				`snapshot=2`,
				`#-----------`,
				`time=161677632`,
				`mem_heap_B=1616423`,
				`mem_heap_extra_B=140569`,
				`mem_stacks_B=0`,
				`heap_tree=detailed`,
				detailedHeap + "2",
			}, "\n"),
			want: &Massif{
				Snapshots: []Snapshot{
					{
						Index: 0,
						Time:  "0",
					},
					{
						Index:           1,
						Time:            "103242501",
						MemoryHeap:      1123587 * unit.Byte,
						MemoryHeapExtra: 89197 * unit.Byte,
						HeapTree:        detailedHeap + "1",
					},
					{
						Index:           2,
						Time:            "161677632",
						MemoryHeap:      1616423 * unit.Byte,
						MemoryHeapExtra: 140569 * unit.Byte,
						HeapTree:        detailedHeap + "2",
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(test.content))
			t.Logf("Test File:\n%s", test.content)
			if err != nil {
				t.Fatalf("Parse(..) = _, %v, want _, nil", err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Parse(..) has diff (want -> got):\n%s", diff)
			}
		})
	}
}

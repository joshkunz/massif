package massifparse_test

import (
	"os"
	"log"
	"fmt"
	"encoding/json"

	"github.com/joshkunz/massifparse"
)

func Example() {
	f, err := os.Open("testdata/example.massif")
	if err != nil {
		log.Fatalf("failed to open testdata/example.massif: %v", err)
	}
	defer f.Close()

	parsed, err := massifparse.Parse(f)
	if err != nil {
		log.Fatalf("failed to parse testdata/example.massif: %v", err)
	}

	// Convert the Massif struct into JSON for display purposes only. This is
	// not required.
	display, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal Massif for display: %v", err)
	}

	fmt.Print(string(display))

	// Output:
	// {
	//   "Description": "--massif-out-file=./example.massif",
	//   "Binary": "example_binary",
	//   "Args": [
	//     "--option-one",
	//     "--option-two",
	//     "--flag=value"
	//   ],
	//   "TimeUnit": "i",
	//   "Snapshots": [
	//     {
	//       "Index": 0,
	//       "Time": "0",
	//       "MemoryHeap": 0,
	//       "MemoryHeapExtra": 0,
	//       "MemoryStack": 0,
	//       "HeapTree": ""
	//     },
	//     {
	//       "Index": 1,
	//       "Time": "103242501",
	//       "MemoryHeap": 8988696,
	//       "MemoryHeapExtra": 713576,
	//       "MemoryStack": 0,
	//       "HeapTree": "Note: This heap is snipped.\nn6: 1123587 (heap allocation functions) malloc/new/new[], --alloc-fns, etc.\n n0: 12558 in 28 places, all below massif's threshold (1.00%)"
	//     },
	//     {
	//       "Index": 2,
	//       "Time": "161677632",
	//       "MemoryHeap": 12931384,
	//       "MemoryHeapExtra": 1124552,
	//       "MemoryStack": 0,
	//       "HeapTree": "Note: This heap is snipped.\nn6: 1616423 (heap allocation functions) malloc/new/new[], --alloc-fns, etc.\n n0: 13839 in 28 places, all below massif's threshold (1.00%)"
	//     },
	//     {
	//       "Index": 3,
	//       "Time": "273947848",
	//       "MemoryHeap": 22066184,
	//       "MemoryHeapExtra": 1915064,
	//       "MemoryStack": 0,
	//       "HeapTree": ""
	//     }
	//   ]
	// }
}
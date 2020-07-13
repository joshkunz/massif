package massifparse

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/martinlindhe/unit"
)

const (
	snapshotSeparator = "#-----------"
)

type Snapshot struct {
	Index           int
	Time            string
	MemoryHeap      unit.Datasize
	MemoryHeapExtra unit.Datasize
	MemoryStack     unit.Datasize
	HeapTree        string
}

type Massif struct {
	Description string
	Binary      string
	Args        []string
	TimeUnit    string
	Snapshots   []Snapshot
}

type parser struct {
	sc    *bufio.Scanner
	m     *Massif
	atEOF bool
	line  int
}

func (p *parser) abort(f string, args ...interface{}) {
	newArgs := append([]interface{}{p.line}, args...)
	panic(fmt.Errorf("[line %d] "+f, newArgs...))
}

type eofContext string

const (
	noAbortOnEOF eofContext = ""
)

func (p *parser) scan(c eofContext) {
	p.atEOF = !p.sc.Scan()
	if p.atEOF {
		err := p.sc.Err()
		if err != nil {
			panic(err)
		}
		if c != noAbortOnEOF {
			p.abort("unexpected end of file while parsing %s", c)
		}
	} else {
		p.line++
	}
}

func (p *parser) text() string {
	return p.sc.Text()
}

func (p *parser) tryParseHeaderField(field string, out *string) bool {
	line := p.text()
	if !strings.HasPrefix(line, field+": ") {
		return false
	}
	*out = strings.SplitN(line, " ", 2)[1]
	p.scan(noAbortOnEOF)
	return true
}

func (p *parser) parseHeader() {
	p.scan(noAbortOnEOF)
	p.tryParseHeaderField("desc", &p.m.Description)

	var command string
	if p.tryParseHeaderField("cmd", &command) {
		p.m.Binary = strings.Fields(command)[0]
		p.m.Args = strings.Fields(command)[1:]
	}

	p.tryParseHeaderField("time_unit", &p.m.TimeUnit)
}

func (p *parser) eatLine(c eofContext, l string) {
	if line := p.text(); !(line == l) {
		p.abort("got line %q, expected %q", line, l)
	}
	p.scan(c)
}

func (p *parser) parseSnapshotVar(name string) string {
	line := p.text()
	fields := strings.SplitN(line, "=", 2)
	if len(fields) < 2 {
		p.abort("got %q, expected variable \"%s=...\"", line, name)
	}
	parsedName, value := fields[0], fields[1]
	if parsedName != name {
		p.abort("got variable %q, but looking for variable %q", parsedName, name)
	}
	return value
}

func (p *parser) parseSnapshotVarInt64(name string) int64 {
	raw := p.parseSnapshotVar(name)
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		p.abort("snapshot variable %q value %q not an integer: %v", name, raw, err)
	}
	return val
}

func (p *parser) parseDetailedHeapTree() string {
	var tree strings.Builder
	var addNewline bool
	// Not using p.scan, because we want to catch EOF here.
	for p.sc.Scan() {
		// If we're about to start the next snapshot, then we're done parsing
		// the heap tree.
		if p.text() == snapshotSeparator {
			return tree.String()
		}
		if addNewline {
			tree.WriteByte('\n')
		}
		tree.WriteString(p.text())
		addNewline = true
	}
	// Need to update p.atEOF, so we can detect partially recognized input.
	p.atEOF = true
	return tree.String()
}

func (p *parser) parseSnapshot() Snapshot {
	var s Snapshot
	p.eatLine("snapshot", snapshotSeparator)

	s.Index = int(p.parseSnapshotVarInt64("snapshot"))
	p.scan("snapshot")

	p.eatLine("snapshot", snapshotSeparator)

	s.Time = p.parseSnapshotVar("time")
	p.scan("snapshot")
	s.MemoryHeap = unit.Datasize(p.parseSnapshotVarInt64("mem_heap_B")) * unit.Byte
	p.scan("snapshot")
	s.MemoryHeapExtra = unit.Datasize(p.parseSnapshotVarInt64("mem_heap_extra_B")) * unit.Byte
	p.scan("snapshot")
	s.MemoryStack = unit.Datasize(p.parseSnapshotVarInt64("mem_stacks_B")) * unit.Byte
	p.scan("snapshot")
	heapType := p.parseSnapshotVar("heap_tree")

	if heapType == "empty" {
		p.scan(noAbortOnEOF)
		return s
	}

	// Note, this assumes we did *not* call p.scan() after parsing "heap_tree"
	s.HeapTree = p.parseDetailedHeapTree()
	return s
}

func (p *parser) parse() (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		err = e.(error)
	}()
	p.parseHeader()
	for p.text() == snapshotSeparator {
		p.m.Snapshots = append(p.m.Snapshots, p.parseSnapshot())
	}
	if !p.atEOF {
		p.abort("trailing unparsable content starting on this line.")
	}
	return p.sc.Err()
}

func Parse(in io.Reader) (*Massif, error) {
	var out Massif
	p := parser{sc: bufio.NewScanner(in), m: &out}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return &out, nil
}

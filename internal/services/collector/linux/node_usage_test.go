//go:build linux
// +build linux

package linux

import (
	"fmt"
	"strings"
	"testing"

	"github.com/CEBEP9HUH/OTUS_Go_Diploma/internal/services/statistic"
	"github.com/stretchr/testify/require"
)

type lineData struct {
	FS    interface{}
	Nodes interface{}
	Used  interface{}
	Free  interface{}
	Use   interface{}
	Mnt   interface{}
}

var nodeHeader = lineData{
	FS:    "Filesystem",
	Nodes: "Inodes",
	Used:  "IUsed",
	Free:  "IFree",
	Use:   "IUse%",
	Mnt:   "Mounted on",
}

func makeTestTable(header lineData, data []lineData) string {
	const (
		tmplHead = "%-15.15s %15.15s %15.15s %15.15s %15.15s %-15.15s"
		tmplData = "%-15.15s %15d %15d %15d %14d%% %-15.15s"
	)
	var table strings.Builder
	h := fmt.Sprintf(tmplHead, header.FS, header.Nodes,
		header.Used, header.Free, header.Use, header.Mnt)
	table.WriteString(h)
	table.WriteString(lineSeparator)
	for _, d := range data {
		l := fmt.Sprintf(tmplData, d.FS, d.Nodes,
			d.Used, d.Free, d.Use, d.Mnt)
		table.WriteString(l)
		table.WriteString(lineSeparator)
	}
	return table.String()
}

func TestTable(t *testing.T) {
	data := []lineData{
		{
			FS:    "tmpfs",
			Nodes: 4102311,
			Used:  1171,
			Free:  4101140,
			Use:   1,
			Mnt:   "/run",
		},
		{
			FS:    "/dev/nvme0n1p4",
			Nodes: 61030400,
			Used:  1131482,
			Free:  59898918,
			Use:   22,
			Mnt:   "/",
		},
		{
			FS:    "tmpfs",
			Nodes: 4102311,
			Used:  45,
			Free:  4102266,
			Use:   1,
			Mnt:   "/dev/shm",
		},
		{
			FS:    "/dev/nvme0n1p3",
			Nodes: 0,
			Used:  0,
			Free:  0,
			Use:   0,
			Mnt:   "/boot/efi",
		},
	}
	expect := statistic.NodeUsage{
		NodeUsage: make(map[string]statistic.FSNodeInfo, len(data)),
	}
	for _, v := range data {
		expect.NodeUsage[v.Mnt.(string)] = statistic.FSNodeInfo{
			FS:    v.FS.(string),
			Usage: float32(v.Use.(int)),
		}
	}
	table := makeTestTable(nodeHeader, data)

	var n nodeUsageLinux
	res, err := n.extractStat(table)
	require.NoError(t, err)
	require.EqualValues(t, expect, res)
}

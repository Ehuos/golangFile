package main

import (
	"bufio"
	"errors"
	"fmt"
	//"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sort"
	"time"
	"io/ioutil"
)

const TARGET = "Swap:"
const FORMAT = "%5d %9s %s\n"

type swap_info struct {
	Pid  int64
	Size int64
	Comm string
}

type swap_infos []*swap_info

func (p swap_infos) Len() int {
	return len(p)
}

func (p swap_infos) Less(i, j int) bool {
	return p[i].Size < p[j].Size
}

func (p swap_infos) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	start := time.Now()

	slist := make(swap_infos, 0)
	getSwap(&slist)
	sort.Sort(&slist)
	fmt.Printf("%5s %9s %s\n", "PID", "SWAP", "COMMAND")
	var total int64 = 0
	for _, v := range slist {
		fmt.Printf(FORMAT, v.Pid, filesize(v.Size), v.Comm)
		total = total+v.Size
	}
	fmt.Printf("Total  %8s\n", filesize(total))
	t1 := time.Now()

	fmt.Printf("Cost time %v\n", t1.Sub(start))
}


func getSwap(list *swap_infos) (err error) {
	f, _ := os.Open("/proc")
	names, err := f.Readdirnames(0)
	if err != nil {
		f.Close()
		return
	}
	for _, name := range names {
		toint, err := strconv.ParseInt(name, 10, 0)
		if err == nil {
			info := &swap_info{
				Pid: toint,
			}
			err = getSwapFor(info)
			if err == nil {
				*list = append(*list, info)
			}
		}
	}
	f.Close()
	return
}

func getSwapFor(info *swap_info) (err error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/cmdline", info.Pid))
	if err != nil {
		return
	}

	buff, err := ioutil.ReadAll(f)
	if err != nil {
		f.Close()
		return
	}

	f.Close()

	var len = len(buff)

	if len == 0 {
		return errors.New("o cmdline")
	}

	len --

	for i := 0; i < len; i++ {
		if buff[i] == 0 {
			buff[i] = 32
		}
	}

	buff = buff[:len-1]
	info.Comm = string(buff)
	f, err = os.Open(fmt.Sprintf("/proc/%d/smaps", info.Pid))
	if err != nil {
		return
	}

	size := getSwapSize(bufio.NewScanner(f))

	if size == 0 {
		return errors.New("wrap 0")
	}

	info.Size = size
	f.Close()
	return
}

var warp_string_buf string

func getSwapSize(r *bufio.Scanner) (size int64) {
	size = 0
	for r.Scan(){
		warp_string_buf = r.Text()
		b := strings.HasPrefix(warp_string_buf, TARGET)
		if b {
			curr_warp_array := strings.Split(warp_string_buf, " ")
			toint, _ := strconv.ParseInt(curr_warp_array[len(curr_warp_array)-2], 10, 0)
			size = size+toint
		}
	}
	return
}

var units string = "KMGT"

func filesize(s int64) string {
	var unit int8 = 0
	var left float32 = float32(s)
	for unit < 3 {
		if left > 1100 {
			left = left/1024
		} else {
			break
		}
		unit++
	}
	return fmt.Sprintf("%.1f", left) + string(units[unit]) + "iB"
}

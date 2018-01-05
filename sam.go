package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f1, _ := os.Open("test_vlans.csv")
	f2, _ := os.Open("test_requests.csv")
	defer f1.Close()
	defer f2.Close()

	r1 := csv.NewReader(f1)
	r2 := csv.NewReader(f2)

	// vlans headers
	_, _ = r1.Read()
	// request headers
	_, _ = r2.Read()

	vlans, _ := r1.ReadAll()
	requests, _ := r2.ReadAll()

	graph := NewGraph(vlans)
	graph.populateGraph()
	for key, val := range graph.nodeMap {
		var keys []int
		for k := range val.primaryDevices {
			keys = append(keys, k)
		}
		for k := range val.secondaryDevices {
			keys = append(keys, k)
		}
		val.commonDevices = mapper(keys)
		fmt.Printf("%#v,%#v\n", key, val)
	}
	output := perfromMapping(graph, requests)
	for _, o := range output {
		fmt.Println(o)
	}
}

type networkGraph struct {
	vlans   [][]string
	nodeMap map[int]*vlanNode
}

func NewGraph(vlans [][]string) *networkGraph {
	return &networkGraph{vlans, make(map[int]*vlanNode)}
}

func (g *networkGraph) populateGraph() {
	for _, v := range g.vlans {

		vlanID, _ := strconv.Atoi(v[2])
		node := &vlanNode{}

		if n, ok := g.nodeMap[vlanID]; ok {
			node = n
		} else {
			node = &vlanNode{
				nodeID:           vlanID,
				primaryDevices:   make(map[int]int),
				secondaryDevices: make(map[int]int),
				commonDevices:    make(map[int]int),
			}
			g.nodeMap[vlanID] = node
		}

		port, _ := strconv.Atoi(v[1])
		deviceID, _ := strconv.Atoi(v[0])
		if port == 1 {
			node.primaryDevices[deviceID] = 1
		} else {
			node.secondaryDevices[deviceID] = 1
		}
	}
}

type vlanNode struct {
	nodeID           int
	primaryDevices   map[int]int
	secondaryDevices map[int]int
	commonDevices    map[int]int
}

// Return a map of devices with count 2
func mapper(arr []int) map[int]int {
	temp := make(map[int]int)
	res := make(map[int]int)
	for _, v := range arr {
		temp[v]++
		if temp[v] == 2 {
			res[v] = 2
		}
	}
	return res
}

func min(t map[int]int) int {
	min := 0
	for key := range t {
		if min == 0 {
			min = key
		} else if min > key {
			min = key
		}
	}
	return min
}
func perfromMapping(graph *networkGraph, requests [][]string) [][]int {
	vID := 1
	res := make([][]int, 0)
	for _, req := range requests {

		r, _ := strconv.Atoi(req[1])
		reqID, _ := strconv.Atoi(req[0])

		for {
			// fmt.Println("---------------------------")
			// fmt.Printf("reqid: %d\t redundancy: %d\t vid: %d\n", reqID, r, vID)
			vNode := graph.nodeMap[vID]

			if len(vNode.primaryDevices) == 0 {
				vID++
				// fmt.Println("next vid", vID)
			}
			vNode = graph.nodeMap[vID]

			if r == 1 {
				if len(vNode.primaryDevices) == 0 {
					// fmt.Println("turning to next vid from vid", vID)
					vID++
				}
				vNode = graph.nodeMap[vID]

				deviceID := min(vNode.commonDevices)
				delete(vNode.primaryDevices, deviceID)
				delete(vNode.secondaryDevices, deviceID)
				delete(vNode.commonDevices, deviceID)

				// res = append(res, []int{reqID, deviceID, 0, vID})
				// res = append(res, []int{reqID, deviceID, 1, vID})
				fmt.Printf("req_id :%d\t device_id: %d\t port:%d\t vid: %d\n", reqID, deviceID, 0, vID)
				fmt.Printf("req_id :%d\t device_id: %d\t port:%d\t vid: %d\n", reqID, deviceID, 1, vID)
				break

			} else if r == 0 {
				if len(vNode.primaryDevices) == 0 {
					// fmt.Println("turning to next vid from vid", vID)
					vID++
				}
				vNode = graph.nodeMap[vID]

				deviceID := min(vNode.primaryDevices)
				delete(vNode.primaryDevices, deviceID)

				// res = append(res, []int{reqID, deviceID, 0, vID})
				fmt.Printf("req_id :%d\t device_id: %d\t port:%d\t vid: %d\n", reqID, deviceID, 1, vID)
				break
			}
		}

	}
	return nil
}

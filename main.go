package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f1, _ := os.Open("vlans.csv")
	f2, _ := os.Open("requests.csv")
	f3, _ := os.Create("output.csv")
	defer f1.Close()
	defer f2.Close()
	defer f3.Close()

	r1 := csv.NewReader(f1)
	r2 := csv.NewReader(f2)

	// skip vlans headers
	_, _ = r1.Read()
	// skip request headers
	_, _ = r2.Read()

	vlans, _ := r1.ReadAll()
	requests, _ := r2.ReadAll()

	graph := NewGraph(vlans)
	graph.populateGraph()

	// Find the common devices for each vlan node
	for _, val := range graph.nodeMap {
		var keys []int
		for k := range val.primaryDevices {
			keys = append(keys, k)
		}
		for k := range val.secondaryDevices {
			keys = append(keys, k)
		}
		val.commonDevices = uniqueCount(keys)
	}

	output := perfromMapping(graph, requests)

	writer := csv.NewWriter(f3)
	// write headers
	writer.Write([]string{"request_id", "device_id", "primary_port", "vlan_id"})
	err := writer.WriteAll(output)
	if err != nil {
		fmt.Println(err)
	}
	writer.Flush()
}

type vlanNode struct {
	nodeID           int
	primaryDevices   map[int]int
	secondaryDevices map[int]int
	commonDevices    map[int]int
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

// Return a map of devices with count 2
func uniqueCount(arr []int) map[int]int {
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

// Return the minimum value from a map
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

func perfromMapping(graph *networkGraph, requests [][]string) [][]string {
	vID := 1
	res := make([][]string, 0)
	var vNode *vlanNode
	for _, req := range requests {

		r, _ := strconv.Atoi(req[1])

		for {
			vNode = graph.nodeMap[vID]

			if r == 1 {
				if len(vNode.commonDevices) == 0 {
					vID++
					continue
				}
				deviceID := min(vNode.commonDevices)
				delete(vNode.primaryDevices, deviceID)
				delete(vNode.secondaryDevices, deviceID)
				delete(vNode.commonDevices, deviceID)

				res = append(res, []string{req[0], strconv.Itoa(deviceID), strconv.Itoa(0), strconv.Itoa(vID)})
				res = append(res, []string{req[0], strconv.Itoa(deviceID), strconv.Itoa(1), strconv.Itoa(vID)})
				break

			} else {
				if len(vNode.primaryDevices) == 0 {
					vID++
					continue
				}

				deviceID := min(vNode.primaryDevices)
				delete(vNode.primaryDevices, deviceID)

				res = append(res, []string{req[0], strconv.Itoa(deviceID), strconv.Itoa(1), strconv.Itoa(vID)})
				break
			}

		}

	}
	return res
}

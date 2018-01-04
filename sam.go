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

	_, _ = r1.Read()
	vlans, _ := r1.ReadAll()
	_, _ = r2.ReadAll()

	graph := NewGraph(vlans)
	graph.populateGraph()
	for keys, val := range graph.nodeMap {
		val.nodesCount = mapper(append(val.primaryDevices, val.secondaryDevices...))
		fmt.Printf("%#v,%#v\n", keys, val)
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
		id, _ := strconv.Atoi(v[0])
		node := &vlanNode{}
		if n, ok := g.nodeMap[id]; ok {
			node = n
		} else {
			node = &vlanNode{nodeID: id}
			g.nodeMap[id] = node
		}
		port, _ := strconv.Atoi(v[1])
		deviceID, _ := strconv.Atoi(v[2])
		if port == 1 {
			node.primaryDevices = append(node.primaryDevices, deviceID)
		} else {
			node.secondaryDevices = append(node.secondaryDevices, deviceID)
		}
	}
}

type vlanNode struct {
	nodeID           int
	primaryDevices   []int
	secondaryDevices []int
	nodesCount       map[int]int
}

func (n *vlanNode) existsPrimarySecondary(deviceID int) {
	
}

func mapper(arr []int) map[int]int {
	res := make(map[int]int)
	for _, v := range arr {
		res[v]++
		/*
			if res[v] == 2 {
				common = append(common, v)
			}
		*/
	}
	return res
}

func (n *vlanNode) eat() {
	n.primaryDevices = n.primaryDevices[1:]
	n.secondaryDevices = n.secondaryDevices[1:]
}
func nonRedundant(n *vlanNode) int {
	return n.primaryDevices[0]
}

func redundant(t map[int]int) int {
	min := 0
	for key, val := range t {
		if val == 2 {
			if min == 0 {
				min = key
			} else if min > key {
				min = key
			}
		}
	}
	// return some execessively high number
	// Todo: make it less whacky
	if min == 0 {
		return 100000
	}
	return min
}
func perfromMapping(graph networkGraph, requests [][]string) {
	// res := make([][]int, 0)
	for _, req := range requests {

		r, _ := strconv.Atoi(req[1])
		for keys := range graph.nodeMap {
			switch r {
			case 1:
				a := redundant()
				b := redundant(t1m)
				c := redundant(t2m)

				// fmt.Println("red values:",a,b,c)

				if a < b && a < c || a == b || a == c {
					fmt.Printf("redundant, device:%d, vland_id:%d\n", 0, a)
					delete(t0m, a)
					t0.eat()
					// fmt.Println(t0)
				} else if b < a && b < c || b == c {
					fmt.Printf("redundant, device:%d, vland_id:%d\n", 1, b)
					delete(t1m, b)
					t1.eat()
					// fmt.Println(t1)
				} else {
					fmt.Printf("redundant, device:%d, vland_id:%d\n", 2, c)
					delete(t2m, c)
					t2.eat()
					// fmt.Println(t2)
				}

			case 0:
				a := nonRedundant(t0)
				b := nonRedundant(t1)
				c := nonRedundant(t2)

				// fmt.Println("non-values:",a,b,c)
				if (a < b || a == b) && (a == c || a < c) {
					fmt.Printf("non, device:%d, vland_id:%d\n", 0, a)
					t0.nodes[0].nodes = t0.nodes[0].nodes[1:]
					// fmt.Println(t0)
				} else if b < a && (b < c || b == c) {
					fmt.Printf("non, device:%d, vland_id:%d\n", 1, b)
					t1.nodes[0].nodes = t1.nodes[0].nodes[1:]
					// fmt.Println(t1)
				} else {
					fmt.Printf("non, device:%d, vland_id:%d\n", 2, c)
					t2.nodes[0].nodes = t2.nodes[0].nodes[1:]
					// fmt.Println(t2)
				}
			}
		}
	}
}

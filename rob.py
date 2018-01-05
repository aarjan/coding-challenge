# -*- coding: utf-8 -*-
"""
    This module consists of algorithm to effectively reserve VLAN IDs
    on a particular device and a port depending on the redundancy criteria
    that is defined during the request.
    Note: The implementation involves graph implementation where each VLAN ID
    is treated as a 'Node(VlanNode)' and devices are treated as its 'Children'.
    VLANNode: This holds the vlan id as a satellite data and has diffrent children attributes
    which is discussed below. Ex: VLANNode(2) #vlan node with id of 2 
    Children of the VLAN Node are as follows:
    
    1. DevicesPrimary: This is a device that has vlan id in its range available for use
    for primary port. Ex: {1, 2, 3}
    2. DevicesSecondary: This is a device that has vlan id in its range available for use
    for secondary port. Ex: {3, 6, 7}
    3. DevicesCommon: This is a device that has a vlan id  in its range available for use for 
    both the primary and secondary port. This makes the process of elimination during mapping
    much more faster. This attribute is only here for making the computaion process easier.
    Ex: {3} #common device.
""" 
import csv

class VLanNode:
    """The VLanNode object represent a single node in a cluster of available 
    vlan ids. Each vlan node has attributes that holds sets of device ids.
    
    :param value: represents the vlan id
    :param devices_primary: set of device ids whose vlan id is available 
                            for use in primary port
    :param devices_secondary: set of device ids whose vlan id is available 
                            for use in secondary port
    :param devices_common: set of device ids whose vlan id is availabe 
                            for use in both the secondary and primary port
    """

    def __init__(self, value):
        self.value = value
        # involves a lot of membership testing which we want to compute in O(1) 
        self.devices_primary = set()
        self.devices_secondary = set()
        self.devices_common = set()
   
    def get_device_list(self, is_primary):
        """Helper method that returns appropriate set of device based on
        weather or not the device has vlan id in its range for primary
        or the secondary port.
        """ 
        return self.devices_primary if is_primary else self.devices_secondary
   
    def exists_primary_secondary(self, device_id):
        """Helper method that checks if vlan id (self.value)
        is available for use on both the primary and secondary port
        """
        return self.devices_primary.\
            intersection(self.devices_secondary).\
            intersection(set([device_id]))
  
    def __repr__(self):
        return 'VLanNode({0})'.format(self.value)


class NetworkGraph:
    """This graph object holds the network of vlan nodes.
    :param vlans: holds the data of vlans (vlans.csv)
    :param id_vlan_node_map: holds the mapping of vlan_id(key) and vlan_node(value)
                            for later use during populating the network graph.
    """
    
    def __init__(self, vlans):
        self.vlans = vlans
        self.id_vlan_node_map = {}
    
    def populate_graph(self):
        """This method populates the network with the help of vlans.csv.
        Summary:
        Step 1: Iterate over the vlans.csv
        Step 2: Get the vlan id
        Step 3: Check to see if the node with that id already exists.
                If it exists, then grab the node else create a new Vlan node
                and put it into `self.id_vlan_node_map` for later access
        Step 4: Find out if the vlan id is for the primary or secondary port
        Step 5: Accordingly grab the device list from the vlan node attribute.
        Step 6: Get the device id and add it to the device list
        Step 7: Also check if the vlan id exists for both the primary and secondary port.
                If it exists, add it to common device list.
        Step 8: Repeat
        """
        
        for vlan_item in self.vlans:
            # vlan id
            _id = int(vlan_item['vlan_id'])
            if _id in self.id_vlan_node_map:
                vlan_node = self.id_vlan_node_map[_id]
            else:
                vlan_node = VLanNode(_id)
                # add to the mapping
                self.id_vlan_node_map[_id] = vlan_node
            is_primary_port = True if vlan_item['primary_port'] == '0' else False
            device_list = vlan_node.get_device_list(is_primary_port)
            device_id = int(vlan_item['device_id'])
            device_list.add(device_id)
            # check to see if this device's vlan id is available for both the primary and secondary
            if vlan_node.exists_primary_secondary(device_id):
                vlan_node.devices_common.add(device_id) 


def perform_mapping(graph, requests, vlans_ids):
    """This function is the core algorithm behind assiging appropriate vlan id
    for the device obeying the rules defined when redundancy is required or not.
    :param graph: Instance of Network Graph
    :param requests: list of requests parsed from requests.csv
    :param vlans_ids: list of all the vlan ids available for all the devices in the network
    Summary:
    Step 1: Iterate over the requests.csv
    Step 2: Start from the lowest available vlan id using while loop
    Step 3: Grab the Vlan Node from the vlan id using `graph.id_vlan_node_map`
    Step 4: If request does not require redundancy("0")
            a. Find the device with the lowest device id from the `vlan_node.devices_secondary`
            b. And if the device exists in devices_common, remove it from devices_common
            since it is no longer available which require "redundancy"
    Step 5: If request require redundancy("1")
            a. Find the device with the lowest device id from the `vlan_node.devices_common`
            b. And if the device exists in devices_primary and devices_secondary, remove it
            since it is no longer availabe for later use
    Step 6: If there is no device available, if bothe Step4 / Step5 fails, get the next availble
            lowest vlan node and repeat Step4/Step5
    """

    for request in requests:
        request_id = request['request_id']
        current_index = 0
        while True:
            
            try:
                current_vlan_id = vlans_ids[current_index]
            except IndexError:
                break
            
            vlan_node = graph.id_vlan_node_map[current_vlan_id]
            
            try:
                if request['redundant'] == '0':
                    device_id = min(vlan_node.devices_secondary)
                    vlan_node.devices_secondary.remove(device_id)
                    if device_id in vlan_node.devices_common:
                        vlan_node.devices_common.remove(device_id)
                    print(request_id, device_id, 1, current_vlan_id)
                    break
                
                elif request['redundant'] == '1':
                    device_id = min(vlan_node.devices_common)
                    vlan_node.devices_common.remove(device_id)
                    if device_id in vlan_node.devices_secondary:
                        vlan_node.devices_secondary.remove(device_id)
                    if device_id in vlan_node.devices_primary:
                        vlan_node.devices_primary.remove(device_id)
                    print(request_id, device_id, 0, current_vlan_id)
                    print(request_id, device_id, 1, current_vlan_id)
                    break

            except ValueError as e:
                print("exception ",e)
                current_index += 1
        

def main():
    import csv
    vlans = list(csv.DictReader(open('test_vlans.csv')))
    g = NetworkGraph(vlans)
    g.populate_graph()
    print(g.id_vlan_node_map)
    vlan_ids = sorted(list(g.id_vlan_node_map.keys()))
    requests = list(csv.DictReader(open('test_requests.csv')))
    for keys in g.id_vlan_node_map:
        print(keys,g.id_vlan_node_map[keys].devices_secondary,g.id_vlan_node_map[keys].devices_primary,g.id_vlan_node_map[keys].devices_common)
    perform_mapping(g, requests, vlan_ids)
    
if __name__=="__main__":
    main()
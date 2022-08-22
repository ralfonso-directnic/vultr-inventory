package main

type Inventory map[string]interface{}

type Hosts struct{
	Hosts []string `json:"hosts"`
}



type ResponseMeta struct {
    Hostvars map[string]map[string]interface{} `json:"hostvars"`
}

type ResponseChildren struct {
	Children []interface{}
}

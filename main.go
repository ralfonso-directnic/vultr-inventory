package main

import (
	"github.com/vultr/govultr/v2"
	"os"
	"context"
	"log"
	"golang.org/x/oauth2"
	"encoding/json"
	"fmt"
	"flag"
	"strings"
	"time"
	"gopkg.in/ini.v1"
)
var host string
var list bool
var debug bool
var indent bool
var refresh bool
var cacheFile string

func main(){

	cacheFile = "./vultr_inventory.cache"

	flag.StringVar(&host,"host", "", "Return Individual Host")
	flag.BoolVar(&list,"list", true, "List hosts")
	flag.BoolVar(&indent,"indent",false,"Indent output")
	flag.BoolVar(&debug,"debug", false, "Enable debug output.")
	flag.BoolVar(&refresh,"refresh",false,"Refresh cache")
	flag.Parse()

	if(cacheExists() && refresh==false){

		dat,_ := os.ReadFile(cacheFile)
		fmt.Print(string(dat))

	}else{

		fetchHosts()

	}


}

func cacheExists() (bool){

	//is it too old?



	if(fileExists(cacheFile)){



		 stat,e := os.Stat(cacheFile)

		 if(e!=nil){}

		mtime := stat.ModTime()

		now := time.Now()

		deadline := mtime.Add(time.Hour * 4)


		if(now.After(deadline)){

			return false
		}



		return true

	}

	return false

}



func fetchHosts() ([]interface{},error) {

    var apiKey string

	cfg, err := ini.Load("./vultr.ini")

	if(err!=nil){



	}else{

		apiKey=cfg.Section("default").Key("key").String()

	}
	

	if(len(apiKey)<1) {

		apiKey = os.Getenv("VULTR_API_KEY")

	}

	if len(apiKey)<1 {

		log.Fatal("VULTR_API_KEY required")

	}

	config := &oauth2.Config{}
	ctx := context.Background()
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: apiKey})
	vultrClient := govultr.NewClient(oauth2.NewClient(ctx, ts))

	// Optional changes
	_ = vultrClient.SetBaseURL("https://api.vultr.com")
	vultrClient.SetUserAgent("vultr-inventory")
	vultrClient.SetRateLimit(500)

	listOptions := &govultr.ListOptions{PerPage:250}


	    resp := Inventory{}
	    hostvars := ResponseMeta{}
	    all := Hosts{}
	    vultr:= Hosts{}

	    hostvars.Hostvars = make(map[string]map[string]interface{})

		i, _, err := vultrClient.Instance.List(context.Background(), listOptions)

		if err != nil {
			return nil, err
		}

		for _,it := range i {


			key := fmt.Sprintf("region_%s",it.Region)

			if obj, ok := resp[key]; ok {

				rk := obj.(Hosts)
				rk.Hosts = append(rk.Hosts,it.Label)
				resp[key] = rk

			}else{

				rk := Hosts{}
				rk.Hosts = append(rk.Hosts,it.Label)
				resp[key] = rk

			}

			key = fmt.Sprintf("os_%s",strings.Replace(strings.ToLower(it.Os)," ","_",-1))

			if obj, ok := resp[key]; ok {

				rk := obj.(Hosts)
				rk.Hosts = append(rk.Hosts,it.Label)
				resp[key] = rk

			}else{

				rk := Hosts{}
				rk.Hosts = append(rk.Hosts,it.Label)
				resp[key] = rk

			}


			for _,t := range it.Tags {


                key = fmt.Sprintf("group_%s",t)


				if obj, ok := resp[key]; ok {

                    rk := obj.(Hosts)
					rk.Hosts = append(rk.Hosts,it.Label)
					resp[key] = rk

				}else{

					rk := Hosts{}
					rk.Hosts = append(rk.Hosts,it.Label)
					resp[key] = rk

				}

			}

			vultr.Hosts = append(vultr.Hosts,it.Label)
			all.Hosts = append(all.Hosts,it.Label)

			if entry,ok := hostvars.Hostvars[it.Label]; ok {

				entry["ansible_host"] = it.MainIP
				hostvars.Hostvars[it.Label] = entry

			}else{

			    entry := make(map[string]interface{})
			    entry["ansible_host"] = it.MainIP
				hostvars.Hostvars[it.Label] = entry

			}


			//entries = append(entries,it)

		}

	    resp["_meta"] = hostvars
		resp["all"] = all
		resp["vultr"] = vultr


	var b []byte
	var e error



	if(indent==true) {

		b,e = json.MarshalIndent(resp, "", " ")

	}else{

		b,e = json.Marshal(resp)

	}

	if(e!=nil){}

		//update cache

	   os.WriteFile(cacheFile,b,0644)
	   fmt.Print(string(b))
		var ent []interface{}
		return ent,nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


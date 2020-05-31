# Warframe Wiki Module Query
[![Go Report Card](https://goreportcard.com/badge/github.com/obowersa/wfwiki)](https://goreportcard.com/report/github.com/obowersa/wfwiki)

## Summary
wfwiki is a go library for retrieving Warframe Fandom Wiki module  as structs

## Details
wfwiki is heavily inspired by [Snekw's work](https://wf.snekw.com/)

Behind the scenes, wfwiki uses a rate limited http request (1 tick per second) to call the Fandom Wiki's API. The JSON
is then parsed through an embeded lua VM to generate a JSON string for the data table, and unmarshall's this into a struct.

Current tables which are supported:
- Weapons
- Warframes
- Modules

## TODO
- [ ] Improve testing and test coverage
- [ ] Implement CICD on Tekton
- [ ] Support for data tables which have name collisions in the same object. Currently duplicate names like SecondaryAttack for a weapon will get dropped
- [ ] Implement simple cache with checking of hash
- [ ] Implement lua worker pool
- [ ] Move GetStats outside of package. Currently it is used for  a temporary test client
- [ ] Create examples folder

## Example
```go
package main
import "github.com/obowersa/wfwiki/pkg/wfwiki"

func main (){
	wf := wfwiki.NewWFWiki()
	wf.GetStats("warframe","Ash")
}
```
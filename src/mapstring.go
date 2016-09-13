
package main

import (
	"encoding/json"
	"fmt"
)

type MapContainer struct {
	M map[string]string
}

func main () {
	j := []byte(`{"a": "b", "c": "d"}`)

	var interfaceVal interface{}
	json.Unmarshal(j, &interfaceVal)

	// At this point, you have an interface{} whose true type is
	// map[string]interface{}. Note fmt.Printf "reads through"
	// the interface{} to the underlying value, using the reflection
	// capabilities:
	fmt.Printf("interfaceVal's type: %T\n", interfaceVal)

	strings := interfaceVal.(map[string]interface{})

	// but the values are still interface{}:
	valueA := strings["a"]
	// this panics with
	// "invalid operation: valueA + "\n" (mismatched types interface {} and string)"
	// fmt.Printf(valueA + "\n")

	// you would need:
	fmt.Printf(valueA.(string) + "\n")
	// which printed b

	// Now what you're looking for is a map[string]string, but we can not
	// go straight there from here:
	// s2 := interfaceVal.(map[string]string)
	// yields a run-time panic of
	// "panic: interface conversion: interface is map[string]interface {}, not map[string]string"

	// I'm not 100% sure, but I believe Go is entirely invariant in computer
	// science terminology (see
	// http://en.wikipedia.org/wiki/Covariance_and_contravariance_%28computer_science%29 ),
	// so for a type assertion to succeed, you need an exact match.
	// This is actually the easiest case to understand, which is probably why it is
	// what Go uses. Suppose my JSON was instead {"a": "b", "c": 1}... what would
	// asserting it as map[string]string _do_ with c? Silently discard? Error?
	// There isn't really a right answer, AT THE LANGUAGE LEVEL. In a moment, we're
	// going to observe that the JSON library has made a different choice, but at
	// a different abstraction level.

	// All that said, you now have a couple of choices. You could:
	var stringMap map[string]string
	json.Unmarshal(j, &stringMap)
	fmt.Printf("String map: %#v\n", stringMap)

	// or
	var mapContainer MapContainer
	json.Unmarshal(j, &mapContainer.M)
	fmt.Printf("mapContainer: %#v\n", mapContainer)

	// Two more notes. To directly get the map container, you need:
	var mapContainer2 MapContainer
	json.Unmarshal([]byte(`{"m": {"a": "b", "c": "d"}}`), &mapContainer2)
	fmt.Printf("mapContainer2: %#v\n", mapContainer2)
	// Note I had to change to a capital M for that, to make it public.

	// And finally, note the JSON library WILL just discard things:
	var mapContainer3 MapContainer
	json.Unmarshal([]byte(`{"m": {"a": "b", "c": "d", "e": 3.2}}`), &mapContainer3)
	fmt.Printf("mapContainer3: %#v\n", mapContainer3)
	// but that is a choice made by the JSON library, rather than the language.
}

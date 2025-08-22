package script

var OnTick func()

//export tick
func tick() {
	OnTick()
}

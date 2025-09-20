package env

func IsRC() bool {
	return isRC()
}

func ModuleSlotIndex() (int, bool) {
	index := moduleSlotIndex()
	if index < 0 {
		return 0, false
	}

	return index, true
}

func ModuleSlotCockpitIndex() (int, bool) {
	index := moduleSlotCockpitIndex()
	if index < 0 {
		return 0, false
	}

	return index, true
}

func ModuleSlotIndexInClassGroup() (int, bool) {
	index := moduleSlotIndexInClassGroup()
	if index < 0 {
		return 0, false
	}

	return index, true
}

//go:wasm-module env
//export is_rc
func isRC() bool

//go:wasm-module env
//export module_slot_index
func moduleSlotIndex() int

//go:wasm-module env
//export module_slot_cockpit_index
func moduleSlotCockpitIndex() int

//go:wasm-module env
//export module_slot_index_in_class_group
func moduleSlotIndexInClassGroup() int

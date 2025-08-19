package vehicle

const (
	VehicleErrorNotFound = VehicleError(256)
	BogieNotFound        = VehicleError(512)
	AxleNotFound         = VehicleError(1024)
	CouplingNotFound     = VehicleError(2048)
	PantographNotFound   = VehicleError(4096)
	UnknownError         = VehicleError(0)

	CouplingFront = 0
	CouplingBack  = 1

	SurfaceTypeGravel = 0
	SurfaceTypeStreet = 1
	SurfaceTypeGrass  = 2

	RailQualitySmooth          = 0
	RailQualityRough           = 1
	RailQualityFroggySmooth    = 2
	RailQualityFroggyRough     = 3
	RailQualityFlatGroove      = 4
	RailQualityHighSpeedSmooth = 5
	RailQualitySmoothDirt      = 6
	RailQualityRoughDirt       = 7
)

type Bogie uint32
type Axle struct {
	index int
	bogie Bogie
}
type Pantograph uint32
type VehicleError uint32

func (e VehicleError) Error() string {
	switch e {
	case VehicleErrorNotFound:
		return "vehicle not found"
	case BogieNotFound:
		return "bogie not found"
	case AxleNotFound:
		return "axle not found"
	case CouplingNotFound:
		return "coupling not found"
	case PantographNotFound:
		return "pantograph not found"
	case UnknownError:
		return "unknown error"
	}

	return "unknown error"
}

func GetBogie(index int) (Bogie, error) {
	ret := bogieIsValid(uint32(index))

	if ret > 255 {
		return 0, VehicleError(ret)
	}

	return Bogie(index), nil
}

func (b Bogie) GetAxle(index int) (Axle, error) {
	ret := axleIsValid(uint32(b), uint32(index))

	if ret > 255 {
		return Axle{}, VehicleError(ret)
	}

	return Axle{index: index, bogie: b}, nil
}

func (a Axle) SetTractionForceNewton(value float32) {
	setTractionForceNewton(uint32(a.bogie), uint32(a.index), value)
}

func (a Axle) SetBrakeForceNewton(value float32) {
	setBrakeForceNewton(uint32(a.bogie), uint32(a.index), value)
}

func (a Axle) RailQuality() uint32 {
	return railQuality(uint32(a.bogie), uint32(a.index))
}

func (a Axle) SurfaceType() uint32 {
	return surfaceType(uint32(a.bogie), uint32(a.index))
}

func (a Axle) InverseRadius() float32 {
	return inverseRadius(uint32(a.bogie), uint32(a.index))
}

func (b Bogie) SetRailBrakeForceNewton(value float32) {
	setRailBrakeForceNewton(uint32(b), value)
}

func GetPantograph(index int) (Pantograph, error) {
	ret := pantographIsValid(uint32(index))

	if ret > 255 {
		return 0, VehicleError(ret)
	}

	return Pantograph(index), nil
}

func (p Pantograph) Height() float64 {
	return pantographHeight(uint32(p))
}

func (p Pantograph) Voltage() float64 {
	return pantographVoltage(uint32(p))
}

func IsCoupled(end int) bool {
	return isCoupled(uint32(end)) == 1
}

//go:wasm-module vehicle
//export is_coupled
func isCoupled(bogie uint32) uint32

//go:wasm-module vehicle
//export bogie_is_valid
func bogieIsValid(bogie uint32) uint32

//go:wasm-module vehicle
//export axle_is_valid
func axleIsValid(bogie, axle uint32) uint32

//go:wasm-module vehicle
//export pantograph_is_valid
func pantographIsValid(end uint32) uint32

//go:wasm-module vehicle
//export rail_quality
func railQuality(bogie, axle uint32) uint32

//go:wasm-module vehicle
//export surface_type
func surfaceType(bogie, axle uint32) uint32

//go:wasm-module vehicle
//export inverse_radius
func inverseRadius(bogie, axle uint32) float32

//go:wasm-module vehicle
//export velocity_vs_ground
func VelocityVsGround() float64

//go:wasm-module vehicle
//export acceleration_vs_ground
func AccelerationVsGround() float32

//go:wasm-module vehicle
//export pantograph_height
func pantographHeight(pantograph uint32) float64

//go:wasm-module vehicle
//export pantograph_voltage
func pantographVoltage(pantograph uint32) float64

//go:wasm-module vehicle
//export set_traction_force_newton
func setTractionForceNewton(bogie, axle uint32, value float32)

//go:wasm-module vehicle
//export set_brake_force_newton
func setBrakeForceNewton(bogie, axle uint32, value float32)

//go:wasm-module vehicle
//export set_rail_brake_force_newton
func setRailBrakeForceNewton(bogie uint32, value float32)

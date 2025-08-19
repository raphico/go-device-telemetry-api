package device

type DeviceType string

const (
	DeviceTypeTemperatureSensor DeviceType = "temperature_sensor"
	DeviceTypeHumiditySensor    DeviceType = "humidity_sensor"
	DeviceTypeMotionSensor      DeviceType = "motion_sensor"
)

var validDeviceTypes = map[DeviceType]struct{}{
	DeviceTypeTemperatureSensor: {},
	DeviceTypeHumiditySensor:    {},
	DeviceTypeMotionSensor:      {},
}

func NewDeviceType(value string) (DeviceType, error) {
	dt := DeviceType(value)
	if _, ok := validDeviceTypes[dt]; !ok {
		return "", ErrInvalidDeviceType
	}
	return dt, nil
}

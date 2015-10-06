package main

import (
	"reflect"

	"github.com/spf13/viper"
	"github.com/spf13/cast"
)

type Device struct {
	name string
	enabled bool
	protocol string
	port int
	rule string
}

type Configuration struct {
	devices []Device
}

func (conf *Configuration) Load(configfile *string) {
	conf.readConfig(configfile)

	conf.devices = make([]Device, 0, 100)
	for key, value := range viper.Get("devices").(map[string]interface{}) {
		d := conf.bindDevice(key, value)
		conf.devices = append(conf.devices, *d)
	}
}

func (conf *Configuration) readConfig(configfile *string) {
	log.Info("Loading config file: ", configfile)

	viper.SetConfigFile(*configfile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Cannot load config file")
	}
}

func (conf *Configuration) bindDevice(key string, value interface{}) *Device {
	if reflect.TypeOf(value).Kind() != reflect.Map {
		return nil
	}

	device := new(Device)
	device.name = key
	device.enabled = cast.ToBool(value.(map[string]interface{})["enabled"])
	device.protocol = cast.ToString(value.(map[string]interface{})["protocol"])
	device.port = cast.ToInt(value.(map[string]interface{})["port"])
	device.rule = cast.ToString(value.(map[string]interface{})["rule"])

	return device
}

func (conf *Configuration) GetDevices() []Device {
	return conf.devices
}

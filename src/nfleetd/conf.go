package main

import (
	"reflect"
	"github.com/spf13/viper"
	"github.com/spf13/cast"
)

type Worker struct {
	thread int
}

type Device struct {
	name string
	enabled bool
	protocol string
	address string
	port int
	rule string
}

type Configuration struct {
	worker *Worker
	devices []Device
}

func (conf *Configuration) Load(configfile *string) {
	conf.readConfig(configfile)

	conf.loadWorker()
	conf.loadDevices()
}

func (conf *Configuration) readConfig(configfile *string) {
	log.Info("Loading config file: ", configfile)

	viper.SetConfigFile(*configfile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Cannot load config file")
	}
}

func (conf *Configuration) loadWorker() {
	conf.worker = new(Worker)
	conf.worker.thread = cast.ToInt(viper.Get("worker.thread"))
}

func (conf *Configuration) loadDevices() {
	conf.devices = make([]Device, 0, 100)
	for key, value := range viper.Get("devices").(map[string]interface{}) {
		d := conf.bindDevice(key, value)
		conf.devices = append(conf.devices, *d)

		log.Debug("Loaded configuration: ", d)
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
	device.address = cast.ToString(value.(map[string]interface{})["address"])
	device.port = cast.ToInt(value.(map[string]interface{})["port"])
	device.rule = cast.ToString(value.(map[string]interface{})["rule"])

	return device
}

func (conf *Configuration) getWorker() *Worker {
	return conf.worker
}

func (conf *Configuration) GetDevices() []Device {
	return conf.devices
}

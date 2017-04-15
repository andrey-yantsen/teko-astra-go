package main

import (
	"flag"

	"fmt"

	astra_l "github.com/andrey-yantsen/teko-astra-go"
)

func main() {
	serial := flag.String("serial", "/dev/ttyAPP1", "serial port device (/dev/...)")
	address := flag.Uint("address", 255, "device address")
	action_find := flag.Bool("find-device", false, "run FindDevice command")
	action_register := flag.Bool("register-device", false, "run RegisterDevice command")
	action_delete := flag.Bool("delete-device", false, "run DeleteDevice command")
	action_createl2net := flag.Uint("create-l2-net", 0, "(re-)create Level2 net")
	action_registerl2dev := flag.Bool("register-l2-device", false, "register new Level2 device")
	action_getstate := flag.Bool("get-state", false, "get RI-M status")
	action_getevents := flag.Bool("get-events", false, "get RI-M events")
	flag.Parse()
	if driver, err := astra_l.Connect(*serial); err != nil {
		panic(err)
	} else {
		device := driver.GetDevice(uint16(*address))
		if *action_find {
			if f, err := device.FindDevice(); err != nil {
				panic(err)
			} else {
				fmt.Printf("FindDevice response: %+v\n", f)
			}
		}
		if *action_delete {
			if err := device.DeleteDevice(); err != nil {
				panic(err)
			} else {
				fmt.Printf("Device %+v deleted\n", *address)
			}
		}
		if *action_register {
			broadcastDevice := driver.GetDevice(255)
			if f, err := broadcastDevice.FindDevice(); err != nil {
				panic(err)
			} else if err := device.RegisterDevice(f.EUI.GetShortDeviceEUI()); err != nil {
				panic(err)
			} else {
				fmt.Printf("Device %+v registered\n", *address)
			}
		}
		if *action_getstate {
			if s, err := device.GetState(); err != nil {
				panic(err)
			} else {
				fmt.Printf("Device state: %+v\n", s)
			}
		}
		if *action_getevents {
			if e, err := device.GetEvents(); err != nil {
				panic(err)
			} else {
				fmt.Printf("Device events: %+v\n", e)
			}
		}
		if *action_createl2net > 0 {
			if *action_createl2net > 3 {
				panic("channel number should be less than 4")
			}
			if err := device.CreateLevel2Net(uint8(*action_createl2net)); err != nil {
				panic(err)
			} else {
				println("Net successfully registered")
			}
		}
		if *action_registerl2dev {
			if dev, err := device.RegisterLevel2Device(0); err != nil {
				panic(err)
			} else {
				fmt.Printf("Registered new L2 device: %+v\n", dev)
			}
		}
	}
}

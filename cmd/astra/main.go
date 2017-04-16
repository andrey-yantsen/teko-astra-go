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
	action_deletel2dev := flag.Uint("delete-l2-device", 0, "delete Level2 device")
	action_deletel2all := flag.Bool("delete-l2-all", false, "delete all Level2 devices")
	action_getstate := flag.Bool("get-state", false, "get RI-M status")
	action_getevents := flag.Bool("get-events", false, "get RI-M events")
	action_getl2netconfig := flag.Bool("get-net-l2-config", false, "get L2 net config")
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
				fmt.Printf("Got %d device events\n", len(e))
				for _, ev := range e {
					switch ev.(type) {
					case astra_l.EventNoLink:
						fmt.Printf("NoLink event: %+v\n", ev)
					case astra_l.EventSStateRtm:
						fmt.Printf("EventSStateRtm event: %+v\n", ev)
					case astra_l.EventSStateRtmLC:
						fmt.Printf("EventSStateRtmLC event: %+v\n", ev)
					case astra_l.EventSStateBrr:
						fmt.Printf("EventSStateBrr event: %+v\n", ev)
					case astra_l.EventSStateRimRtr:
						fmt.Printf("EventSStateRimRtr event: %+v\n", ev)
					case astra_l.EventSStateKeychain:
						fmt.Printf("EventSStateKeychain event: %+v\n", ev)
					case astra_l.EventSStateOtherWithNoData:
						fmt.Printf("EventSStateOtherWithNoData event: %+v\n", ev)
					case astra_l.EventSStateOtherWithPower:
						fmt.Printf("EventSStateOtherWithPower event: %+v\n", ev)
					case astra_l.EventSStateOtherWithTemperature:
						fmt.Printf("EventSStateOtherWithTemperature event: %+v\n", ev)
					case astra_l.EventSStateOtherWithSmoke:
						fmt.Printf("EventSStateOtherWithSmoke event: %+v\n", ev)
					case astra_l.EventRRState:
						fmt.Printf("EventRRState event: %+v\n", ev)
					default:
						fmt.Printf("Got some strange event: %+v\n", ev)
					}
				}
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
		if *action_deletel2dev > 0 {
			if err := device.DeleteLevel2Device(uint16(*action_deletel2dev)); err != nil {
				panic(err)
			} else {
				println("L2 device deleted")
			}
		}
		if *action_deletel2all {
			if err := device.DeleteAllLevel2Devices(); err != nil {
				panic(err)
			} else {
				println("All L2 devices deleted")
			}
		}
		if *action_getl2netconfig {
			if c, err := device.GetNetLevel2Config(); err != nil {
				panic(err)
			} else {
				fmt.Printf("NET config: %+v\n", c)
			}
		}
	}
}

```                                              
   __     ___              ___   _____     ___   
 /'_ `\  / __`\  _______  / __`\/\ '__`\  /'___\ 
/\ \L\ \/\ \L\ \/\______\/\ \L\ \ \ \L\ \/\ \__/ 
\ \____ \ \____/\/______/\ \____/\ \ ,__/\ \____\
 \/___L\ \/___/           \/___/  \ \ \/  \/____/
   /\____/                         \ \_\         
   \_/__/                           \/_/         

```

## what

A golang implementation of the Open Pixel Control protocol.

## usage

```
package main

import("github.com/kellydunn/go-opc")

func main {
     // Setup a new server
     s := opc.NewServer()

     // Register your devices (where r is an implementation of opc.Device)
     s.RegisterDevice(r)

     // Listen for incoming messages and serve them accordingly
     go s.ListenOnPort("tcp", "7890")
     go s.Process()

     // Create a client
     c := opc.NewClient("tcp", "localhost:7890")

     // Make a message!
     // This creates a message to send on channel 0
     // Or according to the OPC spec, a Broadcast Message.
     m := opc.NewMessage(0)  

     // Set pixel #1 to white.
     m.SetPixelColors(1, 255, 255, 255)
     
     // Send the message!
     c.Send(m)

     // The first pixel of all registered devices should be white!
}    

```

## design

The applications of OPC are not currently tied to any single communication model, and it is currently unclear if there is any canonical method of dispatching OPC messages.  So, when using this library, it is encouraged to implement the `opc.Device` interface such that you can further define the details of your devices and how they should be written to.

A very simple implementation of the `opc.Device` interface could be the following:

```
type DummyDevice struct {
     conn net.Conn
     channel uint8
}

// Simple write behavior.  Write the OPC Message over a network connection.
func (d *DummyDevice) Write(m *opc.Message) error {
     _, err := conn.Write(m.byteArray)
     if err != nil {
        return err
     }
     
     return nil
}

// Simple Channel getter.  Return the channel in which to associate this device.
func (d *DummyDevice) Channel() uint8 {
     return channel
}
```

## related work

  - Open Pixel Protocol (http://openpixelcontrol.org/)
  - Fadecandy OPC implementation (https://github.com/scanlime/fadecandy/tree/master/server/src)
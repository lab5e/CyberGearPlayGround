# "Hello CyberGear"

# Hardware required

* Xiaomi CyberGear Motor.
* [CANable-MKS](https://canable.io/getting-started.html) USB to CAN adapter. Adapter is also available from [AliExpress](https://www.aliexpress.com/w/wholesale-MKS-canable-pro.html?spm=a2g0o.productlist.search.0).
* 24V Power Supply.
* [Amass XT30(2+2)-F](https://www.china-amass.com/product/contain/1Yf5h7G4u1927079) cable. Cable with female connector is available from [AliExpress](https://www.aliexpress.com/w/wholesale-XT30(2%2B2)%2525252dF.html?spm=a2g0o.home.search.0)

# Electrical connections

![pinout](pictures/pinout.png)

* Power supply: BAT+ / GND
* MKS Canable: CAN_L / CAN_H ("G" on MKS Canable can be left unconnected)

# Protocol

## CAN
The CyberGear CAN protocol is described in the translated [manual](../manual/Translated%20copy%20%20of%20CyberGear%20micromotor%20instruction%20manual.pdf).

## SLCAN

Testing via a PC/MAC etc can be done via this gocg utility and a Serial Line CAN <-> CAN transciever like [CANable-MKS](https://github.com/makerbase-mks/CANable-MKS)

The frame format is :

| byte | description | example |
|------|-------------|---------|
| 1    | Frame type identifier | 't' : Standard (11-bit) CAN frame.| 
|     |   |  'T' : 'T': Extended (29-bit) CAN frame (use 'T' for talking to the CyberGear)| 
|     |   |  'r' :  Standard CAN remote frame. |     
|     |   |  'R' : Extended CAN remote frame.|     
| 2-17 | CAN id (00000000-1FFFFFFF) | 8 bytes ascii encoded. example: 0xFF will be transmitted as "FF" |
| 18   | DLC (Data length) | Number of "dd" pairs must match the data length |
| [19-2*DLC] | Data | Payload|

Please refer to the translated [manual](../manual/Translated%20copy%20%20of%20CyberGear%20micromotor%20instruction%20manual.pdf) for a description of CANid/ frame header and data encoding.

Before sending CyberGear specific frames, the following SLCAN commands have to be sent:

1. S8\<CR\> - Set CAN bitrate to 1MBit
1. O\<CR\> - Open the CAN channel in normal mode (send and receive)

Before closing the session, it might be a good idea to send

1. C<CR> - Close the CAN channel




# Go

## Currently implemented commands


|command| arguments | example |description|
|-------|----------------|-----|-----|
|open   | \<serialport\> | open COM15 |Opens the serial port for reading and writing.|
|close  |                | close|Closes the currently open serial port.|
|enable | \<motor id\>   |enable 7F| Enable a motor (7F is the default cybergear id).|
|disable| \<motor id\>   | disable 7F|Disables / stops the motor.|
|speed  | \<motor id\> \<speed\>|speed 7F 2.2| Sets motor speed (rad/s). Valid speed settings are in the range [-30, 30]|


## Examples

1. Open serial port, enable motor with CAN id 0x7F (default) and setting the speed to 10 .
```
> open COM15
> enable 7F
> set_speed 7F 10
```

2. Stop motor and close serial port.
```
> disable 7F
> close
```

>**Work in progress**:  
>- Parse motor status response frames (current speed, position, torque etc)
>- set current / position mode
>- set current / speed limit
>- set position
>- change CAN id



 
# Python
(Python example is courtesy of Hans Elias at HH)

* pip install python-can (use socketcan or slcan interface).


# Useful stuff

* [cangaroo](https://github.com/normaldotcom/cangaroo/) open source can bus analyzer with support for transmit/receive of standard and FD frames and DBC decoding of incoming frames.
* [canable](https://canable.io/getting-started.html) Getting started with MKS Canable.

# Protocol implementations

* [M5 Stack](https://github.com/project-sternbergia/cybergear_m5) (C++)
* [CyberGearKit](https://github.com/CmST0us/CyberGearKit) (Swift)

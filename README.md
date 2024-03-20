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

# Go

## Currently implemented commands


|command| arguments | example |description|
|-------|----------------|-----|-----|
|open   | \<serialport\> | open COM15 |Opens the serial port for reading and writing.|
|close  |                | close|Closes the currently open serial port.|
|enable | \<motor id\>   |enable 7F| Enable a motor (7F is the default cybergear id).|
|disable| \<motor id\>   | disable 7F|Disables / stops the motor.|
|speed  | \<motor id\> \<speed\>|speed 7F 2.2| Sets motor speed (rad/s).|

>TODO:
>- read motor parameters (current speed / location / current etc)
>- current / position mode
>- more stuff...


# Python
(Python example is courtesy of Hans Elias at HH)

* pip install python-can (use socketcan or slcan interface).


# Useful stuff

* [cangaroo](https://github.com/normaldotcom/cangaroo/) open source can bus analyzer with support for transmit/receive of standard and FD frames and DBC decoding of incoming frames.
* [canable](https://canable.io/getting-started.html) Getting started with MKS Canable.

# Protocol implementations

* [M5 Stack](https://github.com/project-sternbergia/cybergear_m5) (C++)
* [CyberGearKit](https://github.com/CmST0us/CyberGearKit) (Swift)

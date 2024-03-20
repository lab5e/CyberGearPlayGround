import can
import struct
import math
import time

def id_fields(mode, a, b):
    return (b & 0xff) | ((a & 0xffff) << 8) | ((mode & 0b11111) << 24)

def float_to_u16(val):
    mapped = int(val * (2**16))
    if mapped < 0: mapped = 0
    if mapped > 65535: mapped = 65535
    return mapped

def float_to_i16(val):
    return float_to_u16((val + 1) / 2.0) - 32768

class Motor:
    REG_RUN_MODE = 0x7005
    REG_IQ_REF = 0x7006
    REG_SPD_REF = 0x700a
    REG_LIMIT_TORQUE = 0x700b
    REG_CUR_KP = 0x7010
    REG_CUR_KI = 0x7011
    REG_CUR_FILT_GAIN = 0x7014
    REG_LOC_REF = 0x7016
    REG_LIMIT_SPD = 0x7017
    REG_LIMIT_CUR = 0x7018
    REG_MECH_POS = 0x7019
    REG_IQF = 0x701a
    REG_MECHVEL = 0x701b
    REG_VBUS = 0x701c
    REG_ROTATION = 0x701d
    REG_LOC_KP = 0x701e
    REG_SPD_KP = 0x701f
    REG_SPD_KI = 0x7020

    RUN_MODE_OPERATIONAL = 0
    RUN_MODE_POSITION = 1
    RUN_MODE_SPEED = 2
    RUN_MODE_CURRENT = 3

    def __init__(self, bus, motor_id):
        self.bus = bus
        self.motor_id = motor_id

    def enable(self):
        msg = can.Message(
            arbitration_id=id_fields(3, 0, self.motor_id),
            # arbitration_id=id_fields(0, 0, self.motor_id),
            data=[],
            is_extended_id=True
        )
        try:
            self.bus.send(msg)
            print("Message sent on {}".format(self.bus.channel_info))
            print(msg)
        except can.CanError:
            print("Message NOT sent")

    def disable(self):
        msg = can.Message(
            arbitration_id=id_fields(4, 0, self.motor_id),
            data=[],
            is_extended_id=True
        )
        print(msg)
        self.bus.send(msg)

    def control(self, torque, target_angle, target_velocity, k_p, k_d):
        """
        torque: Torque to apply in Nm. From -12N to 12N.
        target_angle: Angle to move to, from -4pi to 4pi.
        target_velocity: Angular velocity to move to target. From -30rad/s to 30rad/s.
        k_p: Proportional coeficient to Position PID controller. From 0.0 to 500.0.
        k_d: Derivative coeficient to Position PID controller. From 0.0 to 5.0.
        """

        target_angle_num = float_to_u16(((target_angle / math.pi) + 4) / 8)
        target_velocity_num = float_to_u16(((target_velocity / math.pi) + 30) / 60)
        k_p_num = float_to_u16(k_p / 500)
        k_d_num = float_to_u16(k_d / 5)

        msg = can.Message(
            arbitration_id=id_fields(1, float_to_i16(torque * 12), self.motor_id),
            data=struct.pack(
                "<HHHH",
                target_angle_num,
                target_velocity_num,
                k_p_num,
                k_d_num
            ),
            is_extended_id=True
        )
        printf("XXX\n")
        self.bus.send(msg)

    def param_write_typ(self, reg, typ, val):
        msg = can.Message(
            arbitration_id=id_fields(18, 0, self.motor_id),
            data=struct.pack("<Hxx" + typ, reg, val),
            is_extended_id=True
        )
        print("Message:")
        print(msg)
        self.bus.send(msg)

    def param_write_uint(self, reg, val):
        return self.param_write_typ(reg, "I", val)

    def param_write_int(self, reg, val):
        return self.param_write_typ(reg, "i", val)

    def param_write_float(self, reg, val):
        return self.param_write_typ(reg, "f", val)

#with can.Bus(interface='socketcan', channel='can0') as bus:
#with can.interface.Bus(bustype='slcan', channel='COM1', baudrate='500000') as bus:
with can.interface.Bus(bustype="slcan", channel="COM15", baudrate='115200') as bus:
    print(bus)
    print (bus.state)

    m = Motor(bus, 0x7F)

    try:

        time.sleep(1)

        m.enable()

        print("Enabled")


        # Set to current control, set current ref
        #m.param_write_uint(m.REG_RUN_MODE, m.RUN_MODE_CURRENT)
        #m.param_write_float(m.REG_IQ_REF, 0.2)


        # # Set to position control, apply speed limit, go to location
        # m.param_write_uint(m.REG_RUN_MODE, m.RUN_MODE_POSITION)
        # m.param_write_float(m.REG_LIMIT_SPD, 30)
        # m.param_write_float(m.REG_LOC_REF, -3.1415 * 3 * 3.1 * 1)

        # Set to position control, apply speed limit, go to location
        m.param_write_uint(m.REG_RUN_MODE, m.RUN_MODE_SPEED)
        m.param_write_float(m.REG_SPD_REF, 2.1)


        while True:
            time.sleep(1)

    finally:
        m.disable()
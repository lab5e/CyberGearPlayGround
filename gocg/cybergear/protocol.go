package cybergear

import (
	"fmt"
	"math"
)

const MAX_CAN_ID = 0x7F
const EXTENDED_FRAME_TYPE = 'T'

type communicationType int

const (
	COMMUNICATION_FETCH_DEVICE_ID              communicationType = 0  // Get the device ID (communication type 0); get the device ID and 64-bit MCU unique identifier
	COMMUNICATION_MOTION_CONTROL_COMMAND       communicationType = 1  // Operation control mode motor control instructions (communication type 1) are used to send control instructions to the motor
	COMMUNICATION_STATUS_REPORT                communicationType = 2  // Motor feedback data (communication type 2) is used to feedback the motor operating status to the host
	COMMUNICATION_ENABLE_DEVICE                communicationType = 3  // Motor enable operation (communication type 3)
	COMMUNICATION_DISABLE_DEVICE               communicationType = 4  // Motor stopped (communication type 4)
	COMMUNICATION_SET_MECHANICAL_ZERO_POSITION communicationType = 6  // Setting the mechanical zero position of the motor (communication type 6) will set the current motor position to the mechanical zero position (lost after power failure)
	COMMUNICATION_SET_CAN_ID                   communicationType = 7  // Set motor CAN_ID (communication type 7) to change the current motor CAN_ID, which will take effect immediately.
	COMMUNICATION_READ_SINGLE_PARAM            communicationType = 17 // Single parameter reading (communication type 17)
	COMMUNICATION_WRITE_SINGLE_PARAM           communicationType = 18 // Single parameter writing (communication type 18) (lost after power failure)
	COMMUNICATION_ERROR_REPORT                 communicationType = 21 // Fault feedback frame (communication type 21)
)

type configParameter int

const ( // (Translated from chinese)
	CONFIG_WR_NAME             configParameter = 0x0000 // Motor name (Parameter Type: String)
	CONFIG_R_BAR_CODE          configParameter = 0x0001 // BarCode (Parameter Type: String)
	CONFIG_R_BOOT_CODE_VERSION configParameter = 0x1000 // BootCodeVersion (Parameter Type: String)
	CONFIG_R_BOOT_BUILD_DATE   configParameter = 0x1001 // BootBuildDate (Parameter Type: String)
	CONFIG_R_BOOT_BUILD_TIME   configParameter = 0x1002 // BootBuildTime (Parameter Type: String)
	CONFIG_R_APP_CODE_VERSION  configParameter = 0x1003 // AppCodeVersion (Parameter Type: String)
	CONFIG_R_APP_GIT_VERSION   configParameter = 0x1004 // AppGitVersion (Parameter Type: String)
	CONFIG_R_APP_BUILD_DATE    configParameter = 0x1005 // AppBuildDate  (Parameter Type: String)
	CONFIG_R_APP_BUILD_TIME    configParameter = 0x1006 // AppBuildTime	(Parameter Type: ?)
	CONFIG_R_APP_CODE_NAME     configParameter = 0x1007 // AppCodeName (Parameter Type: ?)
	CONFIG_R_ECHO_PARA1        configParameter = 0x2000 // echoPara1 (Parameter Type: ?)
	CONFIG_R_ECHO_PARA2        configParameter = 0x2001 // echoPara2 (Parameter Type: ?)
	CONFIG_R_ECHO_PARA3        configParameter = 0x2002 // echoPara3 (Parameter Type: ?)
	CONFIG_R_ECHO_PARA4        configParameter = 0x2003 // echoPara4 (Parameter Type: ?)
	CONFIG_WR_ECHO_FRE_HZ      configParameter = 0x2004 // echoFreHz (Parameter Type: ?)
	CONFIG_R_MECH_OFFSET       configParameter = 0x2005 // MechOffset Motor magnetic encoder angle offset (Parameter Type float [-7, 7])
	CONFIG_WR_MECH_POS_INIT    configParameter = 0x2006 // MechPos_init Reference angle for initial multi-turn (Parameter Type float [-50, 50])
	CONFIG_WR_LIMIT_TORQUE     configParameter = 0x2007 // Torque limit (Parameter Type float [0, 12])
	CONFIG_WR_I_FW_MAX         configParameter = 0x2008 // I_FW_MAX (Field weakening current value, default 0) (Parameter Type float [0, 33])
	CONFIG_WR_MOTOR_INDEX      configParameter = 0x2009 // Motor index, marks the motor joint position (Parameter Type uint8_t [0, 20])
	CONFIG_WR_CAN_ID           configParameter = 0x200a // CAN_ID (Node ID) (Parameter Type uint8_t [0, 127])
	CONFIG_WR_CAN_MASTER       configParameter = 0x200b // CAN_MASTER can host  id (Parameter Type uint8_t [0, 127])
	CONFIG_WR_CAN_TIMEOUT      configParameter = 0x200c // CAN_TIMEOUT can timeout threshold, default 0 (Parameter Type uint32_t [0, 10000])
	CONFIG_WR_MOTOR_OVER_TEMP  configParameter = 0x200d // Motor overtemp protection value, temp (degree) *10 (Parameter Type uint16_t [0, 1500])
	CONFIG_WR_OVER_TEMP_TIME   configParameter = 0x200e // Overtemperature time (Parameter Type uint32_t [0, 100000])
	CONFIG_WR_GEAR_RATIO       configParameter = 0x200f // GearRatio (Parameter Type float [1, 64])
	CONFIG_WR_TQ_CALI_TYPE     configParameter = 0x2010 // Tq_caliType Torque calibration method setting (Parameter Type uint8_t 0, 1)
	CONFIG_WR_CUR_FILT_GAIN    configParameter = 0x2011 // cur_filt_gain Current filter parameters (Parameter Type float [0, 1])
	CONFIG_WR_CUR_KP           configParameter = 0x2012 // cur_kp current kp (Parameter Type float [0, 200])
	CONFIG_WR_CUR_KI           configParameter = 0x2013 // cur_ki current ki (Parameter Type float [0, 200])
	CONFIG_WR_SPD_KP           configParameter = 0x2014 // spd_kp speed kp (Parameter Type float [0, 200])
	CONFIG_WR_SPD_KI           configParameter = 0x2015 // spd_ki speed ki (Parameter Type float [0, 200])
	CONFIG_WR_LOC_KP           configParameter = 0x2016 // loc_kp location kp (Parameter Type float [0, 200])
	CONFIG_WR_SPD_FILT_GAIN    configParameter = 0x2017 // spd_filt_gain Speed filter parameters (Parameter Type float [0, 1])
	CONFIG_WR_LIMIT_SPD        configParameter = 0x2018 // Position loop speed limit (Parameter Type float [0, 200])
	CONFIG_WR_LIMIT_CUR        configParameter = 0x2019 // Position speed control current limit (Parameter Type float [0, 27])
	CONFIG_R_TIME_USE0         configParameter = 0x3000 // timeUse0 (Parameter Type uint16_t)
	CONFIG_R_TIME_USE1         configParameter = 0x3001 // timeUse1 (Parameter Type uint16_t)
	CONFIG_R_TIME_USE2         configParameter = 0x3002 // timeUse2 (Parameter Type uint16_t)
	CONFIG_R_TIME_USE3         configParameter = 0x3003 // timeUse3 (Parameter Type uint16_t)
	CONFIG_R_ENCODER_RAW       configParameter = 0x3004 // Magnetic encoder raw sampling value (Parameter Type uint16_t)
	CONFIG_R_MCU_TEMP          configParameter = 0x3005 // MCU Internal temperature, *10 (Parameter Type uint16_t)
	CONFIG_R_MOTOR_TEMP        configParameter = 0x3006 // Motor NTC temperature *10 (Parameter Type uint16_t)
	CONFIG_R_VBUS_MV           configParameter = 0x3007 // Bus voltage (mv) (Parameter Type uint16_t)
	CONFIG_R_ADC1_OFFSET       configParameter = 0x3008 // Sampling Channel 1 Zero Current Bias (adc1 offset) (Parameter Type int32_t)
	CONFIG_R_ADC2_OFFSET       configParameter = 0x3009 // Sampling Channel 2 Zero Current Bias (adc2 offset) (Parameter Type int32_t)
	CONFIG_R_ADC1_RAW          configParameter = 0x300a // adc1Raw adc sample value 1 (Parameter Type uint32_t)
	CONFIG_R_ADC2_RAW          configParameter = 0x300b // adc2Raw adc Sample value 2 (Parameter Type uint32_t)
	CONFIG_R_VBUS_V            configParameter = 0x300c // VBUS voltage (Parameter Type float)
	CONFIG_R_CMD_ID            configParameter = 0x300d // cmdId id ring command, A (Parameter Type float)
	CONFIG_R_CMD_IQ            configParameter = 0x300e // cmdIq iq ring command, A (Parameter Type float)
	CONFIG_R_CMD_LOC_REF       configParameter = 0x300f // cmdlocref Position loop command, rad (Parameter Type float)
	CONFIG_R_CMD_SPD_REF       configParameter = 0x3010 // cmdspdref Speed loop command, rad/s (Parameter Type float)
	CONFIG_R_CMD_TORQUE        configParameter = 0x3011 // cmdTorque Torque command, nm (Parameter Type float)
	CONFIG_R_CMD_POS           configParameter = 0x3012 // cmdPos mit protocol angle command (Parameter Type float)
	CONFIG_R_CMD_VEL           configParameter = 0x3013 // cmdVel mit protocol speed command (Parameter Type float)
	CONFIG_R_ROTATION          configParameter = 0x3014 // rotation Number of turns (Parameter Type int16_t)
	CONFIG_R_MOD_POS           configParameter = 0x3015 // Mechanical angle of the motor without counting revolutions, rad (Parameter Type float)
	CONFIG_R_MECH_POS          configParameter = 0x3016 // mechPos Load end lap mechanical angle, rad (Parameter Type float)
	CONFIG_R_MECH_VEL          configParameter = 0x3017 // mechVel Load end speed,rad/s (Parameter Type float)
	CONFIG_R_ELEC_POS          configParameter = 0x3018 // elecPos electrical angle (Parameter Type float)
	CONFIG_R_IA                configParameter = 0x3019 // ia U line current, A (Parameter Type float)
	CONFIG_R_IB                configParameter = 0x301a // ib V line current, A (Parameter Type float)
	CONFIG_R_IC                configParameter = 0x301b // ic W line current, A (Parameter Type float)
	CONFIG_R_TICK              configParameter = 0x301c // tick (Parameter Type uint32_t)
	CONFIG_R_PHASE_ORDER       configParameter = 0x301d // phaseOrder Calibration direction mark (Parameter Type uint8_t)
	CONFIG_R_IQF               configParameter = 0x301e // iqf iq filter value, A (Parameter Type float)
	CONFIG_R_BOARD_TEMP        configParameter = 0x301f // boardTemp board temperature, *10 (Parameter Type int16_t)
	CONFIG_R_IQ                configParameter = 0x3020 // iq original value, A (Parameter Type float)
	CONFIG_R_ID                configParameter = 0x3021 // id original value, A (Parameter Type float)
	CONFIG_R_FAULT_STATUS      configParameter = 0x3022 // fault status value (Parameter Type uint32_t)
	CONFIG_R_WARN_STATUS       configParameter = 0x3023 // warning status value (Parameter Type uint32_t)
	CONFIG_R_DRV_FAULT         configParameter = 0x3024 // Driver chip fault value (Parameter Type uint16_t)
	CONFIG_R_DRV_TEMP          configParameter = 0x3025 // driver chip temperature value, degrees (Parameter Type int16_t)
	CONFIG_R_UQ                configParameter = 0x3026 // Uq q-axis voltage  (Parameter Type float)
	CONFIG_R_UD                configParameter = 0x3027 // Ud d-axis voltage (Parameter Type float)
	CONFIG_R_DTC_U             configParameter = 0x3028 // U phase output duty cycle (Parameter Type float)
	CONFIG_R_DTC_V             configParameter = 0x3029 // V phase output duty cycle (Parameter Type float)
	CONFIG_R_DTC_W             configParameter = 0x302a // W phase output duty cycle (Parameter Type float)
	CONFIG_R_CLOSED_LOOP_V_BUS configParameter = 0x302b // vbus in closed loop (Parameter Type)
	CONFIG_R_CLOSED_LOOP_V_REF configParameter = 0x302c // closed loop vq, vd combined voltage (Parameter Type float)
	CONFIG_R_TORQUE_FDB        configParameter = 0x302d // torque feedback value, nm (Parameter Type float)
	CONFIG_R_RATED_I           configParameter = 0x302e // motor rated current (Parameter Type float)
	CONFIG_R_LIMIT_I           configParameter = 0x302f // Motor limits maximum current (Parameter Type float)
)

type runModeType uint16

const (
	OPEARATION_CONTROL_MODE runModeType = 0x00
	LOCATION_MODE           runModeType = 0x01
	SPEED_MODE              runModeType = 0x02
	CURRENT_MODE            runModeType = 0x03
)

type motorParameterIndex uint16 // Read write

const (
	PARAMETER_RUN_MODE      motorParameterIndex = 0x7005 // R/W
	PARAMETER_IQ_REF        motorParameterIndex = 0x7006 // R/W
	PARAMETER_SPD_REF       motorParameterIndex = 0x700A // R/W
	PARAMETER_IMIT_TORQUE   motorParameterIndex = 0x700B // R/W
	PARAMETER_CUR_KP        motorParameterIndex = 0x7010 // R/W
	PARAMETER_CUR_KI        motorParameterIndex = 0x7011 // R/W
	PARAMETER_CUR_FILT_GAIN motorParameterIndex = 0x7014 // R/W
	PARAMETER_LOC_REF       motorParameterIndex = 0x7016 // R/W
	PARAMETER_LIMIT_SPD     motorParameterIndex = 0x7017 // R/W
	PARAMETER_LIMIT_CUR     motorParameterIndex = 0x7018 // R/W
	/*
		PARAMETER_MECH_POS		motorParameterIndex = 0x7019	// R
		PARAMETER_IQF			motorParameterIndex = 0x701A	// R
		PARAMETER_MECH_VEL		motorParameterIndex = 0x701B	// R
		PARAMETER_MECH_VBUS		motorParameterIndex = 0x701C	// R
		PARAMETER_MECH_ROTATION	motorParameterIndex = 0x701D	// R/W
		PARAMETER_LOC_KP		motorParameterIndex = 0x701E	// R/W
		PARAMETER_SPD_KP		motorParameterIndex = 0x701F	// R/W
		PARAMETER_SPD_KI		motorParameterIndex = 0x7020	// R/W
	*/
)

// CyberGear Communication mode 1 sends control instructions
type motionControl struct {
	id     byte    // Target motor CAN ID
	angle  float32 // Target angle (-4π~4π)
	speed  float32 //Target angular velocity (-30rad/s~30rad/s)
	kp     float32 // Kp (0.0~500.0)
	kd     float32 // Kd (0.0~5.0)
	torque float32 // Torque corresponding (-12Nm~12Nm)
}

// CyberGear Communication mode 17 single parameter reading
type singleParam struct {
	hostId    byte                // Host CAN Id
	motorId   byte                // Motor CAN Id
	parameter motorParameterIndex // 参数索引
	data      [4]byte
}

// CyberGear SLCAN frame
type SLCanFrame struct {
	header   [10]byte
	data     [16]byte
	checksum byte
}

func NewSLCanFrame() SLCanFrame {
	return SLCanFrame{
		header: [10]byte{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30},
		data:   [16]byte{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30},
	}
}

func (f *SLCanFrame) Serialize() []byte {
	buf := []byte{}
	buf = append(buf, f.header[:]...)

	if f.header[9] != 0x30 {
		buf = append(buf, f.data[:(f.header[9]-0x30)*2]...)
	}

	// For some reason, the SLCan implementation for MKS CANable and / or cybergear doesn't seem to be interested in the checksum
	// var checksum byte
	// for _, b := range buf {
	// 	checksum ^= byte(b)
	// }
	// checksumStr := fmt.Sprintf("%02X", checksum)
	// buf = append(buf, checksumStr[1])
	// buf = append(buf, checksumStr[0])
	return buf
}

// 4.1.4 Motor enable operation (communication type 3)
func EnableMotorCmd(hostId byte, motorId byte) (*SLCanFrame, error) {

	if hostId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid host Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	if motorId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid motor Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	hostIdString := fmt.Sprintf("%02X", hostId)
	motorIdString := fmt.Sprintf("%02X", motorId)
	communicationType := fmt.Sprintf("%02X", COMMUNICATION_ENABLE_DEVICE)

	frame := NewSLCanFrame()
	frame.header[0] = 'T' // Extended frame
	frame.header[1] = communicationType[0]
	frame.header[2] = communicationType[1]
	frame.header[5] = hostIdString[0]
	frame.header[6] = hostIdString[1]
	frame.header[7] = motorIdString[0]
	frame.header[8] = motorIdString[1]

	return &frame, nil
}

// 4.1.5 Motor stopped (communication type 4)
func DisableMotorCmd(hostId byte, motorId byte) (*SLCanFrame, error) {
	if hostId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid host Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	if motorId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid motor Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	hostIdString := fmt.Sprintf("%02X", hostId)
	motorIdString := fmt.Sprintf("%02X", motorId)
	communicationType := fmt.Sprintf("%02X", COMMUNICATION_DISABLE_DEVICE)

	frame := NewSLCanFrame()
	frame.header[0] = 'T' // Extended frame
	frame.header[1] = communicationType[0]
	frame.header[2] = communicationType[1]
	frame.header[5] = hostIdString[0]
	frame.header[6] = hostIdString[1]
	frame.header[7] = motorIdString[0]
	frame.header[8] = motorIdString[1]

	return &frame, nil
}

func SetRunMode(hostId byte, motorId byte, mode runModeType) (*SLCanFrame, error) {
	if hostId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid host Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	if motorId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid motor Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	hostIdString := fmt.Sprintf("%02X", hostId)
	motorIdString := fmt.Sprintf("%02X", motorId)
	communicationType := fmt.Sprintf("%02X", COMMUNICATION_WRITE_SINGLE_PARAM)

	frame := NewSLCanFrame()
	frame.header[0] = 'T' // Extended frame
	frame.header[1] = communicationType[0]
	frame.header[2] = communicationType[1]
	frame.header[5] = hostIdString[0]
	frame.header[6] = hostIdString[1]
	frame.header[7] = motorIdString[0]
	frame.header[8] = motorIdString[1]
	frame.header[9] = '8' // DLC

	index := fmt.Sprintf("%04X", PARAMETER_RUN_MODE)

	frame.data[0] = index[2]
	frame.data[1] = index[3]
	frame.data[2] = index[0]
	frame.data[3] = index[1]

	runMode := fmt.Sprintf("%02X", mode)

	frame.data[8] = runMode[0]
	frame.data[9] = runMode[1]
	frame.data[10] = '0'
	frame.data[11] = '0'
	frame.data[12] = '0'
	frame.data[13] = '0'
	frame.data[14] = '0'
	frame.data[15] = '0'

	return &frame, nil
}

// 4.1.9 Single parameter writing (communication type 18) (lost in case of power failure)
//
//	Parameter				Name			Description							Type	bytes	Unit						R/W
//	PARAMETER_IQ_REF 		iq_ref 			Current Mode Iq						float	4 		23~23A 						R/W
//	PARAMETER_SPD_REF		spd_ref 		Speed mode							float 	4		-30~30rad/s					R/W
//	PARAMETER_IMIT_TORQUE	imit_torque		Torque limit 						float 	4		0~12Nm 						R/W
//	PARAMETER_CUR_KP 		cur_kp 			Current Kp							float	4		Default value 0.125			R/W
//	PARAMETER_CUR_KI 		cur_ki 			Current Ki 							float 	4 		Default value 0.0158		R/W
//	PARAMETER_CUR_FILT_GAIN	cur_filt_gain	Current filter coefficient			float	4		0~1.0, default value 0.1	R/W
//	PARAMETER_LOC_REF 		loc_ref 		Position mode angle					float	4 		rad 						R/W
//	PARAMETER_LIMIT_SPD		limit_spd 		Location mode speed limit			float 	4 		0~30rad/s 					R/W
//	PARAMETER_LIMIT_CUR		limit_cur 		Speed Position mode Current limit	float 	4 		0~23A						R/W
func WriteParameterCmd(hostId byte, motorId byte, index motorParameterIndex, data float32) (*SLCanFrame, error) {
	if hostId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid host Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	if motorId > MAX_CAN_ID {
		return nil, fmt.Errorf("invalid motor Id (%d). Max Id is %d", hostId, MAX_CAN_ID)
	}

	hostIdString := fmt.Sprintf("%02X", hostId)
	motorIdString := fmt.Sprintf("%02X", motorId)
	communicationType := fmt.Sprintf("%02X", COMMUNICATION_WRITE_SINGLE_PARAM)

	frame := NewSLCanFrame()
	frame.header[0] = 'T' // Extended frame
	frame.header[1] = communicationType[0]
	frame.header[2] = communicationType[1]
	frame.header[5] = hostIdString[0]
	frame.header[6] = hostIdString[1]
	frame.header[7] = motorIdString[0]
	frame.header[8] = motorIdString[1]
	frame.header[9] = '8' // DLC

	indexString := fmt.Sprintf("%04X", index)

	frame.data[0] = indexString[2]
	frame.data[1] = indexString[3]
	frame.data[2] = indexString[0]
	frame.data[3] = indexString[1]

	data_int := math.Float32bits(data)
	incrediblyWeirdEncoding := fmt.Sprintf("%02X%02X%02X%02X", data_int>>24&0xFF, data_int>>16&0xFF, data_int>>8&0xFF, data_int&0xFF)

	frame.data[8] = incrediblyWeirdEncoding[6]
	frame.data[9] = incrediblyWeirdEncoding[7]
	frame.data[10] = incrediblyWeirdEncoding[4]
	frame.data[11] = incrediblyWeirdEncoding[5]
	frame.data[12] = incrediblyWeirdEncoding[2]
	frame.data[13] = incrediblyWeirdEncoding[3]
	frame.data[14] = incrediblyWeirdEncoding[0]
	frame.data[15] = incrediblyWeirdEncoding[1]

	return &frame, nil
}

// /* 对 CAN ID 区域特定的比特位进行赋值整数
//  * @param: frame 要设置的帧
//  * @param: bit_start 比特开始位
//  * @param: bit_length 比特长度
//  * @param: value 设置值(整数)
//  * */
// CYBERGEARAPI void cyber_gear_set_can_id_int_value(const cgFrame *frame, int bit_start, int bit_length, int value);

// /* 获取 CAN ID 区域特定的比特位数据
//  * @param: frame 要设置的帧
//  * @param: bit_start 比特开始位
//  * @param: bit_length 比特长度
//  * @return: int 值数据
//  * */
// CYBERGEARAPI int cyber_gear_get_can_id_int_value(const cgFrame *frame, int bit_start, int bit_length);

// /* 对 CAN ID 设置通讯类型
//  * @param: frame 要设置的帧
//  * @param: type 通信类型
//  * */
// CYBERGEARAPI void cyber_gear_set_can_id_communication_type(const cgFrame *frame, communicationType type);

// /* 对 CAN ID 设置主机CANID
//  * @param: frame 要设置的帧
//  * @param: value 主机CANID
//  * */
// CYBERGEARAPI void cyber_gear_set_can_id_host_can_id(const cgFrame *frame, int value);

// /* 对 CAN ID 设置目标电机CANID
//  * @param: frame 要设置的帧
//  * @param: value 电机CANID
//  * */
// CYBERGEARAPI void cyber_gear_set_can_id_target_can_id(const cgFrame *frame, int value);

// /* 构造一个参数写入的CAN包 (通信类型18), 参数值为整数
//  * @param: frame 要设置的帧
//  * @param: index 参数index
//  * @param: value 参数值
//  * */
// CYBERGEARAPI void cyber_gear_build_parameter_write_frame_with_int_value(const cgFrame *frame, motorParameterIndex index, int value);

// /* 构造一个参数写入的CAN包 (通信类型18), 参数值为浮点
//  * @param: frame 要设置的帧
//  * @param: index 参数index
//  * @param: value 参数值
//  * */
// CYBERGEARAPI void cyber_gear_build_parameter_write_frame_with_float_value(const cgFrame *frame, motorParameterIndex index, float value);

// /* 构造一个参数读取的CAN包 （通信类型17）
//  * @param: frame 要设置的帧
//  */
// CYBERGEARAPI void cyber_gear_build_parameter_read_frame(const cgFrame *frame, motorParameterIndex index);

// /* 解析一个参数读取的CAN包 （通信类型17）
//  * @param: frame 要设置的帧
//  */
// CYBERGEARAPI singleParam cyber_gear_parse_parameter_read_frame(const cgFrame *frame);

// /* 获取帧的通信类型
//  * @param: frame 解析的帧
//  * @return: communicationType 通信类型
//  * */
// CYBERGEARAPI communicationType cyber_gear_get_can_id_communication_type(const cgFrame * const frame);

// /* 获取帧的目标CAN_ID
//  * @param: frame 要设置的帧
//  * @return: int 目标 CAN_ID
//  * */
// CYBERGEARAPI int cyber_gear_get_can_id_target_id(const cgFrame * const frame);

// /* 获取帧的主机CAN_ID
//  * @param: frame 要设置的帧
//  * @return: int 主机 CAN_ID
//  * */
// CYBERGEARAPI int cyber_gear_get_can_id_host_id(const cgFrame * const frame);

// /* 运控模式电机控制指令 (通信类型 1)用来向电机发送控制指令
//  * @param: control_param 控制参数
//  * */
// CYBERGEARAPI void cyber_gear_build_motion_control_frame(const cgFrame *frame, const motionControl control_param);

// /* 解析通信类型6 构建一个 位置机械零位帧
//  * @param: frame 要解析的帧
//  * @return: communicationType 通信类型
//  * */
// CYBERGEARAPI void cyber_gear_build_set_mechanical_zero_position_frame(const cgFrame * frame);

// /* 解析通信类型7 构建一个 设置 CAN ID 帧
//  * @param: frame 要解析的帧
//  * @return: communicationType 通信类型
//  * */
// CYBERGEARAPI void cyber_gear_build_set_can_id_frame(const cgFrame * frame, int setting_can_id);

// /*  Dump 一个 CyberGear 的 电机运行状态帧 帧 */
// CYBERGEARAPI void cyber_gear_dump_motor_status_frame(const motorStatus status);

package pdufield

//BindTransmitterTLVParameter defines the TLV fields names
//realted with BindTransmitter PDU
type bindTransmiterTLVParameter struct {
	//ScInterfaceVersion ...
	ScInterfaceVersion TLVType
}

//submitSMTLVParameter defines the TLV fileds names related with SubmitSM PDU
type submitSMTLVParameter struct {
	//UserMessageReference ESME assigned message reference number.
	UserMessageReference TLVType

	//SourcePort Indicates the application port number associated with the source address
	// of the message. This parameter should be present for WAP applications.
	SourcePort TLVType

	//SourceAddrSubmit The subcomponent in the destination device which created the user data.
	SourceAddrSubUnit TLVType

	//DestionationPort Indicates the application port number associated with the destination address
	//of the message. This parameter should be present for WAP applications.
	DestinationPort TLVType

	//DestAddrSubmit The subcomponent in the destination device for which the user data is intended.
	DestAddrSubUnit TLVType

	//SarMsgRefNum The reference number for a particular concatenated short message
	SarMsgRefNum TLVType

	//SarTotalSegments Indicates the total number of short messages within the
	//concatenated short message.
	SarTotalSegments TLVType

	//SarSegmentSeqNum Indicates the sequence number of a particular short message
	//fragment within the concatenated short message.
	SarSegmentSeqNum TLVType

	//MoreMessagesToSend MoreMessagesToSend Indicates that there are more
	//messages to follow for the destination SME.
	MoreMessagesToSend TLVType

	//PayloadType defines the type of payload (e.g. WDP WCMP etc.).
	PayloadType TLVType

	//MessagePayload Contains the extended short message user data. Up to 64K octets can be
	//transmitted.
	//Note: The short message data should be inserted in either the short_message
	//or message_payload fields. Both fields should not be used simultaneously.
	//The sm_length field should be set to zero if using the message_payload parameter.
	MessagePayload TLVType

	//PrivacyIndicator Indicates the level of privacy associated with the message
	PrivacyIndicator TLVType

	//CallbackNum ),A callback number associated with the short message
	//This parameter can be included a number of times for multiple callback
	//addresses.
	CallbackNum TLVType

	//CallbackNumPresInd Defines the callback number presentation and screening
	//If this parameter is present and there are multiple instances of the
	//callback_num parameter then this parameter must occur an equal number of
	//instances and the order of occurrence determines the particular
	//callback_num_pres_ind which corresponds to a particular callback_num.
	CallbackNumPresInd TLVType

	//CallbackNumAtag Associates a displayable alphanumeric tag with the callback number
	//If this parameter is present and there are multiple instances of the
	// callback_num parameter then this parameter must occur an equal number
	//of instances and the order of occurrence determines the particular
	//callback_num_atag which corresponds to a particular callback_num.
	CallbackNumAtag TLVType

	//SourceSubAddress The subaddress of the message originator.
	SourceSubAddress TLVType

	//DestSubAddress The subaddress of the message destination.
	DestSubAddress TLVType

	//UserResponseCode A user response code. The actual response codes are implementation specific.
	UserResponseCode TLVType

	//DisplayTime Provides the receiving MS with a display time associated with the message.
	DisplayTime TLVType

	//SMSSingal Indicates the alerting mechanism when the message is received by an MS.
	SMSSignal TLVType

	//MSValidity Indicates validity information for this message to the recipient MS.
	MSValidity TLVType

	//MSMsgWaitFacilities This parameter controls the indication and specifies the message
	// type (of the message associated with the MWI) at the mobile station.
	MSMsgWaitFacilities TLVType

	//NumberOfMessages  Indicates the number of messages stored in a mail box
	NumberOfMessages TLVType

	//AlertOnMsgDelivery  Request an MS alert signal be invoked on message delivery
	AlertOnMsgDelivery TLVType

	//LanguageIndicator  Indicates the language of an alphanumeric text message
	LanguageIndicator TLVType

	//ItsReplyType  The MS user’s reply method to an SMS delivery message received from
	//the network is indicated and controlled by this parameter.
	ItsReplyType TLVType

	//ItsSessionInfo  Session control information for Interactive Teleservice
	ItsSessionInfo TLVType

	//UsedServiceOp  This parameter is used to identify the required USSD Service
	//type when interfacing to a USSD system.
	UsedServiceOp TLVType
}

//SubmitSMMultiTVLParameter defines the TLV fields names for
// submit sm multipart PDU.
type submitSMMultiTLVParameter struct {
	UserMessageReference TLVType
	SourcePort           TLVType
	SourceAddrSubUnit    TLVType
	DestinationPort      TLVType
	DestAddrSubUnit      TLVType
	SarMsgRefNum         TLVType
	SarTotalSegments     TLVType
	SarSegmentSeqNum     TLVType
	PayloadType          TLVType
	MessagePayload       TLVType
	PrivacyIndicator     TLVType
	CallbackNum          TLVType
	CallbackNumPresInd   TLVType
	CallbackNumAtag      TLVType
	SourceSubAddress     TLVType
	DestSubAddress       TLVType
	DisplayTime          TLVType
	SMSSignal            TLVType
	MSValidity           TLVType
	MSMsgWaitFacilities  TLVType
	AlertOnMsgDelivery   TLVType
	LanguageIndicator    TLVType
	DestFlag             TLVType
	SMEAddress           TLVType
	DistributionListName TLVType
}

type deliverSMTLVParameter struct {
	UserMessageReference TLVType
	SourcePort           TLVType
	DestinationPort      TLVType
	SarMsgRefNum         TLVType
	SarTotalSegments     TLVType
	SarSegmentSeqNum     TLVType
	UserResponseCode     TLVType
	PrivacyIndicator     TLVType
	PayloadType          TLVType
	MessagePayload       TLVType
	CallbackNum          TLVType
	SourceSubAddress     TLVType
	DestSubAddress       TLVType
	LanguageIndicator    TLVType
	ItsSessionInfo       TLVType
	NetworErrorCode      TLVType
	MessageState         TLVType
	ReceiptedMessageID   TLVType
}

type dataSMTLVParameter struct {
	SourcePort           TLVType
	SourceAddrUnit       TLVType
	SourceNetworkType    TLVType
	SourceBearerType     TLVType
	SourceTelematicID    TLVType
	DestinationPort      TLVType
	DestAddrSubUnit      TLVType
	DestNetworkType      TLVType
	DestBearerType       TLVType
	DestTelematicsID     TLVType
	SarMsgRefNum         TLVType
	SarTotalSegments     TLVType
	SarSegmentSeqNum     TLVType
	MoreMessagesToSend   TLVType
	QosTimeToLive        TLVType
	PayloadType          TLVType
	MessagePayload       TLVType
	SetDPF               TLVType
	ReceiptedMessageID   TLVType
	MessageState         TLVType
	NetworkErrorCode     TLVType
	UserMessageReference TLVType
	PrivacyInicator      TLVType
	CallbackNum          TLVType
	CallbackNumPresInd   TLVType
	CallbackNumAtag      TLVType
	SourceSubAddress     TLVType
	DestSubAddress       TLVType
	UserResponseCode     TLVType
	DisplayTime          TLVType
	SMSSignal            TLVType
	MSValidity           TLVType
	MsMsgWaitFacilities  TLVType
	NumberOfMessages     TLVType
	AlertOnMsgDelivery   TLVType
	LanguageIndicator    TLVType
	ItsReplyType         TLVType
	ItsSessionInfo       TLVType
}

type dataSMRespTLVParameter struct {
	DeliveryFaiureReason     TLVType
	NetWorkErrorCode         TLVType
	AdditionalStatusInfoText TLVType
	DPFResult                TLVType
}

// BindTransmiterTLVParameter defines the TLV parameters available  for BindTrasmiter PDU
var BindTransmiterTLVParameter *bindTransmiterTLVParameter

// SubmitSMTLVParameter defines the TLV parameters available for SubmitSM PDU
var SubmitSMTLVParameter *submitSMTLVParameter

// SubmitSMMultiTLVParameter defines the TLV parameters available for submitSMMulti PDU
var SubmitSMMultiTLVParameter *submitSMMultiTLVParameter

// DeliverSMTLVParameter defines TLV parameters available for DeliverSM PDU
var DeliverSMTLVParameter *deliverSMTLVParameter

//DataSMTLVParameter defines the TLV parameters available for DataSM PDU
var DataSMTLVParameter *dataSMTLVParameter

// DataSMRespTLVParameter defines the TLV parameters available for DataSMResp PDU
var DataSMRespTLVParameter *dataSMRespTLVParameter

func init() {
	BindTransmiterTLVParameter = &bindTransmiterTLVParameter{
		ScInterfaceVersion: TLVType(1),
	}
	SubmitSMTLVParameter = &submitSMTLVParameter{
		UserMessageReference: TLVType(1),
		SourcePort:           TLVType(2),
		SourceAddrSubUnit:    TLVType(3),
		DestinationPort:      TLVType(4),
		DestAddrSubUnit:      TLVType(5),
		SarMsgRefNum:         TLVType(6),
		SarTotalSegments:     TLVType(7),
		SarSegmentSeqNum:     TLVType(8),
		MoreMessagesToSend:   TLVType(9),
		PayloadType:          TLVType(10),
		MessagePayload:       TLVType(11),
		PrivacyIndicator:     TLVType(12),
		CallbackNum:          TLVType(13),
		CallbackNumPresInd:   TLVType(14),
		CallbackNumAtag:      TLVType(15),
		SourceSubAddress:     TLVType(16),
		DestSubAddress:       TLVType(17),
		UserResponseCode:     TLVType(18),
		DisplayTime:          TLVType(19),
		SMSSignal:            TLVType(20),
		MSValidity:           TLVType(21),
		MSMsgWaitFacilities:  TLVType(22),
		NumberOfMessages:     TLVType(23),
		AlertOnMsgDelivery:   TLVType(24),
		LanguageIndicator:    TLVType(25),
		ItsReplyType:         TLVType(26),
		ItsSessionInfo:       TLVType(27),
		UsedServiceOp:        TLVType(28),
	}
	SubmitSMMultiTLVParameter = &submitSMMultiTLVParameter{
		UserMessageReference: TLVType(1),
		SourcePort:           TLVType(2),
		SourceAddrSubUnit:    TLVType(3),
		DestinationPort:      TLVType(4),
		DestAddrSubUnit:      TLVType(5),
		SarMsgRefNum:         TLVType(6),
		SarTotalSegments:     TLVType(7),
		SarSegmentSeqNum:     TLVType(8),
		PayloadType:          TLVType(10),
		MessagePayload:       TLVType(11),
		PrivacyIndicator:     TLVType(12),
		CallbackNum:          TLVType(13),
		CallbackNumPresInd:   TLVType(14),
		CallbackNumAtag:      TLVType(15),
		SourceSubAddress:     TLVType(16),
		DestSubAddress:       TLVType(17),
		DisplayTime:          TLVType(19),
		SMSSignal:            TLVType(20),
		MSValidity:           TLVType(21),
		MSMsgWaitFacilities:  TLVType(22),
		AlertOnMsgDelivery:   TLVType(24),
		LanguageIndicator:    TLVType(25),
		DestFlag:             TLVType(26),
		SMEAddress:           TLVType(27),
		DistributionListName: TLVType(28),
	}

	DeliverSMTLVParameter = &deliverSMTLVParameter{
		UserMessageReference: TLVType(1),
		SourcePort:           TLVType(2),
		DestinationPort:      TLVType(4),
		SarMsgRefNum:         TLVType(6),
		SarTotalSegments:     TLVType(7),
		SarSegmentSeqNum:     TLVType(8),
		PayloadType:          TLVType(9),
		MessagePayload:       TLVType(10),
		PrivacyIndicator:     TLVType(11),
		CallbackNum:          TLVType(12),
		SourceSubAddress:     TLVType(13),
		DestSubAddress:       TLVType(14),
		LanguageIndicator:    TLVType(15),
		ItsSessionInfo:       TLVType(16),
		NetworErrorCode:      TLVType(17),
		MessageState:         TLVType(18),
		ReceiptedMessageID:   TLVType(19),
	}
	DataSMTLVParameter = &dataSMTLVParameter{
		SourcePort:           TLVType(1),
		SourceAddrUnit:       TLVType(2),
		SourceNetworkType:    TLVType(3),
		SourceBearerType:     TLVType(4),
		SourceTelematicID:    TLVType(5),
		DestinationPort:      TLVType(6),
		DestAddrSubUnit:      TLVType(7),
		DestNetworkType:      TLVType(8),
		DestBearerType:       TLVType(9),
		DestTelematicsID:     TLVType(10),
		SarMsgRefNum:         TLVType(11),
		SarTotalSegments:     TLVType(12),
		SarSegmentSeqNum:     TLVType(13),
		MoreMessagesToSend:   TLVType(14),
		QosTimeToLive:        TLVType(15),
		PayloadType:          TLVType(16),
		MessagePayload:       TLVType(17),
		SetDPF:               TLVType(18),
		ReceiptedMessageID:   TLVType(19),
		MessageState:         TLVType(20),
		NetworkErrorCode:     TLVType(21),
		UserMessageReference: TLVType(22),
		PrivacyInicator:      TLVType(23),
		CallbackNum:          TLVType(24),
		CallbackNumPresInd:   TLVType(25),
		CallbackNumAtag:      TLVType(26),
		SourceSubAddress:     TLVType(27),
		DestSubAddress:       TLVType(28),
		UserResponseCode:     TLVType(29),
		DisplayTime:          TLVType(30),
		SMSSignal:            TLVType(31),
		MSValidity:           TLVType(32),
		MsMsgWaitFacilities:  TLVType(33),
		NumberOfMessages:     TLVType(34),
		AlertOnMsgDelivery:   TLVType(35),
		LanguageIndicator:    TLVType(36),
		ItsReplyType:         TLVType(37),
		ItsSessionInfo:       TLVType(38),
	}

	DataSMRespTLVParameter = &dataSMRespTLVParameter{
		DeliveryFaiureReason:     TLVType(1),
		NetWorkErrorCode:         TLVType(2),
		AdditionalStatusInfoText: TLVType(3),
		DPFResult:                TLVType(4),
	}
}

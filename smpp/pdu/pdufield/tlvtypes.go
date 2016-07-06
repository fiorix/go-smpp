package pdufield

// BindTransmiterTLVParameter defines the TLV parameters available  for BindTrasmiter PDU
const (
	AdditionalStatusInfoText = "AdditionalStatusInfoText"
	AlertOnMsgDelivery       = "AlertOnMsgDelivery"
	CallbackNum              = "CallbackNum"
	CallbackNumPresInd       = "CallbackNumPresInd"
	CallbackNumAtag          = "CallbackNumAtag"
	DeliveryFailureReason    = "DeliveryFailureReason"
	DestAddrSubUnit          = "DestAddrSubUnit"
	DestBearerType           = "DestBearerType"
	DestFlag                 = "DestFlag"
	DestNetworkType          = "DestNetworkType"
	DestSubAddress           = "DestSubAddress"
	DestinationPort          = "DestinationPort"
	DestTelematicsID         = "DestTelematicsID"
	DisplayTime              = "DisplayTime"
	DistributionListName     = "DistributionListName"
	DPFResult                = "DPFResult"
	ItsReplyType             = "ItsReplyType"
	ItsSessionInfo           = "ItsSessionInfo"
	LanguageIndicator        = "LanguageIndicator"
	MessagePayload           = "MessagePayload"
	MessageStateTLV          = "MessageState"
	MoreMessagesToSend       = "MoreMessagesToSend"
	MSMsgWaitFacilities      = "MSMsgWaitFacilities"
	MSValidity               = "MSValidity"
	NumberOfMessages         = "NumberOfMessages"
	NetworkErrorCode         = "NetworkErrorCode"
	PayloadType              = "PayloadType"
	PrivacyIndicator         = "PrivacyIndicator"
	QosTimeToLive            = "QosTimeToLive"
	ReceiptedMessageID       = "ReceiptedMessageID"
	SarMsgRefNum             = "SarMsgRefNum"
	SarTotalSegments         = "SarTotalSegments"
	SarSegmentSeqNum         = "SarSegmentSeqNum"
	ScInterfaceVersion       = "ScInterfaceVersion"
	SetDPF                   = "SetDPF"
	SMEAddress               = "SMEAddress"
	SMSSignal                = "SMSSignal"
	SourceAddrSubUnit        = "SourceAddrSubUnit"
	SourceBearerType         = "SourceBearerType"
	SourceNetworkType        = "SourceNetworkType"
	SourcePort               = "SourcePort"
	SourceSubAddress         = "SourceSubAddress"
	SourceTelematicID        = "SourceTelematicID "
	UsedServiceOp            = "UsedServiceOp"
	UserMessageReference     = "UserMessageReference"
	UserResponseCode         = "UserResponseCode"
)

var BindTransmiterTLVParameter = map[string]TLVType{
	ScInterfaceVersion: TLVType(1),
}

// SubmitSMTLVParameter defines the TLV parameters available for SubmitSM PDU
var SubmitSMTLVParameter = map[string]TLVType{
	// UserMessageReference ESME assigned message reference number.
	UserMessageReference: TLVType(1),

	// SourcePort Indicates the application port number associated with the source address
	// of the message. This parameter should be present for WAP applications.
	SourcePort: TLVType(2),

	// SourceAddrSubmit The subcomponent in the destination device which created the user data.
	SourceAddrSubUnit: TLVType(3),

	// DestinationPort Indicates the application port number associated with the destination address
	// of the message. This parameter should be present for WAP applications.
	DestinationPort: TLVType(4),

	// DestAddrSubmit The subcomponent in the destination device for which the user data is intended.
	DestAddrSubUnit: TLVType(5),

	// SarMsgRefNum The reference number for a particular concatenated short message
	SarMsgRefNum: TLVType(6),

	// SarTotalSegments Indicates the total number of short messages within the
	// concatenated short message.
	SarTotalSegments: TLVType(7),

	// SarSegmentSeqNum Indicates the sequence number of a particular short message
	// fragment within the concatenated short message.
	SarSegmentSeqNum: TLVType(8),

	// MoreMessagesToSend MoreMessagesToSend Indicates that there are more
	// messages to follow for the destination SME.
	MoreMessagesToSend: TLVType(9),

	// PayloadType defines the type of payload (e.g. WDP WCMP etc.),.
	PayloadType: TLVType(10),

	// MessagePayload Contains the extended short message user data. Up to 64K octets can be
	// transmitted.
	// Note: The short message data should be inserted in either the short_message
	// or message_payload fields. Both fields should not be used simultaneously.
	// The sm_length field should be set to zero if using the message_payload parameter.
	MessagePayload: TLVType(11),

	// PrivacyIndicator Indicates the level of privacy associated with the message
	PrivacyIndicator: TLVType(12),

	// CallbackNum ,A callback number associated with the short message
	// This parameter can be included a number of times for multiple callback
	// addresses.
	CallbackNum: TLVType(13),

	// CallbackNumPresInd Defines the callback number presentation and screening
	// If this parameter is present and there are multiple instances of the
	// callback_num parameter then this parameter must occur an equal number of
	// instances and the order of occurrence determines the particular
	// callback_num_pres_ind which corresponds to a particular callback_num.
	CallbackNumPresInd: TLVType(14),
	// CallbackNumAtag Associates a displayable alphanumeric tag with the callback number
	// If this parameter is present and there are multiple instances of the
	// callback_num parameter then this parameter must occur an equal number
	// of instances and the order of occurrence determines the particular
	// callback_num_atag which corresponds to a particular callback_num.
	CallbackNumAtag: TLVType(15),

	// SourceSubAddress The subaddress of the message originator.
	SourceSubAddress: TLVType(16),

	// DestSubAddress The subaddress of the message destination.
	DestSubAddress: TLVType(17),

	// UserResponseCode A user response code. The actual response codes are implementation specific.
	UserResponseCode: TLVType(18),

	// DisplayTime Provides the receiving MS with a display time associated with the message.
	DisplayTime: TLVType(19),

	// SMSSingal Indicates the alerting mechanism when the message is received by an MS.
	SMSSignal: TLVType(20),

	// MSValidity Indicates validity information for this message to the recipient MS.
	MSValidity: TLVType(21),

	// MSMsgWaitFacilities This parameter controls the indication and specifies the message
	// type (of the message associated with the MWI), at the mobile station.
	MSMsgWaitFacilities: TLVType(22),

	// NumberOfMessages  Indicates the number of messages stored in a mail box
	NumberOfMessages: TLVType(23),

	// AlertOnMsgDelivery  Request an MS alert signal be invoked on message delivery
	AlertOnMsgDelivery: TLVType(24),

	// LanguageIndicator  Indicates the language of an alphanumeric text message
	LanguageIndicator: TLVType(25),

	// ItsReplyType  The MS user’s reply method to an SMS delivery message received from
	// the network is indicated and controlled by this parameter.
	ItsReplyType: TLVType(26),

	// ItsSessionInfo  Session control information for Interactive Teleservice
	ItsSessionInfo: TLVType(27),

	// UsedServiceOp  This parameter is used to identify the required USSD Service
	// type when interfacing to a USSD system.
	UsedServiceOp: TLVType(28),
}

// SubmitSMMultiTLVParameter defines the TLV parameters available for submitSMMulti PDU
var SubmitSMMultiTLVParameter = map[string]TLVType{
	// UserMessageReference ,ESME assigned message reference number
	UserMessageReference: TLVType(1),

	// SourcePort  Indicates the application port number associated with the source
	// address of the message. This parameter should be present for WAP applications.
	SourcePort: TLVType(2),

	// SourceAddrSubUnit The subcomponent in the destination device which created the user data
	SourceAddrSubUnit: TLVType(3),

	// DestinationPort Indicates the application port number associated with the destination address of the message
	// This parameter should be present for WAP applications
	DestinationPort: TLVType(4),

	// DestAddrSubUnit  The subcomponent in the destination device for which the user data is intended
	DestAddrSubUnit: TLVType(5),

	// SarMsgRefNum   The reference number for a particular concatenated short message
	SarMsgRefNum: TLVType(6),

	// SarTotalSegments  Indicates the total number of short messages within the concatenated short message
	SarTotalSegments: TLVType(7),

	// SarSegmentSeqNum Indicates the sequence number of a particular short message fragment within
	// the concatenated short message.
	SarSegmentSeqNum: TLVType(8),

	// PayloadType Defines the type of payload (e.g. WDP, WCMP, etc.),
	PayloadType: TLVType(10),

	// MessagePayload  Contains the extended short message user data. Up to 64K octets can be transmitted
	// Note: The short message data should be inserted in either the short_message or message_payload fields
	// Both fields should not be used simultaneously The sm_length field should be set to zero
	// if using the message_payload parameter.
	MessagePayload: TLVType(11),

	//PrivacyIndicator Indicates the level of privacy associated with the message
	PrivacyIndicator: TLVType(12),

	// CallbackNum  A callback number associated with the short message
	// This parameter can be included a number of times for multiple callback addresses.
	CallbackNum: TLVType(13),

	// CallbackNumPresInd Identifies the presentation and screening associated with the callback number
	// If this parameter is present and there are multiple instances of the callback_num parameter
	// then this parameter must occur an equal number of instances and the order of occurrence determines
	// the particular callback_num_pres_ind which corresponds to a particular callback_num.
	CallbackNumPresInd: TLVType(14),

	// CallbackNumAtag Associates a displayable alphanumeric tag with the callback number.
	// If this parameter is present and there are multiple instances of the callback_num parameter
	// then this parameter must occur an equal number of instances and the order of occurrence determines
	// the particular callback_num_atag which corresponds to a particular callback_num
	CallbackNumAtag: TLVType(15),

	// SourceSubAddress  The subaddress of the message originator
	SourceSubAddress: TLVType(16),

	// DestSubAddress The subaddress of the message destination
	DestSubAddress: TLVType(17),

	// DisplayTime Provides the receiving MS based SME with a display time associated with the message
	DisplayTime: TLVType(19),

	// SMSSignal  Indicates the alerting mechanism when the message is received by an MS
	SMSSignal: TLVType(20),

	// MSValidity Indicates validity information for this message to the recipient MS
	MSValidity: TLVType(21),

	// MSMsgWaitFacilities This parameter controls the indication and specifies the message type
	// (of the message associated with the MWI), at the mobile station.
	MSMsgWaitFacilities: TLVType(22),

	// AlertOnMsgDelivery Requests an MS alert signal be invoked on message delivery
	AlertOnMsgDelivery: TLVType(24),

	// LanguageIndicator Indicates the language of an alphanumeric text message.
	LanguageIndicator: TLVType(25),

	// DestFlag  Flag which will identify whether destination address is a Distribution List name or SME address
	DestFlag: TLVType(26),

	// SMEAddress  Depending on dest_flag this could be an SME Address or a Distribution List Name
	SMEAddress: TLVType(27),

	// DistributionListName  Depending on dest_flag this could be an SME Address or a Distribution List Name
	DistributionListName: TLVType(28),
}

// DeliverSMTLVParameter defines TLV parameters available for DeliverSM PDU
var DeliverSMTLVParameter = map[string]TLVType{

	// UserMessageReference A reference assigned by the originating SME to the message.
	// In the case that the deliver_sm is carrying an SMSC delivery receipt, an SME delivery acknowledgement
	// or an SME user acknowledgement (as indicated in the esm_class field),, the user_message_reference parameter
	// is set to the message reference of the original message
	UserMessageReference: TLVType(1),

	// SourcePort Indicates the application port number associated with the source address of the message
	// The parameter should be present for WAP applications.
	SourcePort: TLVType(2),

	// DestinationPort Indicates the application port number associated with the destination address of the message
	// The parameter should be present for WAP applications
	DestinationPort: TLVType(4),

	// SarMsgRefNum  The reference number for a particular concatenated short message
	SarMsgRefNum: TLVType(6),

	// SarTotalSegments Indicates the total number of short messages within the concatenated short message
	SarTotalSegments: TLVType(7),

	// SarSegmentSeqNum  Indicates the sequence number of a particular short message fragment within the
	// concatenated short message
	SarSegmentSeqNum: TLVType(8),

	// PayloadType  Defines the type of payload (e.g. WDP, WCMP, etc.),
	PayloadType: TLVType(9),

	// MessagePayload  Contains the extended short message user data. Up to 64K octets can be transmitted
	// Note: The short message data should be inserted in either the short_message or message_payload fields.
	// Both fields should not be used simultaneously.
	// The sm_length field should be set to zero if using the message_payload parameter.
	MessagePayload: TLVType(10),

	// PrivacyIndicator Indicates a level of privacy associated with the message
	PrivacyIndicator: TLVType(11),

	// CallbackNum A callback number associated with the short message. This parameter can be included a
	// number of times for multiple call back addresses.
	CallbackNum: TLVType(12),

	// SourceSubAddress  The subaddress of the message originator.
	SourceSubAddress: TLVType(13),

	// DestSubAddress The subaddress of the message destination.
	DestSubAddress: TLVType(14),

	// LanguageIndicator  Indicates the language of an alphanumeric text message
	LanguageIndicator: TLVType(15),

	// ItsSessionInfo Session control information for Interactive Teleservice
	ItsSessionInfo: TLVType(16),

	// NetworErrorCode  May be present for Intermediate Notifications and SMSC Delivery Receipts
	NetworkErrorCode: TLVType(17),

	// MessageState Should be present for SMSC Delivery Receipts and Intermediate Notifications
	MessageStateTLV: TLVType(18),

	// ReceiptedMessageID  SMSC message ID of receipted message Should be present for SMSC Delivery Receipts
	// and Intermediate Notifications
	ReceiptedMessageID: TLVType(19),
}

//DataSMTLVParameter defines the TLV parameters available for DataSM PDU
var DataSMTLVParameter = map[string]TLVType{
	// SourcePort Indicates the application port number associated with the source address of the message
	// This parameter should be present for WAP applications
	SourcePort: TLVType(1),

	//SourceAddrSubUnit  The subcomponent in the destination device which created the user data
	SourceAddrSubUnit: TLVType(2),

	// SourceNetworkType The correct network associated with the originating device
	SourceNetworkType: TLVType(3),

	// SourceBearerType The correct bearer type for the delivering the user data to the destination
	SourceBearerType: TLVType(4),

	// SourceTelematicID  The telematics identifier associated with the source
	SourceTelematicID: TLVType(5),

	// DestinationPort  Indicates the application port number associated with the destination address of the message
	// This parameter should be present for WAP applications
	DestinationPort: TLVType(6),

	// DestAddrSubUnit  The subcomponent in the destination device for which the user data is intended
	DestAddrSubUnit: TLVType(7),

	// DestNetworkType The correct network for the destination device
	DestNetworkType: TLVType(8),

	// DestBearerType The correct bearer type for the delivering the user data to the destination
	DestBearerType: TLVType(9),

	// DestTelematicsID  The telematics identifier associated with the destination
	DestTelematicsID: TLVType(10),

	// SarMsgRefNum The reference number for a particular concatenated short message
	SarMsgRefNum: TLVType(11),

	// SarTotalSegments Indicates the total number of short messages within the concatenated short message
	SarTotalSegments: TLVType(12),

	// SarSegmentSeqNum  Indicates the sequence number of a particular short message fragment
	// within the concatenated short message
	SarSegmentSeqNum: TLVType(13),

	// MoreMessagesToSend Indicates that there are more messages to follow for the destination SME
	MoreMessagesToSend: TLVType(14),

	// QosTimeToLive Time to live as a relative time in seconds from submission
	QosTimeToLive: TLVType(15),

	// PayloadType  Defines the type of payload (e.g. WDP, WCMP, etc.),
	PayloadType: TLVType(16),

	// MessagePayload Contains the message user data. Up to 64K octets can be transmitted
	MessagePayload: TLVType(17),

	// SetDPF  Indicator for setting Delivery Pending Flag on delivery failure.
	SetDPF: TLVType(18),

	// ReceiptedMessageID SMSC message ID of message being receipted. Should be present for SMSC Delivery
	// Receipts and Intermediate Notifications
	ReceiptedMessageID: TLVType(19),

	// MessageState Should be present for SMSC Delivery Receipts and Intermediate Notifications
	MessageStateTLV: TLVType(20),

	// NetworkErrorCode May be present for SMSC Delivery Receipts and Intermediate Notifications
	NetworkErrorCode: TLVType(21),

	// UserMessageReference ESME assigned message reference number
	UserMessageReference: TLVType(22),

	// PrivacyInicator  Indicates a level of privacy associated with the message.
	PrivacyIndicator: TLVType(23),

	// CallbackNum  A callback number associated with the short message. This parameter can be included a number
	// of times for multiple call back addresses
	CallbackNum: TLVType(24),

	// CallbackNumPresInd This parameter identifies the presentation and screening associated with the callback number
	// If this parameter is present and there are multiple instances of the callback_num parameter then
	// this parameter must occur an equal number of instances and the order of occurrence determines
	// the particular callback_num_pres_ind which corresponds to a particular callback_num
	CallbackNumPresInd: TLVType(25),

	// CallbackNumAtag  This parameter associates a displayable alphanumeric tag with the callback number.
	// If this parameter is present and there are multiple instances of the callback_num parameter then this
	// parameter must occur an equal number of instances and the order of occurrence determines the particular
	// callback_num_atag which corresponds to a particular callback_num
	CallbackNumAtag: TLVType(26),

	// SourceSubAddress The subaddress of the message originator.
	SourceSubAddress: TLVType(27),

	// DestSubAddress The subaddress of the message destination
	DestSubAddress: TLVType(28),

	// UserResponseCode A user response code. The actual response codes are implementation specific
	UserResponseCode: TLVType(29),

	// DisplayTime  Provides the receiving MS based SME with a display time associated with the message
	DisplayTime: TLVType(30),

	// SMSSignal  Indicates the alerting mechanism when the message is received by an MS
	SMSSignal: TLVType(31),

	// MSValidity  Indicates validity information for this message to the recipient MS
	MSValidity: TLVType(32),

	// MsMsgWaitFacilities  This parameter controls the indication and specifies the message
	// type (of the message associated with the MWI), at the mobile station.
	MSMsgWaitFacilities: TLVType(33),

	// NumberOfMessages  Indicates the number of messages stored in a mail box (e.g. voice mail box),
	NumberOfMessages: TLVType(34),

	// AlertOnMsgDelivery Requests an MS alert signal be invoked on message delivery
	AlertOnMsgDelivery: TLVType(35),

	// LanguageIndicator  Indicates the language of an alphanumeric text message.
	LanguageIndicator: TLVType(36),

	// ItsReplyType  The MS user’s reply method to an SMS delivery message received from the network
	// is indicated and controlled by this parameter
	ItsReplyType: TLVType(37),

	// ItsSessionInf Session control information for Interactive Teleservice
	ItsSessionInfo: TLVType(38),
}

// DataSMRespTLVParameter defines the TLV parameters available for DataSMResp PDU
var DataSMRespTLVParameter = map[string]TLVType{
	// DeliveryFailureReason  Include to indicate reason for delivery failure
	DeliveryFailureReason: TLVType(1),

	// NetWorkErrorCode  Error code specific to a wireless network
	NetworkErrorCode: TLVType(2),

	// AdditionalStatusInfoText ASCII text giving a description of the meaning of the response
	AdditionalStatusInfoText: TLVType(3),

	// DPFResult Indicates whether the Delivery Pending Flag was set
	DPFResult: TLVType(4),
}

package pdufield

// BindTransmiterTLVParameter defines the TLV parameters available  for BindTrasmiter PDU
var BindTransmiterTLVParameter map[string]TLVType

// SubmitSMTLVParameter defines the TLV parameters available for SubmitSM PDU
var SubmitSMTLVParameter map[string]TLVType

// SubmitSMMultiTLVParameter defines the TLV parameters available for submitSMMulti PDU
var SubmitSMMultiTLVParameter map[string]TLVType

// DeliverSMTLVParameter defines TLV parameters available for DeliverSM PDU
var DeliverSMTLVParameter map[string]TLVType

//DataSMTLVParameter defines the TLV parameters available for DataSM PDU
var DataSMTLVParameter map[string]TLVType

// DataSMRespTLVParameter defines the TLV parameters available for DataSMResp PDU
var DataSMRespTLVParameter map[string]TLVType

func init() {

	BindTransmiterTLVParameter = make(map[string]TLVType)
	// ScInterfaceVersion SMPP version supported by SMSC
	BindTransmiterTLVParameter["ScInterfaceVersion"] = TLVType(1)

	SubmitSMTLVParameter = make(map[string]TLVType)

	// UserMessageReference ESME assigned message reference number.
	SubmitSMTLVParameter["UserMessageReference"] = TLVType(1)

	// SourcePort Indicates the application port number associated with the source address
	// of the message. This parameter should be present for WAP applications.
	SubmitSMTLVParameter["SourcePort"] = TLVType(2)

	// SourceAddrSubmit The subcomponent in the destination device which created the user data.
	SubmitSMTLVParameter["SourceAddrSubUnit"] = TLVType(3)

	// DestinationPort Indicates the application port number associated with the destination address
	// of the message. This parameter should be present for WAP applications.
	SubmitSMTLVParameter["DestinationPort"] = TLVType(4)

	// DestAddrSubmit The subcomponent in the destination device for which the user data is intended.
	SubmitSMTLVParameter["DestAddrSubUnit"] = TLVType(5)

	// SarMsgRefNum The reference number for a particular concatenated short message
	SubmitSMTLVParameter["SarMsgRefNum"] = TLVType(6)

	// SarTotalSegments Indicates the total number of short messages within the
	// concatenated short message.
	SubmitSMTLVParameter["SarTotalSegments"] = TLVType(7)

	// SarSegmentSeqNum Indicates the sequence number of a particular short message
	// fragment within the concatenated short message.
	SubmitSMTLVParameter["SarSegmentSeqNumv"] = TLVType(8)

	// MoreMessagesToSend MoreMessagesToSend Indicates that there are more
	// messages to follow for the destination SME.
	SubmitSMTLVParameter["MoreMessagesToSend"] = TLVType(9)

	// PayloadType defines the type of payload (e.g. WDP WCMP etc.).
	SubmitSMTLVParameter["PayloadType"] = TLVType(10)

	// MessagePayload Contains the extended short message user data. Up to 64K octets can be
	// transmitted.
	// Note: The short message data should be inserted in either the short_message
	// or message_payload fields. Both fields should not be used simultaneously.
	// The sm_length field should be set to zero if using the message_payload parameter.
	SubmitSMTLVParameter["MessagePayload"] = TLVType(11)

	// PrivacyIndicator Indicates the level of privacy associated with the message
	SubmitSMTLVParameter["PrivacyIndicator"] = TLVType(12)

	// CallbackNum ,A callback number associated with the short message
	// This parameter can be included a number of times for multiple callback
	// addresses.
	SubmitSMTLVParameter["CallbackNum"] = TLVType(13)

	// CallbackNumPresInd Defines the callback number presentation and screening
	// If this parameter is present and there are multiple instances of the
	// callback_num parameter then this parameter must occur an equal number of
	// instances and the order of occurrence determines the particular
	// callback_num_pres_ind which corresponds to a particular callback_num.
	SubmitSMTLVParameter["CallbackNumPresInd"] = TLVType(14)
	// CallbackNumAtag Associates a displayable alphanumeric tag with the callback number
	// If this parameter is present and there are multiple instances of the
	// callback_num parameter then this parameter must occur an equal number
	// of instances and the order of occurrence determines the particular
	// callback_num_atag which corresponds to a particular callback_num.
	SubmitSMTLVParameter["CallbackNumAtag"] = TLVType(15)

	// SourceSubAddress The subaddress of the message originator.
	SubmitSMTLVParameter["SourceSubAddress"] = TLVType(16)

	// DestSubAddress The subaddress of the message destination.
	SubmitSMTLVParameter["DestSubAddress"] = TLVType(17)

	// UserResponseCode A user response code. The actual response codes are implementation specific.
	SubmitSMTLVParameter["UserResponseCode"] = TLVType(18)

	// DisplayTime Provides the receiving MS with a display time associated with the message.
	SubmitSMTLVParameter["DisplayTime"] = TLVType(19)

	// SMSSingal Indicates the alerting mechanism when the message is received by an MS.
	SubmitSMTLVParameter["SMSSignal"] = TLVType(20)

	// MSValidity Indicates validity information for this message to the recipient MS.
	SubmitSMTLVParameter["MSValidity"] = TLVType(21)

	// MSMsgWaitFacilities This parameter controls the indication and specifies the message
	// type (of the message associated with the MWI) at the mobile station.
	SubmitSMTLVParameter["MSMsgWaitFacilities"] = TLVType(22)

	// NumberOfMessages  Indicates the number of messages stored in a mail box
	SubmitSMTLVParameter["NumberOfMessages"] = TLVType(23)

	// AlertOnMsgDelivery  Request an MS alert signal be invoked on message delivery
	SubmitSMTLVParameter["AlertOnMsgDelivery"] = TLVType(24)

	// LanguageIndicator  Indicates the language of an alphanumeric text message
	SubmitSMTLVParameter["LanguageIndicato"] = TLVType(25)

	// ItsReplyType  The MS user’s reply method to an SMS delivery message received from
	// the network is indicated and controlled by this parameter.
	SubmitSMTLVParameter["ItsReplyType"] = TLVType(26)

	// ItsSessionInfo  Session control information for Interactive Teleservice
	SubmitSMTLVParameter["ItsSessionInfo"] = TLVType(27)

	// UsedServiceOp  This parameter is used to identify the required USSD Service
	// type when interfacing to a USSD system.
	SubmitSMTLVParameter["UsedServiceOp"] = TLVType(28)

	SubmitSMMultiTLVParameter = make(map[string]TLVType)

	// UserMessageReference ,ESME assigned message reference number
	SubmitSMMultiTLVParameter["UserMessageReference"] = TLVType(1)

	// SourcePort  Indicates the application port number associated with the source
	// address of the message. This parameter should be present for WAP applications.
	SubmitSMMultiTLVParameter["SourcePort"] = TLVType(2)

	// SourceAddrSubUnit The subcomponent in the destination device which created the user data
	SubmitSMMultiTLVParameter["SourceAddrSubUnit"] = TLVType(3)

	// DestinationPort Indicates the application port number associated with the destination address of the message
	// This parameter should be present for WAP applications
	SubmitSMMultiTLVParameter["DestinationPort"] = TLVType(4)

	// DestAddrSubUnit  The subcomponent in the destination device for which the user data is intended
	SubmitSMMultiTLVParameter["DestAddrSubUnit"] = TLVType(5)

	// SarMsgRefNum   The reference number for a particular concatenated short message
	SubmitSMMultiTLVParameter["SarMsgRefNum"] = TLVType(6)

	// SarTotalSegments  Indicates the total number of short messages within the concatenated short message
	SubmitSMMultiTLVParameter["SarTotalSegments"] = TLVType(7)

	// SarSegmentSeqNum Indicates the sequence number of a particular short message fragment within
	// the concatenated short message.
	SubmitSMMultiTLVParameter["SarSegmentSeqNum"] = TLVType(8)

	// PayloadType Defines the type of payload (e.g. WDP, WCMP, etc.)
	SubmitSMMultiTLVParameter["PayloadType"] = TLVType(10)

	// MessagePayload  Contains the extended short message user data. Up to 64K octets can be transmitted
	// Note: The short message data should be inserted in either the short_message or message_payload fields
	// Both fields should not be used simultaneously The sm_length field should be set to zero
	// if using the message_payload parameter.
	SubmitSMMultiTLVParameter["MessagePayload"] = TLVType(11)

	//PrivacyIndicator Indicates the level of privacy associated with the message
	SubmitSMMultiTLVParameter["PrivacyIndicator"] = TLVType(12)

	// CallbackNum  A callback number associated with the short message
	// This parameter can be included a number of times for multiple callback addresses.
	SubmitSMMultiTLVParameter["CallbackNum"] = TLVType(13)

	// CallbackNumPresInd Identifies the presentation and screening associated with the callback number
	// If this parameter is present and there are multiple instances of the callback_num parameter
	// then this parameter must occur an equal number of instances and the order of occurrence determines
	// the particular callback_num_pres_ind which corresponds to a particular callback_num.
	SubmitSMMultiTLVParameter["CallbackNumPresInd"] = TLVType(14)

	// CallbackNumAtag Associates a displayable alphanumeric tag with the callback number.
	// If this parameter is present and there are multiple instances of the callback_num parameter
	// then this parameter must occur an equal number of instances and the order of occurrence determines
	// the particular callback_num_atag which corresponds to a particular callback_num
	SubmitSMMultiTLVParameter["CallbackNumAtag"] = TLVType(15)

	// SourceSubAddress  The subaddress of the message originator
	SubmitSMMultiTLVParameter["SourceSubAddress"] = TLVType(16)

	// DestSubAddress The subaddress of the message destination
	SubmitSMMultiTLVParameter["DestSubAddress"] = TLVType(17)

	// DisplayTime Provides the receiving MS based SME with a display time associated with the message
	SubmitSMMultiTLVParameter["DisplayTime"] = TLVType(19)

	// SMSSignal  Indicates the alerting mechanism when the message is received by an MS
	SubmitSMMultiTLVParameter["SMSSignal"] = TLVType(20)

	// MSValidity Indicates validity information for this message to the recipient MS
	SubmitSMMultiTLVParameter["MSValidity"] = TLVType(21)

	// MSMsgWaitFacilities This parameter controls the indication and specifies the message type
	// (of the message associated with the MWI) at the mobile station.
	SubmitSMMultiTLVParameter["MSMsgWaitFacilities"] = TLVType(22)

	// AlertOnMsgDelivery Requests an MS alert signal be invoked on message delivery
	SubmitSMMultiTLVParameter["AlertOnMsgDelivery"] = TLVType(24)

	// LanguageIndicator Indicates the language of an alphanumeric text message.
	SubmitSMMultiTLVParameter["LanguageIndicator"] = TLVType(25)

	// DestFlag  Flag which will identify whether destination address is a Distribution List name or SME address
	SubmitSMMultiTLVParameter["DestFlag"] = TLVType(26)

	// SMEAddress  Depending on dest_flag this could be an SME Address or a Distribution List Name
	SubmitSMMultiTLVParameter["SMEAddress"] = TLVType(27)

	// DistributionListName  Depending on dest_flag this could be an SME Address or a Distribution List Name
	SubmitSMMultiTLVParameter["DistributionListName"] = TLVType(28)

	DeliverSMTLVParameter = make(map[string]TLVType)

	// UserMessageReference A reference assigned by the originating SME to the message.
	// In the case that the deliver_sm is carrying an SMSC delivery receipt, an SME delivery acknowledgement
	// or an SME user acknowledgement (as indicated in the esm_class field), the user_message_reference parameter
	// is set to the message reference of the original message
	DeliverSMTLVParameter["UserMessageReference"] = TLVType(1)

	// SourcePort Indicates the application port number associated with the source address of the message
	// The parameter should be present for WAP applications.
	DeliverSMTLVParameter["SourcePort"] = TLVType(2)

	// DestinationPort Indicates the application port number associated with the destination address of the message
	// The parameter should be present for WAP applications
	DeliverSMTLVParameter["DestinationPort"] = TLVType(4)

	// SarMsgRefNum  The reference number for a particular concatenated short message
	DeliverSMTLVParameter["SarMsgRefNum"] = TLVType(6)

	// SarTotalSegments Indicates the total number of short messages within the concatenated short message
	DeliverSMTLVParameter["SarTotalSegments"] = TLVType(7)

	// SarSegmentSeqNum  Indicates the sequence number of a particular short message fragment within the
	// concatenated short message
	DeliverSMTLVParameter["SarSegmentSeqNum"] = TLVType(8)

	// PayloadType  Defines the type of payload (e.g. WDP, WCMP, etc.)
	DeliverSMTLVParameter["PayloadType"] = TLVType(9)

	// MessagePayload  Contains the extended short message user data. Up to 64K octets can be transmitted
	// Note: The short message data should be inserted in either the short_message or message_payload fields.
	// Both fields should not be used simultaneously.
	// The sm_length field should be set to zero if using the message_payload parameter.
	DeliverSMTLVParameter["MessagePayload"] = TLVType(10)

	// PrivacyIndicator Indicates a level of privacy associated with the message
	DeliverSMTLVParameter["PrivacyIndicator"] = TLVType(11)

	// CallbackNum A callback number associated with the short message. This parameter can be included a
	// number of times for multiple call back addresses.
	DeliverSMTLVParameter["CallbackNum"] = TLVType(12)

	// SourceSubAddress  The subaddress of the message originator.
	DeliverSMTLVParameter["SourceSubAddress"] = TLVType(13)

	// DestSubAddress The subaddress of the message destination.
	DeliverSMTLVParameter["DestSubAddress"] = TLVType(14)

	// LanguageIndicator  Indicates the language of an alphanumeric text message
	DeliverSMTLVParameter["LanguageIndicator"] = TLVType(15)

	// ItsSessionInfo Session control information for Interactive Teleservice
	DeliverSMTLVParameter["ItsSessionInfo"] = TLVType(16)

	// NetworErrorCode  May be present for Intermediate Notifications and SMSC Delivery Receipts
	DeliverSMTLVParameter["NetworErrorCode"] = TLVType(17)

	// MessageState Should be present for SMSC Delivery Receipts and Intermediate Notifications
	DeliverSMTLVParameter["MessageState"] = TLVType(18)

	// ReceiptedMessageID  SMSC message ID of receipted message Should be present for SMSC Delivery Receipts
	// and Intermediate Notifications
	DeliverSMTLVParameter["ReceiptedMessageID"] = TLVType(19)

	DataSMTLVParameter = make(map[string]TLVType)
	// SourcePort Indicates the application port number associated with the source address of the message
	// This parameter should be present for WAP applications
	DataSMTLVParameter["SourcePort"] = TLVType(1)

	//SourceAddrSubUnit  The subcomponent in the destination device which created the user data
	DataSMTLVParameter["SourceAddrSubUnit"] = TLVType(2)

	// SourceNetworkType The correct network associated with the originating device
	DataSMTLVParameter["SourceNetworkType"] = TLVType(3)

	// SourceBearerType The correct bearer type for the delivering the user data to the destination
	DataSMTLVParameter["SourceBearerType"] = TLVType(4)

	// SourceTelematicID  The telematics identifier associated with the source
	DataSMTLVParameter["SourceTelematicID"] = TLVType(5)

	// DestinationPort  Indicates the application port number associated with the destination address of the message
	// This parameter should be present for WAP applications
	DataSMTLVParameter["DestinationPort"] = TLVType(6)

	// DestAddrSubUnit  The subcomponent in the destination device for which the user data is intended
	DataSMTLVParameter["DestAddrSubUnit:"] = TLVType(7)

	// DestNetworkType The correct network for the destination device
	DataSMTLVParameter["DestNetworkType"] = TLVType(8)

	// DestBearerType The correct bearer type for the delivering the user data to the destination
	DataSMTLVParameter["DestBearerType"] = TLVType(9)

	// DestTelematicsID  The telematics identifier associated with the destination
	DataSMTLVParameter["DestTelematicsID"] = TLVType(10)

	// SarMsgRefNum The reference number for a particular concatenated short message
	DataSMTLVParameter["SarMsgRefNum"] = TLVType(11)

	// SarTotalSegments Indicates the total number of short messages within the concatenated short message
	DataSMTLVParameter["SarTotalSegments"] = TLVType(12)

	// SarSegmentSeqNum  Indicates the sequence number of a particular short message fragment
	// within the concatenated short message
	DataSMTLVParameter["SarSegmentSeqNum"] = TLVType(13)

	// MoreMessagesToSend Indicates that there are more messages to follow for the destination SME
	DataSMTLVParameter["MoreMessagesToSend"] = TLVType(14)

	// QosTimeToLive Time to live as a relative time in seconds from submission
	DataSMTLVParameter["QosTimeToLive"] = TLVType(15)

	// PayloadType  Defines the type of payload (e.g. WDP, WCMP, etc.)
	DataSMTLVParameter["PayloadType"] = TLVType(16)

	// MessagePayload Contains the message user data. Up to 64K octets can be transmitted
	DataSMTLVParameter["MessagePayload"] = TLVType(17)

	// SetDPF  Indicator for setting Delivery Pending Flag on delivery failure.
	DataSMTLVParameter["SetDPF"] = TLVType(18)

	// ReceiptedMessageID SMSC message ID of message being receipted. Should be present for SMSC Delivery
	// Receipts and Intermediate Notifications
	DataSMTLVParameter["ReceiptedMessageID"] = TLVType(19)

	// MessageState Should be present for SMSC Delivery Receipts and Intermediate Notifications
	DataSMTLVParameter["MessageState"] = TLVType(20)

	// NetworkErrorCode May be present for SMSC Delivery Receipts and Intermediate Notifications
	DataSMTLVParameter["NetworkErrorCode"] = TLVType(21)

	// UserMessageReference ESME assigned message reference number
	DataSMTLVParameter["UserMessageReference"] = TLVType(22)

	// PrivacyInicator  Indicates a level of privacy associated with the message.
	DataSMTLVParameter["PrivacyInicator"] = TLVType(23)

	// CallbackNum  A callback number associated with the short message. This parameter can be included a number
	// of times for multiple call back addresses
	DataSMTLVParameter["CallbackNum"] = TLVType(24)

	// CallbackNumPresInd This parameter identifies the presentation and screening associated with the callback number
	// If this parameter is present and there are multiple instances of the callback_num parameter then
	// this parameter must occur an equal number of instances and the order of occurrence determines
	// the particular callback_num_pres_ind which corresponds to a particular callback_num
	DataSMTLVParameter["CallbackNumPresInd"] = TLVType(25)

	// CallbackNumAtag  This parameter associates a displayable alphanumeric tag with the callback number.
	// If this parameter is present and there are multiple instances of the callback_num parameter then this
	// parameter must occur an equal number of instances and the order of occurrence determines the particular
	// callback_num_atag which corresponds to a particular callback_num
	DataSMTLVParameter["CallbackNumAtag"] = TLVType(26)

	// SourceSubAddress The subaddress of the message originator.
	DataSMTLVParameter["SourceSubAddress"] = TLVType(27)

	// DestSubAddress The subaddress of the message destination
	DataSMTLVParameter["DestSubAddress"] = TLVType(28)

	// UserResponseCode A user response code. The actual response codes are implementation specific
	DataSMTLVParameter["UserResponseCode"] = TLVType(29)

	// DisplayTime  Provides the receiving MS based SME with a display time associated with the message
	DataSMTLVParameter["DisplayTime"] = TLVType(30)

	// SMSSignal  Indicates the alerting mechanism when the message is received by an MS
	DataSMTLVParameter["SMSSignal"] = TLVType(31)

	// MSValidity  Indicates validity information for this message to the recipient MS
	DataSMTLVParameter["MSValidity"] = TLVType(32)

	// MsMsgWaitFacilities  This parameter controls the indication and specifies the message
	// type (of the message associated with the MWI) at the mobile station.
	DataSMTLVParameter["MsMsgWaitFacilities"] = TLVType(33)

	// NumberOfMessages  Indicates the number of messages stored in a mail box (e.g. voice mail box)
	DataSMTLVParameter["NumberOfMessages"] = TLVType(34)

	// AlertOnMsgDelivery Requests an MS alert signal be invoked on message delivery
	DataSMTLVParameter["AlertOnMsgDelivery"] = TLVType(35)

	// LanguageIndicator  Indicates the language of an alphanumeric text message.
	DataSMTLVParameter["LanguageIndicator"] = TLVType(36)

	// ItsReplyType  The MS user’s reply method to an SMS delivery message received from the network
	// is indicated and controlled by this parameter
	DataSMTLVParameter["ItsReplyType"] = TLVType(37)

	// ItsSessionInf Session control information for Interactive Teleservice
	DataSMTLVParameter["ItsSessionInfo"] = TLVType(38)

	DataSMRespTLVParameter = make(map[string]TLVType)

	// DeliveryFailureReason  Include to indicate reason for delivery failure
	DataSMRespTLVParameter["DeliveryFailureReason"] = TLVType(1)

	// NetWorkErrorCode  Error code specific to a wireless network
	DataSMRespTLVParameter["NetWorkErrorCode"] = TLVType(2)

	// AdditionalStatusInfoText ASCII text giving a description of the meaning of the response
	DataSMRespTLVParameter["AdditionalStatusInfoText"] = TLVType(3)

	// DPFResult Indicates whether the Delivery Pending Flag was set
	DataSMRespTLVParameter["DPFResult"] = TLVType(4)

}

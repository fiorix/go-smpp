package smpp34

import "fmt"

type CMDStatus uint32

type CMDId uint32

type Header struct {
	Length   uint32
	Id       CMDId
	Status   CMDStatus
	Sequence uint32
}

func NewPduHeader(l uint32, id CMDId, status CMDStatus, seq uint32) *Header {
	return &Header{l, id, status, seq}
}

func (s CMDId) Error() string {
	switch s {
	case GENERIC_NACK:
		return fmt.Sprint("GENERIC_NACK")
	case BIND_RECEIVER:
		return fmt.Sprint("BIND_RECEIVER")
	case BIND_RECEIVER_RESP:
		return fmt.Sprint("BIND_RECEIVER_RESP")
	case BIND_TRANSMITTER:
		return fmt.Sprint("BIND_TRANSMITTER")
	case BIND_TRANSMITTER_RESP:
		return fmt.Sprint("BIND_TRANSMITTER_RESP")
	case QUERY_SM:
		return fmt.Sprint("QUERY_SM")
	case QUERY_SM_RESP:
		return fmt.Sprint("QUERY_SM_RESP")
	case SUBMIT_SM:
		return fmt.Sprint("SUBMIT_SM")
	case SUBMIT_SM_RESP:
		return fmt.Sprint("SUBMIT_SM_RESP")
	case DELIVER_SM:
		return fmt.Sprint("DELIVER_SM")
	case DELIVER_SM_RESP:
		return fmt.Sprint("DELIVER_SM_RESP")
	case UNBIND:
		return fmt.Sprint("UNBIND")
	case UNBIND_RESP:
		return fmt.Sprint("UNBIND_RESP")
	case REPLACE_SM:
		return fmt.Sprint("REPLACE_SM")
	case REPLACE_SM_RESP:
		return fmt.Sprint("REPLACE_SM_RESP")
	case CANCEL_SM:
		return fmt.Sprint("CANCEL_SM")
	case CANCEL_SM_RESP:
		return fmt.Sprint("CANCEL_SM_RESP")
	case BIND_TRANSCEIVER:
		return fmt.Sprint("BIND_TRANSCEIVER")
	case BIND_TRANSCEIVER_RESP:
		return fmt.Sprint("BIND_TRANSCEIVER_RESP")
	case OUTBIND:
		return fmt.Sprint("OUTBIND")
	case ENQUIRE_LINK:
		return fmt.Sprint("ENQUIRE_LINK")
	case ENQUIRE_LINK_RESP:
		return fmt.Sprint("ENQUIRE_LINK_RESP")
	case SUBMIT_MULTI:
		return fmt.Sprint("SUBMIT_MULTI")
	case SUBMIT_MULTI_RESP:
		return fmt.Sprint("SUBMIT_MULTI_RESP")
	case ALERT_NOTIFICATION:
		return fmt.Sprint("ALERT_NOTIFICATION")
	case DATA_SM:
		return fmt.Sprint("DATA_SM")
	case DATA_SM_RESP:
		return fmt.Sprint("DATA_SM_RESP")
	default:
		return fmt.Sprint("Unknown PDU Type. ID:", uint32(s))
	}
}

func (s CMDStatus) Error() string {
	switch s {
	default:
		return fmt.Sprint("Unknown Status:", uint32(s))
	case ESME_ROK:
		return fmt.Sprint("No Error")
	case ESME_RINVMSGLEN:
		return fmt.Sprint("Message Length is invalid")
	case ESME_RINVCMDLEN:
		return fmt.Sprint("Command Length is invalid")
	case ESME_RINVCMDID:
		return fmt.Sprint("Invalid Command ID")
	case ESME_RINVBNDSTS:
		return fmt.Sprint("Incorrect BIND Status for given command")
	case ESME_RALYBND:
		return fmt.Sprint("ESME Already in Bound State")
	case ESME_RINVPRTFLG:
		return fmt.Sprint("Invalid Priority Flag")
	case ESME_RINVREGDLVFLG:
		return fmt.Sprint("Invalid Registered Delivery Flag")
	case ESME_RSYSERR:
		return fmt.Sprint("System Error")
	case ESME_RINVSRCADR:
		return fmt.Sprint("Invalid Source Address")
	case ESME_RINVDSTADR:
		return fmt.Sprint("Invalid Dest Addr")
	case ESME_RINVMSGID:
		return fmt.Sprint("Message ID is invalid")
	case ESME_RBINDFAIL:
		return fmt.Sprint("Bind Failed")
	case ESME_RINVPASWD:
		return fmt.Sprint("Invalid Password")
	case ESME_RINVSYSID:
		return fmt.Sprint("Invalid System ID")
	case ESME_RCANCELFAIL:
		return fmt.Sprint("Cancel SM Failed")
	case ESME_RREPLACEFAIL:
		return fmt.Sprint("Replace SM Failed")
	case ESME_RMSGQFUL:
		return fmt.Sprint("Message Queue Full")
	case ESME_RINVSERTYP:
		return fmt.Sprint("Invalid Service Type")
	case ESME_RINVNUMDESTS:
		return fmt.Sprint("Invalid number of destinations")
	case ESME_RINVDLNAME:
		return fmt.Sprint("Invalid Distribution List name")
	case ESME_RINVDESTFLAG:
		return fmt.Sprint("Destination flag is invalid")
	case ESME_RINVSUBREP:
		return fmt.Sprint("Invalid 'submit with replace' request")
	case ESME_RINVESMCLASS:
		return fmt.Sprint("Invalid esm_class field data")
	case ESME_RCNTSUBDL:
		return fmt.Sprint("Cannot Submit to Distribution List")
	case ESME_RSUBMITFAIL:
		return fmt.Sprint("submit_sm or submit_multi failed")
	case ESME_RINVSRCTON:
		return fmt.Sprint("Invalid Source address TON")
	case ESME_RINVSRCNPI:
		return fmt.Sprint("Invalid Source address NPI")
	case ESME_RINVDSTTON:
		return fmt.Sprint("Invalid Destination address TON")
	case ESME_RINVDSTNPI:
		return fmt.Sprint("Invalid Destination address NPI")
	case ESME_RINVSYSTYP:
		return fmt.Sprint("Invalid system_type field")
	case ESME_RINVREPFLAG:
		return fmt.Sprint("Invalid replace_if_present flag")
	case ESME_RINVNUMMSGS:
		return fmt.Sprint("Invalid number of messages")
	case ESME_RTHROTTLED:
		return fmt.Sprint("Throttling error (ESME has exceeded allowed message limit")
	case ESME_RINVSCHED:
		return fmt.Sprint("Invalid Scheduled Delivery Time")
	case ESME_RINVEXPIRY:
		return fmt.Sprint("Invalid message validity period (Expiry time)")
	case ESME_RINVDFTMSGID:
		return fmt.Sprint("Predefined Message Invalid or Not Found")
	case ESME_RX_T_APPN:
		return fmt.Sprint("ESME Receiver Temporary App Error Code")
	case ESME_RX_P_APPN:
		return fmt.Sprint("ESME Receiver Permanent App Error Code")
	case ESME_RX_R_APPN:
		return fmt.Sprint("ESME Receiver Reject Message Error Code")
	case ESME_RQUERYFAIL:
		return fmt.Sprint("Query_sm request failed")
	case ESME_RINVOPTPARSTREAM:
		return fmt.Sprint("Error in the optional part of the PDU Body.")
	case ESME_ROPTPARNOTALLWD:
		return fmt.Sprint("Optional Parameter not allowed")
	case ESME_RINVPARLEN:
		return fmt.Sprint("Invalid Parameter Length.")
	case ESME_RMISSINGOPTPARAM:
		return fmt.Sprint("Expected Optional Parameter missing")
	case ESME_RINVOPTPARAMVAL:
		return fmt.Sprint("Invalid Optional Parameter Value")
	case ESME_RDELIVERYFAILURE:
		return fmt.Sprint("Delivery Failure (used for data_sm_resp)")
	case ESME_RUNKNOWNERR:
		return fmt.Sprint("Unknown Error")
	}
}

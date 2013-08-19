package smpp34

const (
	// SMPP Protocol Version
	SMPP_VERSION = 0x34

	// Sequence number start/end
	SEQUENCE_NUM_START = 0x00000001
	SEQUENCE_NUM_END   = 0x7FFFFFFF
)

const (
	// ESME Error Constants
	ESME_ROK              = 0x00000000 // OK!
	ESME_RINVMSGLEN       = 0x00000001 // Message Length is invalid
	ESME_RINVCMDLEN       = 0x00000002 // Command Length is invalid
	ESME_RINVCMDID        = 0x00000003 // Invalid Command ID
	ESME_RINVBNDSTS       = 0x00000004 // Incorrect BIND Status for given com-
	ESME_RALYBND          = 0x00000005 // ESME Already in Bound State
	ESME_RINVPRTFLG       = 0x00000006 // Invalid Priority Flag
	ESME_RINVREGDLVFLG    = 0x00000007 // Invalid Registered Delivery Flag
	ESME_RSYSERR          = 0x00000008 // System Error
	ESME_RINVSRCADR       = 0x0000000A // Invalid Source Address
	ESME_RINVDSTADR       = 0x0000000B // Invalid Dest Addr
	ESME_RINVMSGID        = 0x0000000C // Message ID is invalid
	ESME_RBINDFAIL        = 0x0000000D // Bind Failed
	ESME_RINVPASWD        = 0x0000000E // Invalid Password
	ESME_RINVSYSID        = 0x0000000F // Invalid System ID
	ESME_RCANCELFAIL      = 0x00000011 // Cancel SM Failed
	ESME_RREPLACEFAIL     = 0x00000013 // Replace SM Failed
	ESME_RMSGQFUL         = 0x00000014 // Message Queue Full
	ESME_RINVSERTYP       = 0x00000015 // Invalid Service Type
	ESME_RINVNUMDESTS     = 0x00000033 // Invalid number of destinations
	ESME_RINVDLNAME       = 0x00000034 // Invalid Distribution List name
	ESME_RINVDESTFLAG     = 0x00000040 // Destination flag is invalid
	ESME_RINVSUBREP       = 0x00000042 // Invalid 'submit with replace' request
	ESME_RINVESMCLASS     = 0x00000043 // Invalid esm_class field data
	ESME_RCNTSUBDL        = 0x00000044 // Cannot Submit to Distribution List
	ESME_RSUBMITFAIL      = 0x00000045 // submit_sm or submit_multi failed
	ESME_RINVSRCTON       = 0x00000048 // Invalid Source address TON
	ESME_RINVSRCNPI       = 0x00000049 // Invalid Source address NPI
	ESME_RINVDSTTON       = 0x00000050 // Invalid Destination address TON
	ESME_RINVDSTNPI       = 0x00000051 // Invalid Destination address NPI
	ESME_RINVSYSTYP       = 0x00000053 // Invalid system_type field
	ESME_RINVREPFLAG      = 0x00000054 // Invalid replace_if_present flag
	ESME_RINVNUMMSGS      = 0x00000055 // Invalid number of messages
	ESME_RTHROTTLED       = 0x00000058 // Throttling error (ESME has exceeded allowed message limits)
	ESME_RINVSCHED        = 0x00000061 // Invalid Scheduled Delivery Time
	ESME_RINVEXPIRY       = 0x00000062 // Invalid message validity period (Expiry time)
	ESME_RINVDFTMSGID     = 0x00000063 // Predefined Message Invalid or Not Found
	ESME_RX_T_APPN        = 0x00000064 // ESME Receiver Temporary App Error Code
	ESME_RX_P_APPN        = 0x00000065 // ESME Receiver Permanent App Error Code
	ESME_RX_R_APPN        = 0x00000066 // ESME Receiver Reject Message Error Code
	ESME_RQUERYFAIL       = 0x00000067 // Query_sm request failed
	ESME_RINVOPTPARSTREAM = 0x000000C0 // Error in the optional part of the PDU Body
	ESME_ROPTPARNOTALLWD  = 0x000000C1 // Optional Parameter not allowed
	ESME_RINVPARLEN       = 0x000000C2 // Invalid Parameter Length
	ESME_RMISSINGOPTPARAM = 0x000000C3 // Expected Optional Parameter missing
	ESME_RINVOPTPARAMVAL  = 0x000000C4 // Invalid Optional Parameter Value
	ESME_RDELIVERYFAILURE = 0x000000FE // Delivery Failure (used for data_sm_resp)
	ESME_RUNKNOWNERR      = 0x000000FF // Unknown Error
)

const (
	// PDU Types
	GENERIC_NACK          = 0x80000000
	BIND_RECEIVER         = 0x00000001
	BIND_RECEIVER_RESP    = 0x80000001
	BIND_TRANSMITTER      = 0x00000002
	BIND_TRANSMITTER_RESP = 0x80000002
	QUERY_SM              = 0x00000003
	QUERY_SM_RESP         = 0x80000003
	SUBMIT_SM             = 0x00000004
	SUBMIT_SM_RESP        = 0x80000004
	DELIVER_SM            = 0x00000005
	DELIVER_SM_RESP       = 0x80000005
	UNBIND                = 0x00000006
	UNBIND_RESP           = 0x80000006
	REPLACE_SM            = 0x00000007
	REPLACE_SM_RESP       = 0x80000007
	CANCEL_SM             = 0x00000008
	CANCEL_SM_RESP        = 0x80000008
	BIND_TRANSCEIVER      = 0x00000009
	BIND_TRANSCEIVER_RESP = 0x80000009
	OUTBIND               = 0x0000000B
	ENQUIRE_LINK          = 0x00000015
	ENQUIRE_LINK_RESP     = 0x80000015
	SUBMIT_MULTI          = 0x00000021
	SUBMIT_MULTI_RESP     = 0x80000021
	ALERT_NOTIFICATION    = 0x00000102
	DATA_SM               = 0x00000103
	DATA_SM_RESP          = 0x80000103
)

package buffer

type CPacket struct {
	Code int
}

type CmdCode interface {
	Write(buf *ByteBuf)
	Read(buf *ByteBuf)
	Copy() CmdCode
}
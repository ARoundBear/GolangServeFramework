package znet

import (
	"ZinxLearning/zinx/utils"
	"ZinxLearning/zinx/ziface"
	"bytes"
	"encoding/binary"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	return 8
}

func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(binaryData)

	msg := &Message{}

	err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	err = binary.Read(dataBuf, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, err //errors.New("too Large msg data recv!")
	}
	return msg, nil
}

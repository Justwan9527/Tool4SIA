package SIA

import (
	"fmt"

	"errors"

	"time"

	"github.com/sigurn/crc16"
)

type Sia struct {
	ID    string
	Rrcvr string
	Lpref string
	Acct  string
	Data  string
	ISeq  uint16
}

func (obj *Sia) GenSiaData() ([]byte, error) {
	if obj.ID != "SIA-DCS" {
		//fmt.Println("Unsupported ID:", obj.ID)

		return []byte{}, errors.New("Unsupported ID:" + obj.ID)
	}

	//check sia field
	/*
		if len(obj.Rrcvr) == 0 || len(obj.Rrcvr) > 6 {

		}
	*/

	//
	//ID Token, SIA-DCS or ADM-CID
	var id string = "SIA-DCS"

	//fmt.Println("id:", id)

	//seq from 0001-9999, four ASCII characters
	obj.ISeq++
	var seq string = fmt.Sprintf("%04d", obj.ISeq)
	//fmt.Println("seq:", seq)

	//Rrcvr. receiver number consists of an ASCII "R", follow by 1-6 HEX ASCII digits.
	var Rrcvr string = obj.Rrcvr

	Rrcvr = "R" + Rrcvr

	//fmt.Println("Rrcvr:", Rrcvr)

	//Lpref. account prefix consists of an ASCII "L", follow by 1-6 HEX ASCII disits.
	var LPref string = obj.Lpref
	LPref = "L" + LPref
	//fmt.Println("LPref:", LPref)

	//acct. account number consists of an ASCII "#", followed by 3-16 ASCII characters representing hexadecimal digits

	var acct string = obj.Acct
	//fmt.Println("acct:", acct)

	acct = "#" + acct

	//[data] all data is in ASCII characters and the bracket characters "[" and "]" are included. "NF1234/NPAZone123"

	var data string = "[" + acct + "|" + obj.Data + "]"
	//fmt.Println("data:", data)

	//timestamp. the format of the timestamp is:<_HH:MM:SS,MM-DD-YYYY>. the braces are not part of the transmitted message. But the dunerscore, colon, commaand hyphen characters are included
	currentTime := time.Now()
	var year, month, day int = currentTime.Year(), (int)(currentTime.Month()), currentTime.Day()
	var hour, minute, second int = currentTime.Hour(), currentTime.Minute(), currentTime.Second()

	//formattedTime := currentTime.Format("16:08:34,03-06-2024")
	formattedTime := fmt.Sprintf("%02d:%02d:%02d,%02d-%02d-%04d", hour, minute, second, day, month, year)
	var timestamp string = "_" + formattedTime
	//fmt.Println("timestamp:", timestamp)

	var data4crc string = "\"" + id + "\"" + seq + Rrcvr + LPref + acct + data + timestamp
	//fmt.Println("data4crc:", data4crc, "length:", len(data4crc))

	//length of the message
	var lLLL int = len(data4crc)

	//convert 0LLL to 3 hex digits(in ASCII)
	var sLLL string = fmt.Sprintf("0%03x", lLLL)
	//fmt.Println("sLLL:", sLLL)
	//calculate crc

	table := crc16.MakeTable(crc16.CRC16_ARC)
	crc := crc16.Checksum([]byte(data4crc), table)
	c := fmt.Sprintf("%X", crc)
	//fmt.Println("CRC16==", c)

	//complete SIA data
	var siaData string = "\r" + c + sLLL + data4crc + "\n"
	//fmt.Println("siaData:", siaData)

	return []byte(siaData), nil

}

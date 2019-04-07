package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/monnand/dhkx"
	log "github.com/sirupsen/logrus"
	"math/big"
	"net"
)

func (sess *ElinkSession) handlePacket(conn net.Conn, packet Elink) error {
	log.Debugln("<<<<<<<<<<<<<<<<<<<<handlePacket>>>>>>>>>>>>>>>>>>>>")
	log.Debugln("handlePacket start, and sess is ", sess)
	log.Debugln(packet.Flag, packet.Len)
	log.Debugf("\n%s\n", hex.Dump(packet.Data))
	base := Base{}
	if sess.key == nil {
		err := json.Unmarshal(packet.Data, &base)
		if err != nil {
			log.Errorln("json unmarshal error:", err)
			return err
		}
	} else {
		//解密
		log.Debugln("cipher data:")
		log.Debugf("\n%s\n", hex.Dump(packet.Data))
		log.Debugln("decrypt start")
		data, err := sess.Decrypt(packet.Data)
		if err != nil {
			log.Errorln("decrypt error:", err)
			return err
		}
		log.Debugln("plain data:")
		log.Debugf("\n%s\n", hex.Dump(data))
		packet.DeCryData = data
		err = json.Unmarshal(data, &base)
		if err != nil {
			return err
		}
	}

	switch base.Type {
	case ElinkTypeKeyNgReq:
		log.Println("recv keyngreq")
		handleErr(sess.handleKeyNgReq(packet))
	case ElinkTypeKeyNgAck:
		log.Println("recv keyngack")
	case ElinkTypeDH:
		log.Println("recv dh")
		handleErr(sess.handleDH(packet))
		log.Printf("after recv dh", sess)
	case ElinkTypeDevReg:
		log.Println("recv dev_reg")
		handleErr(sess.handleDevReg(packet))
		handleErr(sess.getStatus())
	case ElinkTypeKeepAlive:
		log.Println("recv keepalive")
		handleErr(sess.handleKeepAlive(packet))
	case ElinkTypeAck:
		log.Println("recv ack")
	case ElinkTypeCfg:
		log.Println("recv cfg")
	case ElinkTypeGetStatus:
		log.Println("recv get_status")
	case ElinkTypeStatus:
		log.Println("recv status")
		handleErr(sess.handleStatus(packet))
	case ElinkTypeRealDevInfo:
		log.Println("recv real dev info")
		handleErr(sess.handleRealDevInfo(packet))
	default:
		log.Errorln("unknown type")
	}
	return nil
}

func handleErr(err error) {
	if err != nil {
		log.Errorln(err)
	}
}

func (sess *ElinkSession) handleKeyNgReq(packet Elink) error {
	j := KeyNgReq{}
	err := json.Unmarshal(packet.Data, &j)
	if err != nil {
		log.Println("json unmarshal keyngreq packet error")
		return err
	}
	log.Println("%+v", j)
	support := false
	for _, mode := range j.KeyModeList {
		if mode.KeyMode == "dh" {
			support = true
		}
	}

	sess.PerMac = j.Mac
	sess.Seq.RecvSeq = j.Seq

	if !support {
		return errors.New("not support key mode")
	}

	ackJ := KeyNgAck{
		Base{
			Type: ElinkTypeKeyNgAck,
			Seq:  j.Seq,
			Mac:  sess.Mac,
		},
		"dh",
	}

	data, err := json.Marshal(ackJ)
	if err != nil {
		return err
	}

	elink := Elink{
		Data: data,
	}
	_, err = sess.Write(serialPacket(elink))
	if err != nil {
		return err
	}

	return nil
}

func (sess *ElinkSession) handleDH(packet Elink) error {
	dhj := DH{}
	err := json.Unmarshal(packet.Data, &dhj)
	if err != nil {
		return err
	}

	err = sess.updateSeq(dhj.Seq)
	if err != nil {
		return err
	}

	gBytes, _ := base64.StdEncoding.DecodeString(dhj.Data.DHG)
	pBytes, _ := base64.StdEncoding.DecodeString(dhj.Data.DHP)
	kBytes, _ := base64.StdEncoding.DecodeString(dhj.Data.DHKey)

	var gInt big.Int
	var pInt big.Int
	gInt.SetBytes(gBytes)
	pInt.SetBytes(pBytes)

	group := dhkx.CreateGroup(&pInt, &gInt)
	bobPrivateKey, err := group.GeneratePrivateKey(rand.Reader)
	if err != nil {
		return err
	}

	alicePubKey := dhkx.NewPublicKey(kBytes)
	encKey, err := group.ComputeKey(alicePubKey, bobPrivateKey)
	if err != nil {
		return err
	}

	sess.key = encKey.Bytes()
	pub := bobPrivateKey.Bytes()
	pubBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(pub)))
	base64.StdEncoding.Encode(pubBase64, pub)

	dhAck := DH{
		Base{
			"dh",
			dhj.Seq,
			sess.Mac,
		},
		DHData{
			DHKey: string(pubBase64),
			DHG:   dhj.Data.DHG,
			DHP:   dhj.Data.DHP,
		},
	}

	data, err := json.Marshal(dhAck)
	if err != nil {
		return err
	}

	elink := Elink{
		Data: data,
	}
	_, err = sess.Write(serialPacket(elink))
	if err != nil {
		return err
	}

	log.Println("sent dh")
	log.Debugf("\n%s\n", hex.Dump(serialPacket(elink)))
	log.Printf("sess is \n%+v\n", sess)

	return nil
}

func (sess *ElinkSession) handleDevReg(packet Elink) error {
	j := DevReg{}
	err := json.Unmarshal(packet.DeCryData, &j)
	if err != nil {
		return err
	}

	sess.DevData = j.Data
	err = sess.updateSeq(j.Seq)
	if err != nil {
		return err
	}

	return sess.sendAck()
}

func (sess ElinkSession) getStatus() error {
	j := GetStatus{
		Base{
			Type: "get_status",
			Seq:  sess.Seq.SendSeq + 1,
			Mac:  sess.Mac,
		},
		[]Get{
			{Name: StatusNameWiFi},
			{Name: StatusNameWiFiSwitch},
			{Name: StatusNameLedSwitch},
			{Name: StatusNameWiFiTimer},
		},
	}

	log.Println("sent get_status")
	return sess.writePacket(j)
}

func (sess ElinkSession) getAPData() error {
	j := GetStatus{
		Base{
			Type: "get_ap_data",
			Seq:  sess.Seq.SendSeq + 1,
			Mac:  sess.Mac,
		},
		[]Get{
			{Name: APDataNameCPURate},
			{Name: APDataNameMem},
			{Name: APDataNameUploadSpeed},
			{Name: APDataNameDownloadSpeed},
			{Name: APDataNameOnlineTime},
			{Name: APDataNameNum},
			{Name: APDataNameChannel},
			{Name: APDataNameLoad},
		},
	}

	log.Println("sent get_ap_data")

	return sess.writePacket(j)
}

func (sess ElinkSession) getDevInfo() error {
	j := GetStatus{
		Base{
			Type: "get_real_devinfo",
			Seq:  sess.Seq.SendSeq + 1,
			Mac:  sess.Mac,
		},
		[]Get{},
	}

	log.Println("sent get_real_devinfo")
	return sess.writePacket(j)
}

func (sess ElinkSession) wifiConfig(cfg WiFiSet) error {
	j := WifiConfig{
		Base: Base{
			Type: "cfg",
			Mac:  sess.Mac,
			Seq:  sess.Seq.SendSeq + 1,
		},
		//Status: sess.Status,
		Set: cfg,
	}

	log.Println("sent wifi config")

	return sess.writePacket(j)
}

func (sess ElinkSession) switchConfig(cfg SwitchSet) error {
	j := SwitchConfig{
		Base{
			Type: "cfg",
			Mac:  sess.Mac,
			Seq:  sess.Seq.SendSeq + 1,
		},
		cfg,
	}

	log.Println("sent switch config")

	return sess.writePacket(j)
}

func (sess ElinkSession) upgradeConfig(cfg UpgradeSet) error {
	j := UpgradeConfig{
		Base{
			Type: "cfg",
			Mac:  sess.Mac,
			Seq:  sess.Seq.SendSeq + 1,
		},
		cfg,
	}

	log.Println("sent switch config")

	return sess.writePacket(j)
}

func (sess ElinkSession) reboot() error {
	j := Reboot{
		Base{
			Type: "cfg",
			Mac:  sess.Mac,
			Seq:  sess.Seq.SendSeq + 1,
		},
		RebootSet{
			"reboot",
		},
	}

	log.Println("sent reboot config")

	return sess.writePacket(j)
}

func (sess ElinkSession) reset() error {
	j := Reboot{
		Base{
			Type: "cfg",
			Mac:  sess.Mac,
			Seq:  sess.Seq.SendSeq + 1,
		},
		RebootSet{
			"reset",
		},
	}

	log.Println("sent reboot config")

	return sess.writePacket(j)
}

func (sess *ElinkSession) handleStatus(packet Elink) error {
	j := Status{}
	err := json.Unmarshal(packet.DeCryData, &j)
	if err != nil {
		return err
	}
	sess.Status = j.Status
	log.Printf("status is \n%+v\n", j)
	return sess.sendAck()
}

func (sess ElinkSession) handleKeepAlive(packet Elink) error {
	j := KeepAlive{}
	err := json.Unmarshal(packet.DeCryData, &j)
	if err != nil {
		return err
	}

	err = sess.updateSeq(j.Seq)
	if err != nil {
		return err
	}

	return sess.sendAck()
}

func (sess *ElinkSession) handleRealDevInfo(packet Elink) error {
	j := RealDevInfo{}
	err := json.Unmarshal(packet.DeCryData, &j)
	if err != nil {
		return err
	}

	sess.RealDevInfo = j.RealDev
	err = sess.updateSeq(j.Seq)
	if err != nil {
		return err
	}

	return sess.sendAck()
}

func (sess ElinkSession) sendAck() error {
	j := Ack{
		Type: "ack",
		Mac:  sess.Mac,
		Seq:  sess.Seq.RecvSeq,
	}

	data, err := json.Marshal(j)
	if err != nil {
		return err
	}

	cipherData, err := sess.Encrypt(data)
	if err != nil {
		return err
	}

	packet := Elink{
		Data: cipherData,
	}

	log.Println("sent ack")
	log.Debugf("\n%s\n", hex.Dump(data))
	_, err = sess.Write(serialPacket(packet))
	if err != nil {
		return err
	}
	return nil
}

func (sess ElinkSession) Decrypt(ciphertext []byte) ([]byte, error) {
	// no iv ?
	block, err := aes.NewCipher(sess.key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	//iv := ciphertext[:aes.BlockSize]
	//ciphertext = ciphertext[aes.BlockSize:]
	iv := make([]byte, aes.BlockSize)

	mod := cipher.NewCBCDecrypter(block, iv)

	mod.CryptBlocks(ciphertext, ciphertext)

	i := len(ciphertext) - 1
	for ; i >= 0; i-- {
		if ciphertext[i] != 0x00 {
			break
		}
	}

	return ciphertext[:i+1], nil
}

func (sess ElinkSession) Encrypt(plaintext []byte) ([]byte, error) {
	// no iv ?
	if len(plaintext)%aes.BlockSize != 0 {
		padLen := aes.BlockSize - len(plaintext)%aes.BlockSize
		plaintext = append(plaintext, make([]byte, padLen)...)
	}

	block, err := aes.NewCipher(sess.key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	/*	iv := ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil , err
		}*/
	iv := make([]byte, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext[aes.BlockSize:], nil
}

func (sess *ElinkSession) updateSeq(newSeq int) error {
	if newSeq-sess.Seq.RecvSeq != 1 {
		return errors.New("seq error")
	} else {
		sess.Seq.RecvSeq = newSeq
	}
	return nil
}

func (sess *ElinkSession) writePacket(j interface{}) error {
	data, err := json.Marshal(j)
	if err != nil {
		return err
	}

	cipherData, err := sess.Encrypt(data)
	if err != nil {
		return err
	}

	packet := Elink{
		Data: cipherData,
	}

	log.Debugf("\n%s\n", hex.Dump(data))
	_, err = sess.Write(serialPacket(packet))
	if err != nil {
		return err
	}

	sess.Seq.SendSeq += 1
	return nil
}

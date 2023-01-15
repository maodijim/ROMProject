package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/des"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"reflect"
	"time"

	Cmd "ROMProject/Cmds"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func Decompress(data []byte) ([]byte, error) {
	f := bytes.NewReader(data)
	z, err := zlib.NewReader(f)
	if err != nil {
		log.Warnf("failed to decompress data %s", err)
		return data, err
	}
	result, err := ioutil.ReadAll(z)
	return result, err
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// w := zlib.NewWriter(&b)
	w, err := zlib.NewWriterLevel(&b, zlib.NoCompression)
	_, err = w.Write(data)
	w.Close()
	return b.Bytes(), err
}

func EncryptDesEcb(data, cipherKey []byte) ([]byte, error) {
	block, err := des.NewCipher(cipherKey)
	if err != nil {
		return nil, err
	}
	blockSzie := block.BlockSize()
	// Zero Padding
	padding := blockSzie - len(data)%blockSzie
	padtxt := bytes.Repeat([]byte{0}, padding)
	data = append(data, padtxt...)
	crypted := make([]byte, len(data))
	dst := crypted
	for len(data) > 0 {
		block.Encrypt(dst, data[:blockSzie])
		data = data[blockSzie:]
		dst = dst[blockSzie:]
	}
	return crypted, nil
}

func DecryptDesEcb(data, cipherKey []byte) ([]byte, error) {
	block, err := des.NewCipher(cipherKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		return nil, errors.New("input not full block")
	}
	decrypted := make([]byte, len(data))
	dst := decrypted
	for len(data) > 0 {
		block.Decrypt(dst, data[:blockSize])
		data = data[blockSize:]
		dst = dst[blockSize:]
	}
	// Zero UnPadding
	decrypted = bytes.TrimFunc(decrypted, func(r rune) bool {
		return r == rune(0)
	})
	return decrypted, nil
}

func isEncrypt(body []byte) bool {
	if len(body) >= 1 && body[0]&2 == 2 {
		return true
	}
	return false
}

func isCompress(body []byte) bool {
	if len(body) >= 1 && body[0]&1 == 1 {
		return true
	}
	return false
}

func isClientCMD(id1, id2 int) bool {
	return id1 == 99 && id2 == 1
}

func IsValidFlag(flag byte) bool {
	for _, f := range TcpFlag {
		if f[0] == flag {
			return true
		}
	}
	return false
}

func parseCmdNonce(body []byte) (hasNonce bool, hasData bool, nonceSize int) {
	if len(body) > 2 {
		nonceSize = int(body[2])
	} else {
		return false, false, 0
	}
	startPoint := nonceSize + 4
	if (nonceSize == 44 || startPoint == 50 || nonceSize == 50) && body[3] == 0 {
		hasNonce = true
	} else {
		hasData = true
	}
	if len(body) > startPoint {
		hasData = true
	}
	return hasNonce, hasData, nonceSize
}

func GetContentSize(flag, lengthHeader []byte) int {
	contentLength := int(binary.LittleEndian.Uint16(lengthHeader))
	if isEncrypt(flag) {
		if contentLength%8 != 0 {
			contentLength = 8 - contentLength%8 + contentLength
		}
	}
	return contentLength
}

const (
	ErrParseWireType    = "proto: cannot parse reserved wire type"
	ErrMisMatchEndGroup = "proto: mismatching end group marker"
	ErrInvalidField     = "proto: invalid field number"
)

func ParseCmd(data []byte, param proto.Message) (err error) {
	hasNonce, hasData, nonceSize := parseCmdNonce(data)
	startPoint := int(math.Min(float64(nonceSize+4), float64(len(data))))
	if hasNonce {
		err = proto.Unmarshal(data[startPoint:], param)
	} else if hasData {
		err = proto.Unmarshal(data[2:], param)
	}
	if err != nil {
		switch err.Error() {
		case ErrParseWireType, ErrMisMatchEndGroup, ErrInvalidField, io.ErrUnexpectedEOF.Error():
			log.Debugf("parse cmd error: %s", err)
		default:
			log.Debugf("parse cmd error: %s", err)
		}
	}

	if hasNonce {
		printNonceStr(data, startPoint)
	}
	return err
}

func ConstructBody(cmdId, cmdParamId int32, flag, body, nonce, cipherKey []byte) []byte {
	// newBodyLength := len(body) + 2
	var newBody []byte
	if nonce != nil {
		newBody = make([]byte, 4)
		newBody[0] = byte(cmdId)
		newBody[1] = byte(cmdParamId)
		newBody[2] = byte(len(nonce))
		newBody[3] = 0
		newBody = append(newBody[:], nonce[:]...)
		newBody = append(newBody[:], body[:]...)
	} else {
		newBody = make([]byte, 2)
		newBody[0] = byte(cmdId)
		newBody[1] = byte(cmdParamId)
		newBody = append(newBody[:], body[:]...)
	}
	if isCompress(flag) {
		result, err := Compress(newBody)
		if err != nil {
			log.Printf("failed to compress body: %s", err)
			return body
		}
		newBody = result
	}
	if isEncrypt(flag) {
		result, err := EncryptDesEcb(newBody, cipherKey)
		if err != nil {
			log.Printf("failed to encrypt body: %s", err)
			return newBody
		}
		newBody = result
	}
	return newBody
}

func GetTimeNow(inMili bool) int64 {
	if inMili {
		return time.Now().UnixNano() / int64(time.Millisecond)
	}
	return time.Now().Unix()
}

func RandomSleepTime(max, min int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func PrintTranslateMsgResult(cmdParamName string, err error, param proto.Message) {
	log.Infof("Parsing %s", cmdParamName)
	if err != nil {
		log.Warnf("failed to parse protobuf: %s", err)
	}
	log.Infof("unmarshal: %v", param)
}

func ParseBody(body, cipherKey []byte) [][]byte {
	bodyLength := len(body)
	offset := 0
	headerLength := 3
	var cmdList [][]byte
	for bodyLength > offset {
		if len(body[offset:]) > headerLength {
			// log.Printf("starting at position %d with header %x", offset, body[offset:offset+headerLength])
			contentLength := int(binary.LittleEndian.Uint16(body[offset:][1:]))
			start := offset + headerLength
			end := start
			var newBody []byte
			if isEncrypt(body[offset:]) {
				if contentLength%8 != 0 {
					contentLength = 8 - contentLength%8 + contentLength
				}
				end += contentLength
				if end > bodyLength {
					offset = bodyLength
					continue
				}
				newBody = body[start:end]
				decrypted, err := DecryptDesEcb(newBody, cipherKey)
				if err != nil {
					log.Warnf("failed to decrypt: %s", err)
					offset = end
					continue
				}
				newBody = decrypted
			} else {
				end += contentLength
				newBody = body[start:end]
			}
			if isCompress(body[offset:]) {
				newBody, _ = Decompress(newBody)
			}
			offset = end
			cmdList = append(cmdList, newBody)
			if len(newBody) > 2 && isClientCMD(int(newBody[0]), int(newBody[1])) {
				// log.Printf("Client CMD found contiune parsing...")
				cmdList = append(cmdList, doSubParse(newBody)...)
			}
		}
	}
	return cmdList
}

func doSubParse(body []byte) [][]byte {
	bodyLength := len(body)
	var cmdList [][]byte
	offset := 4
	for offset < bodyLength && bodyLength-offset > 2 {
		contentLength := int(binary.LittleEndian.Uint16(body[offset:]))
		if contentLength > bodyLength {
			offset = bodyLength
			continue
		}
		start := offset + 2
		end := start + contentLength
		if end > bodyLength {
			end = bodyLength
		}
		newBody := body[start:end]
		cmdList = append(cmdList, newBody)
		offset = end
	}
	return cmdList
}

func printNonceStr(body []byte, startPoint int) {
	noStr := &Cmd.Nonce{}
	err1 := proto.Unmarshal(body[4:startPoint], noStr)
	ts := time.Unix(int64(noStr.GetTimestamp()), 0)
	log.Printf("date: %s, Nonce Str: %v; err: %v", ts, noStr, err1)
}

func getNonce(includeTime bool, currentIndex *uint32) []byte {
	currentTime := int64(0)
	if includeTime {
		currentTime = GetTimeNow(false)
	}
	*currentIndex += 1
	sign := fmt.Sprintf("%d_%d_!^ro&", currentTime, currentTime)
	signStr := fmt.Sprintf("%x", sha1.Sum([]byte(sign)))
	nonce := &Cmd.Nonce{
		Index: currentIndex,
		Sign:  &signStr,
	}
	if includeTime {
		newT := uint32(currentTime)
		nonce.Timestamp = &newT
	}
	pOut, _ := proto.Marshal(nonce)
	return pOut
}

func GetNpcAttrValByType(attrs []*Cmd.UserAttr, attrType Cmd.EAttrType) (attrVal int32) {
	for _, attr := range attrs {
		if attr.GetType() == attrType {
			attrVal = attr.GetValue()
		}
	}
	return attrVal
}

func GetNpcDataValByType(datas []*Cmd.UserData, dataType Cmd.EUserDataType) (dataVal uint64) {
	for _, data := range datas {
		if data.GetType() == dataType {
			dataVal = data.GetValue()
		}
	}
	return dataVal
}

func GetMemberDataByType(datas []*Cmd.MemberData, dataType Cmd.EMemberData) (dataVal uint64) {
	for _, data := range datas {
		if data.GetType() == dataType {
			dataVal = data.GetValue()
		}
	}
	return dataVal
}

func GetRandom(min, max int) int {
	return rand.Intn(max-min) + min
}

func SliceContain(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if reflect.TypeOf(a) == reflect.TypeOf(e) && a == e {
			return true
		}
	}
	return false
}

func StrSliceContain(strList []string, e string) bool {
	for _, a := range strList {
		if a == e {
			return true
		}
	}
	return false
}

func Uint64SliceContains(s []*uint64, e uint64) bool {
	for _, a := range s {
		if *a == e {
			return true
		}
	}
	return false
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func MapContains[T comparable](s map[T]bool, e T) bool {
	_, ok := s[e]
	return ok
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomZhCharacterName(length int) string {
	var name string
	lastChar := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	fourthChar := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	thirdChar := []string{"名", "轩", "汉", "姝", "逸", "宫", "布", "菲", "星", "托"}
	secChar := []string{"问", "叶", "时", "季", "乐", "月", "荦", "醉", "梦", "生"}
	firstChar := []string{"明", "猎", "登", "高", "长", "伪", "罗", "净", "空", "竹"}
	rand.Seed(time.Now().UnixNano())
	name += firstChar[rand.Intn(len(firstChar))]
	rand.Seed(time.Now().UnixNano())
	name += secChar[rand.Intn(len(secChar))]
	rand.Seed(time.Now().UnixNano())
	name += thirdChar[rand.Intn(len(thirdChar))]
	rand.Seed(time.Now().UnixNano())
	name += fourthChar[rand.Intn(len(fourthChar))]
	rand.Seed(time.Now().UnixNano())
	name += lastChar[rand.Intn(len(lastChar))]
	return name
}

func GetAttrPointReq(currentPoint int32) int32 {
	return int32(math.Floor(float64(currentPoint)/9) + 2)
}

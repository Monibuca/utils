package codec

import (
	"io"

	"github.com/Monibuca/utils/v3"
)

const (
	// FLV Tag Type
	FLV_TAG_TYPE_AUDIO  = 0x08
	FLV_TAG_TYPE_VIDEO  = 0x09
	FLV_TAG_TYPE_SCRIPT = 0x12
)

var (
	Codec2SoundFormat = map[string]byte{
		"aac":  10,
		"pcma": 7,
		"pcmu": 8,
	}
	// 音频格式. 4 bit
	SoundFormat = map[byte]string{
		0:  "Linear PCM, platform endian",
		1:  "ADPCM",
		2:  "MP3",
		3:  "Linear PCM, little endian",
		4:  "Nellymoser 16kHz mono",
		5:  "Nellymoser 8kHz mono",
		6:  "Nellymoser",
		7:  "G.711 A-law logarithmic PCM",
		8:  "G.711 mu-law logarithmic PCM",
		9:  "reserved",
		10: "AAC",
		11: "Speex",
		14: "MP3 8Khz",
		15: "Device-specific sound"}

	// 采样频率. 2 bit
	SoundRate = map[byte]int{
		0: 5500,
		1: 11000,
		2: 22000,
		3: 44000}

	// 量化精度. 1 bit
	SoundSize = map[byte]string{
		0: "8Bit",
		1: "16Bit"}

	// 音频类型. 1bit
	SoundType = map[byte]string{
		0: "Mono",
		1: "Stereo"}

	// 视频帧类型. 4bit
	FrameType = map[byte]string{
		1: "keyframe (for AVC, a seekable frame)",
		2: "inter frame (for AVC, a non-seekable frame)",
		3: "disposable inter frame (H.263 only)",
		4: "generated keyframe (reserved for server use only)",
		5: "video info/command frame"}

	// 视频编码类型. 4bit
	CodecID = map[byte]string{
		1:  "JPEG (currently unused)",
		2:  "Sorenson H.263",
		3:  "Screen video",
		4:  "On2 VP6",
		5:  "On2 VP6 with alpha channel",
		6:  "Screen video version 2",
		7:  "AVC",
		12: "H265"}
)

var FLVHeader = []byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0, 0, 0, 9, 0, 0, 0, 0}

func WriteFLVTag(w io.Writer, t byte, timestamp uint32, payload []byte) (err error) {
	head := utils.GetSlice(11)
	defer utils.RecycleSlice(head)
	tail := utils.GetSlice(4)
	defer utils.RecycleSlice(tail)
	head[0] = t
	dataSize := uint32(len(payload))
	utils.BigEndian.PutUint32(tail, dataSize+11)
	utils.BigEndian.PutUint24(head[1:], dataSize)
	head[4] = byte(timestamp >> 16)
	head[5] = byte(timestamp >> 8)
	head[6] = byte(timestamp)
	head[7] = byte(timestamp >> 24)
	if _, err = w.Write(head); err != nil {
		return
	}
	// Tag Data
	if _, err = w.Write(payload); err != nil {
		return
	}
	if _, err = w.Write(tail); err != nil { // PreviousTagSizeN(4)
		return
	}
	return
}
func ReadFLVTag(r io.Reader) (t byte, timestamp uint32, payload []byte, err error) {
	head := utils.GetSlice(11)
	defer utils.RecycleSlice(head)
	if _, err = io.ReadFull(r, head); err != nil {
		return
	}
	t = head[0]
	dataSize := utils.BigEndian.Uint24(head[1:])
	timestamp = (uint32(head[7]) << 24) | (uint32(head[4]) << 16) | (uint32(head[5]) << 8) | uint32(head[6])
	payload = make([]byte, int(dataSize))
	if _, err = io.ReadFull(r, payload); err == nil {
		t := utils.GetSlice(4)
		_, err = io.ReadFull(r, t)
		utils.RecycleSlice(t)
	}
	return
}

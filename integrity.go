package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("zero-byte file")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil {
		return fmt.Errorf("cannot read: %w", err)
	}
	buf = buf[:n]

	format, checker := detectFormat(buf, info.Size())
	if checker == nil {
		return fmt.Errorf("unrecognized audio format")
	}

	if err := checkExtension(path, format); err != nil {
		return err
	}

	return checker()
}

func detectFormat(buf []byte, fileSize int64) (string, func() error) {
	if len(buf) < 12 {
		return "", nil
	}

	// FLAC: starts with "fLaC"
	if buf[0] == 0x66 && buf[1] == 0x4C && buf[2] == 0x61 && buf[3] == 0x43 {
		return "flac", func() error { return checkFLAC(buf) }
	}

	// Ogg: starts with "OggS"
	if buf[0] == 0x4F && buf[1] == 0x67 && buf[2] == 0x67 && buf[3] == 0x53 {
		return "ogg", func() error { return checkOgg(buf) }
	}

	// MP3 with ID3v2: starts with "ID3"
	if buf[0] == 0x49 && buf[1] == 0x44 && buf[2] == 0x33 {
		return "mp3", func() error { return checkMP3ID3(buf, fileSize) }
	}

	// MP3 raw frame: sync word FF Ex or FF Fx (11 bits set)
	if buf[0] == 0xFF && buf[1]&0xE0 == 0xE0 {
		return "mp3", func() error { return checkMP3Frame(buf) }
	}

	// M4A/AAC: "ftyp" at offset 4
	if buf[4] == 0x66 && buf[5] == 0x74 && buf[6] == 0x79 && buf[7] == 0x70 {
		return "m4a", func() error { return checkM4A(buf, fileSize) }
	}

	// WAV: "RIFF" at 0, "WAVE" at 8
	if buf[0] == 0x52 && buf[1] == 0x49 && buf[2] == 0x46 && buf[3] == 0x46 &&
		buf[8] == 0x57 && buf[9] == 0x41 && buf[10] == 0x56 && buf[11] == 0x45 {
		return "wav", func() error { return checkWAV(buf, fileSize) }
	}

	return "", nil
}

var extToFormat = map[string]string{
	".flac": "flac",
	".mp3":  "mp3",
	".ogg":  "ogg",
	".opus": "ogg",
	".m4a":  "m4a",
	".m4b":  "m4a",
	".mp4":  "m4a",
	".aac":  "m4a",
	".wav":  "wav",
	".wave": "wav",
}

func checkExtension(path string, detected string) error {
	ext := strings.ToLower(filepath.Ext(path))
	expected, known := extToFormat[ext]
	if !known {
		return nil // unknown extension, skip mismatch check
	}
	if expected != detected {
		return fmt.Errorf("extension %s but content is %s", ext, detected)
	}
	return nil
}

// FLAC: "fLaC" + STREAMINFO block (type 0, length 34, sane fields)
func checkFLAC(buf []byte) error {
	if len(buf) < 42 {
		return fmt.Errorf("flac header truncated")
	}

	// First metadata block must be STREAMINFO (type 0)
	blockType := buf[4] & 0x7F
	if blockType != 0 {
		return fmt.Errorf("flac first block is %d, expected STREAMINFO (0)", blockType)
	}

	// STREAMINFO length must be 34
	blockLen := int(buf[5])<<16 | int(buf[6])<<8 | int(buf[7])
	if blockLen != 34 {
		return fmt.Errorf("flac STREAMINFO length %d, expected 34", blockLen)
	}

	minBlock := binary.BigEndian.Uint16(buf[8:10])
	maxBlock := binary.BigEndian.Uint16(buf[10:12])
	if minBlock < 16 {
		return fmt.Errorf("flac min block size %d < 16", minBlock)
	}
	if maxBlock < minBlock {
		return fmt.Errorf("flac max block size %d < min %d", maxBlock, minBlock)
	}

	// Sample rate: 20 bits at bytes 18-20
	sampleRate := int(buf[18])<<12 | int(buf[19])<<4 | int(buf[20])>>4
	if sampleRate == 0 {
		return fmt.Errorf("flac sample rate is 0")
	}

	return nil
}

// MP3 with ID3v2 tag: validate tag header, check size fits in file
func checkMP3ID3(buf []byte, fileSize int64) error {
	if len(buf) < 10 {
		return fmt.Errorf("mp3 ID3 header truncated")
	}

	version := buf[3]
	if version != 2 && version != 3 && version != 4 {
		return fmt.Errorf("mp3 ID3v2.%d unsupported", version)
	}

	// Size bytes must be synchsafe (each < 0x80)
	for i := 6; i < 10; i++ {
		if buf[i] >= 0x80 {
			return fmt.Errorf("mp3 ID3 size byte %d invalid: 0x%02X", i-6, buf[i])
		}
	}

	tagSize := int64(buf[6])<<21 | int64(buf[7])<<14 | int64(buf[8])<<7 | int64(buf[9])
	tagSize += 10 // add header size
	if tagSize >= fileSize {
		return fmt.Errorf("mp3 ID3 tag size %d >= file size %d (no audio data)", tagSize, fileSize)
	}

	return nil
}

// MP3 raw MPEG frame: validate sync word fields
func checkMP3Frame(buf []byte) error {
	if len(buf) < 4 {
		return fmt.Errorf("mp3 frame header truncated")
	}

	// Version: bits 4-3 of byte 1
	version := (buf[1] >> 3) & 0x03
	if version == 1 {
		return fmt.Errorf("mp3 reserved MPEG version")
	}

	// Layer: bits 2-1 of byte 1
	layer := (buf[1] >> 1) & 0x03
	if layer == 0 {
		return fmt.Errorf("mp3 reserved layer")
	}

	// Bitrate index: upper 4 bits of byte 2
	bitrate := buf[2] >> 4
	if bitrate == 0x0F {
		return fmt.Errorf("mp3 invalid bitrate index")
	}

	// Sample rate index: bits 3-2 of byte 2
	sampleRate := (buf[2] >> 2) & 0x03
	if sampleRate == 0x03 {
		return fmt.Errorf("mp3 reserved sample rate")
	}

	return nil
}

// Ogg container: validate page header, then identify Vorbis or Opus payload
func checkOgg(buf []byte) error {
	if len(buf) < 28 {
		return fmt.Errorf("ogg page header truncated")
	}

	if buf[4] != 0x00 {
		return fmt.Errorf("ogg stream version %d, expected 0", buf[4])
	}

	if buf[5]&0x02 == 0 {
		return fmt.Errorf("ogg first page missing BOS flag")
	}

	numSegments := int(buf[26])
	segTableEnd := 27 + numSegments
	if len(buf) < segTableEnd {
		return fmt.Errorf("ogg segment table truncated")
	}

	remaining := buf[segTableEnd:]

	// Vorbis: payload starts with 01 "vorbis"
	if len(remaining) >= 7 &&
		remaining[0] == 0x01 &&
		remaining[1] == 'v' && remaining[2] == 'o' && remaining[3] == 'r' &&
		remaining[4] == 'b' && remaining[5] == 'i' && remaining[6] == 's' {
		return checkVorbis(remaining[7:])
	}

	// Opus: payload starts with "OpusHead"
	if len(remaining) >= 8 &&
		remaining[0] == 'O' && remaining[1] == 'p' && remaining[2] == 'u' && remaining[3] == 's' &&
		remaining[4] == 'H' && remaining[5] == 'e' && remaining[6] == 'a' && remaining[7] == 'd' {
		return checkOpus(remaining[8:])
	}

	return fmt.Errorf("ogg payload is neither Vorbis nor Opus")
}

func checkVorbis(buf []byte) error {
	// Need: version(4) + channels(1) + sample_rate(4) = 9 bytes
	if len(buf) < 9 {
		return fmt.Errorf("vorbis identification header truncated")
	}

	version := binary.LittleEndian.Uint32(buf[0:4])
	if version != 0 {
		return fmt.Errorf("vorbis version %d, expected 0", version)
	}

	channels := buf[4]
	if channels == 0 {
		return fmt.Errorf("vorbis channel count is 0")
	}

	sampleRate := binary.LittleEndian.Uint32(buf[5:9])
	if sampleRate == 0 {
		return fmt.Errorf("vorbis sample rate is 0")
	}

	return nil
}

func checkOpus(buf []byte) error {
	// Need: version(1) + channels(1) = 2 bytes
	if len(buf) < 2 {
		return fmt.Errorf("opus header truncated")
	}

	// Major version (upper nibble) must be 0
	if buf[0]>>4 != 0 {
		return fmt.Errorf("opus major version %d, expected 0", buf[0]>>4)
	}

	if buf[1] == 0 {
		return fmt.Errorf("opus channel count is 0")
	}

	return nil
}

// M4A/AAC: validate ftyp box structure
func checkM4A(buf []byte, fileSize int64) error {
	if len(buf) < 12 {
		return fmt.Errorf("m4a ftyp box truncated")
	}

	boxSize := int64(binary.BigEndian.Uint32(buf[0:4]))
	if boxSize < 16 {
		return fmt.Errorf("m4a ftyp box size %d < 16", boxSize)
	}
	if boxSize > fileSize {
		return fmt.Errorf("m4a ftyp box size %d > file size %d", boxSize, fileSize)
	}

	// Major brand (bytes 8-11) should be printable ASCII
	for i := 8; i < 12; i++ {
		if buf[i] < 0x20 || buf[i] > 0x7E {
			return fmt.Errorf("m4a major brand contains non-ASCII byte 0x%02X", buf[i])
		}
	}

	// Check next box header is valid if we have enough data
	if int64(len(buf)) > boxSize+7 {
		nextType := buf[boxSize+4 : boxSize+8]
		for _, b := range nextType {
			if b < 0x20 || b > 0x7E {
				return fmt.Errorf("m4a second box type contains non-ASCII byte 0x%02X", b)
			}
		}
	}

	return nil
}

// WAV: validate RIFF/WAVE header and fmt chunk
func checkWAV(buf []byte, fileSize int64) error {
	if len(buf) < 36 {
		return fmt.Errorf("wav header truncated")
	}

	// RIFF size should roughly match file size
	riffSize := int64(binary.LittleEndian.Uint32(buf[4:8]))
	if riffSize > fileSize {
		return fmt.Errorf("wav RIFF size %d > file size %d", riffSize, fileSize)
	}

	// fmt chunk at offset 12
	if buf[12] != 0x66 || buf[13] != 0x6D || buf[14] != 0x74 || buf[15] != 0x20 {
		return fmt.Errorf("wav missing fmt chunk at expected offset")
	}

	// Audio format
	audioFmt := binary.LittleEndian.Uint16(buf[20:22])
	switch audioFmt {
	case 1, 3, 6, 7, 0xFFFE: // PCM, float, A-law, mu-law, extensible
	default:
		return fmt.Errorf("wav unknown audio format %d", audioFmt)
	}

	// Channels
	channels := binary.LittleEndian.Uint16(buf[22:24])
	if channels == 0 || channels > 8 {
		return fmt.Errorf("wav channel count %d out of range 1-8", channels)
	}

	// Sample rate
	sampleRate := binary.LittleEndian.Uint32(buf[24:28])
	if sampleRate < 8000 || sampleRate > 384000 {
		return fmt.Errorf("wav sample rate %d out of range 8000-384000", sampleRate)
	}

	// Bits per sample
	bps := binary.LittleEndian.Uint16(buf[34:36])
	switch bps {
	case 8, 16, 24, 32:
	default:
		return fmt.Errorf("wav bits per sample %d not 8/16/24/32", bps)
	}

	return nil
}

package favicon

import (
	"bytes"
	"encoding/binary"
)

// https://en.wikipedia.org/wiki/BMP_file_format

type BitmapFileHeader struct {
	Signature [2]byte // The header field used to identify the BMP and DIB file is 0x42 0x4D in hexadecimal, same as BM in ASCII.
	FileSize  uint32  // The size of the BMP file in bytes
	_         uint16  // Reserved; actual value depends on the application that creates the image
	_         uint16  // Reserved; actual value depends on the application that creates the image
	Offset    uint32  // The offset, i.e. starting address, of the byte where the bitmap image data (pixel array) can be found.
}

// Windows BITMAPINFOHEADER - there are other possibilities, but this looks to be
// the most common possibility.
type DIBHeader struct {
	Size                 uint32 // the size of this header, in bytes (40)
	Width                uint32 // the bitmap width in pixels (signed integer)
	Height               uint32 // the bitmap height in pixels (signed integer)
	ColorPlanes          uint16 // the number of color planes (must be 1)
	BitsPerPixel         uint16 // the number of bits per pixel, which is the color depth of the image. Typical values are 1, 4, 8, 16, 24 and 32.
	Compression          uint32 // the compression method being used. See the next table for a list of possible values
	ImageSize            uint32 // the image size. This is the size of the raw bitmap data; a dummy 0 can be given for BI_RGB bitmaps.
	HorizontalResolution uint32 // the horizontal resolution of the image. (pixel per metre, signed integer)
	VerticalResolution   uint32 // the vertical resolution of the image. (pixel per metre, signed integer)
	ColorCount           uint32 // the number of colors in the color palette, or 0 to default to 2^n
	ImportantColorCount  uint32 // the number of important colors used, or 0 when every color is important; generally ignored
}

func isBMP(icoBytes []byte) (bool, error) {
	r := bytes.NewReader(icoBytes)
	r.Seek(0, 0)

	dibHeader := DIBHeader{}

	err := binary.Read(r, binary.LittleEndian, &dibHeader)
	if err != nil {
		return false, err
	}

	if dibHeader.Size == 40 && dibHeader.ColorPlanes == 1 {
		return true, nil
	}

	return false, nil
}

func appendBitmapFileHeader(entry *IconDirEntry, icoBytes []byte) (*bytes.Buffer, error) {
	r := bytes.NewReader(icoBytes)
	r.Seek(0, 0)

	bitmapFileHeader := BitmapFileHeader{
		Signature: [2]byte{'B', 'M'},
		FileSize:  14 + entry.Size, // header + current file size
		Offset:    14 + 40,         // header + dib
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, bitmapFileHeader); err != nil {
		return nil, err
	}

	outz := append(buf.Bytes(), icoBytes...)

	buf2 := new(bytes.Buffer)
	if err := binary.Write(buf2, binary.LittleEndian, outz); err != nil {
		return nil, err
	}

	return buf2, nil
}

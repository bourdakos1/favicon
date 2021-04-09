package favicon

import (
	"encoding/binary"
	"fmt"
	"image"
	_ "image/png" // WTF is this magic....
	"io/ioutil"
	"os"
	"path/filepath"
)

// https://en.wikipedia.org/wiki/ICO_(file_format)

type Favicon struct {
	file    *os.File
	iconDir IconDir
	entries []IconDirEntry
}

type IconDir struct {
	_          uint16 // Reserved. Must always be 0.
	Type       uint16 // Specifies image type: 1 for icon (.ICO) image, 2 for cursor (.CUR) image. Other values are invalid.
	ImageCount uint16 // Specifies number of images in the file.
}

type IconDirEntry struct {
	Width        byte   // Specifies image width in pixels. Can be any number between 0 and 255. Value 0 means image width is 256 pixels.
	Height       byte   // Specifies image height in pixels. Can be any number between 0 and 255. Value 0 means image height is 256 pixels.
	ColorCount   byte   // Specifies number of colors in the color palette. Should be 0 if the image does not use a color palette.
	_            byte   // Reserved. Should be 0.
	ColorPlanes  uint16 // Specifies color planes. Should be 0 or 1.
	BitsPerPixel uint16 // Specifies bits per pixel.
	Size         uint32 // Specifies the size of the image's data in bytes.
	Offset       uint32 // Specifies the offset of BMP or PNG data from the beginning of the ICO file
}

func Pack(files []string, output string) {
	iconDir := IconDir{
		0,
		1,
		uint16(len(files)),
	}

	err := os.MkdirAll(filepath.Dir(output), os.ModeDir|os.ModePerm)
	if err != nil {
		panic(err)
	}

	buf, err := os.Create(output)
	defer buf.Close()
	if err != nil {
		panic(err)
	}

	if err := binary.Write(buf, binary.LittleEndian, iconDir); err != nil {
		panic(err)
	}

	var offset = uint32(6 + (len(files) * 16))
	for _, path := range files {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			panic(err)
		}

		fi, err := file.Stat()
		if err != nil {
			panic(err)
		}

		image, _, err := image.DecodeConfig(file)

		size := uint32(fi.Size())

		iconDirEntry := IconDirEntry{
			byte(image.Width),
			byte(image.Height),
			0,
			0,
			0,
			0,
			size,
			offset,
		}
		if err := binary.Write(buf, binary.LittleEndian, iconDirEntry); err != nil {
			panic(err)
		}
		offset += size
	}

	for _, path := range files {
		// file, err := os.Open(path)
		// defer file.Close()
		file, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if err := binary.Write(buf, binary.LittleEndian, file); err != nil {
			panic(err)
		}
	}

	buf.Sync()
}

func New(file *os.File) (*Favicon, error) {
	favicon := Favicon{}
	favicon.file = file

	err := readIconDir(file, &favicon)
	if err != nil {
		return nil, err
	}

	err = readEntries(file, &favicon)
	if err != nil {
		return nil, err
	}

	return &favicon, err
}

func readIconDir(file *os.File, favicon *Favicon) error {
	err := binary.Read(file, binary.LittleEndian, &favicon.iconDir)
	if err != nil {
		return err
	}
	return nil
}

func readEntries(file *os.File, favicon *Favicon) error {
	for i := 0; i < int(favicon.iconDir.ImageCount); i++ {
		iconDirEntry := IconDirEntry{}
		err := binary.Read(file, binary.LittleEndian, &iconDirEntry)
		if err != nil {
			return err
		}
		favicon.entries = append(favicon.entries, iconDirEntry)
	}
	return nil
}

func (favicon Favicon) SaveAsPNGs(nameRoot string) error {
	for _, entry := range favicon.entries {
		favicon.file.Seek(int64(entry.Offset), 0)
		bytes := make([]byte, entry.Size)

		_, err := favicon.file.Read(bytes)
		if err != nil {
			return err
		}

		mightBeBMP, _ := isBMP(bytes)
		if mightBeBMP {
			fmt.Println("Warning: File might be BMP.")
		}

		exportName := fmt.Sprintf("%s_%dx%d.png", nameRoot, entry.Width, entry.Height)
		err = ioutil.WriteFile(exportName, bytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (favicon Favicon) SaveAsBMPs(nameRoot string) error {
	for _, entry := range favicon.entries {
		favicon.file.Seek(int64(entry.Offset), 0)
		bytes := make([]byte, entry.Size)

		_, err := favicon.file.Read(bytes)
		if err != nil {
			return err
		}

		bmpBytes, err := appendBitmapFileHeader(&entry, bytes)
		if err != nil {
			return err
		}

		exportName := fmt.Sprintf("%s_%dx%d.bmp", nameRoot, entry.Width, entry.Height)
		err = ioutil.WriteFile(exportName, bmpBytes.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

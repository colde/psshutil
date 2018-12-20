package widevine

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/colde/psshutil/fileHandling"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
)

func Parse(f *os.File, size int64) {
	dataSize, err := fileHandling.ReadFromFile(f, 4)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	sizeInt := int64(binary.BigEndian.Uint32(dataSize))

	buf, err := fileHandling.ReadFromFile(f, sizeInt)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	widevineHeader := &WidevinePsshData{}
	err = proto.Unmarshal(buf, widevineHeader)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	key_ids := widevineHeader.GetKeyId()

	fmt.Println("Widevine Content ID:", string(widevineHeader.GetContentId()))
	for _, key_id := range key_ids {
		fmt.Println("Widevine Key ID:", base64.StdEncoding.EncodeToString(key_id))
	}
	fmt.Println("Widevine provider ID:", string(widevineHeader.GetProvider()))
}

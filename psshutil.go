/*
psshutil implements function to manipulate and use pssh boxes in isobmff files
*/
package main

import (
	"os"
  "log"
  "encoding/binary"
  "github.com/nu7hatch/gouuid"
)

func main() {
	var totalSize int64

	f, e := os.Open("video.mp4")
	if e != nil {
		log.Fatalf(e.Error())
	}
	defer f.Close()

  fi ,err := f.Stat()
  if err != nil {
		log.Fatalln(err.Error())
	}
  totalSize = fi.Size()

  loopAtoms(f, totalSize, 0)
}

func parsePssh(f *os.File, box string, size int64) {
  log.Println("Parsing PSSH")

  // Full box header
  _, err := readFromFile(f, 4)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }

  systemID, err := readFromFile(f, 16)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }

  systemUUID, err := uuid.Parse(systemID)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }
  log.Println("Found SystemID: ", systemUUID)
  switch systemUUID.String() {
  case "edef8ba9-79d6-4ace-a3c8-27dcd51d21ed":
    log.Println("Found Widevine")
  case "9a04f079-9840-4286-ab92-e65be0885f95":
    log.Println("Found PlayReady")
  default:
    log.Println("Unable to detect DRM system")
  }
}

func loopAtoms(f *os.File, totalSize int64, offset int64) {
  var pos int64

  for totalSize > pos {
    size, box, err := readHeader(f)
    if err != nil {
      log.Fatalln(err.Error())
    }

    sizeInt := int64(binary.BigEndian.Uint32(size))

    if string(box) == "moov" {
      loopAtoms(f, sizeInt - 8, pos + 8)
      pos += sizeInt
    } else {
      //log.Println(size, string(box))
      if string(box) == "pssh" {
        parsePssh(f, string(box), sizeInt)
      }
      pos += sizeInt
      seek := pos + offset
      f.Seek(seek, 0)
    }
  }
  return
}

func readFromFile(f *os.File, size int64) ([]byte,  error) {
	buf := make([]byte, size)
	_, err := f.Read(buf)

	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}

	return buf, nil
}

func readHeader(f *os.File)  ([]byte,  []byte,  error){
  buf, err := readFromFile(f, 8)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, nil, err
	}


	return buf[0:4], buf[4:8], nil
}

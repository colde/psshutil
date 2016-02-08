/*
psshutil implements function to manipulate and use pssh boxes in isobmff files
*/
package main

import (
	"os"
  "log"
  "encoding/binary"
  "github.com/nu7hatch/gouuid"
  "github.com/colde/psshutil/widevine"
  "github.com/colde/psshutil/playready"
  "github.com/colde/psshutil/fileUtility"
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

  // Full box header
  _, err := fileUtility.ReadFromFile(f, 4)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }

  systemID, err := fileUtility.ReadFromFile(f, 16)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }

  systemUUID, err := uuid.Parse(systemID)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }
  switch systemUUID.String() {
  case "edef8ba9-79d6-4ace-a3c8-27dcd51d21ed":
    log.Println("Found Widevine", systemUUID)
    // Size determined to be size - 8 (box header), 4 (fullbox header), 16 (systemid)
    widevine.Parse(f, size-28)
  case "9a04f079-9840-4286-ab92-e65be0885f95":
    log.Println("Found PlayReady", systemUUID)
    // Size determined to be size - 8 (box header), 4 (fullbox header), 16 (systemid)
    playready.Parse(f, size-28)
  default:
    log.Println("Found unknown DRM system", systemUUID)
  }
}

func loopAtoms(f *os.File, totalSize int64, offset int64) {
  var pos int64

  for totalSize > pos {
    size, box, err := fileUtility.ReadHeader(f)
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

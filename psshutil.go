/*
psshutil implements function to manipulate and use pssh boxes in isobmff files
*/
package main

import (
	"os"
  "log"
  "fmt"
  "flag"
  "encoding/binary"
  "github.com/nu7hatch/gouuid"
  "github.com/colde/psshutil/widevine"
  "github.com/colde/psshutil/playready"
  "github.com/colde/psshutil/fileHandling"
)

func main() {
  var fileName = flag.String("i", "", "Input file for reading/parsing")
  flag.Parse()

  if *fileName == "" {
    fmt.Println("Usage: psshutil -i <video.mp4>")
    os.Exit(0)
  }

	var totalSize int64

	f, e := os.Open(*fileName)
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
  _, err := fileHandling.ReadFromFile(f, 4)
  if err != nil {
    log.Fatalln(err.Error())
    return
  }

  systemID, err := fileHandling.ReadFromFile(f, 16)
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
    fmt.Println("Found Widevine", systemUUID)
    // Size determined to be size - 8 (box header), 4 (fullbox header), 16 (systemid)
    widevine.Parse(f, size-28)
  case "9a04f079-9840-4286-ab92-e65be0885f95":
    fmt.Println("Found Microsoft PlayReady", systemUUID)
    // Size determined to be size - 8 (box header), 4 (fullbox header), 16 (systemid)
    playready.Parse(f, size-28)
	case "f239e769-efa3-4850-9c16-a903c6932efb":
		fmt.Println("Found Adobe Primetime DRM, version 4", systemUUID)
	case "5e629af5-38da-4063-8977-97ffbd9902d4":
		fmt.Println("Found Marlin DRM")
	case "adb41c24-2dbf-4a6d-958b-4457c0d27b95":
		fmt.Println("Found Nagra MediaAccess PRM 3.0")
	case "a68129d3-575b-4f1a-9cba-3223846cf7c3":
		fmt.Println("Cisco/NDS VideoGuard Everywhere DRM")
	case "9a27dd82-fde2-4725-8cbc-4234aa06ec09":
		fmt.Println("Found Verimatrix VCAS")
	case "1f83e1e8-6ee9-4f0d-ba2f-5ec4e3ed1a66":
		fmt.Println("Found Arris SecureMedia")
	case "644fe7b5-260f-4fad-949a-0762ffb054b4":
		fmt.Println("Found CMLA (OMA DRM)")
	case "6a99532d-869f-5922-9a91-113ab7b1e2f3":
		fmt.Println("Found MobiTV DRM")
	case "35bf197b-530e-42d7-8b65-1b4bf415070f":
		fmt.Println("Found DivX DRM Series 5")
	case "b4413586-c58c-ffb0-94a5-d4896c1af6c3":
		fmt.Println("Found Viaccess-Orca DRM")
	case "80a6be7e-1448-4c37-9e70-d5aebe04c8d2":
		fmt.Println("Found Irdeto Content Protection for DASH")
	case "dcf4e3e3-62f1-5818-7ba6-0a6fe33ff3dd":
		fmt.Println("Found DigiCAP SmartXess for DASH")
  default:
    fmt.Println("Found unknown DRM system", systemUUID)
  }
	fmt.Println()
}

func loopAtoms(f *os.File, totalSize int64, offset int64) {
  var pos int64

  for totalSize > pos {
    size, box, err := fileHandling.ReadHeader(f)
    if err != nil {
      log.Fatalln(err.Error())
    }

    sizeInt := int64(binary.BigEndian.Uint32(size))

    if string(box) == "moov" {
      loopAtoms(f, sizeInt - 8, pos + 8)
      pos += sizeInt
    } else {
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

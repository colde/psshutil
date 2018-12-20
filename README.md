# psshutil

A simple library for parsing PSSH headers from **ISOBMFF** boxes

Status
------
Parses basic information for PlayReady and Widevine, including KeyIDs.
Will signal presence for a number of DRM systems.

Build and install
-----------------
Make sure you have a working **go** environment

	go get github.com/colde/psshutil
	go install github.com/colde/psshutil

Usage
-----
  psshutil -i video.mp4

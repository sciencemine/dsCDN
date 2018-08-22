package server

/*
 * The controller is responsible for handling the main handler functions for the REST endpoints
 * It deals with query parameters and setting the correct headers and sending back the data to
 * the client. Files starting with the prefix con_ are controller files that have been refactored
 * out for readability
 */

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/textproto"
	"strconv"
)

//Controller ...
type Controller struct {
	Repository Repository
}

//writeAssetToPart writes an asset to a mime multipart body
func writeAssetToPart(assetObj Asset, writer *multipart.Writer) error {
	//Send a request to get the asset file
	url := assetObj.URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Could not get the asset from the url: ", err)
		return err
	}
	defer resp.Body.Close()

	//Copy the asset to a byte buffer to find the size of the asset
	buf := &bytes.Buffer{}
	nRead, err := io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Println("Could not copy the responsBody to the buffer: ", err)
		return err
	}

	//Create the MIMEHeader and create a new part
	h := make(textproto.MIMEHeader)
	h.Set("Content-Type", assetObj.MIMEType)
	h.Set("Content-Length", strconv.FormatInt(nRead, 10)) //len(data)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	h.Set("Content-Name", assetObj.ID.Hex())
	fw, err := writer.CreatePart(h)
	if err != nil {
		log.Fatalln("Could not create form field: ", err)
		return err
	}

	//Create the encoding writer and write the data to it.
	w := quotedprintable.NewWriter(fw)
	w.Binary = true //Tell it to treat the data as binary
	w.Write(buf.Bytes())

	return nil
}

//writeDescToParts is a wrapper function for write asset to parts that will loop
//through a description array and add all the assets to the multipart message
func writeDescToParts(arr []desc, writer *multipart.Writer, seen map[string]bool) {
	for _, item := range arr {
		if assetObj, ok := item.Asset.(Asset); ok {
			id := assetObj.ID.Hex()
			if vis, _ := seen[id]; !vis {
				if err := writeAssetToPart(assetObj, writer); err != nil {
					continue
				}
				seen[id] = true
			}
		}
	}
}

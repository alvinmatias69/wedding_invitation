package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

func Start() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("GET /", fs)

	http.HandleFunc("GET /image.jpg", func(w http.ResponseWriter, r *http.Request) {
		if authorization := r.Header.Get("Authorization"); len(authorization) == 0 {
			w.Header().Add("wWW-authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("salah"))
			return
		}

		mediaCtx, _ := jpegstructure.NewJpegMediaParser().ParseFile("./static/image.jpg")
		segmentList := mediaCtx.(*jpegstructure.SegmentList)
		builder, _ := segmentList.ConstructExifBuilder()
		ifd0Ib, _ := exif.GetOrCreateIbFromRootIb(builder, "IFD0")
		_ = ifd0Ib.SetStandardWithName("Artist", "https://medium.com/@kleinc./modifying-exif-metadata-with-go-268c22bf654e")
		_ = ifd0Ib.SetStandardWithName("XPKeywords", "value")
		segmentList.SetExif(builder)

		w.WriteHeader(http.StatusOK)
		segmentList.Write(w)
	})

	http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	http.HandleFunc("GET /kontol", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("kontol"))
	})

	fmt.Printf("starting server in port: %v\n", "8080")
	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}

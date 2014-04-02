package fpm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/evanphx/fpm/go/fpm/ar"
)

type DebPackage struct {
	Data *MetaData
	Path string
}

type DebPackageWriter struct {
	controlBuf bytes.Buffer
	controlTar *tar.Writer

	dataBuf bytes.Buffer
	dataTar *tar.Writer
}

func (d *DebPackageWriter) AddControl(name string, data []byte) error {
	hdr := &tar.Header{
		Name: name,
		Size: int64(len(data)),
	}

	if err := d.controlTar.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := d.controlTar.Write(data); err != nil {
		return err
	}

	return nil
}

func (d *DebPackageWriter) AddData(f File) error {
	hdr := &tar.Header{
		Name: f.Name(),
		Size: f.Size(),
	}

	if err := d.dataTar.WriteHeader(hdr); err != nil {
		return err
	}

	r, err := f.Data()
	if err != nil {
		return err
	}

	if _, err := io.Copy(d.dataTar, r); err != nil {
		return err
	}

	r.Close()

	return nil
}

var debMagic = []byte("2.0\n")

func (d *DebPackage) Write(input *DirPackage) error {
	tw := &DebPackageWriter{}
	tw.controlTar = tar.NewWriter(&tw.controlBuf)
	tw.dataTar = tar.NewWriter(&tw.dataBuf)

	files := input.Files()

	if err := tw.AddControl("debian-binary", debMagic); err != nil {
		return err
	}

	if err := tw.AddControl("control", d.control(files)); err != nil {
		return err
	}

	tw.controlTar.Close()

	for _, file := range files {
		if err := tw.AddData(file); err != nil {
			return err
		}
	}

	tw.dataTar.Close()

	out, err := os.OpenFile(d.Path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	archive := ar.NewWriter(out)
	archive.WriteGlobalHeader()

	buf := new(bytes.Buffer)

	gw := gzip.NewWriter(buf)
	_, err = io.Copy(gw, bytes.NewReader(tw.controlBuf.Bytes()))
	if err != nil {
		return err
	}

	gw.Close()

	header := new(ar.Header)
	header.ModTime = time.Now()
	header.Mode = 0644
	header.Name = "control.tar.gz"
	header.Size = int64(buf.Len())

	if err := archive.WriteHeader(header); err != nil {
		return err
	}

	archive.Write(buf.Bytes())

	buf = new(bytes.Buffer)

	gw = gzip.NewWriter(buf)
	if err != nil {
		return err
	}

	io.Copy(gw, bytes.NewReader(tw.dataBuf.Bytes()))
	gw.Close()

	header = new(ar.Header)
	header.ModTime = time.Now()
	header.Mode = 0644
	header.Name = "data.tar.gz"
	header.Size = int64(buf.Len())

	if err := archive.WriteHeader(header); err != nil {
		return err
	}

	archive.Write(buf.Bytes())

	return out.Close()
}

func (d *DebPackage) control(files []File) []byte {
	total := int64(0)

	for _, f := range files {
		total += f.Size()
	}

	d.Data.Extra["installed_size"] = total

	var buf bytes.Buffer

	err := debTemplate.Execute(&buf, d.Data)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

var debTemplate *template.Template

func init() {
	t, err := template.New("deb_control").Parse(debTemplateData)
	if err != nil {
		panic(err)
	}

	debTemplate = t
}

const debTemplateData = `
Package: {{ .Name }}
Version: {{ .Version }}
License: {{ .License }}
Vendor: {{ .Vendor }}
Architecture: {{ .Architecture }}
Maintainer: {{ .Maintainer  }}
Installed-Size: {{ .Extra.installed_size }}
Section: {{ .Section }}
Priority: {{ .Extra.priority }}
Homepage: {{ .URL }}
Description: {{ .Description }}
`

package promexport

import (
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
)

type Result struct {
	Hostname string
	Zone     string
	Status   string
	Qtime    int
}

type Promwriter struct {
	Wr       *os.File
	Hostname string
	Filename string
}

func NewPromwriter(fname string) (*Promwriter, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	p := Promwriter{Hostname: hostname, Filename: fname}
	p.Wr, err = os.Create(fname + ".tmp")
	if err != nil {
		return nil, err
	}
	p.writeheader()
	if err != nil {
		return nil, err
	}
	// r := Result{Hostname: hostname, Zone: "co.za", Status: "OK", Qtime: 360}
	return &p, nil

}

func (p *Promwriter) writeheader() error {
	_, err := p.Wr.Write([]byte("# HELP epp_check Connection time to get an epp header\n"))
	_, err = p.Wr.Write([]byte("# TYPE epp_check gauge\n"))
	return err
}

// Closes the filehandle and renames to final location to achieve the atomic creation needed for TextExporter
func (p *Promwriter) Close() error {
	p.Wr.Close()
	return os.Rename(p.Filename+".tmp", p.Filename)
}

func (p *Promwriter) Writeresult(zone string, status string, querytime int) error {
	res := Result{Hostname: p.Hostname, Zone: zone, Status: status, Qtime: querytime}
	tmpl := "epp_check{zone=\"{{.Zone}}\", status=\"{{.Status}}\", hostname=\"{{.Hostname}}\"} {{.Qtime}}\n"
	t, err := template.New("prometheus").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	logrus.Info("Write result", res)
	return t.Execute(p.Wr, res)
}

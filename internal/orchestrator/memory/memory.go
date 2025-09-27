package memory

import "fmt"

type DocMemory struct {
	mem map[string]string
}

func New() *DocMemory {
	return &DocMemory{mem: make(map[string]string)}
}

func (d *DocMemory) Write(docName, content string) {
	d.mem[docName] = content
}

func (d *DocMemory) Read(docName string) (string, error) {
	doc, ok := d.mem[docName]
	if !ok {
		return "", fmt.Errorf("%s does not exist in memory", docName)
	}
	return doc, nil
}

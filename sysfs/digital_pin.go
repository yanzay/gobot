package sysfs

import (
	"fmt"
	"os"
	"strconv"
)

const (
	IN       = "in"
	OUT      = "out"
	HIGH     = 1
	LOW      = 0
	GPIOPATH = "/sys/class/gpio"
)

type DigitalPin interface {
	Unexport() error
	Export() error
	Read() (int, error)
	Direction(string) error
	Write(int) error
}

type digitalPin struct {
	pin   string
	label string
}

// NewDigitalPin returns a DigitalPin given the pin number and an optional sysfs pin label.
// If no label is supplied the default label will prepend "gpio" to the pin number,
// eg. a pin number of 10 will have a label of "gpio10"
func NewDigitalPin(pin int, v ...string) DigitalPin {
	d := &digitalPin{pin: strconv.Itoa(pin)}
	if len(v) > 0 {
		d.label = v[0]
	} else {
		d.label = "gpio" + d.pin
	}

	return d
}

// Direction sets the direction for the pin
func (d *digitalPin) Direction(dir string) error {
	_, err := writeFile(fmt.Sprintf("%v/%v/direction", GPIOPATH, d.label), []byte(dir))
	return err
}

// Write writes to the pin
func (d *digitalPin) Write(b int) error {
	_, err := writeFile(fmt.Sprintf("%v/%v/value", GPIOPATH, d.label), []byte(strconv.Itoa(b)))
	return err
}

// Read reads the current value of the pin
func (d *digitalPin) Read() (n int, err error) {
	buf, err := readFile(fmt.Sprintf("%v/%v/value", GPIOPATH, d.label))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(buf[0]))
}

// Export exports the pin for use by the operating system
func (d *digitalPin) Export() error {
	_, err := writeFile(GPIOPATH+"/export", []byte(d.pin))
	return err
}

// Unexport unexports the pin and releases the pin from the operating system
func (d *digitalPin) Unexport() error {
	_, err := writeFile(GPIOPATH+"/unexport", []byte(d.pin))
	return err
}

func writeFile(path string, data []byte) (i int, err error) {
	file, err := OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 2)
	_, err = file.Read(buf)
	return buf, err
}

package mtsWsdl

import (
	"encoding/xml"
	"fmt"
	"time"
)

const formatWithoutZone = "2006-01-02T15:04:05.9999999"

var UtcOffset = time.Hour * -4 // because we use Izhevsk time zone for MTS

type TimeWithoutZone struct {
	time.Time
}

func (c *TimeWithoutZone) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		parse, err = time.Parse(formatWithoutZone, v)
		if err != nil {
			return err
		}
		parse = parse.Add(UtcOffset)
	}
	*c = TimeWithoutZone{parse}
	return nil
}

func (c *TimeWithoutZone) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	result := c.UTC().Add(-UtcOffset)
	dateString := fmt.Sprintf("%s", result.Format(formatWithoutZone))
	return e.EncodeElement(dateString, start)
}

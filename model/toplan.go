package model

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func ToPlan(rawdata []byte) (*Plan, error) {
	data, err := charmap.ISO8859_1.NewDecoder().Bytes(rawdata)
	if err != nil {
		data = rawData
		log.Printf("Error decoding plan input to utf-8: %v\n", err)
	}

	rawplan := string(data)

	headStart := strings(rawplan, "<head>")
	headEnd := strings(rawplan, "</head>") + 7
	if headStart == -1 || headEnd-7 == -1 {
		log.Printf("No head found in raw plan.\n")
	} else {
		rawplan = append(rawplan[:headStart], rawplan[headEnd:])
	}

	re := regexp.MustCompile("(</?html>)|(</?body>)|(</?p>)|(</?br>)|(</?CENTER>)")
	rawplan = re.ReplaceAllString(rawplan, "")

	decoder := xml.NewDecoder(strings.NewReader(rawplan))
	if err = moveToNext("font", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"font\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	token, err := decoder.RawToken()
	charData, _ := token.(xml.CharData)
	createdString := string(charData[7:])
	loc, _ := time.LoadLocation("Europe/Berlin")
	created := time.ParseInLocation("02.06.2006 15:04", createdString, loc)

	if err = moveToNext("div", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for \"div\" with day of first part: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	token, err = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	dayString := string(charData)
	dayFormat := "2.6.2006"
	day1 := time.ParseInLocation(dayFormat, strings.Split(" ")[0], loc)

	if err = moveToNext("table", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for second \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for third \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"tr\" of the first plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for second \"tr\" of the first plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	firstPartSubstitutions := make([]Substitution, 0, 20)
	token, err = decoder.RawToken()
	_, ok := token.(xml.StartElement)
	for ok && err == nil {
		append(firstPartSubstitutions, readSubstitution(decoder))
		token, err = decoder.RawToken()
		_, ok = token.(xml.StartElement)
	}

	if err = moveToNext("div", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for \"div\" with day of second part: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	token, err = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	dayString := string(charData)
	dayFormat := "2.6.2006"
	day2 := time.ParseInLocation(dayFormat, strings.Split(" ")[0], loc)

	if err = moveToNext("table", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for fourth \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for fifth \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for sixth \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"tr\" of the second plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for second \"tr\" of the second plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	secondPartSubstitutions := make([]Substitution, 0, 20)
	token, err = decoder.RawToken()
	_, ok = token.(xml.StartElement)
	for ok && err == nil {
		append(secondPartSubstitutions, readSubstitution(decoder))
		token, err = decoder.RawToken()
		_, ok = token.(xml.StartElement)
	}

	parts := []Part{
		Part{Day: day1, Substitutions: firstPartSubstitutions},
		Part{Day: day2, Substitutions: secondPartSubstitutions}}
	return &Plan{
		Created: created,
		Parts:   parts}, nil
}

func readSubstitution(decoder *xml.Decoder) Substitution {
	_, _ := decoder.RawToken()
	token, _ := decoder.RawToken()
	charData, _ := token.(xml.CharData)
	class := string(charData)
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	token, _ = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	periodString := string(charData)
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	token, _ = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	substTeacherShort := string(charData)
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	token, _ = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	instdTeacherShort := string(charData)
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	token, _ = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	instdSubjectShort := string(charData)
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	token, _ = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	kind := string(charData)
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	_, _ = decoder.RawToken()
	_, _ = decoder.RawToken()

	_, _ = decoder.RawToken()
	token, _ = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	text := string(charData)

	return Substitution{
		Class:        class,
		Period:       periodString,
		SubstTeacher: Teacher{Short: substTeacherShort},
		InstdTeacher: Teacher{Short: instdTeacherShort},
		InstdSubject: Teacher{Short: instdSubjectShort},
		Kind:         kind,
		Text:         text}
}

func moveToNext(elementName string, decoder *xml.Decoder) error {
	for {
		token, err := decoder.RawToken()
		if err != nil {
			return err
		}
		if startElement, ok := token.(xml.StartElement); ok {
			if startElement.Name.Lokal == elementName {
				return nil
			}
		}
	}
}

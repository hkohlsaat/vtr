package model

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func ToPlan(uploadReader io.Reader) (*Plan, error) {
	plan, err := decodePlan(uploadReader)
	if err != nil {
		return plan, err
	}

	refine(plan)
	return plan, nil
}

func decodePlan(uploadReader io.Reader) (*Plan, error) {
	encReader := charmap.ISO8859_1.NewDecoder().Reader(uploadReader)

	decoder := xml.NewDecoder(encReader)
	decoder.Entity = xml.HTMLEntity

	if err := moveToNext("font", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"font\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	token, err := decoder.RawToken()
	charData, _ := token.(xml.CharData)
	createdString := string(charData[7:])
	loc, _ := time.LoadLocation("Europe/Berlin")
	created, _ := time.ParseInLocation("02.01.2006 15:04", createdString, loc)

	if err = moveToNext("div", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for \"div\" with day of first part: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	token, err = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	dayString := string(charData)
	dayFormat := "2.1.2006"
	day1, _ := time.ParseInLocation(dayFormat, strings.Split(dayString, " ")[0], loc)

	if err = moveToNext("table", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for second \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for third \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"tr\" of the first plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for second \"tr\" of the first plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", false, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for end of second \"tr\" of the first plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}

	decoder.RawToken() // Skip xml.CharData.

	firstPartSubstitutions := make([]Substitution, 0, 20)
	token, err = decoder.RawToken()
	_, ok := token.(xml.StartElement)
	for ok && err == nil {
		firstPartSubstitutions = append(firstPartSubstitutions, readSubstitution(decoder))
		decoder.RawToken()
		token, err = decoder.RawToken()
		_, ok = token.(xml.StartElement)
	}

	if err = moveToNext("div", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for \"div\" with day of second part: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	token, err = decoder.RawToken()
	charData, _ = token.(xml.CharData)
	dayString = string(charData)
	day2, _ := time.ParseInLocation(dayFormat, strings.Split(dayString, " ")[0], loc)

	if err = moveToNext("table", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for fourth \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for fifth \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("table", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for sixth \"table\": %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for first \"tr\" of the second plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", true, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for second \"tr\" of the second plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}
	if err = moveToNext("tr", false, decoder); err != nil {
		err = errors.New(fmt.Sprintf("Error searching for end of second \"tr\" of the second plan table: %v\n", err))
		log.Println(err)
		return &Plan{}, err
	}

	decoder.RawToken() // Skip xml.CharData.

	secondPartSubstitutions := make([]Substitution, 0, 20)
	token, err = decoder.RawToken()
	_, ok = token.(xml.StartElement)
	for ok && err == nil {
		secondPartSubstitutions = append(secondPartSubstitutions, readSubstitution(decoder))
		decoder.RawToken()
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
	var (
		class             = readInformation(decoder) // Read class.
		periodString      = readInformation(decoder) // Read period(s).
		substTeacherShort = readInformation(decoder) // Read substitution teacher.
		instdTeacherShort = readInformation(decoder) // Read instead teacher.
		instdSubjectShort = readInformation(decoder) // Read instead subject.
		kind              = readInformation(decoder) // Read kind.
		_                 = readInformation(decoder) // Skip "Vtr. von".
		text              = readInformation(decoder) // Read text.
	)
	// Close table row.
	decoder.RawToken()

	return Substitution{
		Class:        class,
		Period:       periodString,
		SubstTeacher: Teacher{Short: substTeacherShort},
		InstdTeacher: Teacher{Short: instdTeacherShort},
		InstdSubject: Subject{Short: instdSubjectShort},
		Kind:         kind,
		Text:         text}
}

func readInformation(decoder *xml.Decoder) string {
	decoder.RawToken()
	token, _ := decoder.RawToken()
	if startElement, ok := token.(xml.StartElement); ok && startElement.Name.Local == "span" {
		token, _ = decoder.RawToken()
		defer decoder.RawToken()
	}
	charData, _ := token.(xml.CharData)
	information := string(charData)
	decoder.RawToken()
	return information
}

func moveToNext(elementName string, se bool, decoder *xml.Decoder) error {
	for {
		token, err := decoder.RawToken()
		if err != nil {
			return err
		}
		if se {
			if startElement, ok := token.(xml.StartElement); ok {
				if startElement.Name.Local == elementName {
					return nil
				}
			}
		} else {
			if endElement, ok := token.(xml.EndElement); ok {
				if endElement.Name.Local == elementName {
					return nil
				}
			}
		}
	}
}

func refine(plan *Plan) {
	const nbsp = "\u00A0"
	for p, part := range plan.Parts {
		for s, substitution := range part.Substitutions {
			if substitution.Period == nbsp {
				plan.Parts[p].Substitutions[s].Period = ""
			}
			if substitution.Class == nbsp {
				plan.Parts[p].Substitutions[s].Class = ""
			}
			if substitution.SubstTeacher.Short == nbsp || substitution.SubstTeacher.Short == "???" ||
				substitution.SubstTeacher.Short == "+" || substitution.SubstTeacher.Short == "---" {
				plan.Parts[p].Substitutions[s].SubstTeacher.Short = ""
			} else {
				plan.Parts[p].Substitutions[s].SubstTeacher.Read()
			}
			if substitution.InstdTeacher.Short == nbsp {
				plan.Parts[p].Substitutions[s].InstdTeacher.Short = ""
			} else {
				plan.Parts[p].Substitutions[s].InstdTeacher.Read()
			}
			if substitution.InstdSubject.Short == nbsp {
				plan.Parts[p].Substitutions[s].InstdSubject.Short = ""
			} else {
				plan.Parts[p].Substitutions[s].InstdSubject.Read()
			}
			if substitution.Kind == nbsp {
				plan.Parts[p].Substitutions[s].Kind = ""
			} else if substitution.Kind == "Statt-Vertretung" {
				plan.Parts[p].Substitutions[s].Kind = "Vertretung"
			}
			if substitution.Text == nbsp {
				plan.Parts[p].Substitutions[s].Text = ""
			} else {
				re := regexp.MustCompile("Aufg(\\.|(abe)) [A-Za-z]{2,3}")
				task := re.FindString(substitution.Text)
				if task != "" {
					provider := strings.Trim(task[len(task)-3:], " ")
					taskProvider := Teacher{Short: provider}
					taskProvider.Read()
					plan.Parts[p].Substitutions[s].TaskProvider = taskProvider
				}
			}
		}
	}
}

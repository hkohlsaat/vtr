package model

type Subject struct {
	Short      string `gorm:"unique_index"`
	Name       string
	SplitClass bool
}

func (subject *Subject) Create() {
	if db.NewRecord(*subject) {
		db.Create(subject)
	}
}

func ReadAllSubjects() []Subject {
	var subjects []Subject
	db.Find(&subjects)
	return subjects
}

func (subject *Subject) Read() {
	var readSubject Subject
	db.Where(subject).First(&readSubject)
	if readSubject.Short == subject.Short {
		*subject = readSubject
	}
}

func (subject *Subject) Exists() bool {
	subject.Read()
	if subject.Name == "" {
		return false
	}
	return true
}

func (subject *Subject) Update() {
	subject.UpdateShort(subject.Short)
}
func (subject *Subject) UpdateShort(short string) {
	if !subject.SplitClass {
		subject.Read()
		db.Model(subject).Update("splitClass", "false")
	}
	db.Table("subjects").Where(&Subject{Short: short}).Update(*subject)
}

func (subject *Subject) Delete() {
	db.Where(subject).Delete(Subject{})
}

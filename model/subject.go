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

func (subject *Subject) Read() {
	var readSubject Subject
	db.Where(subject).First(&readSubject)
	if readSubject.Short == subject.Short {
		*subject = readSubject
	}
}

func (subject *Subject) Update() {
	db.Table("subjects").Where(&Subject{Short: subject.Short}).Update(*subject)
}

func (subject *Subject) Delete() {
	db.Where(subject).Delete(Subject{})
}

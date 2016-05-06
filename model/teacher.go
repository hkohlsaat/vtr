package model

type Teacher struct {
	Short string `gorm:"unique_index"`
	Name  string
	Sex   rune
}

func (teacher *Teacher) Create() {
	if db.NewRecord(*teacher) {
		db.Create(teacher)
	}
}

func ReadAllTeachers() []Teacher {
	var teachers []Teacher
	db.Find(&teachers)
	return teachers
}

func (teacher *Teacher) Read() {
	var readTeacher Teacher
	db.Where(teacher).First(&readTeacher)
	if readTeacher.Short == teacher.Short {
		*teacher = readTeacher
	}
}

func (teacher *Teacher) Exists() bool {
	teacher.Read()
	if teacher.Name == "" {
		return false
	}
	return true
}

func (teacher *Teacher) Update() {
	db.Table("teachers").Where(&Teacher{Short: teacher.Short}).Update(*teacher)
}
func (teacher *Teacher) UpdateShort(short string) {
	db.Table("teachers").Where(&Teacher{Short: short}).Update(*teacher)
}

func (teacher *Teacher) Delete() {
	db.Where(teacher).Delete(Teacher{})
}

package models


// AssignmentType defines types like HW, Exam
type LocalAssignmentType struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	Color    string `gorm:"not null"`
	NotionID string
}

/*func Get_Local_AssignmentType_byName(name string, db *gorm.DB) *AssignmentType {
	assignmentType := &AssignmentType{}
	err := db.Where("name = ?", name).First(assignmentType).Error
	if err != nil {
		log.Fatalln("Error getting assignment type with name: ", err)
		return nil
	}
	return assignmentType
}*/

func (a *LocalAssignmentType) ToMap() map[string]string {
	return map[string]string{
		"id":		a.NotionID,
		"name":		a.Name,
		"color":	a.Color,
	}
}

// AssignmentStatus defines statuses like Not Started, In Progress
type LocalAssignmentStatus struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	Color    string `gorm:"not null"`
	NotionID string
}

/*func Get_AssignmentStatus_byName(name string, db *gorm.DB) *AssignmentStatus {
	assignmentStatus := &AssignmentStatus{}
	err := db.Where("name = ?", name).First(assignmentStatus).Error
	if err != nil {
		log.Fatalln("Error getting assignment status with name: ", err)
		return nil
	}
	return assignmentStatus
}*/

func (a *LocalAssignmentStatus) ToMap() map[string]string {
	return map[string]string{
		"id":		a.NotionID,
		"name":		a.Name,
		"color":	a.Color,
	}
}

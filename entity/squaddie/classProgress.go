package squaddie

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility"
)

// ClassProgress tracks the ClassProgress's current class and any levels they have taken so far.
type ClassProgress struct {
	BaseClassID         string                          `json:"base_class" yaml:"base_class"`
	CurrentClassID      string                          `json:"current_class" yaml:"current_class"`
	ClassLevelsConsumed map[string]*ClassLevelsConsumed `json:"class_levels" yaml:"class_levels"`
}

// AddClass gives the ClassProgress a new class it can gain levels in, if it wasn't already added.
func (classProgress *ClassProgress) AddClass(classReference *squaddieclass.ClassReference) {
	if classProgress.ClassLevelsConsumed[classReference.ID] != nil {
		return
	}

	classProgress.ClassLevelsConsumed[classReference.ID] = &ClassLevelsConsumed{
		ClassID:        classReference.ID,
		ClassName:      classReference.Name,
		LevelsConsumed: []string{},
	}
}

// GetLevelCountsByClass returns a mapping of class names to levels gained.
func (classProgress *ClassProgress) GetLevelCountsByClass() map[string]int {
	count := map[string]int{}
	for classID, progress := range classProgress.ClassLevelsConsumed {
		count[classID] = len(progress.LevelsConsumed)
	}

	return count
}

// MarkLevelUpBenefitAsConsumed makes the ClassProgress remember it used this benefit to level up already.
func (classProgress *ClassProgress) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string) {
	classProgress.ClassLevelsConsumed[benefitClassID].LevelsConsumed = append(classProgress.ClassLevelsConsumed[benefitClassID].LevelsConsumed, benefitID)
}

// SetClass changes the ClassProgress's CurrentClassID to the given classID.
//   It also sets the BaseClass if it hasn't been already.
//   Raises an error if classID has not been added to the squaddie yet.
func (classProgress *ClassProgress) SetClass(classID string) error {
	if _, exists := classProgress.ClassLevelsConsumed[classID]; !exists {
		newError := fmt.Errorf(`cannot switch to unknown class "%s"`, classID)
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	if classProgress.BaseClassID == "" {
		classProgress.BaseClassID = classID
	}

	classProgress.CurrentClassID = classID
	return nil
}

// SetBaseClassIfNoBaseClass sets the BaseClass if it hasn't been already.
func (classProgress *ClassProgress) SetBaseClassIfNoBaseClass(classID string) {
	if classProgress.BaseClassID == "" {
		classProgress.BaseClassID = classID
	}
}

// IsClassLevelAlreadyUsed returns true if a LevelUpBenefit with the given SquaddieID has already been used.
func (classProgress *ClassProgress) IsClassLevelAlreadyUsed(benefitID string) bool {
	return classProgress.anyClassLevelsConsumed(func(classID string, progress *ClassLevelsConsumed) bool {
		return progress.IsLevelAlreadyConsumed(benefitID)
	})
}

// HasAddedClass returns true if the ClassProgress has already added a class with the name classIDToFind
func (classProgress *ClassProgress) HasAddedClass(classIDToFind string) bool {
	return classProgress.anyClassLevelsConsumed(func(classID string, progress *ClassLevelsConsumed) bool {
		return classID == classIDToFind
	})
}

// anyClassLevelsConsumed returns true if any of the squaddie's class levels consumed satisfies a given condition.
func (classProgress *ClassProgress) anyClassLevelsConsumed(condition func(classID string, progress *ClassLevelsConsumed) bool) bool {
	for classID, progress := range classProgress.ClassLevelsConsumed {
		if condition(classID, progress) {
			return true
		}
	}
	return false
}

// GetBaseClassID returns the base class ID.
func (classProgress *ClassProgress) GetBaseClassID() string {
	return classProgress.BaseClassID
}

// GetCurrentClassID returns the current class ID.
func (classProgress *ClassProgress) GetCurrentClassID() string {
	return classProgress.CurrentClassID
}

// GetClassLevelsConsumed returns the current class ID.
func (classProgress *ClassProgress) GetClassLevelsConsumed() *map[string]*ClassLevelsConsumed  {
	return &classProgress.ClassLevelsConsumed
}

// HasSameClassesAs sees if the other class progress has the same fields.
func (classProgress *ClassProgress) HasSameClassesAs(other *ClassProgress) bool {
	if classProgress.GetBaseClassID() != other.GetBaseClassID() { return false }
	if classProgress.GetCurrentClassID() != other.GetCurrentClassID() { return false }

	otherClassLevelsConsumed := *other.GetClassLevelsConsumed()
	if len(*classProgress.GetClassLevelsConsumed()) != len(otherClassLevelsConsumed) { return false }

	classLevelsConsumedByClassID := map[string]bool{}
	for classLevelsConsumedClassID := range *classProgress.GetClassLevelsConsumed() {
		classLevelsConsumedByClassID[classLevelsConsumedClassID] = false
	}

	for classID, classLevelsConsumed := range otherClassLevelsConsumed {
		_, exists := classLevelsConsumedByClassID[classID]
		if !exists { return false }
		if !classProgress.ClassLevelsConsumed[classID].HasSameConsumptionAs(classLevelsConsumed) { return false }
		classLevelsConsumedByClassID[classID] = true
	}

	for _, wasFound := range classLevelsConsumedByClassID {
		if wasFound == false { return false }
	}

	return true
}

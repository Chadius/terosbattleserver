package squaddie

import (
	"fmt"
	"github.com/cserrant/terosBattleServer/entity/power"
	"github.com/cserrant/terosBattleServer/entity/squaddieclass"
	"github.com/cserrant/terosBattleServer/utility"
)

// Squaddie is the base unit you can deploy and control on a field.
type Squaddie struct {
	ID                  string                    `json:"id" yaml:"id"`
	Name                string                    `json:"name" yaml:"name"`
	Affiliation         string                    `json:"affiliation" yaml:"affiliation"`
	BaseClassID         string                    `json:"base_class" yaml:"base_class"`
	CurrentClass        string                    `json:"current_class" yaml:"current_class"`
	ClassLevelsConsumed map[string]*ClassProgress `json:"class_levels" yaml:"class_levels"`
	CurrentHitPoints    int                       `json:"current_hit_points" yaml:"current_hit_points"`
	MaxHitPoints        int                       `json:"max_hit_points" yaml:"max_hit_points"`
	Aim                 int                       `json:"aim" yaml:"aim"`
	Strength            int                       `json:"strength" yaml:"strength"`
	Mind                int                       `json:"mind" yaml:"mind"`
	Dodge               int                       `json:"dodge" yaml:"dodge"`
	Deflect             int                       `json:"deflect" yaml:"deflect"`
	CurrentBarrier      int                       `json:"current_barrier" yaml:"current_barrier"`
	MaxBarrier          int                       `json:"max_barrier" yaml:"max_barrier"`
	Armor               int                       `json:"armor" yaml:"armor"`
	Movement            Movement                  `json:"movement" yaml:"movement"`
	PowerReferences     []*power.Reference        `json:"powers" yaml:"powers"`
}

// NewSquaddie generates a squaddie with maxed out health.
func NewSquaddie(name string) *Squaddie {
	newSquaddie := Squaddie{
		ID:                  utility.StringWithCharset(8, "abcdefgh0123456789"),
		Name:                name,
		Affiliation:         "Player",
		CurrentHitPoints:    0,
		MaxHitPoints:        5,
		Aim:                 0,
		Strength:            0,
		Mind:                0,
		Dodge:               0,
		Deflect:             0,
		CurrentBarrier:      0,
		MaxBarrier:          0,
		Armor:               0,
		ClassLevelsConsumed: map[string]*ClassProgress{},
		Movement:            Movement{
			Distance:        3,
			Type:            Foot,
			HitAndRun:       false,
		},
	}
	newSquaddie.SetHPToMax()
	return &newSquaddie
}

// CheckSquaddieForErrors makes sure the created squaddie doesn't have an error.
func CheckSquaddieForErrors(newSquaddie *Squaddie) (newError error) {
	if newSquaddie.Affiliation != "Player" {
		return fmt.Errorf("Squaddie has unknown affiliation: '%s'", newSquaddie.Affiliation)
	}

	return nil
}

// SetHPToMax restores the Squaddie's HitPoints.
func (squaddie *Squaddie) SetHPToMax() {
	squaddie.CurrentHitPoints = squaddie.MaxHitPoints
}

// SetBarrierToMax restores the Squaddie's Barrier.
func (squaddie *Squaddie) SetBarrierToMax() {
	squaddie.CurrentBarrier = squaddie.MaxBarrier
}

// GetDefensiveStatsAgainstPhysical calculates how this squaddie can defend against physical attacks.
func (squaddie *Squaddie) GetDefensiveStatsAgainstPhysical() (evasion, barrierDamageReduction, armorDamageReduction int) {
	return squaddie.Dodge, squaddie.CurrentBarrier, squaddie.Armor
}

// GetDefensiveStatsAgainstSpell calculates how this squaddie can defend against spell attacks.
func (squaddie *Squaddie) GetDefensiveStatsAgainstSpell() (evasion, barrierDamageReduction, armorDamageReduction int) {
	return squaddie.Deflect, squaddie.CurrentBarrier, 0
}

// GetOffensiveStatsWithPhysical calculates the squaddie's bonuses with physical attacks.
func (squaddie *Squaddie) GetOffensiveStatsWithPhysical() (toHitBonus, damageBonus int) {
	return squaddie.Aim, squaddie.Strength
}

// GetOffensiveStatsWithSpell calculates the squaddie's bonuses with Spell attacks.
func (squaddie *Squaddie) GetOffensiveStatsWithSpell() (toHitBonus, damageBonus int) {
	return squaddie.Aim, squaddie.Mind
}

// AddInnatePower gives the Squaddie access to the power.
//  Raises an error if the squaddie already has the power.
func (squaddie *Squaddie) AddInnatePower(newPower *power.Power) error {
	if ContainsPowerID(squaddie.PowerReferences, newPower.ID) {
		return fmt.Errorf(`squaddie "%s" already has innate power with ID "%s"`, squaddie.Name, newPower.ID)
	}

	squaddie.PowerReferences = append(squaddie.PowerReferences, &power.Reference{Name: newPower.Name, ID: newPower.ID})
	return nil
}

// GetInnatePowerIDNames returns a list of all the powers the squaddie has access to.
func (squaddie *Squaddie) GetInnatePowerIDNames() []*power.Reference {
	powerIDNames := []*power.Reference{}
	for _, reference := range squaddie.PowerReferences {
		powerIDNames = append(powerIDNames, &power.Reference{Name: reference.Name, ID: reference.ID})
	}
	return powerIDNames
}

// AddClass gives the Squaddie a new class it can gain levels in.
func (squaddie *Squaddie) AddClass(class *squaddieclass.Class) {
	squaddie.ClassLevelsConsumed[class.ID] = &ClassProgress{
		ClassID:        class.ID,
		ClassName:      class.Name,
		LevelsConsumed: []string{},
	}
}

// GetLevelCountsByClass returns a mapping of class names to levels gained.
func (squaddie *Squaddie) GetLevelCountsByClass() map[string]int {
	count := map[string]int{}
	for classID, progress := range squaddie.ClassLevelsConsumed {
		count[classID] = len(progress.LevelsConsumed)
	}

	return count
}

// MarkLevelUpBenefitAsConsumed makes the Squaddie remember it used this benefit to level up already.
func (squaddie *Squaddie) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string)  {
	squaddie.ClassLevelsConsumed[benefitClassID].LevelsConsumed = append(squaddie.ClassLevelsConsumed[benefitClassID].LevelsConsumed, benefitID)
}

// GainLevelInClass is an alias for MarkLevelUpBenefitAsConsumed
func (squaddie *Squaddie) GainLevelInClass(classID, levelID string)  {
	squaddie.MarkLevelUpBenefitAsConsumed(classID, levelID)
}

// SetClass changes the Squaddie's CurrentClass to the given classID.
//   It also sets the BaseClass if it hasn't been already.
//   Raises an error if classID has not been added to the squaddie yet.
func (squaddie *Squaddie) SetClass(classID string) error {
	if _, exists := squaddie.ClassLevelsConsumed[classID]; !exists {
		return fmt.Errorf(`cannot switch "%s" to unknown class "%s"`, squaddie.Name, classID)
	}

	if squaddie.BaseClassID == "" {
		squaddie.BaseClassID = classID
	}

	squaddie.CurrentClass = classID
	return nil
}

// SetBaseClassIfNoBaseClass sets the BaseClass if it hasn't been already.
func (squaddie *Squaddie) SetBaseClassIfNoBaseClass(classID string) {
	if squaddie.BaseClassID == "" {
		squaddie.BaseClassID = classID
	}
}

// IsClassLevelAlreadyUsed returns true if a LevelUpBenefit with the given ID has already been used.
func (squaddie *Squaddie) IsClassLevelAlreadyUsed(benefitID string) bool {
	return squaddie.anyClassLevelsConsumed(func(classID string, progress *ClassProgress) bool {
		return progress.IsLevelAlreadyConsumed(benefitID)
	})
}

// HasAddedClass returns true if the Squaddie has already added a class with the name classIDToFind
func (squaddie *Squaddie) HasAddedClass(classIDToFind string) bool {
	return squaddie.anyClassLevelsConsumed(func(classID string, progress *ClassProgress) bool {
		return classID == classIDToFind
	})
}

// anyClassLevelsConsumed returns true if any of the squaddie's class levels consumed satisfies a given condition.
func (squaddie *Squaddie) anyClassLevelsConsumed(condition func(classID string, progress *ClassProgress) bool) bool {
	for classID, progress := range squaddie.ClassLevelsConsumed {
		if condition(classID, progress) {
			return true
		}
	}
	return false
}

// ClearInnatePowers removes all of the squaddie's powers.
func (squaddie *Squaddie) ClearInnatePowers() {
	squaddie.PowerReferences = []*power.Reference{}
}

// RemovePowerByID removes the power with the given ID from the squaddie's
func (squaddie *Squaddie) RemovePowerByID(powerToRemoveID string) {
	if ContainsPowerID(squaddie.PowerReferences, powerToRemoveID) == false {
		return
	}
	powerIndexToDelete := IndexOfPowerID(squaddie.PowerReferences, powerToRemoveID)
	squaddie.PowerReferences = append(squaddie.PowerReferences[:powerIndexToDelete], squaddie.PowerReferences[powerIndexToDelete+1:]...)
}

// ClearTemporaryPowerReferences empties the temporary references to powers.
func (squaddie *Squaddie) ClearTemporaryPowerReferences() {
	squaddie.PowerReferences = []*power.Reference{}
}

// GetMovementDistancePerRound Returns the distance the Squaddie can travel.
func (squaddie *Squaddie) GetMovementDistancePerRound() int {
	return squaddie.Movement.Distance
}

// GetMovementType returns the Squaddie's movement type
func (squaddie *Squaddie) GetMovementType() MovementType {
	return squaddie.Movement.Type
}

// CanHitAndRun indicates if the Squaddie can move after attacking.
func (squaddie *Squaddie) CanHitAndRun() bool {
	return squaddie.Movement.HitAndRun
}

// IndexOfPowerID returns the index of the reference to a power with the given ID.
//   returns -1 if it cannot be found.
func IndexOfPowerID(references []*power.Reference, powerID string) int {
	for index, reference := range references {
		if reference.ID == powerID {
			return index
		}
	}
	return -1
}

// ContainsPowerID returns true if the squaddie has a reference to a power with the given ID.
func ContainsPowerID(references []*power.Reference, powerID string) bool {
	for _, reference := range references {
		if reference.ID == powerID {
			return true
		}
	}
	return false
}

// FilterPowerID returns a list of power references that satisfy the condition.
func FilterPowerID(references []*power.Reference, condition func(*power.Reference) bool) []*power.Reference {
	selectedReferences := []*power.Reference{}
	for _, reference := range references {
		if condition(reference) == true {
			selectedReferences = append(selectedReferences, reference)
		}
	}
	return selectedReferences
}

// AnyPowerID returns true if one Power Reference satisfies the condition.
func AnyPowerID(references []*power.Reference, condition func(*power.Reference) bool) bool {
	for _, reference := range references {
		if condition(reference) == true {
			return true
		}
	}
	return false
}

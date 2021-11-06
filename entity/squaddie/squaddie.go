package squaddie

import (
	"fmt"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/utility"
)

// Affiliation describes the Squaddie's allegiance
type Affiliation string

// Squaddie Affiliation constants
const (
	Player  Affiliation = "Player"
	Enemy   Affiliation = "Enemy"
	Ally    Affiliation = "Ally"
	Neutral Affiliation = "Neutral"
)

// Squaddie represents a person, creature or thing that can take actions on a battlefield.
type Squaddie struct {
	Identification  Identification  `json:"identification" yaml:"identification"`
	ClassProgress   ClassProgress   `json:"class_progress" yaml:"class_progress"`
	Defense         Defense         `json:"defense" yaml:"defense"`
	Offense         Offense         `json:"offense" yaml:"offense"`
	Movement        Movement        `json:"movement" yaml:"movement"`
	PowerCollection PowerCollection `json:"powers" yaml:"powers"`
}

// NewSquaddie generates a squaddie with maxed out health.
func NewSquaddie(name string) *Squaddie {
	newSquaddie := Squaddie{
		Identification: Identification{
			SquaddieID:          "squaddie_" + utility.StringWithCharset(8, "abcdefgh0123456789"),
			SquaddieName:        name,
			SquaddieAffiliation: Player,
		},
		ClassProgress: ClassProgress{
			ClassLevelsConsumed: map[string]*ClassLevelsConsumed{},
		},
		Defense: Defense{
			SquaddieCurrentHitPoints: 0,
			SquaddieMaxHitPoints:     5,
			SquaddieDodge:            0,
			SquaddieDeflect:          0,
			SquaddieCurrentBarrier:   0,
			SquaddieMaxBarrier:       0,
			SquaddieArmor:            0,
		},
		Offense: Offense{
			SquaddieAim:      0,
			SquaddieStrength: 0,
			SquaddieMind:     0,
		},
		Movement: Movement{
			SquaddieMovementDistance:     3,
			SquaddieMovementType:         Foot,
			SquaddieMovementCanHitAndRun: false,
		},
	}
	newSquaddie.Defense.SetHPToMax()
	return &newSquaddie
}

// CheckSquaddieForErrors makes sure the created squaddie doesn't have an error.
func CheckSquaddieForErrors(newSquaddie *Squaddie) (newError error) {
	if !newSquaddie.Identification.HasValidAffiliation() {
		newError := fmt.Errorf("squaddie %s has unknown affiliation: '%s'", newSquaddie.ID(), newSquaddie.Affiliation())
		utility.Log(newError.Error(), 0, utility.Error)
		return newError
	}

	return nil
}

// ID delegates.
func (s *Squaddie) ID() string {
	return s.Identification.ID()
}

// Affiliation delegates.
func (s *Squaddie) Affiliation() Affiliation {
	return s.Identification.Affiliation()
}

// Name delegates.
func (s *Squaddie) Name() string {
	return s.Identification.Name()
}

// MaxHitPoints delegates.
func (s *Squaddie) MaxHitPoints() int {
	return s.Defense.MaxHitPoints()
}

// Dodge delegates.
func (s *Squaddie) Dodge() int {
	return s.Defense.Dodge()
}

// Deflect delegates.
func (s *Squaddie) Deflect() int {
	return s.Defense.Deflect()
}

// MaxBarrier delegates.
func (s *Squaddie) MaxBarrier() int {
	return s.Defense.MaxBarrier()
}

// Armor delegates.
func (s *Squaddie) Armor() int {
	return s.Defense.Armor()
}

// CurrentHitPoints delegates.
func (s *Squaddie) CurrentHitPoints() int {
	return s.Defense.CurrentHitPoints()
}

// CurrentBarrier delegates.
func (s *Squaddie) CurrentBarrier() int {
	return s.Defense.CurrentBarrier()
}

// Aim delegates.
func (s *Squaddie) Aim() int {
	return s.Offense.Aim()
}

// Strength delegates.
func (s *Squaddie) Strength() int {
	return s.Offense.Strength()
}

// Mind delegates.
func (s *Squaddie) Mind() int {
	return s.Offense.Mind()
}

//MovementDistance delegates.
func (s *Squaddie) MovementDistance() int {
	return s.Movement.MovementDistance()
}

// MovementType delegates.
func (s *Squaddie) MovementType() MovementType {
	return s.Movement.MovementType()
}

// MovementCanHitAndRun delegates.
func (s *Squaddie) MovementCanHitAndRun() bool {
	return s.Movement.CanHitAndRun()
}

// HasSameStatsAs returns true if other's stats matches this one.
//   The comparison ignores the ID.
func (s *Squaddie) HasSameStatsAs(other *Squaddie) bool {
	if s.Name() != other.Name() {
		return false
	}
	if s.Affiliation() != other.Affiliation() {
		return false
	}

	if !s.hasSameDefenseAs(other) {
		return false
	}
	if !s.hasSameOffenseAs(other) {
		return false
	}
	if !s.hasSameMovementAs(other) {
		return false
	}
	return true
}

func (s *Squaddie) hasSameMovementAs(other *Squaddie) bool {
	if s.MovementType() != other.MovementType() {
		return false
	}
	if s.MovementDistance() != other.MovementDistance() {
		return false
	}
	if s.MovementCanHitAndRun() != other.MovementCanHitAndRun() {
		return false
	}
	return true
}

func (s *Squaddie) hasSameOffenseAs(other *Squaddie) bool {
	if s.Aim() != other.Aim() {
		return false
	}
	if s.Strength() != other.Strength() {
		return false
	}
	if s.Mind() != other.Mind() {
		return false
	}
	return true
}

func (s *Squaddie) hasSameDefenseAs(other *Squaddie) bool {
	if s.MaxHitPoints() != other.MaxHitPoints() {
		return false
	}
	if s.Dodge() != other.Dodge() {
		return false
	}
	if s.Deflect() != other.Deflect() {
		return false
	}
	if s.MaxBarrier() != other.MaxBarrier() {
		return false
	}
	if s.Armor() != other.Armor() {
		return false
	}
	if s.CurrentHitPoints() != other.CurrentHitPoints() {
		return false
	}
	if s.CurrentBarrier() != other.CurrentBarrier() {
		return false
	}
	return true
}

// ClearPowerReferences delegates.
func (s *Squaddie) ClearPowerReferences() {
	s.PowerCollection.ClearPowerReferences()
}

// HasPowerWithID delegates.
func (s *Squaddie) HasPowerWithID(powerID string) bool {
	return s.PowerCollection.HasPowerWithID(powerID)
}

// HasEquippedPower delegates.
func (s *Squaddie) HasEquippedPower() bool {
	return s.PowerCollection.HasEquippedPower()
}

// EquipPower delegates.
func (s *Squaddie) EquipPower(powerID string) {
	s.PowerCollection.EquipPower(powerID)
}

// GetEquippedPowerID delegates.
func (s *Squaddie) GetEquippedPowerID() string {
	return s.PowerCollection.GetEquippedPowerID()
}

// AddPowerReference delegates.
func (s *Squaddie) AddPowerReference(reference *power.Reference) {
	s.PowerCollection.AddPowerReference(reference)
}

// RemovePowerReferenceByPowerID delegates.
func (s *Squaddie) RemovePowerReferenceByPowerID(powerID string) {
	s.PowerCollection.RemovePowerReferenceByPowerID(powerID)
}

// GetLevelCountsByClass delegates.
func (s *Squaddie) GetLevelCountsByClass() map[string]int {
	return s.ClassProgress.GetLevelCountsByClass()
}

// MarkLevelUpBenefitAsConsumed delegates.
func (s *Squaddie) MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID string) {
	s.ClassProgress.MarkLevelUpBenefitAsConsumed(benefitClassID, benefitID)
}

// SetClass delegates.
func (s *Squaddie) SetClass(classID string) error {
	return s.ClassProgress.SetClass(classID)
}

// SetBaseClassIfNoBaseClass delegates.
func (s *Squaddie) SetBaseClassIfNoBaseClass(classID string) {
	s.ClassProgress.SetBaseClassIfNoBaseClass(classID)
}

// IsClassLevelAlreadyUsed delegates.
func (s *Squaddie) IsClassLevelAlreadyUsed(benefitID string) bool {
	return s.ClassProgress.IsClassLevelAlreadyUsed(benefitID)
}

// HasAddedClass delegates.
func (s *Squaddie) HasAddedClass(classIDToFind string) bool {
	return s.ClassProgress.HasAddedClass(classIDToFind)
}

// AddClass delegates.
func (s *Squaddie) AddClass(classReference *squaddieclass.ClassReference) {
	s.ClassProgress.AddClass(classReference)
}
package levelup_test

import (
	"github.com/chadius/terosbattleserver/entity/levelupbenefit"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/entity/squaddieclass"
	"github.com/chadius/terosbattleserver/usecase/levelup"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	. "gopkg.in/check.v1"
)

type SquaddieUsesLevelUpBenefitSuite struct {
	mageClass              *squaddieclass.Class
	statBooster            *levelupbenefit.LevelUpBenefit
	teros                  *squaddie.Squaddie
	improveAllMovement     *levelupbenefit.LevelUpBenefit
	upgradeToLightMovement *levelupbenefit.LevelUpBenefit

	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
}

var _ = Suite(&SquaddieUsesLevelUpBenefitSuite{})

func (suite *SquaddieUsesLevelUpBenefitSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieclass.ClassBuilder().WithID("ffffffff").WithName("Mage").Build()
	suite.teros = squaddie.NewSquaddieBuilder().Teros().WithName("teros").Strength(1).Mind(2).Dodge(3).Deflect(4).Barrier(6).Armor(7).AddClassByReference(suite.mageClass.GetReference()).Build()
	suite.teros.Defense.SetBarrierToMax()

	suite.statBooster, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("deadbeef").
		WithClassID(suite.mageClass.ID()).
		Dodge(4).
		Deflect(3).
		Barrier(2).
		Armor(1).
		Aim(7).
		Strength(6).
		Mind(5).
		Build()

	suite.improveAllMovement, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("aaaaaaa0").
		WithClassID(suite.mageClass.ID()).
		MovementType(squaddie.Fly).
		CanHitAndRun().
		MovementDistance(1).
		Build()

	suite.upgradeToLightMovement, _ = levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("aaaaaaa1").
		WithClassID(suite.mageClass.ID()).
		MovementType(squaddie.Light).
		Build()

	suite.improveSquaddieStrategy = &levelup.ImproveSquaddieClass{}
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestIncreaseStats(checker *C) {
	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.statBooster, suite.teros)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.MaxHitPoints(), Equals, 5)
	checker.Assert(suite.teros.Aim(), Equals, 7)
	checker.Assert(suite.teros.Strength(), Equals, 7)
	checker.Assert(suite.teros.Mind(), Equals, 7)
	checker.Assert(suite.teros.Dodge(), Equals, 7)
	checker.Assert(suite.teros.Deflect(), Equals, 7)
	checker.Assert(suite.teros.MaxBarrier(), Equals, 8)
	checker.Assert(suite.teros.Armor(), Equals, 8)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieRecordsLevel(checker *C) {
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID()), Equals, false)
	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.statBooster, suite.teros)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{suite.mageClass.ID(): 1})
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID()), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorForNonexistentClass(checker *C) {
	mushroomClassLevel, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("deadbeeg").
		WithClassID("bad SquaddieID").
		Build()
	err := suite.improveSquaddieStrategy.ImproveSquaddie(mushroomClassLevel, suite.teros)
	checker.Assert(err.Error(), Equals, `squaddie "teros" cannot add levels to unknown class "bad SquaddieID"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestRaiseAnErrorIfReusingLevel(checker *C) {
	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.statBooster, suite.teros)
	checker.Assert(err, IsNil)
	checker.Assert(suite.teros.GetLevelCountsByClass(), DeepEquals, map[string]int{"ffffffff": 1})
	checker.Assert(suite.teros.IsClassLevelAlreadyUsed(suite.statBooster.ID()), Equals, true)

	err = suite.improveSquaddieStrategy.ImproveSquaddie(suite.statBooster, suite.teros)
	checker.Assert(err.Error(), Equals, `teros already consumed LevelUpBenefit - class:"ffffffff" id:"deadbeef"`)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestUsingLevelSetsBaseClassIfBaseClassIsUnset(checker *C) {
	checker.Assert(suite.teros.BaseClassID(), Equals, "")
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.statBooster, suite.teros)
	checker.Assert(suite.teros.BaseClassID(), Equals, suite.mageClass.ID())
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieChangeMovement(checker *C) {
	startingMovement := suite.teros.Movement.MovementDistance()

	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.improveAllMovement, suite.teros)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.Movement.MovementDistance(), Equals, startingMovement+1)
	checker.Assert(suite.teros.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.Movement.CanHitAndRun(), Equals, true)
}

func (suite *SquaddieUsesLevelUpBenefitSuite) TestSquaddieCannotDowngradeMovement(checker *C) {
	startingMovement := suite.teros.Movement.MovementDistance()
	suite.improveSquaddieStrategy.ImproveSquaddie(suite.improveAllMovement, suite.teros)

	err := suite.improveSquaddieStrategy.ImproveSquaddie(suite.upgradeToLightMovement, suite.teros)
	checker.Assert(err, IsNil)

	checker.Assert(suite.teros.Movement.MovementDistance(), Equals, startingMovement+1)
	checker.Assert(suite.teros.Movement.MovementType(), Equals, squaddie.MovementType(squaddie.Fly))
	checker.Assert(suite.teros.Movement.CanHitAndRun(), Equals, true)
}

type SquaddieChangePowersWithLevelUpBenefitsSuite struct {
	mageClass               *squaddieclass.Class
	teros                   *squaddie.Squaddie
	powerRepo               *powerrepository.Repository
	squaddieRepo            *squaddie.Repository
	repos                   *repositories.RepositoryCollection
	gainPower               levelupbenefit.LevelUpBenefit
	upgradePower            levelupbenefit.LevelUpBenefit
	spear                   *power.Power
	spearLevel2             *power.Power
	improveSquaddieStrategy levelup.ImproveSquaddieStrategy
}

var _ = Suite(&SquaddieChangePowersWithLevelUpBenefitsSuite{})

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) SetUpTest(checker *C) {
	suite.mageClass = squaddieclass.ClassBuilder().WithID("ffffffff").WithName("Mage").Build()
	suite.teros = squaddie.NewSquaddieBuilder().Teros().AddClassByReference(suite.mageClass.GetReference()).Build()
	suite.teros.Defense.SetBarrierToMax()

	suite.powerRepo = powerrepository.NewPowerRepository()

	suite.spear = power.NewPowerBuilder().Spear().WithID("spearlvl1").Build()

	suite.teros.AddPowerReference(&power.Reference{
		Name:    "spear",
		PowerID: "spearlvl1",
	})

	suite.spearLevel2 = power.NewPowerBuilder().Spear().WithID("spearlvl2").Build()
	newPowers := []*power.Power{suite.spear, suite.spearLevel2}
	suite.powerRepo.AddSlicePowerSource(newPowers)

	newLevel, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("aaab1234").
		WithClassID(suite.mageClass.ID()).
		BigLevel().
		GainPower(suite.spear.PowerID, "spear").
		Build()

	suite.gainPower = *newLevel

	upgradePowerLevel, _ := levelupbenefit.NewLevelUpBenefitBuilder().
		WithID("aaab1235").
		WithClassID(suite.mageClass.ID()).
		BigLevel().
		GainPower(suite.spearLevel2.PowerID, "spear").
		LosePower(suite.spear.PowerID).
		Build()

	suite.upgradePower = *upgradePowerLevel

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{suite.teros})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}
	suite.improveSquaddieStrategy = &levelup.ImproveSquaddieClass{}
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieGainPowers(checker *C) {
	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.gainPower, suite.teros)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetCopyOfPowerReferences()
	checker.Assert(len(attackIDNamePairs), Equals, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.spear.PowerID)
}

func (suite *SquaddieChangePowersWithLevelUpBenefitsSuite) TestSquaddieLosePowers(checker *C) {
	suite.improveSquaddieStrategy.ImproveSquaddie(&suite.gainPower, suite.teros)
	suite.teros.PowerCollection.GetCopyOfPowerReferences()

	err := suite.improveSquaddieStrategy.ImproveSquaddie(&suite.upgradePower, suite.teros)
	checker.Assert(err, IsNil)

	attackIDNamePairs := suite.teros.PowerCollection.GetCopyOfPowerReferences()
	checker.Assert(attackIDNamePairs, HasLen, 1)
	checker.Assert(attackIDNamePairs[0].Name, Equals, "spear")
	checker.Assert(attackIDNamePairs[0].PowerID, Equals, suite.spearLevel2.PowerID)
}

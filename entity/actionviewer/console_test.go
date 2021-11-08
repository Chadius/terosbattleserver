package actionviewer_test

import (
	"github.com/chadius/terosbattleserver/entity/actionviewer"
	"github.com/chadius/terosbattleserver/entity/power"
	"github.com/chadius/terosbattleserver/entity/powerrepository"
	"github.com/chadius/terosbattleserver/entity/powerusagescenario"
	"github.com/chadius/terosbattleserver/entity/squaddie"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
	"github.com/chadius/terosbattleserver/usecase/powercommit"
	"github.com/chadius/terosbattleserver/usecase/powerequip"
	"github.com/chadius/terosbattleserver/usecase/repositories"
	"github.com/chadius/terosbattleserver/utility/testutility"
	powerBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/power"
	squaddieBuilder "github.com/chadius/terosbattleserver/utility/testutility/builder/squaddie"
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ConsoleViewerSuite struct {
	teros   *squaddie.Squaddie
	bandit  *squaddie.Squaddie
	bandit2 *squaddie.Squaddie
	lini    *squaddie.Squaddie

	blot         *power.Power
	axe          *power.Power
	healingStaff *power.Power

	powerRepo    *powerrepository.Repository
	squaddieRepo *squaddie.Repository
	repos        *repositories.RepositoryCollection

	forecastBlotOnBandit *powerattackforecast.Forecast
	resultBlotOnBandit   *powercommit.Result

	forecastBlotOnMultipleBandits *powerattackforecast.Forecast
	resultBlotOnMultipleBandits   *powercommit.Result

	forecastHealingStaffOnTeros *powerattackforecast.Forecast
	resultHealingStaffOnTeros   *powercommit.Result

	viewer *actionviewer.ConsoleActionViewer
}

var _ = Suite(&ConsoleViewerSuite{})

func (suite *ConsoleViewerSuite) SetUpTest(checker *C) {
	suite.teros = squaddieBuilder.Builder().Teros().Build()
	suite.bandit = squaddieBuilder.Builder().Bandit().Build()
	suite.bandit2 = squaddieBuilder.Builder().Bandit().WithName("Bandit2").WithID("banditID2").Build()
	suite.lini = squaddieBuilder.Builder().Lini().Build()

	suite.blot = powerBuilder.Builder().Blot().WithName("Blot").DealsDamage(0).Build()
	suite.axe = powerBuilder.Builder().Axe().Build()
	suite.healingStaff = powerBuilder.Builder().HealingStaff().WithName("Healing Staff").Build()

	suite.squaddieRepo = squaddie.NewSquaddieRepository()
	suite.squaddieRepo.AddSquaddies([]*squaddie.Squaddie{
		suite.teros,
		suite.bandit,
		suite.bandit2,
		suite.lini,
	})

	suite.powerRepo = powerrepository.NewPowerRepository()
	suite.powerRepo.AddSlicePowerSource([]*power.Power{
		suite.blot,
		suite.axe,
		suite.healingStaff,
	})

	suite.repos = &repositories.RepositoryCollection{
		SquaddieRepo: suite.squaddieRepo,
		PowerRepo:    suite.powerRepo,
	}

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.teros,
		[]*power.Reference{
			suite.blot.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	powerequip.LoadAllOfSquaddieInnatePowers(
		suite.bandit2,
		[]*power.Reference{
			suite.axe.GetReference(),
		},
		suite.repos,
	)

	suite.forecastBlotOnBandit = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultBlotOnBandit = &powercommit.Result{
		Forecast: suite.forecastBlotOnBandit,
	}

	suite.forecastBlotOnMultipleBandits = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.teros.ID(),
			PowerID:         suite.blot.ID(),
			Targets:         []string{suite.bandit.ID(), suite.bandit2.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}
	suite.resultBlotOnMultipleBandits = &powercommit.Result{
		Forecast: suite.forecastBlotOnMultipleBandits,
	}

	suite.forecastHealingStaffOnTeros = &powerattackforecast.Forecast{
		Setup: powerusagescenario.Setup{
			UserID:          suite.lini.ID(),
			PowerID:         suite.healingStaff.PowerID,
			Targets:         []string{suite.teros.ID()},
			IsCounterAttack: false,
		},
		Repositories: &repositories.RepositoryCollection{
			SquaddieRepo: suite.squaddieRepo,
			PowerRepo:    suite.powerRepo,
		},
	}

	suite.resultHealingStaffOnTeros = &powercommit.Result{
		Forecast: suite.forecastHealingStaffOnTeros,
	}

	suite.viewer = &actionviewer.ConsoleActionViewer{}
}

func (suite *ConsoleViewerSuite) SetUpTerosAttacksBanditsAndSuffersCounterAttack() {
	suite.resultBlotOnMultipleBandits.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Aim(suite.teros.Aim()).Mind(3).Build()

	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().
		Deflect(suite.bandit.Deflect()).
		HitPoints(suite.bandit.MaxHitPoints()).
		Barrier(1).
		Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.axe = powerBuilder.Builder().CloneOf(suite.axe).WithID(suite.axe.ID()).CanCounterAttack().DealsDamage(3).Build()
	suite.powerRepo.AddPower(suite.axe)

	suite.bandit.Offense = *squaddieBuilder.OffenseBuilder().Strength(0).Build()
	powerequip.SquaddieEquipPower(suite.bandit2, suite.axe.PowerID, suite.repos)

	suite.forecastBlotOnMultipleBandits.CalculateForecast()
	suite.resultBlotOnMultipleBandits.Commit()
}

func (suite *ConsoleViewerSuite) SetUpLiniHealsTeros() {
	suite.resultHealingStaffOnTeros.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Defense.SetHPToMax()
	suite.teros.Defense.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini.Offense = *squaddieBuilder.OffenseBuilder().Mind(1).Build()

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	suite.resultHealingStaffOnTeros.Commit()
}

func (suite *ConsoleViewerSuite) TestShowPowerHitTargetAndDamage(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Aim(suite.teros.Aim()).Mind(3).Build()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 3 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowWhenPowerMisses(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysMissDieRoller{}

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) misses Bandit\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowWhenPowerCriticallyHits(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}
	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Mind(3).Build()
	suite.blot = powerBuilder.Builder().CloneOf(suite.blot).WithID(suite.blot.ID()).CriticalDealsDamage(1).CriticalHitThresholdBonus(9000).Build()
	suite.powerRepo.AddPower(suite.blot)

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil, &output)
	checker.Assert(output.String(), Equals, "Teros (Blot) CRITICALLY hits Bandit, for 4 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowCounterattacks(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.axe = powerBuilder.Builder().CloneOf(suite.axe).WithID(suite.axe.ID()).CanCounterAttack().DealsDamage(2).Build()
	suite.powerRepo.AddPower(suite.axe)

	suite.bandit.Offense = *squaddieBuilder.OffenseBuilder().Strength(0).Build()
	powerequip.SquaddieEquipPower(suite.bandit, suite.axe.PowerID, suite.repos)

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 0 damage\nBandit (axe) counters Teros, for 2 damage\n---\n")
}

func (suite *ConsoleViewerSuite) TestIndicateIfItIsAKillingBlow(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}
	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Mind(suite.bandit.MaxHitPoints() * 2).Build()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()

	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, felling\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowPowerBarrierBurn(checker *C) {
	suite.resultBlotOnBandit.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Mind(3).Build()
	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().Barrier(1).Build()
	suite.bandit.Defense.SetBarrierToMax()

	suite.forecastBlotOnBandit.CalculateForecast()
	suite.resultBlotOnBandit.Commit()
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnBandit, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowMultipleTargets(checker *C) {
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultBlotOnMultipleBandits, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"---\n",
	)
}

func (suite *ConsoleViewerSuite) TestShowPowerHealingEffects(checker *C) {
	suite.resultHealingStaffOnTeros.DieRoller = &testutility.AlwaysHitDieRoller{}

	suite.teros.Defense.SetHPToMax()
	suite.teros.Defense.ReduceHitPoints(suite.teros.MaxHitPoints() - 1)
	suite.lini.Offense = *squaddieBuilder.OffenseBuilder().Mind(1).Build()

	suite.forecastHealingStaffOnTeros.CalculateForecast()
	suite.resultHealingStaffOnTeros.Commit()
	var output strings.Builder
	suite.viewer.PrintResult(suite.resultHealingStaffOnTeros, suite.repos, nil, &output)

	checker.Assert(output.String(), Equals, "Lini (Healing Staff) heals Teros, for 4 healing\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowTargetStatusVerbosity(checker *C) {
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()
	var counterAttackOutput strings.Builder
	suite.viewer.PrintResult(
		suite.resultBlotOnMultipleBandits,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowTargetStatus: true,
		},
		&counterAttackOutput,
	)

	checker.Assert(counterAttackOutput.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"   Bandit: 3/5 HP, 0 barrier\n"+
			"   Bandit2: 2/5 HP\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"   Teros: 2/5 HP\n"+
			"---\n",
	)

	suite.SetUpLiniHealsTeros()
	var output strings.Builder
	suite.viewer.PrintResult(
		suite.resultHealingStaffOnTeros,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowTargetStatus: true,
		},
		&output,
	)

	checker.Assert(output.String(), Equals, "Lini (Healing Staff) heals Teros, for 4 healing\n   Teros: 5/5 HP\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowRollsVerbosity(checker *C) {
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()
	var counterAttackOutput strings.Builder
	suite.viewer.PrintResult(
		suite.resultBlotOnMultipleBandits,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowRolls: true,
		},
		&counterAttackOutput,
	)

	checker.Assert(counterAttackOutput.String(), Equals,
		"Teros (Blot) hits Bandit, for 2 damage + 1 barrier burn\n"+
			"- also hits Bandit2, for 3 damage\n"+
			"   Teros rolls 999 + 0 = 999, Bandit rolls -999 + 0 = -999\n"+
			"   Teros rolls 999 + 0 = 999, Bandit2 rolls -999 + 0 = -999\n"+
			"Bandit2 (axe) counters Teros, for 3 damage\n"+
			"   Bandit2 rolls 999 + -1 = 998, Teros rolls -999 + 0 = -999\n"+
			"---\n",
	)

	suite.SetUpLiniHealsTeros()
	var healingOutput strings.Builder
	suite.viewer.PrintResult(
		suite.resultHealingStaffOnTeros,
		suite.repos,
		&actionviewer.ConsoleActionViewerVerbosity{
			ShowRolls: true,
		},
		&healingOutput,
	)

	checker.Assert(healingOutput.String(), Equals, "Lini (Healing Staff) heals Teros, for 4 healing\n   Auto-hit\n---\n")
}

func (suite *ConsoleViewerSuite) TestShowForecastChanceToHitAndHealing(checker *C) {
	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().
		Aim(2).
		Build()

	suite.bandit2.Defense = *squaddieBuilder.DefenseBuilder().
		Deflect(2).
		Barrier(20).
		Build()
	suite.bandit2.Defense.SetBarrierToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()
	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals,
		"Teros (Blot) vs Bandit: +2 (30/36), for 2 damage + 1 barrier burn\n"+
			"- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn\n"+
			"Bandit2 (axe) counters Teros: -1 (15/36), for 3 damage\n",
	)

	suite.SetUpLiniHealsTeros()
	var healingOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastHealingStaffOnTeros,
		suite.repos,
		&healingOutput,
	)

	checker.Assert(healingOutput.String(), Equals, "Lini (Healing Staff) heals Teros, for 4 healing\n")
}

func (suite *ConsoleViewerSuite) TestShowForecastChanceToCriticallyHitAndGuaranteedMiss(checker *C) {
	suite.teros.Offense = *squaddieBuilder.OffenseBuilder().Aim(2).Build()
	suite.blot = powerBuilder.Builder().CloneOf(suite.blot).WithID(suite.blot.ID()).CriticalDealsDamage(1).CriticalHitThresholdBonus(1).Build()
	suite.powerRepo.AddPower(suite.blot)
	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().Deflect(-200).Build()
	suite.bandit2.Defense = *squaddieBuilder.DefenseBuilder().Deflect(2).Barrier(20).Build()
	suite.bandit2.Defense.SetBarrierToMax()
	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals, "Teros (Blot) vs Bandit: +202 (36/36), for 2 damage + 1 barrier burn\n"+
		" crit: 36/36, for 3 damage + 1 barrier burn\n"+
		"- also vs Bandit2: +0 (21/36) for NO DAMAGE + 3 barrier burn\n"+
		" crit: 1/36 for NO DAMAGE + 4 barrier burn\n"+
		"Bandit2 (axe) counters Teros: -1 (15/36), for 3 damage\n",
	)
}

func (suite *ConsoleViewerSuite) TestShowForecastAttackIsFatal(checker *C) {

	suite.bandit.Defense = *squaddieBuilder.DefenseBuilder().HitPoints(2).Build()
	suite.bandit.Defense.SetHPToMax()

	suite.bandit2.Defense = *squaddieBuilder.DefenseBuilder().Deflect(2).HitPoints(2).Build()
	suite.bandit2.Defense.SetHPToMax()

	suite.SetUpTerosAttacksBanditsAndSuffersCounterAttack()

	var forecastOutput strings.Builder
	suite.viewer.PrintForecast(
		suite.forecastBlotOnMultipleBandits,
		suite.repos,
		&forecastOutput,
	)

	checker.Assert(forecastOutput.String(), Equals, "Teros (Blot) vs Bandit: +0 (21/36), FATAL\n- also vs Bandit2: -2 (10/36), FATAL\nBandit2 (axe) counters Teros: -1 (15/36), for 3 damage\n")
}

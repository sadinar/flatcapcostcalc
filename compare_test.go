package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateCostPerMythic(t *testing.T) {
	metric := calculateBaseGoldCostPerMythic(1, 1000)
	assert.Equal(t, float64(1000), metric)

	metric = calculateBaseGoldCostPerMythic(2, 1000)
	assert.Equal(t, float64(500), metric)

	metric = calculateBaseGoldCostPerMythic(2, 200000000)
	assert.Equal(t, float64(100000000), metric)
}

func TestPerformAscendedFuseEarlyReturnLessThanThree(t *testing.T) {
	pc := PurchaseConfiguration{}
	petList := map[string]uint64{}
	pc.performAscendedFuse(petList)
	assert.Empty(t, petList)

	petList[Ascended] = 2
	pc.performAscendedFuse(petList)
	assert.Empty(t, petList[Mythical])
	assert.Equal(t, uint64(2), petList[Ascended])
}

func TestPerformAscendedFuseCalculation(t *testing.T) {
	pc := PurchaseConfiguration{}
	petList := map[string]uint64{Epic: 1, Legendary: 2000, Ascended: 4, Mythical: 1}
	pc.performAscendedFuse(petList)
	assert.Equal(t, uint64(1), petList[Epic])
	assert.Equal(t, uint64(2000), petList[Legendary])
	assert.Equal(t, uint64(1), petList[Ascended])
	assert.Equal(t, uint64(2), petList[Mythical])

	petList = map[string]uint64{Ascended: 30}
	pc.performAscendedFuse(petList)
	assert.Equal(t, uint64(10), petList[Mythical])
}

func TestPerformProdigiousFuseEarlyReturnLessThanThree(t *testing.T) {
	pc := PurchaseConfiguration{FuseLuckPercentage: 0.25}
	petList := map[string]uint64{}
	pc.performProdigiousFuse(petList)
	assert.Empty(t, petList)

	petList[Prodigious] = 2
	pc.performProdigiousFuse(petList)
	assert.Empty(t, petList[Mythical])
	assert.Empty(t, petList[Ascended])
	assert.Equal(t, uint64(2), petList[Prodigious])

	pc.FuseLuckPercentage = 0
	pc.performProdigiousFuse(petList)
	assert.Empty(t, petList[Mythical])
	assert.Empty(t, petList[Ascended])
	assert.Equal(t, uint64(2), petList[Prodigious])
}

func TestPerformProdigiousFuseCalculation(t *testing.T) {
	pc := PurchaseConfiguration{FuseLuckPercentage: 0.25}
	petList := map[string]uint64{Epic: 5678, Legendary: 2, Prodigious: 5, Ascended: 9, Mythical: 1}
	pc.performProdigiousFuse(petList)
	assert.Equal(t, uint64(5678), petList[Epic])
	assert.Equal(t, uint64(2), petList[Legendary])
	assert.Equal(t, uint64(2), petList[Prodigious])
	assert.Equal(t, uint64(10), petList[Ascended])
	assert.Equal(t, uint64(1), petList[Mythical])

	petList = map[string]uint64{Prodigious: 16}
	pc.FuseLuckPercentage = 0.08
	pc.performProdigiousFuse(petList)
	assert.Equal(t, uint64(1), petList[Prodigious])
	assert.Equal(t, uint64(5), petList[Ascended])
}

func TestPerformBaseFiveFuseEarlyReturnLessThanFive(t *testing.T) {
	pc := PurchaseConfiguration{}
	petList := map[string]uint64{}
	pc.performBaseFiveFuse(petList, Epic)
	assert.Empty(t, petList)

	petList = map[string]uint64{Epic: 1, Legendary: 100}
	pc.FuseLuckPercentage = 0.14
	pc.performBaseFiveFuse(petList, Epic)
	assert.Equal(t, 2, len(petList))
	assert.Equal(t, uint64(1), petList[Epic])
	assert.Equal(t, uint64(100), petList[Legendary])

	petList = map[string]uint64{Epic: 100, Legendary: 2}
	pc.FuseLuckPercentage = 0.25
	pc.performBaseFiveFuse(petList, Legendary)
	assert.Equal(t, 2, len(petList))
	assert.Equal(t, uint64(100), petList[Epic])
	assert.Equal(t, uint64(2), petList[Legendary])

	pc = PurchaseConfiguration{TypeBuying: Rare}
	petList = map[string]uint64{}
	pc.performBaseFiveFuse(petList, Rare)
	assert.Empty(t, petList)
}

func TestPerformBaseFiveFuseCalculation(t *testing.T) {
	pc := PurchaseConfiguration{FuseLuckPercentage: 0.25}
	petList := map[string]uint64{Epic: 100}
	pc.performBaseFiveFuse(petList, Epic)
	assert.Equal(t, 3, len(petList))
	assert.Equal(t, uint64(0), petList[Epic])
	assert.Equal(t, uint64(15), petList[Legendary])
	assert.Equal(t, uint64(5), petList[Prodigious])

	petList[Epic] = 26
	petList[Legendary] = 1
	petList[Prodigious] = 1
	pc.FuseLuckPercentage = 0.01
	pc.performBaseFiveFuse(petList, Epic)
	assert.Equal(t, 3, len(petList))
	assert.Equal(t, uint64(1), petList[Epic])
	assert.Equal(t, uint64(6), petList[Legendary])
	assert.Equal(t, uint64(1), petList[Prodigious])

	petList[Epic] = 100
	petList[Legendary] = 104
	petList[Prodigious] = 1
	petList[Ascended] = 1
	pc.FuseLuckPercentage = 0.25
	pc.performBaseFiveFuse(petList, Legendary)
	assert.Equal(t, 4, len(petList))
	assert.Equal(t, uint64(100), petList[Epic])
	assert.Equal(t, uint64(4), petList[Legendary])
	assert.Equal(t, uint64(16), petList[Prodigious])
	assert.Equal(t, uint64(6), petList[Ascended])

	petList = map[string]uint64{Rare: 50}
	pc.performBaseFiveFuse(petList, Rare)
	assert.Equal(t, uint64(0), petList[Rare])
	assert.Equal(t, uint64(8), petList[Epic])
	assert.Equal(t, uint64(2), petList[Legendary])
	assert.Equal(t, 3, len(petList))
}

func TestCalculateMaxUpgradedPetsCyclesThroughAllPetTypes(t *testing.T) {
	pc := PurchaseConfiguration{
		EggLuckPercentage:  0,
		FuseLuckPercentage: 0,
	}
	pets := map[string]uint64{Epic: 94, Legendary: 53, Prodigious: 2, Ascended: 0, Mythical: 0}
	pc.calculateMaxUpgradedPets(pets)
	assert.Equal(t, uint64(4), pets[Epic])
	assert.Equal(t, uint64(1), pets[Legendary])
	assert.Equal(t, uint64(1), pets[Prodigious])
	assert.Equal(t, uint64(2), pets[Ascended])
	assert.Equal(t, uint64(1), pets[Mythical])

	pets[Epic] = 4
	pets[Legendary] = 3
	pets[Prodigious] = 4
	pets[Ascended] = 6
	pets[Mythical] = 0
	pc.calculateMaxUpgradedPets(pets)
	assert.Equal(t, uint64(4), pets[Epic])
	assert.Equal(t, uint64(3), pets[Legendary])
	assert.Equal(t, uint64(1), pets[Prodigious])
	assert.Equal(t, uint64(1), pets[Ascended])
	assert.Equal(t, uint64(2), pets[Mythical])
}

func TestCalculateMaxUpgradedPetsUsesLuckPercentages(t *testing.T) {
	pc := PurchaseConfiguration{
		EggLuckPercentage:  0.25,
		FuseLuckPercentage: 0.25,
	}
	pets := map[string]uint64{Epic: 1000}
	pc.calculateMaxUpgradedPets(pets)
	assert.Equal(t, uint64(0), pets[Epic])
	assert.Equal(t, uint64(0), pets[Legendary])
	assert.Equal(t, uint64(1), pets[Prodigious])
	assert.Equal(t, uint64(1), pets[Ascended])
	assert.Equal(t, uint64(14), pets[Mythical])

	delete(pets, Legendary)
	delete(pets, Prodigious)
	delete(pets, Ascended)
	delete(pets, Mythical)
	pc.EggLuckPercentage = 0.1
	pc.FuseLuckPercentage = 0.1
	pets[Epic] = 1000
	pc.calculateMaxUpgradedPets(pets)
	assert.Equal(t, uint64(0), pets[Epic])
	assert.Equal(t, uint64(0), pets[Legendary])
	assert.Equal(t, uint64(2), pets[Prodigious])
	assert.Equal(t, uint64(1), pets[Ascended])
	assert.Equal(t, uint64(7), pets[Mythical])
}

func TestCalculateHatchedPetsReturnsHatchCountWhenHatchLuckZero(t *testing.T) {
	pc := PurchaseConfiguration{TypeBuying: Prodigious}
	hatchedPets := pc.calculateHatchedPets(1000)
	assert.Equal(t, 1, len(hatchedPets))
	assert.Equal(t, uint64(1000), hatchedPets[Prodigious])

	pc.TypeBuying = Legendary
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 1, len(hatchedPets))
	assert.Equal(t, uint64(1000), hatchedPets[Legendary])

	pc.TypeBuying = Epic
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 1, len(hatchedPets))
	assert.Equal(t, uint64(1000), hatchedPets[Epic])
}

func TestCalculateHatchedPetsReturnsEpicsAndLegendaryWhenHatchingEpicWithLuck(t *testing.T) {
	pc := PurchaseConfiguration{TypeBuying: Epic, EggLuckPercentage: 0.25}
	hatchedPets := pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(750), hatchedPets[Epic])
	assert.Equal(t, uint64(250), hatchedPets[Legendary])

	pc.EggLuckPercentage = 0.1
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(900), hatchedPets[Epic])
	assert.Equal(t, uint64(100), hatchedPets[Legendary])

	pc.EggLuckPercentage = 0.01
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(990), hatchedPets[Epic])
	assert.Equal(t, uint64(10), hatchedPets[Legendary])
}

func TestCalculateHatchedPetsReturnsLegendaryAndProdigiousWhenHatchingLegendaryWithLuck(t *testing.T) {
	pc := PurchaseConfiguration{TypeBuying: Legendary, EggLuckPercentage: 0.24}
	hatchedPets := pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(760), hatchedPets[Legendary])
	assert.Equal(t, uint64(240), hatchedPets[Prodigious])

	pc.EggLuckPercentage = 0.13
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(870), hatchedPets[Legendary])
	assert.Equal(t, uint64(130), hatchedPets[Prodigious])

	pc.EggLuckPercentage = 0.05
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(950), hatchedPets[Legendary])
	assert.Equal(t, uint64(50), hatchedPets[Prodigious])
}

func TestCalculateHatchedPetsReturnsProdigiousAndAscendedWhenHatchingProdigiousWithLuck(t *testing.T) {
	pc := PurchaseConfiguration{TypeBuying: Prodigious, EggLuckPercentage: 0.22}
	hatchedPets := pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(780), hatchedPets[Prodigious])
	assert.Equal(t, uint64(220), hatchedPets[Ascended])

	pc.EggLuckPercentage = 0.09
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(910), hatchedPets[Prodigious])
	assert.Equal(t, uint64(90), hatchedPets[Ascended])

	pc.EggLuckPercentage = 0.02
	hatchedPets = pc.calculateHatchedPets(1000)
	assert.Equal(t, 2, len(hatchedPets))
	assert.Equal(t, uint64(980), hatchedPets[Prodigious])
	assert.Equal(t, uint64(20), hatchedPets[Ascended])
}

func TestCalculateHatchedPetsHandlesRareCorrectly(t *testing.T) {
	pc := PurchaseConfiguration{TypeBuying: Rare, EggLuckPercentage: 0.07}
	hatchedPets := pc.calculateHatchedPets(100)
	assert.Equal(t, uint64(93), hatchedPets[Rare])
	assert.Equal(t, uint64(7), hatchedPets[Epic])
}

func TestCalculateTotalPetsRejectsInvalidConfig(t *testing.T) {
	pc := PurchaseConfiguration{TypeBuying: "some_junk"}
	_, err := pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid pet type", err.Error())
	pc.TypeBuying = "epic1"
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid pet type", err.Error())

	pc = PurchaseConfiguration{TypeBuying: Epic, EggLuckPercentage: -0.01}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid egg luck percentage", err.Error())
	pc = PurchaseConfiguration{TypeBuying: Legendary, EggLuckPercentage: 0.33}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid egg luck percentage", err.Error())

	pc = PurchaseConfiguration{TypeBuying: Epic, FuseLuckPercentage: -1}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid fuse luck percentage", err.Error())
	pc = PurchaseConfiguration{TypeBuying: Prodigious, FuseLuckPercentage: 0.26}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid fuse luck percentage", err.Error())

	pc = PurchaseConfiguration{TypeBuying: Epic, AchievementGoldBonus: -1}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid achievement gold bonus percentage", err.Error())
	pc = PurchaseConfiguration{TypeBuying: Prodigious, AchievementGoldBonus: 0.36}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid achievement gold bonus percentage", err.Error())

	pc = PurchaseConfiguration{TypeBuying: Epic, FriendGoldBonus: -0.01}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid friend gold bonus percentage", err.Error())
	pc = PurchaseConfiguration{TypeBuying: Prodigious, FriendGoldBonus: 0.31}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid friend gold bonus percentage", err.Error())

	pc = PurchaseConfiguration{TypeBuying: Epic, CaveGoldBonus: -0.1}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid cave gold bonus percentage", err.Error())
	pc = PurchaseConfiguration{TypeBuying: Prodigious, CaveGoldBonus: 1.01}
	_, err = pc.CalculateTotalPets()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid cave gold bonus percentage", err.Error())
}

func TestCalculateTotalPetsNoOposOnNoMoneySpent(t *testing.T) {
	pc := PurchaseConfiguration{MoneySpending: 0, TypeBuying: Prodigious}
	pets, err := pc.CalculateTotalPets()
	assert.Nil(t, err)
	assert.Empty(t, pets)
}

func TestGetCheaperEggsPrices(t *testing.T) {
	pc := PurchaseConfiguration{}
	priceList := pc.getCheaperEggsPrices()
	assert.Equal(t, uint64(650000), priceList[Epic])
	assert.Equal(t, uint64(3000000), priceList[Legendary])
	assert.Equal(t, uint64(10000000), priceList[Prodigious])
}

func TestGetEvenCheaperEggsPrices(t *testing.T) {
	pc := PurchaseConfiguration{}
	priceList := pc.getEvenCheaperEggsPrices()
	assert.Equal(t, uint64(550000), priceList[Epic])
	assert.Equal(t, uint64(2500000), priceList[Legendary])
	assert.Equal(t, uint64(8000000), priceList[Prodigious])
}

func TestCalculateTotalPets(t *testing.T) {
	pc := PurchaseConfiguration{
		MoneySpending:      1000000000,
		TypeBuying:         Legendary,
		EggLuckPercentage:  0.19,
		FuseLuckPercentage: 0.06,
	}
	pc.PetPrices = pc.getCheaperEggsPrices()

	pets, err := pc.CalculateTotalPets()
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), pets[Legendary])
	assert.Equal(t, uint64(0), pets[Prodigious])
	assert.Equal(t, uint64(0), pets[Ascended])
	assert.Equal(t, uint64(15), pets[Mythical])

	pc.EggLuckPercentage = 0
	pc.FuseLuckPercentage = 0
	pets, err = pc.CalculateTotalPets()
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), pets[Legendary])
	assert.Equal(t, uint64(0), pets[Prodigious])
	assert.Equal(t, uint64(1), pets[Ascended])
	assert.Equal(t, uint64(7), pets[Mythical])

	pc.EggLuckPercentage = 0.32
	pc.FuseLuckPercentage = 0.25
	pets, err = pc.CalculateTotalPets()
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), pets[Legendary])
	assert.Equal(t, uint64(2), pets[Prodigious])
	assert.Equal(t, uint64(1), pets[Ascended])
	assert.Equal(t, uint64(26), pets[Mythical])
}

func TestSetSpendableGold(t *testing.T) {
	pc := PurchaseConfiguration{
		TypeBuying:           Legendary,
		AchievementGoldBonus: 0.03,
		CaveGoldBonus:        0.17,
		FriendGoldBonus:      0.2,
		HasCoinBonusPass:     true,
		HasDoubleCoinBoost:   true,
	}
	pc.PetPrices = pc.getCheaperEggsPrices()
	pc.setSpendableGold(100)

	expectedGold := uint64(float64(100) * (1 + 0.03 + 0.17 + 0.2 + 1 + 0.5))
	assert.Equal(t, expectedGold, pc.MoneySpending)
}

func TestConstructor(t *testing.T) {
	pc := New(100, 0, CheaperPriceTable, Legendary, 0.1, 0.11, 0.15, 0.13, 0.2, true, false)
	assert.Equal(t, uint64(247), pc.MoneySpending)
	assert.Equal(t, pc.getCheaperEggsPrices(), pc.PetPrices)
	assert.Equal(t, Legendary, pc.TypeBuying)
	assert.Equal(t, float64(0.1), pc.EggLuckPercentage)
	assert.Equal(t, float64(0.11), pc.FuseLuckPercentage)
	assert.Equal(t, float64(0.15), pc.AchievementGoldBonus)
	assert.Equal(t, float64(0.13), pc.CaveGoldBonus)
	assert.Equal(t, float64(0.2), pc.FriendGoldBonus)
	assert.True(t, pc.HasDoubleCoinBoost)
	assert.False(t, pc.HasCoinBonusPass)

	pc = New(100, 0, EvenCheaperPriceTable, Prodigious, 0.01, 0.02, 0.25, 0.02, 0.1, false, true)
	assert.Equal(t, uint64(187), pc.MoneySpending)
	assert.Equal(t, pc.getEvenCheaperEggsPrices(), pc.PetPrices)
	assert.Equal(t, Prodigious, pc.TypeBuying)
	assert.Equal(t, float64(0.01), pc.EggLuckPercentage)
	assert.Equal(t, float64(0.02), pc.FuseLuckPercentage)
	assert.Equal(t, float64(0.25), pc.AchievementGoldBonus)
	assert.Equal(t, float64(0.02), pc.CaveGoldBonus)
	assert.Equal(t, float64(0.1), pc.FriendGoldBonus)
	assert.False(t, pc.HasDoubleCoinBoost)
	assert.True(t, pc.HasCoinBonusPass)
}

func TestSanityCheckGoldConstants(t *testing.T) {
	assert.True(t, OneHundredMillion < OneBillion)
	assert.True(t, OneBillion < OneTrillion)
}

func TestPetsOverTimeIsZeroWhenNoGoldPerSecondProvided(t *testing.T) {
	pc := New(100, 0, CheaperPriceTable, Legendary, 0.1, 0.11, 0.15, 0.13, 0.2, true, false)
	petsPerHour, err := pc.CalculateMythicPetsPerHour()

	assert.Nil(t, err)
	assert.Equal(t, float32(0), petsPerHour)
}

func TestPetsOverTimeHandlesSmallNumberGracefully(t *testing.T) {
	pc := New(100, 100, CheaperPriceTable, Legendary, 0.1, 0.11, 0.15, 0.13, 0.2, true, false)
	petsPerHour, err := pc.CalculateMythicPetsPerHour()

	assert.Nil(t, err)
	assert.Equal(t, float32(0), petsPerHour)
}

func TestPetsOverTimeHandlesRealisticNumberGracefully(t *testing.T) {
	pc := New(100, 1000000, CheaperPriceTable, Legendary, 0.1, 0.11, 0.15, 0.13, 0.2, true, false)
	petsPerHour, err := pc.CalculateMythicPetsPerHour()

	assert.Nil(t, err)
	assert.Equal(t, float32(49.77), petsPerHour)
}

func TestPetsOverTimeHandlesLargeNumberGracefully(t *testing.T) {
	pc := New(100, OneBillion, CheaperPriceTable, Legendary, 0.1, 0.11, 0.15, 0.13, 0.2, true, false)
	petsPerHour, err := pc.CalculateMythicPetsPerHour()

	assert.Nil(t, err)
	assert.Equal(t, float32(50245.33), petsPerHour)
}

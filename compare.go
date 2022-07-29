package main

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	pc1 := New(
		OneBillion,
		EvenCheaperPriceTable,
		Epic,
		0.30,  // egg luck
		0.25,  // fuse luck
		0.0,   // achievement coin bonus
		0,     // cave coin bonus
		0.0,   // friend coin bonus
		false, // has temporary double coin boost
		false, // has permanent 1.5x coin game pass
	)
	pc1TotalPets, err := pc1.CalculateTotalPets()
	if err != nil {
		panic(err)
	}

	pc2 := New(
		OneBillion,
		EvenCheaperPriceTable,
		Prodigious,
		0.25, // egg luck
		0.06, // fuse luck
		0.0,  // achievement coin bonus
		0,    // cave coin bonus
		0.0,  // friend coin bonus
		false,
		false,
	)
	pc2TotalPets, err := pc2.CalculateTotalPets()
	if err != nil {
		panic(err)
	}

	printComparison(pc1, pc2, pc1TotalPets, pc2TotalPets)
}

func getCostsPerMythic(mythicCount int, displayCost int) (baseCostPerMythic float32, displayedCostPerMythic float32) {
	return calculateBaseGoldCostPerMythic(mythicCount, OneBillion), calculateBaseGoldCostPerMythic(mythicCount, displayCost)
}

func calculateBaseGoldCostPerMythic(mythicCount, moneySpent int) float32 {
	return float32(moneySpent) / float32(mythicCount)
}

func printComparison(pc1, pc2 PurchaseConfiguration, pc1Pets, pc2Pets map[string]int) {
	pc1BaseCostPerMythic, pc1DisplayedCostPerMythic := getCostsPerMythic(pc1Pets[Mythical], pc1.MoneySpending)
	pc2BaseCostPerMythic, pc2DisplayedCostPerMythic := getCostsPerMythic(pc2Pets[Mythical], pc2.MoneySpending)

	p := message.NewPrinter(language.English)
	_, _ = p.Printf("Base gold: %d.\n", OneBillion)
	_, _ = p.Printf("Setup 1's displayed gold before spending: %d.\n", pc1.MoneySpending)
	_, _ = p.Printf("Setup 2's displayed gold before spending: %d.\n", pc2.MoneySpending)
	if pc1BaseCostPerMythic > pc2BaseCostPerMythic {
		_, _ = p.Printf("Setup 2's cost per mythic is better.\n")
		_, _ = p.Printf("Base gold setup 2 vs setup 1: %d per mythic vs %d\n", int(pc2BaseCostPerMythic), int(pc1BaseCostPerMythic))
		_, _ = p.Printf("Displayed gold setup 2 vs setup 1: %d per mythic vs %d\n", int(pc2DisplayedCostPerMythic), int(pc1DisplayedCostPerMythic))
	} else if pc1BaseCostPerMythic < pc2BaseCostPerMythic {
		_, _ = p.Printf("Setup 1's cost per mythic is better.\n")
		_, _ = p.Printf("Base gold setup 1 vs setup 2: %d per mythic vs %d\n", int(pc1BaseCostPerMythic), int(pc2BaseCostPerMythic))
		_, _ = p.Printf("Displayed gold setup 1 vs setup 2: %d per mythic vs %d\n", int(pc1DisplayedCostPerMythic), int(pc2DisplayedCostPerMythic))
	} else {
		fmt.Println("Both setups produce the same number of mythic pets")
	}

	fmt.Println()
	fmt.Println("setup 1 totals:")
	fmt.Println(pc1Pets)

	fmt.Println()
	fmt.Println("setup 2 totals:")
	fmt.Println(pc2Pets)
}

const Rare = "rare"
const Epic = "epic"
const Legendary = "legendary"
const Prodigious = "prodigious"
const Ascended = "ascended"
const Mythical = "mythical"
const CheaperPriceTable = "cheaper"
const EvenCheaperPriceTable = "even_cheaper"
const OneBillion = 1000000000

type PurchaseConfiguration struct {
	MoneySpending        int
	PetPrices            map[string]int
	TypeBuying           string
	EggLuckPercentage    float32
	FuseLuckPercentage   float32
	AchievementGoldBonus float32
	CaveGoldBonus        float32
	FriendGoldBonus      float32
	HasDoubleCoinBoost   bool
	HasCoinBonusPass     bool
}

func New(baseGold int, priceTable, typeBuying string, eggLuckPercentage, fuseLuckPercentage, achievementGoldBonus, caveGoldBonus, friendGoldBonus float32, hasDoubleBoost, hasCoinPass bool) PurchaseConfiguration {
	pc := PurchaseConfiguration{
		TypeBuying:           typeBuying,
		EggLuckPercentage:    eggLuckPercentage,
		FuseLuckPercentage:   fuseLuckPercentage,
		AchievementGoldBonus: achievementGoldBonus,
		CaveGoldBonus:        caveGoldBonus,
		FriendGoldBonus:      friendGoldBonus,
		HasDoubleCoinBoost:   hasDoubleBoost,
		HasCoinBonusPass:     hasCoinPass,
	}
	pc.setSpendableGold(baseGold)
	if priceTable == EvenCheaperPriceTable {
		pc.PetPrices = pc.getEvenCheaperEggsPrices()
	} else {
		pc.PetPrices = pc.getCheaperEggsPrices()
	}
	return pc
}

func (pc *PurchaseConfiguration) CalculateTotalPets() (map[string]int, error) {
	err := pc.validate()
	if err != nil {
		return make(map[string]int, 0), err
	}

	if pc.MoneySpending == 0 {
		return make(map[string]int, 0), nil
	}

	eggsPurchased := pc.MoneySpending / pc.PetPrices[pc.TypeBuying]
	hatchedPetCounts := pc.calculateHatchedPets(eggsPurchased)
	finalPetList := pc.calculateMaxUpgradedPets(hatchedPetCounts)

	return finalPetList, nil
}

func (pc *PurchaseConfiguration) getCheaperEggsPrices() map[string]int {
	return map[string]int{
		Rare:       140000,
		Epic:       650000,
		Legendary:  3000000,
		Prodigious: 10000000,
	}
}

func (pc *PurchaseConfiguration) getEvenCheaperEggsPrices() map[string]int {
	return map[string]int{
		Rare:       120000,
		Epic:       550000,
		Legendary:  2500000,
		Prodigious: 8000000,
	}
}

func (pc *PurchaseConfiguration) setSpendableGold(baseGold int) {
	coinMultiplier := 1 + pc.CaveGoldBonus + pc.AchievementGoldBonus + pc.FriendGoldBonus
	if pc.HasDoubleCoinBoost {
		coinMultiplier += 1
	}
	if pc.HasCoinBonusPass {
		coinMultiplier += 0.5
	}
	gold := float32(baseGold) * coinMultiplier

	pc.MoneySpending = int(gold)
}

func (pc *PurchaseConfiguration) validate() error {
	if pc.MoneySpending < 0 {
		return fmt.Errorf("invalid gold amount")
	}

	isValidType := pc.TypeBuying == Rare ||
		pc.TypeBuying == Epic ||
		pc.TypeBuying == Legendary ||
		pc.TypeBuying == Prodigious
	if !isValidType {
		return fmt.Errorf("invalid pet type")
	}

	if pc.EggLuckPercentage < 0 || pc.EggLuckPercentage > 0.32 {
		return fmt.Errorf("invalid egg luck percentage")
	}

	if pc.FuseLuckPercentage < 0 || pc.FuseLuckPercentage > 0.25 {
		return fmt.Errorf("invalid fuse luck percentage")
	}

	if pc.AchievementGoldBonus < 0 || pc.AchievementGoldBonus > 0.35 {
		return fmt.Errorf("invalid achievement gold bonus percentage")
	}

	if pc.CaveGoldBonus < 0 || pc.CaveGoldBonus > 1 {
		return fmt.Errorf("invalid cave gold bonus percentage")
	}

	if pc.FriendGoldBonus < 0 || pc.FriendGoldBonus > 0.3 {
		return fmt.Errorf("invalid friend gold bonus percentage")
	}

	return nil
}

func (pc *PurchaseConfiguration) calculateHatchedPets(eggsHatched int) map[string]int {
	if pc.EggLuckPercentage == 0 {
		return map[string]int{pc.TypeBuying: eggsHatched}
	}

	upgradedPetsHatched := int(float32(eggsHatched) * pc.EggLuckPercentage)
	basePetsHatched := eggsHatched - upgradedPetsHatched
	switch pc.TypeBuying {
	case Rare:
		return map[string]int{
			Rare: basePetsHatched,
			Epic: upgradedPetsHatched,
		}
	case Epic:
		return map[string]int{
			Epic:      basePetsHatched,
			Legendary: upgradedPetsHatched,
		}
	case Legendary:
		return map[string]int{
			Legendary:  basePetsHatched,
			Prodigious: upgradedPetsHatched,
		}
	case Prodigious:
		return map[string]int{
			Prodigious: basePetsHatched,
			Ascended:   upgradedPetsHatched,
		}
	default:
		return make(map[string]int, 0)
	}
}

func (pc *PurchaseConfiguration) calculateMaxUpgradedPets(hatchedPetCounts map[string]int) map[string]int {
	pc.performBaseFiveFuse(hatchedPetCounts, Rare)
	pc.performBaseFiveFuse(hatchedPetCounts, Epic)
	pc.performBaseFiveFuse(hatchedPetCounts, Legendary)
	pc.performProdigiousFuse(hatchedPetCounts)
	pc.performAscendedFuse(hatchedPetCounts)

	return hatchedPetCounts
}

func (pc *PurchaseConfiguration) performBaseFiveFuse(hatchedPetCounts map[string]int, petRarity string) {
	if hatchedPetCounts[petRarity] < 5 {
		return
	}

	fuseCount := hatchedPetCounts[petRarity] / 5
	hatchedPetCounts[petRarity] = hatchedPetCounts[petRarity] % 5

	upgradedCount := int(float32(fuseCount) * pc.FuseLuckPercentage)
	standardCount := fuseCount - upgradedCount

	switch petRarity {
	case Rare:
		hatchedPetCounts[Epic] += standardCount
		hatchedPetCounts[Legendary] += upgradedCount
	case Epic:
		hatchedPetCounts[Legendary] += standardCount
		hatchedPetCounts[Prodigious] += upgradedCount
	case Legendary:
		hatchedPetCounts[Prodigious] += standardCount
		hatchedPetCounts[Ascended] += upgradedCount
	}
}

func (pc *PurchaseConfiguration) performProdigiousFuse(hatchedPetCounts map[string]int) {
	if hatchedPetCounts[Prodigious] < 3 {
		return
	}

	fuseCount := hatchedPetCounts[Prodigious] / 3
	hatchedPetCounts[Prodigious] = hatchedPetCounts[Prodigious] % 3

	upgradedCount := int(float32(fuseCount) * pc.FuseLuckPercentage)
	standardCount := fuseCount - upgradedCount

	hatchedPetCounts[Ascended] += standardCount
	hatchedPetCounts[Mythical] += upgradedCount
}

func (pc *PurchaseConfiguration) performAscendedFuse(hatchedPetCounts map[string]int) {
	if hatchedPetCounts[Ascended] < 3 {
		return
	}

	fuseCount := hatchedPetCounts[Ascended] / 3
	hatchedPetCounts[Ascended] = hatchedPetCounts[Ascended] % 3

	hatchedPetCounts[Mythical] += fuseCount
}

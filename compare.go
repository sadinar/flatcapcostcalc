package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"math/big"
)

const Rare = "rare"
const Epic = "epic"
const Legendary = "legendary"
const Prodigious = "prodigious"
const Ascended = "ascended"
const Mythical = "mythical"
const CheaperPriceTable = "cheaper"
const EvenCheaperPriceTable = "even_cheaper"
const OneMillion = uint64(1000000)
const OneBillion = uint64(1000000000)
const OneHundredMillion = uint64(100000000)
const OneTrillion = uint64(1000000000000)
const ManualHatchSpeed = 1

func main() {
	performTimeRestrictedComparison()
}

/**
For instructions, see https://github.com/sadinar/flatcapcostcalc
*/

func performTimeRestrictedComparison() {
	trc := NewTimeRestrictedCalculator(
		[]uint64{24 * 1}, // hours spent hatching
		900*OneMillion,   // gold per minute
		80,               // calcify chance
		3.86,             // generate per second
		0.25+0.07,        // egg luck
		0.25,             // fuse luck
		2.00,             // shiny wall luck
		1.09,             // shiny achievement
		1.0526683,        // experts luck
		0.00040*1.5,      // metallic percent
		Legendary,        // type generating
		Legendary,        // type hatching
	)
	trc.Calculate()
}

type TimeRestricted struct {
	HoursOfGenerating     uint64
	HoursOfManualHatching uint64
	GoldPerMinute         uint64
	GeneratePerSecond     float64
	MetallicChance        float64
	calcifyChance         float64
	GenerationHatcher     PetHatcher
	ManualHatcher         PetHatcher
	msgPrinter            MsgPrinter
}

type MsgPrinter interface {
	printMessage(msg string)
}

type StdOutPrinter struct{}

func (sop *StdOutPrinter) printMessage(msg string) {
	fmt.Println(msg)
}

func NewTimeRestrictedCalculator(hoursOfHatching []uint64, goldPerMinute, calcifyChance uint64, generatePerSecond, eggLuckPercentage, fuseLuckPercentage, shinyWallLuck, shinyAchievementLuck, expertsLuck, metallicChance float64, typeGenerating, typeHatching string) TimeRestricted {
	if shinyWallLuck < 1.00 {
		shinyWallLuck = 1.00
	}
	if len(hoursOfHatching) == 1 {
		hoursOfHatching = append(hoursOfHatching, hoursOfHatching[0])
	}

	return TimeRestricted{
		HoursOfGenerating:     hoursOfHatching[0],
		HoursOfManualHatching: hoursOfHatching[1],
		GeneratePerSecond:     generatePerSecond,
		MetallicChance:        metallicChance,
		calcifyChance:         float64(calcifyChance) / 100,
		GoldPerMinute:         goldPerMinute,
		msgPrinter:            &StdOutPrinter{},
		GenerationHatcher: PetHatcher{
			TypeBuying:           typeGenerating,
			PriceTable:           EvenCheaperPriceTable,
			EggLuckPercentage:    eggLuckPercentage,
			FuseLuckPercentage:   fuseLuckPercentage,
			ShinyWallLuck:        shinyWallLuck,
			ShinyAchievementLuck: shinyAchievementLuck,
			ExpertsLuck:          expertsLuck,
		},
		ManualHatcher: PetHatcher{
			TypeBuying:           typeHatching,
			PriceTable:           EvenCheaperPriceTable,
			EggLuckPercentage:    eggLuckPercentage,
			FuseLuckPercentage:   fuseLuckPercentage,
			ShinyWallLuck:        shinyWallLuck,
			ShinyAchievementLuck: shinyAchievementLuck,
			ExpertsLuck:          expertsLuck,
		},
	}
}

func (ts *TimeRestricted) Calculate() error {
	moneyEarned := ts.calculateMoneyEarned()
	ts.setMoneySpending()
	if ts.GenerationHatcher.MoneySpending+ts.ManualHatcher.MoneySpending > moneyEarned {
		p := message.NewPrinter(language.English)
		msg := p.Sprintf(
			"%d out of %d available",
			moneyEarned,
			ts.GenerationHatcher.MoneySpending+ts.ManualHatcher.MoneySpending,
		)
		return fmt.Errorf("not enough money to hatch all those eggs! " + msg)
	}

	generatedPets, err := ts.GenerationHatcher.HatchPets()
	if err != nil {
		return err
	}

	manuallyHatchedPets, err := ts.ManualHatcher.HatchPets()
	if err != nil {
		return err
	}

	crystalCount := uint64(float64(generatedPets[Mythical]) * (ts.calcifyChance + 1.0))
	crystalCount += uint64(float64(manuallyHatchedPets[Mythical]) * (ts.calcifyChance + 1.0))

	shinyScore := ts.GenerationHatcher.GetShinyScore()
	shinyScore += ts.ManualHatcher.GetShinyScore()

	totalPetCount := ts.GenerationHatcher.GetTotalHatchedPetCount() + ts.ManualHatcher.GetTotalHatchedPetCount()
	multiMetallicCount, _ := FindReasonableProbability(totalPetCount, ts.MetallicChance/100)

	moneyLeftOver := moneyEarned - ts.GenerationHatcher.MoneySpending - ts.ManualHatcher.MoneySpending

	genSpeedUpgrades, err := ts.getGenerationSpeedUpgradeCount(crystalCount)
	if err != nil {
		return err
	}

	calcifyUpgrades, err := ts.calculateCalcifyUpgradeCount(crystalCount)
	if err != nil {
		return err
	}

	p := message.NewPrinter(language.English)
	msg := p.Sprintf("\nMetrics for %d hours generating and %d hours manually hatching:\n", ts.HoursOfGenerating, ts.HoursOfManualHatching)
	msg += p.Sprintf("mythic crystals: %d\n", crystalCount)
	msg += p.Sprintf("pet score gained: %d\n", generatedPets[Mythical]+manuallyHatchedPets[Mythical])
	msg += p.Sprintf("shiny score gained: %d\n", shinyScore)
	msg += p.Sprintf("\nmoney spent hatching: %d\n", ts.GenerationHatcher.MoneySpending+ts.ManualHatcher.MoneySpending)
	msg += p.Sprintf("money left over: %d\n", moneyLeftOver)
	msg += p.Sprintf("\nShiny wall upgrades possible: %d\n", ts.getShinyWallUpgradeCount(moneyLeftOver))
	msg += p.Sprintf("Gen speed upgrades possible: %d\n", genSpeedUpgrades)
	msg += p.Sprintf("Calcify upgrades possible: %d\n\n", calcifyUpgrades)
	totalProbability := 0.0
	for i := uint64(0); i <= multiMetallicCount; i++ {
		singleProbability := BinomialProbability(totalPetCount, i, ts.MetallicChance/100)
		totalProbability += singleProbability
		msg += p.Sprintf("probability of exactly %d metallic: %.2f%%\n", i, singleProbability*100)
	}
	msg += p.Sprintf("probability of %d+ metallic: %.2f%%\n", multiMetallicCount+1, (1-totalProbability)*100)

	ts.msgPrinter.printMessage(msg)
	return nil
}

func (ts *TimeRestricted) calculateMoneyEarned() uint64 {
	if ts.HoursOfGenerating >= ts.HoursOfManualHatching {
		return ts.GoldPerMinute * 60 * ts.HoursOfGenerating
	}

	return ts.GoldPerMinute * 60 * ts.HoursOfManualHatching
}

func (ts *TimeRestricted) calculateMetallicChance(totalPetCount uint64) float64 {
	oneRollNoMetallicChance := float64(1) - ts.MetallicChance/100
	noMetallicChance := math.Pow(oneRollNoMetallicChance, float64(totalPetCount))
	return 1 - noMetallicChance
}

func (ts *TimeRestricted) setMoneySpending() {
	prices := ts.GenerationHatcher.getEvenCheaperEggsPrices()
	ts.GenerationHatcher.MoneySpending = uint64(
		float64(prices[ts.GenerationHatcher.TypeBuying]) * ts.GeneratePerSecond * 60 * 60 * float64(ts.HoursOfGenerating),
	)
	ts.ManualHatcher.MoneySpending = uint64(
		float64(prices[ts.ManualHatcher.TypeBuying]) * ManualHatchSpeed * 60 * 60 * float64(ts.HoursOfManualHatching),
	)
}

func (ts *TimeRestricted) getShinyWallUpgradeCount(gold uint64) uint64 {
	if math.Round(ts.GenerationHatcher.ShinyWallLuck*100)/100 >= 2.00 {
		return 0
	}

	currentCost := uint64(
		math.Round((ts.GenerationHatcher.ShinyWallLuck-1)*100+1),
	) * OneBillion
	upgradeCount := uint64(0)
	for {
		if gold < currentCost || currentCost > 100*OneBillion {
			break
		}
		gold -= currentCost
		currentCost += OneBillion
		upgradeCount++
	}

	return upgradeCount
}

func (ts *TimeRestricted) getGenerationSpeedUpgradeCount(crystals uint64) (uint64, error) {
	if ts.GeneratePerSecond > 4.99 {
		return 0, nil
	}

	speed := ts.GeneratePerSecond
	currentCost, err := ts.calculateCurrentSpeedUpgradeCost()
	if err != nil {
		return 0, err
	}

	upgradeCount := uint64(0)
	for {
		if crystals < currentCost || speed >= 5.00 {
			break
		}
		crystals -= currentCost
		increase, err := ts.getGenerationSpeedCostIncrease(ts.GeneratePerSecond)
		if err != nil {
			return 0, err
		}
		currentCost += increase
		speed += 0.01
		speed = math.Round(speed*100) / 100
		upgradeCount++
	}

	return upgradeCount, nil
}

func (ts *TimeRestricted) getGenerationSpeedCostIncrease(currentSpeed float64) (uint64, error) {
	if currentSpeed < 0.25 {
		return 0, fmt.Errorf("invalid generation speed provided")
	}

	switch {
	case currentSpeed < 0.42:
		compareValue := uint64(math.Round(currentSpeed * 100))

		if compareValue%2 == 0 {
			return 0, nil
		}
		return 1, nil
	case currentSpeed < 0.59:
		compareValue := uint64(math.Round(currentSpeed * 100))

		if compareValue%2 == 0 {
			return 1, nil
		}
		return 0, nil
	case currentSpeed < 0.74:
		compareValue := uint64(math.Round(currentSpeed * 100))

		if compareValue%2 == 0 {
			return 0, nil
		}
		return 1, nil
	case currentSpeed < 0.91:
		compareValue := uint64(math.Round(currentSpeed * 100))

		if compareValue%2 == 0 {
			return 1, nil
		}
		return 0, nil
	case currentSpeed < 1.00:
		compareValue := uint64(math.Round(currentSpeed * 100))

		if compareValue%2 == 0 {
			return 0, nil
		}
		return 1, nil
	case currentSpeed < 2.00:
		return 2, nil
	case currentSpeed < 3.00:
		compareValue := uint64(math.Round(currentSpeed * 100))

		if compareValue%2 == 0 {
			return 13, nil
		}
		return 12, nil
	case currentSpeed < 4.00:
		return 45, nil
	case currentSpeed < 5.00:
		return 50, nil
	default:
		return 0, fmt.Errorf("invalid generation speed provided")
	}
}

func (ts *TimeRestricted) calculateCurrentSpeedUpgradeCost() (uint64, error) {
	upgradeCost := uint64(10)
	currentAdjustedSpeed := uint64(math.Round(ts.GeneratePerSecond * 100))
	for speed := uint64(26); speed <= currentAdjustedSpeed; speed++ {
		costIncrease, err := ts.getGenerationSpeedCostIncrease(float64(speed-1) / 100)
		if err != nil {
			return 0, err
		}
		upgradeCost += costIncrease
	}

	return upgradeCost, nil
}

func (ts *TimeRestricted) calculateCalcifyUpgradeCount(crystals uint64) (uint64, error) {
	currentCalcifyChance := uint64(math.Round(ts.calcifyChance * 100))
	upgradeCount := uint64(0)
	upgradeCost := ts.calculateCalcificationSunkCost()
	for {
		if crystals == 0 || currentCalcifyChance == 100 {
			break
		}

		increase, err := ts.getCalcifyUpgradeCost(currentCalcifyChance)
		if err != nil {
			return 0, err
		}

		upgradeCost += increase
		if upgradeCost > crystals {
			break
		}
		crystals -= upgradeCost
		upgradeCount++
		currentCalcifyChance++
	}

	return upgradeCount, nil
}

func (ts *TimeRestricted) calculateCalcificationSunkCost() uint64 {
	totalCost := uint64(0)
	for upgradeCount := uint64(0); upgradeCount < uint64(math.Round(ts.calcifyChance*100)); upgradeCount++ {
		incrementalCost, _ := ts.getCalcifyUpgradeCost(upgradeCount)
		totalCost += incrementalCost
	}

	return totalCost
}

func (ts *TimeRestricted) getCalcifyUpgradeCost(currentValue uint64) (uint64, error) {
	switch {
	case currentValue >= 0 && currentValue < 50:
		return uint64(50), nil
	case currentValue < 100:
		return uint64(100), nil
	default:
		return 0, fmt.Errorf("invalid calcification chance: %d: must be a number from 0 to 100", currentValue)
	}
}

func Factorial(factor uint64) uint64 {
	if factor == 1 {
		return 1
	}
	return factor * Factorial(factor-1)
}

func TotalCombinations(trials uint64, successes uint64) *big.Int {
	if successes == 0 {
		return big.NewInt(1)
	}

	numerator := big.NewInt(int64(trials))
	for i := trials - 1; i > trials-successes; i-- {
		numerator.Mul(numerator, big.NewInt(int64(i)))
	}

	return numerator.Div(numerator, big.NewInt(int64(Factorial(successes))))
}

func BinomialProbability(trials uint64, successes uint64, pSuccess float64) float64 {
	combinations := TotalCombinations(trials, successes)
	fCombinations := big.NewFloat(0)
	_, _, err := fCombinations.Parse(combinations.String(), 10)
	if err != nil {
		panic(err)
	}

	pFail := float64(1) - pSuccess
	s := math.Pow(pSuccess, float64(successes))
	f := math.Pow(pFail, float64(trials-successes))
	t := fCombinations.Mul(fCombinations, big.NewFloat(s))
	t = t.Mul(t, big.NewFloat(f))
	probability, _ := t.Float64()

	return probability
}

func FindReasonableProbability(trials uint64, pSuccess float64) (uint64, float64) {
	totalProbability := float64(0)
	successCount := uint64(0)
	for {
		if totalProbability >= 0.95 {
			break
		}
		p := BinomialProbability(trials, successCount, pSuccess)
		totalProbability += p
		successCount++
	}

	return successCount - 1, totalProbability
}

type PetHatcher struct {
	MoneySpending        uint64
	TypeBuying           string
	PriceTable           string
	petPrices            map[string]uint64
	EggLuckPercentage    float64
	FuseLuckPercentage   float64
	ShinyWallLuck        float64
	ShinyAchievementLuck float64
	ExpertsLuck          float64
	allHatchedPets       map[string]uint64
}

func (ph *PetHatcher) getCheaperEggsPrices() map[string]uint64 {
	return map[string]uint64{
		Rare:       140000,
		Epic:       650000,
		Legendary:  3000000,
		Prodigious: 10000000,
	}
}

func (ph *PetHatcher) getEvenCheaperEggsPrices() map[string]uint64 {
	return map[string]uint64{
		Rare:       120000,
		Epic:       550000,
		Legendary:  2500000,
		Prodigious: 8000000,
	}
}

func (ph *PetHatcher) getShinyPetValues() map[string]uint64 {
	return map[string]uint64{
		Rare:       5,
		Epic:       10,
		Legendary:  15,
		Prodigious: 20,
		Ascended:   30,
		Mythical:   40,
	}
}

func (ph *PetHatcher) HatchPets() (map[string]uint64, error) {
	err := ph.validate()
	if err != nil {
		return make(map[string]uint64, 0), err
	}

	if ph.MoneySpending == 0 {
		return make(map[string]uint64, 0), nil
	}

	err = ph.setPetPrices()
	if err != nil {
		return make(map[string]uint64, 0), nil
	}

	eggsPurchased := ph.MoneySpending / ph.petPrices[ph.TypeBuying]
	hatchedPetCounts := ph.calculateHatchedPets(eggsPurchased)
	ph.allHatchedPets = make(map[string]uint64, 0)
	for eggType, count := range hatchedPetCounts {
		ph.allHatchedPets[eggType] = count
	}

	finalPetList := ph.calculateMaxUpgradedPets(hatchedPetCounts)

	return finalPetList, nil
}

func (ph *PetHatcher) setPetPrices() error {
	switch ph.PriceTable {
	case CheaperPriceTable:
		ph.petPrices = ph.getCheaperEggsPrices()
	case EvenCheaperPriceTable:
		ph.petPrices = ph.getEvenCheaperEggsPrices()
	default:
		return fmt.Errorf("unknown price table specified: " + ph.PriceTable)
	}

	return nil
}

func (ph *PetHatcher) validate() error {
	isValidType := ph.TypeBuying == Rare ||
		ph.TypeBuying == Epic ||
		ph.TypeBuying == Legendary ||
		ph.TypeBuying == Prodigious
	if !isValidType {
		return fmt.Errorf("invalid pet type")
	}

	if ph.EggLuckPercentage < 0 || ph.EggLuckPercentage > 0.35 {
		return fmt.Errorf("invalid egg luck percentage")
	}

	if ph.FuseLuckPercentage < 0 || ph.FuseLuckPercentage > 0.25 {
		return fmt.Errorf("invalid fuse luck percentage")
	}

	return nil
}

func (ph *PetHatcher) calculateHatchedPets(eggsHatched uint64) map[string]uint64 {
	if ph.EggLuckPercentage == 0 {
		return map[string]uint64{ph.TypeBuying: eggsHatched}
	}

	upgradedPetsHatched := uint64(float64(eggsHatched) * ph.EggLuckPercentage)
	basePetsHatched := eggsHatched - upgradedPetsHatched
	switch ph.TypeBuying {
	case Rare:
		return map[string]uint64{
			Rare: basePetsHatched,
			Epic: upgradedPetsHatched,
		}
	case Epic:
		return map[string]uint64{
			Epic:      basePetsHatched,
			Legendary: upgradedPetsHatched,
		}
	case Legendary:
		return map[string]uint64{
			Legendary:  basePetsHatched,
			Prodigious: upgradedPetsHatched,
		}
	case Prodigious:
		return map[string]uint64{
			Prodigious: basePetsHatched,
			Ascended:   upgradedPetsHatched,
		}
	default:
		return make(map[string]uint64, 0)
	}
}

func (ph *PetHatcher) calculateMaxUpgradedPets(hatchedPetCounts map[string]uint64) map[string]uint64 {
	ph.performBaseFiveFuse(hatchedPetCounts, Rare)
	ph.performBaseFiveFuse(hatchedPetCounts, Epic)
	ph.performBaseFiveFuse(hatchedPetCounts, Legendary)
	ph.performProdigiousFuse(hatchedPetCounts)
	ph.performAscendedFuse(hatchedPetCounts)

	return hatchedPetCounts
}

func (ph *PetHatcher) performBaseFiveFuse(hatchedPetCounts map[string]uint64, petRarity string) {
	if hatchedPetCounts[petRarity] < 5 {
		return
	}

	fuseCount := hatchedPetCounts[petRarity] / 5
	hatchedPetCounts[petRarity] = hatchedPetCounts[petRarity] % 5

	upgradedCount := uint64(float64(fuseCount) * ph.FuseLuckPercentage)
	standardCount := fuseCount - upgradedCount

	switch petRarity {
	case Rare:
		hatchedPetCounts[Epic] += standardCount
		ph.allHatchedPets[Epic] += standardCount

		hatchedPetCounts[Legendary] += upgradedCount
		ph.allHatchedPets[Legendary] += upgradedCount
	case Epic:
		hatchedPetCounts[Legendary] += standardCount
		ph.allHatchedPets[Legendary] += standardCount

		hatchedPetCounts[Prodigious] += upgradedCount
		ph.allHatchedPets[Prodigious] += upgradedCount
	case Legendary:
		hatchedPetCounts[Prodigious] += standardCount
		ph.allHatchedPets[Prodigious] += standardCount

		hatchedPetCounts[Ascended] += upgradedCount
		ph.allHatchedPets[Ascended] += upgradedCount
	}
}

func (ph *PetHatcher) performProdigiousFuse(hatchedPetCounts map[string]uint64) {
	if hatchedPetCounts[Prodigious] < 3 {
		return
	}

	fuseCount := hatchedPetCounts[Prodigious] / 3
	hatchedPetCounts[Prodigious] = hatchedPetCounts[Prodigious] % 3

	upgradedCount := uint64(float64(fuseCount) * ph.FuseLuckPercentage)
	standardCount := fuseCount - upgradedCount

	hatchedPetCounts[Ascended] += standardCount
	ph.allHatchedPets[Ascended] += standardCount

	hatchedPetCounts[Mythical] += upgradedCount
	ph.allHatchedPets[Mythical] += upgradedCount
}

func (ph *PetHatcher) performAscendedFuse(hatchedPetCounts map[string]uint64) {
	if hatchedPetCounts[Ascended] < 3 {
		return
	}

	fuseCount := hatchedPetCounts[Ascended] / 3
	hatchedPetCounts[Ascended] = hatchedPetCounts[Ascended] % 3

	hatchedPetCounts[Mythical] += fuseCount
	ph.allHatchedPets[Mythical] += fuseCount
}

func (ph *PetHatcher) GetShinyScore() uint64 {
	if len(ph.allHatchedPets) == 0 {
		return 0
	}

	shinyPetValues := ph.getShinyPetValues()
	totalShinyChance := float64(1) / float64(1000)
	totalShinyChance *= ph.ShinyWallLuck * ph.ShinyAchievementLuck * ph.ExpertsLuck
	totalShinyScore := uint64(0)
	for eggType, count := range ph.allHatchedPets {
		shinyPets := totalShinyChance * float64(count)
		totalShinyScore += uint64(shinyPets) * shinyPetValues[eggType]
	}

	return totalShinyScore
}

func (ph *PetHatcher) GetTotalHatchedPetCount() uint64 {
	total := uint64(0)
	for _, subCount := range ph.allHatchedPets {
		total += subCount
	}

	return total
}

package main

import (
	"fmt"
	"math"
)

func compoundInterest(principal float64, interestRate float64, years float64) float64 {
	compoundedInterest := principal * (1 + (interestRate / 100.0))
	for i := 1; i < int(years); i++ {
		compoundedInterest += compoundedInterest * ((interestRate / 100.0) / 12.0)
	}
	return compoundedInterest
}

func main() {
	var Age int
	var Life_Expectancy int
	var Saving_perMonth float64
	var interestRate float64
	var Actual_Standard_Of_Living_net float64
	var Actual_Wealth float64
	var ThresholdPurchasingPower float64
	var principal float64
	var previousSaving_perMonth float64
	var previousReferenceAcceptable_Purchasing_Power float64
	var Found bool = false
	var Hunting bool = false
	var Disabled bool = false
	var Inconclusive bool = false
	var interestRateTrend string
	fmt.Print("Age: ")
	fmt.Scanln(&Age)
	fmt.Print("Life expectancy: ")
	fmt.Scanln(&Life_Expectancy)
	fmt.Print("Savings target per month: ")
	fmt.Scanln(&Saving_perMonth)
	fmt.Print("Observed interest rate: ")
	fmt.Scanln(&interestRate)
	fmt.Print("Actual standard of Living (NET): ")
	fmt.Scanln(&Actual_Standard_Of_Living_net)
	fmt.Print("Acceptable Threshold in Purchasing Power (-percent): ")
	fmt.Scanln(&ThresholdPurchasingPower)
	fmt.Print("Existing Actual Wealth: ")
	fmt.Scanln(&Actual_Wealth)
	fmt.Print("Will the interest rate climb? y[N] default: ")
	fmt.Scanln(&interestRateTrend)
	PREpreviousSaving_perMonth := Saving_perMonth
	PREReference_acceptable_Purchasing_Power := 100.0
	PREpreviousStandardofLiving := Actual_Standard_Of_Living_net
	for j := Age; j <= Life_Expectancy; j++ {
		SaveStartPurchasingPower := 0.0
		previousSaving_perMonth = PREpreviousSaving_perMonth
		Saving_perMonth = previousSaving_perMonth
		Reference_acceptable_Purchasing_Power := PREReference_acceptable_Purchasing_Power
		previousReferenceAcceptable_Purchasing_Power = Reference_acceptable_Purchasing_Power
		previousStandardofLiving := PREpreviousStandardofLiving
		Actual_Standard_Of_Living_net = previousStandardofLiving
		if (interestRate >= 0.25 && j > Age) && (interestRateTrend == "N" || interestRateTrend == "") {
			if interestRate != 0.25 {
				interestRate -= 0.25 // progressive return to normal. 0.25 base points a year
			}
		} else {
			if j > Age {
				interestRate += 0.25
			}
		}
		for {
			principal = Saving_perMonth
			years_to_Retirement := float64(65 - j)
			//fmt.Println("-------------------NEW SCENARIO---------PROTECTING STANDARD OF LIVING!--------CTRL+C to exit---")
			cutoff := false
			RDSP_rate := 3.5                       //this varies according to situations.
			Tax_Bracket := 1 - 0.1495              //0.0879 //0.1495 //14.95% SECOND BRACKET 2024
			Years_To_Death := Life_Expectancy - 65 //for a man in Canada
			Standard_of_Living_at_Retirement := Actual_Standard_Of_Living_net
			RRQ_Disability_Max_Monthly := 1728.0 / 2               //50% the maximum
			Canada_Pension_Plan_Disability_Max_2021 := 1046.66 / 2 //50% the maximum
			Yearly_Inflation := 1 + (interestRate / 100)           // 5% in 2024 --> anticipated to return around 2%

			//Cost_Of_living := Actual_Standard_Of_Living_net * ((years_to_Retirement + float64(Years_To_Death)) / 24) * 2 // Cost of living index doubles every 24 years
			Insurer_Inflation := 0.02 //2%

			var Raw_Capital float64
			var RDSP_GOVT_MATCHED float64
			//var LivingCapital float64

			//the Governement of Canada's share over RDSP
			for i := 1; RDSP_GOVT_MATCHED <= 100000; i++ {
				Yearly_Capital_RDSP := (principal * 12) * RDSP_rate
				RDSP_GOVT_MATCHED += Yearly_Capital_RDSP - (principal * 12)
				Raw_Capital += Yearly_Capital_RDSP
			}

			//the capital invested monthly
			Raw_Capital += (principal * 12) * years_to_Retirement

			//the compounded interests
			for i := 1; i < int(years_to_Retirement); i++ {
				principal += compoundInterest(principal, interestRate, years_to_Retirement) - principal
				RRQ_Disability_Max_Monthly = RRQ_Disability_Max_Monthly * Yearly_Inflation
				Canada_Pension_Plan_Disability_Max_2021 = Canada_Pension_Plan_Disability_Max_2021 * Yearly_Inflation
				if i == 1 {
					principal += Actual_Wealth
				}
			}

			Retirement_Capital := Raw_Capital + principal
			PreRetirement_Capital := Retirement_Capital

			if Found {
				//fmt.Printf("Savings by the age of 65: %.2f$\n", Retirement_Capital)
			}

			for i := 1; i <= Years_To_Death; i++ {
				Retirement_Capital = Retirement_Capital * (1 - (interestRate / 100))
				Retirement_Capital = Retirement_Capital - (Standard_of_Living_at_Retirement * Tax_Bracket)
				//LivingCapital = (Standard_of_Living_at_Retirement * Tax_Bracket) + ((RRQ_Disability_Max_Monthly * 12 * Yearly_Inflation) * Tax_Bracket) + ((Canada_Pension_Plan_Disability_Max_2021 * 12 * Yearly_Inflation) * Tax_Bracket)
				if Retirement_Capital < 0 {
					//fmt.Printf("Standard of life at %d: %.2f$\n", 65+i, LivingCapital)
					Life_Expectancy = 65 + i
					break
				}
				//RRQ_Disability_Max_Monthly = RRQ_Disability_Max_Monthly * Yearly_Inflation
				//Canada_Pension_Plan_Disability_Max_2021 = Canada_Pension_Plan_Disability_Max_2021 * Yearly_Inflation
			}
			if Found {
				//fmt.Printf("Savings by the age of %d: %.2f$\n", Life_Expectancy+1, Retirement_Capital)

				//fmt.Printf("Adjusted Cost of Living at %d: %.2f$\n", Life_Expectancy, Cost_Of_living)

				//fmt.Printf("Adjusted Reality at %d: %.2f percent actual standards of living\n", Life_Expectancy, ((LivingCapital / Cost_Of_living) * 100))

				//fmt.Println("Cost of living index doubles every 24 years average on 3% inflation historically.")

				//fmt.Printf("Projections using %.1f inflation leads to %.2f percent difference in standard of living at %d\n", (1-Yearly_Inflation)*-100, 100-((LivingCapital/Cost_Of_living)*100), Life_Expectancy)

				//fmt.Printf("This model suggest a saving of %.2f monthly\n", Saving_perMonth)
			}
			years_to_double_cost_of_Living := math.Log(2.0) / math.Log(1+(interestRate/100))
			//years_to_double_cost_of_Living = 72 / interestRate //tripple cost of living
			//fmt.Printf("-------------------------------------THIS PROJECTIONS AT %.2f percent INTEREST RATE-----------%.2f years to DOUBLE Cost of Living-------\n", interestRate, years_to_double_cost_of_Living)

			AdjustedRevenue := Actual_Standard_Of_Living_net
			AdjustedCostOfLiving := Actual_Standard_Of_Living_net
			RDSP_SPENT := 0.0
			DictRetirement := make(map[int][]float64)

			for i := 1; i <= (int(years_to_Retirement) + int(Years_To_Death)); i++ {
				AdjustedRevenue = AdjustedRevenue + (Actual_Standard_Of_Living_net * Insurer_Inflation)
				AdjustedCostOfLiving = AdjustedCostOfLiving + Actual_Standard_Of_Living_net/years_to_double_cost_of_Living
				if (Saving_perMonth*12 < AdjustedCostOfLiving-AdjustedRevenue) && (cutoff == false) {
					if Found {
						//fmt.Printf("|Year %d[%d]<---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|**%.2f$ SAVED MONTHLY  -->CUT-OFF<--\n", 2024+i, j+i, AdjustedRevenue, AdjustedCostOfLiving, (AdjustedCostOfLiving/AdjustedRevenue)*Saving_perMonth)
					}
					DictRetirement[2024+i] = append(DictRetirement[2024+i], AdjustedRevenue, AdjustedCostOfLiving, Saving_perMonth, 0, 0)
					cutoff = true
				} else {
					if cutoff {
						if i >= int(years_to_Retirement-7) && RDSP_SPENT < 100000 && Age+i < 65 {
							if Found {
								//fmt.Printf("|Year %d[%d]<---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|**INJECTING %.2f$ FROM RDSP RETIREMENT CAPITAL\n", 2024+i, j+i, AdjustedRevenue, AdjustedCostOfLiving, AdjustedCostOfLiving-AdjustedRevenue-(Saving_perMonth*12))
							}
							DictRetirement[2024+i] = append(DictRetirement[2024+i], AdjustedRevenue, AdjustedCostOfLiving, Saving_perMonth, AdjustedCostOfLiving-AdjustedRevenue-(Saving_perMonth*12), 0)
							RDSP_SPENT += (AdjustedCostOfLiving - AdjustedRevenue - (Saving_perMonth * 12))
						} else {
							if i >= int(years_to_Retirement) {
								Insurer_Inflation = 1 + ((float64(i) - years_to_Retirement) * 0.01) /// GOVT of Canada 1% inflation adjustment per year
								ElderlyRevenue := ((RRQ_Disability_Max_Monthly * 12 * Tax_Bracket) * Insurer_Inflation) + ((Canada_Pension_Plan_Disability_Max_2021 * 12 * Tax_Bracket) * Insurer_Inflation)
								PreRetirement_Capital = PreRetirement_Capital - (AdjustedCostOfLiving - ElderlyRevenue - (Saving_perMonth * 12))
								if PreRetirement_Capital < 0 {
									break
								}
								if Found {
									//fmt.Printf("|Year %d[%d]<---Elderly Pension ---> %.2f$<---Adjusted Cost of living--->%.2f$---|**INJECTING %.2f$ FROM RETIREMENT CAPITAL[%.2f left]\n", 2024+i, j+i, ElderlyRevenue, AdjustedCostOfLiving, AdjustedCostOfLiving-ElderlyRevenue-(Saving_perMonth*12), PreRetirement_Capital)
								}
								DictRetirement[2024+i] = append(DictRetirement[2024+i], ElderlyRevenue, AdjustedCostOfLiving, Saving_perMonth, 0, AdjustedCostOfLiving-ElderlyRevenue-(Saving_perMonth*12))
							} else {
								if Found {
									//fmt.Printf("|Year %d[%d]<---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|**%.2f percent PURCHASING POWER [insufficent funding]\n", 2024+i, j+i, AdjustedRevenue, AdjustedCostOfLiving, AdjustedRevenue/AdjustedCostOfLiving*100)
								}
								DictRetirement[2024+i] = append(DictRetirement[2024+i], AdjustedRevenue, AdjustedCostOfLiving, Saving_perMonth, 0, 0)
							}
						}
					} else {
						if Found {
							//fmt.Printf("|Year %d[%d]<---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|**%.2f$ SAVED MONTHLY\n", 2024+i, j+i, AdjustedRevenue, AdjustedCostOfLiving, (AdjustedCostOfLiving/AdjustedRevenue)*Saving_perMonth)
						}
						DictRetirement[2024+i] = append(DictRetirement[2024+i], AdjustedRevenue, AdjustedCostOfLiving, Saving_perMonth, 0, 0)
					}
				}
			}
			if Found {
				//fmt.Println("------------------------------------------CLEARER-----------------------------------")
			}
			var i int
			var StartPurchasingPower float64 = 0.0
			for i = 0; i < len(DictRetirement); i++ {
				var NewPurchasingPower float64
				Year := 2025 + i
				RetirementData := DictRetirement[Year]
				AdjustedRevenue, AdjustedCostOfLiving, Saving_perMonth, RDSP_Adjustment, Pension_Adjustments := RetirementData[0], RetirementData[1], RetirementData[2], RetirementData[3], RetirementData[4]
				Cash := AdjustedRevenue - (Saving_perMonth * 12) + RDSP_Adjustment + Pension_Adjustments
				NewPurchasingPower = Cash / AdjustedCostOfLiving * 100
				if i == 0 {
					StartPurchasingPower = NewPurchasingPower
				}
				if Found {
					//fmt.Printf("|Year %d[%d]<---Cash---> %.2f$<---Adjusted Cost of living--->%.2f$---|******%.2f percent PURCHASING POWER\n", Year, Year-2024+j, Cash, AdjustedCostOfLiving, NewPurchasingPower)
				}
				if i > 0 && (StartPurchasingPower-NewPurchasingPower) > ThresholdPurchasingPower {
					break
				}
			}
			if Found {
				//fmt.Printf("This model suggest a saving of %.2f monthly for a interest rate of %.2f\n", Saving_perMonth, interestRate)
				break
			}
			if i == len(DictRetirement) {
				if StartPurchasingPower >= Reference_acceptable_Purchasing_Power && !Hunting { // some quality of life. 75% reference.
					Found = true
				} else {
					if Actual_Standard_Of_Living_net > 0 {
						Saving_perMonth = previousSaving_perMonth
						Actual_Standard_Of_Living_net -= 1000 // at a lower standard of living.
						//fmt.Println(Actual_Standard_Of_Living_net)
					} else {
						if StartPurchasingPower >= Reference_acceptable_Purchasing_Power {
							Hunting = false
							Found = true
						} else {
							if !Hunting && !Disabled {
								Saving_perMonth = previousSaving_perMonth
								Actual_Standard_Of_Living_net = previousStandardofLiving
								Hunting = true
								//fmt.Println("Hunting...")
							} else {
								Saving_perMonth = previousSaving_perMonth
								Actual_Standard_Of_Living_net = previousStandardofLiving
								Reference_acceptable_Purchasing_Power -= 1
								//fmt.Println("Lowering...Purchasing Power to ", Reference_acceptable_Purchasing_Power)
								if Reference_acceptable_Purchasing_Power == 0 {
									interestRate -= 0.25
									Saving_perMonth = previousSaving_perMonth
									Actual_Standard_Of_Living_net = previousStandardofLiving
									Reference_acceptable_Purchasing_Power = previousReferenceAcceptable_Purchasing_Power
									if interestRate <= 0 {
										//fmt.Println("Combination error...Sorry!.")
										//os.Exit(1)
										fmt.Printf("Age:%d<-->INCONCLUSIVE\n", j)
										Inconclusive = true
										break //nothing on that year for that combination.
									}
								}
							}

						}

					}
				}

			} else {
				if Hunting {
					Saving_perMonth -= 10
					//fmt.Println(Saving_perMonth)
					if Saving_perMonth <= 0 {
						Reference_acceptable_Purchasing_Power -= 1
						Saving_perMonth = previousSaving_perMonth
						Actual_Standard_Of_Living_net = previousStandardofLiving
						//fmt.Printf("Lowering purchasing power to %.1f\n", Reference_acceptable_Purchasing_Power)
						if Reference_acceptable_Purchasing_Power == 0 {
							Reference_acceptable_Purchasing_Power = previousReferenceAcceptable_Purchasing_Power
							Hunting = false // it failed to find it at that interest...it's going back.
							Disabled = true
						}
					}
				} else {
					Saving_perMonth += 10
					//fmt.Println(Saving_perMonth)
				}
			}
			SaveStartPurchasingPower = StartPurchasingPower
		}
		if !Inconclusive {
			fmt.Printf("Age:%d<-Monthly retirement savings->%.2f<-at interest rate: %.2f--->Start Purchasing Power of %.2f percent \n", j, Saving_perMonth, interestRate, SaveStartPurchasingPower)
		} else {
			Inconclusive = false
		}
		Found = false
	}
}

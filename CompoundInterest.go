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
	fmt.Print("Age: ")
	fmt.Scanln(&Age)
	years_to_Retirement := float64(65 - Age)
	fmt.Print("Life expectancy: ")
	fmt.Scanln(&Life_Expectancy)
	for {
		fmt.Println("-------------------NEW SCENARIO---------PROTECTING STANDARD OF LIVING!-----------")
		var Saving_perMonth float64
		var interestRate float64
		var Actual_Standard_Of_Living_net float64
		cutoff := false
		fmt.Print("Savings target per month: ")
		fmt.Scanln(&Saving_perMonth)
		principal := Saving_perMonth
		fmt.Print("Observed interest rate: ")
		fmt.Scanln(&interestRate)
		RDSP_rate := 3.5                       //this varies according to situations.
		Tax_Bracket := 1 - 0.1495              //0.0879 //0.1495 //14.95% SECOND BRACKET 2024
		Years_To_Death := Life_Expectancy - 65 //for a man in Canada
		fmt.Print("Actual standard of Living (NET): ")
		fmt.Scanln(&Actual_Standard_Of_Living_net)
		Standard_of_Living_at_Retirement := Actual_Standard_Of_Living_net
		RRQ_Disability_Max_Monthly := 1728.0 / 2               //50% the maximum
		Canada_Pension_Plan_Disability_Max_2021 := 1046.66 / 2 //50% the maximum
		Yearly_Inflation := 1.045                              // 5% in 2024 --> anticipated to return around 2%

		Cost_Of_living := Actual_Standard_Of_Living_net * ((years_to_Retirement + float64(Years_To_Death)) / 24) * 2 // Cost of living index doubles every 24 years
		Insurer_Inflation := 0.02                                                                                    //2%

		Inheritance_Money := 90000.0 //or lottery money?!?

		var Raw_Capital float64
		var RDSP_GOVT_MATCHED float64
		var LivingCapital float64

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
			if i == 15 {
				principal += Inheritance_Money
			}
		}

		Retirement_Capital := Raw_Capital + principal

		fmt.Printf("Savings by the age of 65: %.2f$\n", Retirement_Capital)

		for i := 1; i <= Years_To_Death; i++ {
			Retirement_Capital = Retirement_Capital * (1 - (interestRate / 100))
			Retirement_Capital = Retirement_Capital - (Standard_of_Living_at_Retirement * Tax_Bracket)
			LivingCapital = (Standard_of_Living_at_Retirement * Tax_Bracket) + ((RRQ_Disability_Max_Monthly * 12 * Yearly_Inflation) * Tax_Bracket) + ((Canada_Pension_Plan_Disability_Max_2021 * 12 * Yearly_Inflation) * Tax_Bracket)
			if Retirement_Capital < 0 {
				//LivingCapital = Retirement_Capital
				fmt.Printf("Standard of life at %d: %.2f$\n", 65+i, LivingCapital)
				Life_Expectancy = 65 + i
				break
			}
			RRQ_Disability_Max_Monthly = RRQ_Disability_Max_Monthly * Yearly_Inflation
			Canada_Pension_Plan_Disability_Max_2021 = Canada_Pension_Plan_Disability_Max_2021 * Yearly_Inflation
		}

		fmt.Printf("Savings by the age of %d: %.2f$\n", Life_Expectancy+1, Retirement_Capital)

		fmt.Printf("Adjusted Cost of Living at %d: %.2f$\n", Life_Expectancy, Cost_Of_living)

		fmt.Printf("Adjusted Reality at %d: %.2f percent actual standards of living\n", Life_Expectancy, ((LivingCapital / Cost_Of_living) * 100))

		fmt.Println("Cost of living index doubles every 24 years average on 3% inflation historically.")

		fmt.Printf("Projections using %.1f inflation leads to %.2f percent difference in standard of living at %d\n", (1-Yearly_Inflation)*-100, 100-((LivingCapital/Cost_Of_living)*100), Life_Expectancy)
		fmt.Printf("This model suggest a saving of %.2f monthly\n", Saving_perMonth)

		years_to_double_cost_of_Living := math.Log(2.0) / math.Log(1+(interestRate/100))
		//years_to_double_cost_of_Living = 72 / interestRate //tripple cost of living
		fmt.Printf("-------------------------------------THIS PROJECTIONS AT %.2f percent INTEREST RATE-----------%.2f years to DOUBLE Cost of Living-------\n", interestRate, years_to_double_cost_of_Living)

		AdjustedRevenue := Actual_Standard_Of_Living_net
		AdjustedCostOfLiving := Actual_Standard_Of_Living_net
		RDSP_SPENT := 0.0

		for i := 1; i <= int(years_to_Retirement); i++ {
			AdjustedRevenue = AdjustedRevenue + (Actual_Standard_Of_Living_net * Insurer_Inflation)
			AdjustedCostOfLiving = AdjustedCostOfLiving + Actual_Standard_Of_Living_net/years_to_double_cost_of_Living
			if (Saving_perMonth*12 < AdjustedCostOfLiving-AdjustedRevenue) && (cutoff == false) {
				fmt.Printf("|Year %d <---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|******%.2f$ SAVED MONTHLY  -->CUT-OFF<--\n", 2024+i, AdjustedRevenue, AdjustedCostOfLiving, (AdjustedCostOfLiving/AdjustedRevenue)*Saving_perMonth)
				cutoff = true
			} else {
				if cutoff {
					if i >= int(years_to_Retirement-7) && RDSP_SPENT < 100000 {
						fmt.Printf("|Year %d <---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|******INJECTING %.2f$ FROM RDSP RETIREMENT CAPITAL\n", 2024+i, AdjustedRevenue, AdjustedCostOfLiving, AdjustedCostOfLiving-AdjustedRevenue-(Saving_perMonth*12))
						//AdjustedRevenue += (AdjustedCostOfLiving - AdjustedRevenue)
						RDSP_SPENT += AdjustedCostOfLiving - AdjustedRevenue - (Saving_perMonth * 12)
					} else {
						fmt.Printf("|Year %d <---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|******%.2f percent PURCHASING POWER\n", 2024+i, AdjustedRevenue, AdjustedCostOfLiving, AdjustedRevenue/AdjustedCostOfLiving*100)
					}
				} else {
					fmt.Printf("|Year %d <---Adjusted Revenue---> %.2f$<---Adjusted Cost of living--->%.2f$---|******%.2f$ SAVED MONTHLY\n", 2024+i, AdjustedRevenue, AdjustedCostOfLiving, (AdjustedCostOfLiving/AdjustedRevenue)*Saving_perMonth)
				}
			}
		}
	}

}
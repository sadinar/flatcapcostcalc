# Flattened Collect All Pets (CAP) Pet Cost Calculator

This is a copy of the code from https://github.com/sadinar/capcostcalc condensed down to a single file so it
can run in free go sandboxes such as https://go.dev/play/

To run it in a sandbox, copy the entire contents of the compare.go file and paste it into the input. Once
all the text is copied, run the program.

## Time-restricted Mode

Assumes the user has automatic generation unlocked and is both hatching and generating pets non-stop.
To get started, use the time-restricted constructor in the calculation package to create a new calculator instance.
Then, call the Calculate method to see what the setup will produce.

The input numbers match what is shown in-game and only the metallic chance requires calculation. Multiply the
achievement bonus shown in-game by the metallic luck value shown in-game, then use the result. Alernately, let the
application do the math as shown below.

For example:
```go
tr1 := calculation.NewTimeRestrictedCalculator(
24*1,                         // hours spent hatching
450*2*calculation.OneMillion, // gold per minute
33,                           // calcify chance
3.22,                         // generate per second
0.32,                         // egg luck
0.25,                         // fuse luck
1.74,                         // shiny wall luck
1.09,                         // shiny achievement
1.022399,                     // experts luck
0.0001*1.3,                   // metallic chance
calculation.Epic,             // type generating
calculation.Epic,             // type manually hatching
)
tr1.Calculate()
```

### Input Values

- **_Hours Spent Hatching_**: Whole number greater than 0. Calculation will hatch and generate as many eggs as possible during this many hours
- _**Gold per minute**_: Whole number greater than 0. Gold per minute as shown in-game. Calculation will subtract the cost of eggs from this to determine how many shiny wall upgrades can be purchased once all hatching is complete
- **_Calcify Chance_**: Whole number from 0 to 100 exactly as shown in-game in the automation area. Used to determine how many, on average, mythic shards are produced in the given time period
- _**Generate per second**_: Fractional number from 0.25 to 5.00 in increments of 0.01. Automation station's eggs per second as shown in-game. Will be used to determine how many eggs are generated in the given time period
- _**Egg Luck**_: Cave egg luck plus achievement egg luck as a fractional number from 0 to 0.35. Used to determine how many eggs, on average, were hatched one tier higher
- **_Fuse Luck_**: Cave fuse upgrade as a fraction from 0 to 0.25. Used to calculate how many fuses, on average, will upgrade one tier
- **_Shiny Wall Luck_**: Fractional number from 1.00 to 2.00. Shiny wall luck as shown in-game in the grotto which is used to determine how many, on average, shiny pets are hatched as well as how many additional upgrades can be bought.
- **_Shiny Achievement_**: Fractional number from 1.00 to 1.10 matching the in-game number shown for the achievement. Used to calculate how many shiny pets are hatched, on average
- **_Experts Luck_**: Cave achievement as a fractional number greater than or equal to 1.00
- **_Metallic Chance_**: Fractional number >= 0.00 exactly as shown in-game in the grove. Used to determine the odds of hatching at least one metallic in the given time period
- **_Type Generating_**: Type of egg produced by the automation station from the list of available types. Epic, Legendary, etc. See constants in pet_hatcher.go
- **_Type manually Hatching_**: Type of egg being hatched by standing in front of an egg and using it. Picked from the same list of constants as generation type

## Gold-restricted Mode

Removed from this simplified version. To use the early game calculator, see the full project at https://github.com/sadinar/capcostcalc
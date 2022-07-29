# Flattened Collect All Pets (CAP) Pet Cost Calculator

This is a copy of the code from https://github.com/sadinar/capcostcalc condensed down to a single file so it
can run in free go sandboxes such as https://go.dev/play/

To run it in a sandbox, copy the entire contents of the compare.go file and paste it into the input. Once
all the text is copied, run the program.

To change the two configurations being compared, modify the configuration constructors and the baseGold value.
For example, to decide if max +% fuse luck bonus is better or worse than max +% gold bonuses, compare
the results of these two configurations replacing the values to align with your character:

`baseGold := OneBillion // Base coins being spent before any coin bonuses are applied`

+25% fuse luck:
```
New(
    baseGold,            
    EvenCheaperPriceTable, // Price of the eggs being bought
    Epic,                  // Type of egg being hatched (assumes all hatches are the same type)
    0.29,                  // egg luck
    0.25,                  // fuse luck
    0.3,                   // achievement coin bonus
    0,                     // cave coin bonus
    0.3,                   // friend coin bonus
    true,                  // 2x coin boost
    true,                  // 1.5x coin pass
)
```

+100% gold:
```
New(
    baseGold,
    EvenCheaperPriceTable,
    Prodigious,
    0.29, // egg luck
    0,    // fuse luck
    0.3,  // achievement coin bonus
    1,    // cave coin bonus (1 == +100%)
    0.3,  // friend coin bonus
    true, // 2x coin boost
    true, // 1.5x coin boost game pass
)
```

Note that at higher fuse luck values (~8%+), changing the type being bought to a lower tier like `Epic` will result
in more mythical pets.

## Input Values

- Base gold spent
- Pet price table
- Pet type being hatched
- Egg luck
- Fuse luck
- Achievement coin bonus
- Cave coin bonus
- Friend coin bonus
- Has double coin boost
- Has 1.5x coin game pass

**_Base gold spent_**: Amount of money earned **_before_** any +% gold modifiers are applied. This is **_not_** the
gold amount shown on the display in-game and is not necessary to calculate for input. This value is used for
comparison because cave gold bonuses apply to this value, not the combined total. Any reasonably large number will do.

**_Pet price table_**: Prices of the eggs from Rare through Prodigious.

**_Pet type being hatched_**: Value can be Rare, Epic, Legendary, or Prodigious and assumes all eggs hatched are the
same type.

**_Egg luck_**: Cave egg luck bonus plus achievement egg luck bonus as a number between 0 and 1.32 in
increments of 0.01. For example, with +10% achievement egg luck and +5% cave egg luck, it would be 0.15.

_**Fuse luck**_: Cave fuse luck bonus as a number between 0 and 0.25 in increments of 0.01.

**_Achievement coin bonus_**: Coin bonus from achievements as a number between 0 and 0.35 in increments of 0.05.

**_Cave coin bonus_**: Coin bonus from the cave as a number between 0 and 1 in increments of 0.01.

**_Friend coin bonus_**: Coin bonus from playing with friends. Values should be 0, 0.1, 0.2, or 0.3.

**_Has double coin boost_**: True if 2x coin boost is active and false if it is not.

**_Has 1.5x coin game pass_**: True if user purchased the permanent 1.5x coin boost and false otherwise.
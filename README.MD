# wordle game
**Validation:**
- Must a word in Dictionary
- Must include English alphabet only and case-insensitive

**Playing rule:**  
- **Hit** - The input letter in the correct spot  
- **Present** - The input letter in the incorrect sport   
- **Miss** - The input letter is not included  

**Playing Rule:**
- **Win**: Guess within the maximum round
- **Loss**: Guess without the maximum round
- **GameOver** / **Win** - show the **answer**

---
### How to start the game
1. Clone the project
```
git clone https://github.com/RyanTokManMokMTM/wordle-game.git
```

2. Download the package
```
go mod download
```

3. (Optinal) Update the config setting
```
core/etc
```

4. Go into `core` folder and run the command
```
go run wordlegame.go
```
5. Enjoy!

---

### File structure
```
.
├── client     # client interface & implementation       
├── common     # Tools and utilities                 
├── etc        # Config file  
├── game       # Game Start entry point         
├── internal   # Core Logic             
└── README.md  # Doc
```

### Core implementation
#### guessingWordChecking  logic
> Using a array to indicated there is/are alphabets included in the **guessing word**
- Array size: [a-z][A-Z] , total 52
- Index calculation: according to the ascii code table, `'a'-'a'=0` ,so `arr[0]` =`'a'`, `arr[1]`= `'b'`, `arr[2]` = `'c'` ,`a[26]`=`'A'`, `a[27]`=`'B'`,etc

```
size := len(input)
for i := 0; i < size; i++ {
    //Get the index in couter array by letter(input(i))
    index := utils.LetterIndexOf(rune(in[i]))
    if counter[index] == 0 {
        //The counter do not include this letter
        continue
    }
    
    //The counter include this letter
    if in[i] == guessingWord[i] {
        counter[index] -= 1
        temp[i] = '0'
    }
}
```
---
#### hostCheatingGameChecking logic
1. Comparing user input and the candidate list - `processGameResultAndScore`
2. According to the result(How many `HIT`/`PRESENT`/`MISS`) from `processGameResultAndScore`, and finding current lowest score of the candidate list
3. According to the lowest score, filtering the score is larger than the lowest score and update the candidate list
4. Due to coexisting result issue, filtering results which is/are not the same result.
   > For example: `_?___`, `__?__`. Both are 1 present score, but they are not the same, so we need to select either one of them
5. Continue this process until only one candidate is left.
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
> git clone https://github.com/RyanTokManMokMTM/wordle-game.git

2. Download the package
> go mod download

3. Go into `core` folder and run the command
> go run wordlegame.go

4. Enjoy!

---

### File structure
├── client     # client interface & implementation       
├── common     # Tools and utilities                 
├── etc        # Config file  
├── game       # Game Start entry point         
├── internal   # Core Logic             
└── README.md  # Doc

### Core implementation
#### guessingWordChecking   
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
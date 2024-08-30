package game

//
//func Start(c gameClient.IClient) {
//	reader := bufio.NewReader(os.Stdin)
//	for {
//		fmt.Println("Select a mode : \n1:start a game\n2:exit")
//		selectedMode, _ := reader.ReadString('\n')
//		selectedMode = strings.Replace(selectedMode, "\n", "", -1)
//		flag, err := strconv.Atoi(selectedMode)
//		if err != nil {
//			fmt.Println("Please input an number")
//			continue
//		}
//		switch flag {
//		case 1:
//			logic.GameStart(c)
//		case 2:
//			fmt.Println("Thank you for playing")
//			return
//		default:
//			fmt.Println("Mode not supported")
//			continue
//		}
//	}
//
//}

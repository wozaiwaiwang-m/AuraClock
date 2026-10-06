package magpie

import (
	"fmt"
)

func AskBool(str ...string) bool {
	// func AskBool(str string) bool {
	var as string
	for _, ss := range str {
		fmt.Printf(ss)
	}
	// fmt.Printf(str)
	for {
		fmt.Println("\t\t", "输入（y或n）：")
		fmt.Scan(&as)
		if as == "y" {
			return true
		} else if as == "n" {
			return false
		} else {
			fmt.Print("请重新输入")
		}
	}
}

func AskInt(str ...string) int {
	// func AskInt(str string) int {
	var as int
	for _, ss := range str {
		fmt.Printf(ss)
	}
	// fmt.Printf(str)
	for {
		fmt.Println("\t\t", "输入一个整数：")
		fmt.Scan(&as)
		//待添加判断用户输入的值的类型的代码
		// if as <= 0 {
		// 	fmt.Print("请重新输入")
		// } else {
		// 	return as
		// }
		return as
	}
}

func AskFloat(str ...string) float64 {
	// func AskFloat(str string) float64 {
	var as float64
	for _, ss := range str {
		fmt.Printf(ss)
	}
	// fmt.Printf(str)
	for {
		fmt.Println("\t\t", "输入一个浮点数：")
		fmt.Scan(&as)
		//待添加判断用户输入的值的类型的代码
		// if as <= 0.0 {
		// 	fmt.Print("请重新输入")
		// } else {
		// 	return as
		// }
		return as
	}
}

func AskString(str ...string) string {
	// func AskString(str string) string {
	var as string
	for _, ss := range str {
		fmt.Printf(ss)
	}
	// fmt.Printf(str)
	for {
		fmt.Println("\t\t", "输入一个字符串：")
		fmt.Scan(&as)
		//待添加判断用户输入的值的类型的代码
		// if as <= 0.0 {
		// 	fmt.Print("请重新输入")
		// } else {
		// 	return as
		// }
		return as
	}
}

func AskArray(sli []string, str ...string) (int, string) {
	// var array [int]string
	// var serial []string
	// var as int
	// for _, aa := range sli {
	// 	serial = append(serial, aa)
	// }
	var as int
	for _, ss := range str {
		fmt.Println(ss)
	}
	for aa, bb := range sli {
		fmt.Print("| [", aa+1, "：", bb, "] ")
	}
	//sli = append([]string{""}, sli...)
	// fmt.Printf(str)
	fmt.Println("\t\t", "输入序号：")
	for {
		fmt.Scan(&as)

		if as == 0 || as > len(sli) {
			fmt.Println("请重新输入") //错误处理
		} else {
			return as, sli[as-1]
		}
	}

}

// func RepeatAsk(aa string, ask func(string)) {
// 	bb := AskInt("需要输入几次？")
// 	for range bb {
// 		return ask()
// 	}
// }

// 下面的代码是AI写的，我还不懂，但会用

func RepeatedlyAsk[T any](askFunc func(...string) T, prompts ...string) ([]T, error) {
	var n int
	fmt.Println("需要重复输入几次？")
	fmt.Scan(&n)
	if n <= 0 {
		return nil, fmt.Errorf("重复次数必须为正整数，当前输入: %d", n)
	}

	results := make([]T, 0, n)
	for i := 0; i < n; i++ {
		if len(prompts) > 0 {
			fmt.Printf("第%d次询问:\n", i+1)
		}
		result := askFunc(prompts...)
		results = append(results, result)
	}
	return results, nil
}

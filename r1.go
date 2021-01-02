package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	gernerate    []string
	userProcess  string
	userterminal string
	userInput    string
	//force out
	errout bool

	//PredictS has 1 rule.
	//PredictS []rule
	//PredictC has 2 rules.
	//PredictC []rule
	//PredictA has 2 rules.
	//PredictA []rule
	//PredictB has 2 rules.
	//PredictB []rule
	//PredictQ has 2 rules.
	//PredictQ []rule

	//PredictProg has 1 rule.
	PredictProg []rule
	//PredictDcls has 2 rules.
	PredictDcls []rule
	//PredictDcl has 2 rules.
	PredictDcl []rule
	//PredictStmts has 2 rules.
	PredictStmts []rule
	//PredictStmt is in case 2 and has 2 rules.
	PredictStmt []rule
	//PredictExpr is in case 2 and has 3 rules.
	PredictExpr []rule
	//PredictVal is in case 2 and has 3 rules.
	PredictVal []rule
)

type rule struct {
	num     uint
	predict []string
	context string
}

func init() {
	PredictProg = []rule{
		rule{
			num:     1,
			predict: []string{"floatdcl", "intdcl", "id", "print", "$"},
			context: "Dcls Stmts $",
		},
	}
	PredictDcls = []rule{
		rule{
			num:     2,
			predict: []string{"floatdcl", "intdcl"},
			context: "Dcl Dcls",
		},
		rule{
			num:     3,
			predict: []string{"id", "print", "$"},
			context: "L",
		},
	}
	PredictDcl = []rule{
		rule{
			num:     4,
			predict: []string{"floatdcl"},
			context: "floatdcl id",
		},
		rule{
			num:     5,
			predict: []string{"intdcl"},
			context: "intdcl id",
		},
	}
	PredictStmts = []rule{
		rule{
			num:     6,
			predict: []string{"id", "print"},
			context: "Stmt Stmts",
		},
		rule{
			num:     7,
			predict: []string{"$"},
			context: "L",
		},
	}
	PredictStmt = []rule{
		rule{
			num:     8,
			predict: []string{"id"},
			context: "id assign Val Expr",
		},
		rule{
			num:     9,
			predict: []string{"print"},
			context: "print id",
		},
	}
	PredictExpr = []rule{
		rule{
			num:     10,
			predict: []string{"plus"},
			context: "plus Val Expr",
		},
		rule{
			num:     11,
			predict: []string{"minus"},
			context: "minus Val Expr",
		},
		rule{
			num:     12,
			predict: []string{"id", "print", "$"},
			context: "L",
		},
	}
	PredictVal = []rule{
		rule{
			num:     13,
			predict: []string{"id"},
			context: "id",
		},
		rule{
			num:     14,
			predict: []string{"inum"},
			context: "inum",
		},
		rule{
			num:     15,
			predict: []string{"fnum"},
			context: "L",
		},
	}
}

func main() {
	var (
		way   string
		input string
	)
	for way != "exit" {
		fmt.Println("input,file,exit")
		fmt.Scanln(&way)
		switch {
		case way == "input":
			fmt.Print("input:")
			in := bufio.NewReader(os.Stdin)
			input, _ = in.ReadString('\n')
			restart(input)
		case way == "file":
			fmt.Print("filePath:")
			var path string
			fmt.Scanln(&path)
			input, err := getFileContext(path)
			if err != nil {
				fmt.Println(err)
			}
			restart(input)
		case way == "exit":
			fmt.Println("Exit!")
		default:
			continue
		}
	}

}

func getFileContext(filePath string) (string, error) {
	filePath = strings.Trim(filePath, " \n\t")
	context, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(context), nil
}

func reset() {
	userInput = ""
	userProcess = ""
	errout = false
}

func restart(newUserInput string) {
	reset()
	newUserInput = strings.Trim(newUserInput, " \n\t\r")
	userInput = newUserInput
	userProcess = newUserInput
	gernerate = strings.Split(userInput, " ")
	userterminal = gernerate[0]

	Prog()
}

func Prog() {
	switch {
	case contains(PredictProg[0].predict, userterminal):
		fmt.Printf("%d ", PredictProg[0].num)
		Dcls()
		Stmts()
		match("$")
		if !errout {
			fmt.Println("Accept")
		}
	default:
		fmt.Printf("Error(Prog vs. %s)\n", userterminal)
		errout = true
	}

}

func Dcl() {
	if errout {
		return
	}
	switch {
	case contains(PredictDcl[0].predict, userterminal):
		fmt.Printf("%d ", PredictDcl[0].num)
		match("floatdcl")
		match("id")
	case contains(PredictDcl[1].predict, userterminal):
		fmt.Printf("%d ", PredictDcl[1].num)
		match("intdcl")
		match("id")
	default:
		fmt.Printf("Error(Dcl vs. %s)\n", userterminal)
		errout = true
	}
}

func Dcls() {
	if errout {
		return
	}
	switch {
	case contains(PredictDcls[0].predict, userterminal):
		fmt.Printf("%d ", PredictDcls[0].num)
		Dcl()
		Dcls()
	case contains(PredictDcls[1].predict, userterminal):
		fmt.Printf("%d ", PredictDcls[1].num)
		return
	default:
		fmt.Printf("Error(Dcls vs. %s)\n", userterminal)
		errout = true
	}
}

func Stmts() {
	if errout {
		return
	}
	switch {
	case contains(PredictStmts[0].predict, userterminal):
		fmt.Printf("%d ", PredictStmts[0].num)
		Stmt()
		Stmts()
	case contains(PredictStmts[1].predict, userterminal):
		fmt.Printf("%d ", PredictStmts[1].num)
		return
	default:
		fmt.Printf("Error(Stmts vs. %s)\n", userterminal)
		errout = true
	}
}

func Stmt() {
	if errout {
		fmt.Println("out")
		return
	}
	switch {
	case contains(PredictStmt[0].predict, userterminal):
		fmt.Printf("%d ", PredictStmt[0].num)
		match("id")
		match("assign")
		Val()
		Expr()
	case contains(PredictStmt[1].predict, userterminal):
		fmt.Printf("%d ", PredictStmt[1].num)
		match("print")
		match("id")
	default:
		fmt.Printf("Error(Stmt vs. %s)\n", userterminal)
		errout = true
	}
}

func Val() {
	if errout {
		return
	}
	switch {
	case contains(PredictVal[0].predict, userterminal):
		fmt.Printf("%d ", PredictVal[0].num)
		match("id")
	case contains(PredictVal[1].predict, userterminal):
		fmt.Printf("%d ", PredictVal[1].num)
		match("inum")
	case contains(PredictVal[2].predict, userterminal):
		fmt.Printf("%d ", PredictVal[2].num)
		match("fnum")
	default:
		fmt.Printf("Error(Val vs. %s)\n", userterminal)
		errout = true
	}
}

func Expr() {
	if errout {
		return
	}
	switch {
	case contains(PredictExpr[0].predict, userterminal):
		fmt.Printf("%d ", PredictExpr[0].num)
		match("plus")
		Val()
		Expr()
	case contains(PredictExpr[1].predict, userterminal):
		fmt.Printf("%d ", PredictExpr[1].num)
		match("minus")
		Val()
		Expr()
	case contains(PredictExpr[2].predict, userterminal):
		fmt.Printf("%d ", PredictExpr[2].num)
		return
	default:
		fmt.Printf("Error(Expr vs. %s)\n", userterminal)
		errout = true
	}
}

func match(terminal string) {
	if errout {
		return
	}
	if terminal == userterminal {
		if len(gernerate) == 1 {
			return
		}
		gernerate = gernerate[1:]
		userterminal = gernerate[0]
		return
	}
	fmt.Printf("Error(Expected %s)\n", terminal)
	errout = true
	return
}

// IsPredict return search result
// Is peek in predict?
func contains(set []string, peek string) bool {
	for _, terminal := range set {
		if peek == terminal {
			return true
		}
	}
	return false
}

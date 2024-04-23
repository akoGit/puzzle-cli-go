package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"strconv"
	"log"
	"math/rand"
	"time"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	_ "github.com/mattn/go-sqlite3"
)



func fenToBoard(fen string) [][]string {
	rows := strings.Split(fen, "/")
	board := make([][]string, 8)

	for i, row := range rows {
		board[i] = make([]string, 8)
		col := 0
		for _, char := range row {
			if col >= 8 {
				break
			}
			switch char {
			case 'R':
				board[i][col] = "♜"
			case 'N':
				board[i][col] = "♞"
			case 'B':
				board[i][col] = "♝"
			case 'Q':
				board[i][col] = "♛"
			case 'K':
				board[i][col] = "♚"
			case 'P':
				board[i][col] = "♟"
			case 'r':
				board[i][col] = "♖"
			case 'n':
				board[i][col] = "♘"
			case 'b':
				board[i][col] = "♗"
			case 'q':
				board[i][col] = "♕"
			case 'k':
				board[i][col] = "♔"
			case 'p':
				board[i][col] = "♙"
			default:
				count, _ := strconv.Atoi(string(char))
				for j := 0; j < count; j++ {
					board[i][col] = " "
					col++
				}
				continue
			}
			col++
		}
	}

	return board
}


func GetRandomFEN(databasePath string) (string, string, error) {

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return "","", err
	}
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	rows, err := db.Query("SELECT FEN, Moves FROM mateIn1 ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	var FEN, moves string

	if rows.Next() {
		err := rows.Scan(&FEN, &moves)
		if err != nil {
			return "","", err
		}
	}

	return FEN, moves, nil

}

func main() {

	re := lipgloss.NewRenderer(os.Stdout)

	labelStyle := re.NewStyle().Bold(true).Foreground(lipgloss.Color("#7a9db2"))
	
	fen, moves, err := GetRandomFEN("mateIn1.db")

	fmt.Println(fen + "\n", moves)

	if err != nil {
		log.Fatal(err)
	}

	board := fenToBoard(fen)
	


	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(board...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
	})



	ranks := labelStyle.Render(strings.Join([]string{" A", "B", "C", "D", "E", "F", "G", "H  "}, "   "))
	files := labelStyle.Render(strings.Join([]string{" 8", "7", "6", "5", "4", "3", "2", "1 "}, "\n\n "))

	fmt.Println(lipgloss.JoinVertical(lipgloss.Right,
			lipgloss.JoinHorizontal(lipgloss.Center, files, t.Render()),
	ranks) + "\n")
	 
}

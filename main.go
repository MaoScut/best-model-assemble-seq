package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaoScut/best-model-assemble-seq/structs"
)

func main() {
	unitDataDir := "/home/madp/workspace/best-model-assemble-seq/data/unit-data"
	boardsDataPath := "/home/madp/workspace/best-model-assemble-seq/data/boards.json"
	units, err := readUnitData(unitDataDir)
	if err != nil {
		log.Fatal("readUnitData", err)
	}
	boards, err := readBoardsData(boardsDataPath)
	if err != nil {
		log.Fatal("readBoardsData", err)
	}
	input := Input{
		Boards: boards,
		Units:  units,
	}
	allSolutions := getAllSolutions(input.Units)
	solutionsWithScore := []Solution{}
	for _, s := range allSolutions {
		sws, err := calculateSolutionScore(s, boards)
		if err != nil {
			log.Fatal("calculateSolutionScore", err)
		}
		solutionsWithScore = append(solutionsWithScore, sws)
	}
	maxScoreSolutionIndex := 0
	var maxScore float32
	for index, s := range solutionsWithScore {
		if s.Score > maxScore {
			maxScoreSolutionIndex = index
			maxScore = s.Score
		}
	}
	// b, _ := json.Marshal(solutionsWithScore[maxScoreSolutionIndex])
	// fmt.Printf("%s\n", b)
	// return
	fmt.Println(solutionsWithScore[maxScoreSolutionIndex].String())
}

type Input struct {
	Boards []structs.Board
	Units  []structs.Unit
}

type Solution struct {
	Score float32
	Seq   []SolutionItem
}

type SolutionItem struct {
	Unit            structs.Unit
	BoardSituations []BoardSituationItem
}

type BoardSituationItem struct {
	BoardName  string
	TotalCount int
	RestCount  int
	UseRate    float32
}

func (s *Solution) String() string {
	partNames := []string{}
	for _, item := range s.Seq {
		partNames = append(partNames, item.Unit.Name)
	}
	str := "assemble seq overview\n"
	str += strings.Join(partNames, "->")
	str += "detail\n\n"
	for _, item := range s.Seq {
		str += fmt.Sprintf("after finish %s\n", item.Unit.Name)
		str += fmt.Sprintf("the board situation is\n")
		fullUsedCount := 0
		for idx, sit := range item.BoardSituations {
			str += fmt.Sprintf("%s:%.2f", sit.BoardName, sit.UseRate)
			if idx != len(item.BoardSituations)-1 {
				str += ", "
			}
			if sit.UseRate >= 0.9 {
				fullUsedCount++
			}
		}
		str += "\n"
		str += fmt.Sprintf("full used board count: %d/%d\n", fullUsedCount, len(item.BoardSituations))
		str += "\n\n"
	}
	return str
}

func readUnitData(dir string) (arr []structs.Unit, err error) {
	dirObjs, err := os.ReadDir(dir)
	if err != nil {
		err = fmt.Errorf("os.ReadDir: %w", err)
		return
	}
	for _, obj := range dirObjs {
		unit := structs.Unit{}
		b, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, obj.Name()))
		if err != nil {
			err = fmt.Errorf("os.ReadFile: %w", err)
			return nil, err
		}
		unit.Name = obj.Name()
		str := string(b)
		partsArr := strings.Split(str, "\n")
		for _, p := range partsArr {
			if p != "" {
				reg := regexp.MustCompile(`([a-z]+)(\d+)`)
				matchResult := reg.FindStringSubmatch(p)
				if len(matchResult) == 0 {
					err = fmt.Errorf("invalid input, filename: %s, str: %s", obj.Name(), p)
					return nil, err
				}
				id, err := strconv.Atoi(matchResult[2])
				if err != nil {
					err = fmt.Errorf("strconv: %w", err)
					return nil, err
				}
				unit.Parts = append(unit.Parts, structs.Part{
					Board: matchResult[1],
					Id:    id,
				})
			}
		}
		arr = append(arr, unit)
	}
	return
}

func readBoardsData(path string) (arr []structs.Board, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("os.ReadFile: %w", err)
		return
	}
	arr = []structs.Board{}
	err = json.Unmarshal(b, &arr)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return
	}
	return
}

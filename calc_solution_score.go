package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/MaoScut/best-model-assemble-seq/structs"
)

func calculateSolutionScore(s Solution, boards []structs.Board) (newS Solution, err error) {
	boardsMap := map[string]*BoardsMapValue{}
	for _, b := range boards {
		boardsMap[b.Name] = &BoardsMapValue{
			Max:           b.PartsCount,
			PartSituation: map[int]bool{},
		}
		for i := 0; i < b.PartsCount; i++ {
			boardsMap[b.Name].PartSituation[i+1] = true
		}
	}
	newS = Solution{}
	for index, seqItem := range s.Seq {
		for _, part := range seqItem.Unit.Parts {
			board, ok := boardsMap[part.Board]
			if !ok {
				boardsMapBytes, _ := json.Marshal(boardsMap)
				err = fmt.Errorf("board: %s not found, boardsMap: %s", part.Board, boardsMapBytes)
				return
			}
			delete(board.PartSituation, part.Id)
		}
		newSeq := SolutionItem{
			Unit:            seqItem.Unit,
			BoardSituations: []BoardSituationItem{},
		}
		for name, v := range boardsMap {
			var useRate float32
			restPartsCount := len(v.PartSituation)
			if restPartsCount == 0 {
				useRate = 1
			} else if restPartsCount < 5 {
				useRate = 0.9
			} else {
				useRate = 1 - float32(restPartsCount)/float32(v.Max)
			}
			boardScore := float32(len(s.Seq)-index) * useRate
			newS.Score += boardScore
			newSeq.BoardSituations = append(newSeq.BoardSituations, BoardSituationItem{
				BoardName:  name,
				TotalCount: v.Max,
				RestCount:  restPartsCount,
				UseRate:    useRate,
			})
		}
		sort.Slice(newSeq.BoardSituations, func(i, j int) bool {
			return newSeq.BoardSituations[i].BoardName < newSeq.BoardSituations[j].BoardName
		})
		newS.Seq = append(newS.Seq, newSeq)
	}
	return
}

type BoardsMapValue struct {
	Max           int
	PartSituation map[int]bool
}

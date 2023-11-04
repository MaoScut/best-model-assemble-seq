package main

import "github.com/MaoScut/best-model-assemble-seq/structs"

func getAllSolutions(units []structs.Unit) (ss []Solution) {
	result := &Result{}
	_getAll([]structs.Unit{}, units, result)
	for _, s := range result.arr {
		tmp := Solution{
			Score: 0,
			Seq:   []SolutionItem{},
		}
		for _, u := range s {
			tmp.Seq = append(tmp.Seq, SolutionItem{
				Unit:            u,
				BoardSituations: []BoardSituationItem{},
			})
		}
		ss = append(ss, tmp)
	}
	return
}

func _getAll(pre, rest []structs.Unit, result *Result) {
	if len(rest) == 0 {
		result.arr = append(result.arr, pre)
		return
	}
	for index, u := range rest {
		newRest := remove(rest, index)
		newPre := []structs.Unit{}
		newPre = append(newPre, pre...)
		newPre = append(newPre, u)
		_getAll(newPre, newRest, result)
	}
}

type Result struct {
	arr [][]structs.Unit
}

func remove(old []structs.Unit, index int) (new []structs.Unit) {
	for i := range old {
		if i != index {
			new = append(new, old[i])
		}
	}
	return
}
